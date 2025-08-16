package graph

import (
	"fmt"
	"swimpeek/internal/lanedump"
)

type NodeType string

const (
	// Top-level resources
	ApplicationNode NodeType = "application"
	ComponentNode   NodeType = "component"
	PlaybookNode    NodeType = "playbook"
	ConnectorNode   NodeType = "connector"
	WorkflowNode    NodeType = "workflow"

	// Trigger events
	FlowEventNode      NodeType = "flow_event"
	WebhookNode        NodeType = "webhook"
	PlaybookButtonNode NodeType = "playbook_button"
	RecordEventNode    NodeType = "record_event"
	CronEventNode      NodeType = "cron_event"

	// Actions
	RecordActionNode         NodeType = "record_action"
	RecordSearchActionNode   NodeType = "record_search_action"
	RecordUpdateActionNode   NodeType = "record_update_action"
	RecordCreateActionNode   NodeType = "record_create_action"
	ConnectorActionNode      NodeType = "connector_action"
	ComponentActionNode      NodeType = "component_action"
	EmitEventActionNode      NodeType = "emit_event_action"
	TransformationActionNode NodeType = "transformation_action"
	WhileLoopAction          NodeType = "while_loop_action"
	ForEachLoopAction        NodeType = "for_each_loop_action"
	PythonActionNode         NodeType = "python_action"
	CreateVarsActionNode     NodeType = "create_vars_action"
	UpdateVarsActionNode     NodeType = "update_vars_action"
	ParallelActionNode       NodeType = "parallel_action"
	ConditionalActionNode    NodeType = "condition_action"
	HTTPActionNode           NodeType = "http_action"
)

type EdgeType string

const (
	EmittedByEdge        EdgeType = "emitted_by"
	CalledByEdge         EdgeType = "called_by"
	AccessedByEdge       EdgeType = "accessed_by"
	TriggersWorkflowEdge EdgeType = "triggers_workflow"
	HasEventEdge         EdgeType = "has_event"
	HasActionEdge        EdgeType = "has_action"
	WorkflowEdge         EdgeType = "workflow"
	EntrypointEdge       EdgeType = "entrypoint"
	UnreachableEdge      EdgeType = "unreachable"
	OnSuccessEdge        EdgeType = "on_success"
	OnFailureEdge        EdgeType = "on_failure"
	OnCompleteEdge       EdgeType = "on_complete"
	ElseEdge             EdgeType = "else"
	IfEdge               EdgeType = "if"
)

type ResourceNodes struct {
	AppsById       map[string]*Node
	ComponentsById map[string]*Node
	PlaybooksById  map[string]*Node
	ConnectorsById map[string]*Node
	TriggersById   map[string]*Node
}

// createNodes creates nodes for top-level resources in the dump.
func createNodes(laneState *lanedump.LaneState) *ResourceNodes {
	groups := &ResourceNodes{
		AppsById:       createAppNodes(laneState),
		ComponentsById: createComponentNodes(laneState),
		PlaybooksById:  createPlaybookNodes(laneState),
		ConnectorsById: createConnectorNodes(laneState),
	}

	return groups
}

// createAppNodes creates nodes for each application.
func createAppNodes(state *lanedump.LaneState) map[string]*Node {
	nodes := make(map[string]*Node, len(state.ApplicationsById))

	for appId, app := range state.ApplicationsById {
		label := fmt.Sprintf("[%s] %s", app.Acronym, app.Name)
		meta := newMeta(appId, ApplicationNode, label, "")
		nodes[appId] = newNode(meta)
	}

	return nodes
}

// createComponentNodes creates nodes for each component.
func createComponentNodes(state *lanedump.LaneState) map[string]*Node {
	nodes := make(map[string]*Node, len(state.ComponentsById))

	for compId, comp := range state.ComponentsById {
		meta := newMeta(compId, ComponentNode, comp.Name, comp.Description)
		nodes[compId] = newNode(meta)
	}

	return nodes
}

// createPlaybookNodes creates nodes for each playbook.
func createPlaybookNodes(state *lanedump.LaneState) map[string]*Node {
	nodes := make(map[string]*Node, len(state.PlaybooksById))

	for pbId, pb := range state.PlaybooksById {
		meta := newMeta(pbId, PlaybookNode, pb.Name, pb.Description)
		nodes[pbId] = newNode(meta)
	}

	return nodes
}

// createConnectorNodes creates nodes for each connector.
func createConnectorNodes(state *lanedump.LaneState) map[string]*Node {
	nodes := make(map[string]*Node, len(state.ConnectorsById))

	for _, conn := range state.ConnectorsById {
		meta := newMeta(conn.Meta.Manifest.Name, ConnectorNode, conn.Meta.Manifest.Title, conn.Meta.Manifest.Product)
		nodes[conn.Meta.Manifest.Name] = newNode(meta) // Connectors are referenced by their manifest name
	}

	return nodes
}
