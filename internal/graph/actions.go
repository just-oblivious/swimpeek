package graph

import (
	"fmt"
	"slices"
	"strings"
	"swimpeek/pkg/laneclient"
)

// createActionNodes creates nodes for workflow actions.
func createActionNodes(warns *Warnings, graph *Graph, actions map[string]laneclient.PlaybookAction) (map[string]*Node, error) {
	actNodes := make(map[string]*Node, len(actions))

	for actId, action := range actions {
		actNode, err := createActionNode(warns, graph, action, actId)
		if err != nil {
			warns.Add(fmt.Errorf("failed to create action node for %s: %w", actId, err))
		}
		actNodes[actId] = actNode
	}

	return actNodes, nil
}

// createActionNode creates a node for a workflow action and links it to resources referenced outside the workflow.
func createActionNode(warns *Warnings, graph *Graph, action laneclient.PlaybookAction, actId string) (*Node, error) {
	// Actions with no special handling
	simpleTypes := map[string]NodeType{
		"transformation":  TransformationActionNode,
		"python":          PythonActionNode,
		"createVariables": CreateVarsActionNode,
		"updateVariables": UpdateVarsActionNode,
		"http":            HTTPActionNode,
		"conditional":     ConditionalActionNode,
		"parallelGroup":   ParallelActionNode,
	}
	if nodeType, exists := simpleTypes[action.Type]; exists {
		return newNode(newMeta(actId, nodeType, action.Title, action.Description)), nil
	}

	switch action.Type {
	case "loop":
		// Determine the type of loop (while or for each)
		switch action.Loop.Type {
		case "while":
			return newNode(newMeta(actId, WhileLoopAction, action.Title, action.Description)), nil
		case "for":
			return newNode(newMeta(actId, ForEachLoopAction, action.Title, action.Description)), nil
		default:
			return nil, fmt.Errorf("unknown loop type %s", action.Loop.Type)
		}

	case "emitEvent":
		// Flow event emit action
		sensorName, err := reflectEmitAction(action.Inputs)
		if err != nil {
			return nil, fmt.Errorf("failed to get sensor name for emitEvent: %w", err)
		}
		emitNode := newNode(newMeta(actId, EmitEventActionNode, action.Title, action.Description))
		sensNode, exists := graph.Resources.TriggersById[sensorName]
		if !exists {
			warns.Add(fmt.Errorf("emitEvent action %s references unknown sensor %s", actId, sensorName))
			return sensNode, nil
		}
		newEdge(sensNode, emitNode, EmittedByEdge, nil)
		return emitNode, nil

	case "connector":
		// Component reference
		componentRef, isComponent := strings.CutPrefix(action.Action, "$playbook.component_")

		if isComponent {
			conActionNode := newNode(newMeta(actId, ComponentActionNode, action.Title, action.Description))
			compId := strings.TrimSuffix(componentRef, "_playbook")
			compNode, exists := graph.Resources.ComponentsById[compId]
			if !exists {
				warns.Add(fmt.Errorf("connector action %s references unknown component %s", actId, compId))
				return conActionNode, nil
			}
			newEdge(compNode, conActionNode, CalledByEdge, nil)
			return conActionNode, nil
		}

		// Connector reference
		conActionNode := newNode(newMeta(actId, ConnectorActionNode, action.Title, action.Description))
		connectorRef, _, _ := strings.Cut(action.Action, ".")
		if connectorRef == "" {
			warns.Add(fmt.Errorf("connector action %s has no connector reference", actId))
			return conActionNode, nil
		}
		connectorNode, exists := graph.Resources.ConnectorsById[connectorRef]
		if !exists {
			warns.Add(fmt.Errorf("connector action %s references unknown connector %s", actId, connectorRef))
			return conActionNode, nil
		}
		newEdge(connectorNode, conActionNode, CalledByEdge, nil)
		return conActionNode, nil

	case "recordAction":
		// Determine the subtype of the record action (create, update, or search)
		recordActionType := RecordActionNode
		switch action.RecordActionType {
		case "create":
			recordActionType = RecordCreateActionNode
		case "patch":
			recordActionType = RecordUpdateActionNode
		case "search":
			recordActionType = RecordSearchActionNode
		case "delete":
			recordActionType = RecordDeleteActionNode
		case "upsert":
			recordActionType = RecordUpsertActionNode
		default:
			warns.Add(fmt.Errorf("recordAction %s has unknown recordActionType %s", actId, action.RecordActionType))
		}

		recNode := newNode(newMeta(actId, recordActionType, action.Title, action.Description))

		// Lookup the referenced application
		appId, err := reflectRecordActionAppId(action.Inputs)
		if err != nil {
			warns.Add(fmt.Errorf("recordAction %s reference error: %w", actId, err))
			return recNode, nil
		}
		appNode, exists := graph.Resources.AppsById[appId]
		if !exists {
			warns.Add(fmt.Errorf("recordAction %s references unknown application %s", actId, appId))
			return recNode, nil
		}
		newEdge(appNode, recNode, AccessedByEdge, nil)
		return recNode, nil
	}

	return nil, fmt.Errorf("unknown action type %s", action.Type)
}

// chainActions chains actions starting from the given entrypoints.
func chainActions(warns *Warnings, graph *Graph, source *Node, actions map[string]laneclient.PlaybookAction, entryPoints ...string) error {
	actNodes, err := createActionNodes(warns, graph, actions)
	if err != nil {
		return fmt.Errorf("failed to create action nodes for %s: %w", source.Meta.Id, err)
	}

	// Start walking from the entry points.
	visited := make(map[string]bool)

	for _, entryPoint := range entryPoints {
		entryNode, exists := actNodes[entryPoint]
		if !exists {
			warns.Add(fmt.Errorf("entry point %s not found in action nodes for %s", entryPoint, source.Meta.Id))
			continue
		}
		newEdge(source, entryNode, EntrypointEdge, nil)
		linkActions(warns, graph, entryNode, actions, actNodes, visited)
	}

	// Link unreachable actions to the source
	for actId, actNode := range actNodes {
		if _, exists := visited[actId]; !exists {
			warns.Add(fmt.Errorf("unreachable action %s in %s: %s (%s)", actId, source.Meta.Type, source.Meta.Id, source.Meta.Label))
			newEdge(source, actNode, UnreachableEdge, nil)
		}
	}

	return nil

}

// linkActions recursively links actions.
func linkActions(warns *Warnings, graph *Graph, source *Node, actions map[string]laneclient.PlaybookAction, actNodes map[string]*Node, visited map[string]bool) {
	// Find the action node for the source
	action, exists := actions[source.Meta.Id]
	if !exists {
		warns.Add(fmt.Errorf("action %s not found", source.Meta.Id))
		return
	}

	linkFn := func(edgeType EdgeType, nextActionIds ...string) {
		slices.Sort(nextActionIds)
		for _, nextId := range slices.Compact(nextActionIds) {
			nextNode, exists := actNodes[nextId]
			if !exists {
				warns.Add(fmt.Errorf("next action %s not found in action nodes for %s", nextId, source.Meta.Id))
				continue
			}

			// Create an edge from the source to the next action and continue linking
			newEdge(source, nextNode, edgeType, nil)

			// Skip nodes that have already been visited
			if visited[nextNode.Meta.Id] {
				continue
			}

			linkActions(warns, graph, nextNode, actions, actNodes, visited)
		}

	}

	// Loop actions and parallels: link the inner action chain
	if len(action.Entrypoints) > 0 {
		if err := chainActions(warns, graph, source, action.Actions, action.Entrypoints...); err != nil {
			warns.Add(fmt.Errorf("failed to chain inner actions for %s: %w", source.Meta.Id, err))
		}
	}

	// Conditional actions: link the branches
	if source.Meta.Type == ConditionalActionNode {
		for _, condition := range action.Conditions {
			if condition.Action == "" {
				continue
			}
			linkFn(IfEdge, condition.Action)
		}
		if action.Else != "" {
			linkFn(ElseEdge, action.Else)
		}

		// Do not follow the continuation flows for conditional actions, only the branches defined under "actions" and "else" should be followed.
		// Continuation flows for a conditional action sometimes contain phantom or duplicate actions that cannot be reached
		visited[source.Meta.Id] = true
		return
	}

	// Traverse continuation flows (on-success, on-failure, on-complete)
	for continuationType, actionMaps := range map[EdgeType][]map[string]any{
		OnSuccessEdge:  action.OnSuccess,
		OnFailureEdge:  action.OnFailure,
		OnCompleteEdge: action.OnComplete,
	} {
		nextActionIds := make([]string, 0, len(actionMaps))
		for _, actionMap := range actionMaps {
			for id := range actionMap {
				nextActionIds = append(nextActionIds, id)
			}
		}

		linkFn(continuationType, nextActionIds...)
	}

	visited[source.Meta.Id] = true
}
