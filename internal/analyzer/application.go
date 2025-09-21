package analyzer

import (
	"fmt"
	"swimpeek/internal/graph"
	"swimpeek/pkg/laneclient"
)

type TriggerAction struct {
	Trigger  *graph.Node
	Playbook *graph.Node
	Workflow *graph.Node
	Enabled  bool
}

type ApplicationTriggersResult struct {
	ButtonTriggers      []TriggerAction
	RecordEventTriggers []TriggerAction
}

type AccessAction struct {
	Action        *graph.Node
	Playbook      *graph.Node
	Component     *graph.Node
	Workflow      *graph.Node
	Enabled       bool
	InspectionErr error // Error encountered during field-level inspection of an action, if any.

}

// IsPlaybookAction returns true if this access action is part of a component.
func (a *AccessAction) IsComponentAction() bool {
	return a.Component != nil
}

// ApplicationTriggers analyzes which playbook-workflows can be triggered by the given application.
func (a *Analyzer) ApplicationTriggers(appNode *graph.Node) *ApplicationTriggersResult {

	createTriggersFn := func(trigs map[*graph.Node]bool) []TriggerAction {
		triggerActions := make([]TriggerAction, 0, len(trigs))
		for _, trig := range SortSetByLabel(trigs) {
			wf := a.GetWorkflowForTrigger(trig)
			pb := a.GetPlaybookForWorkflow(wf)
			wfResource := a.GetWorkflowResource(wf)
			wfEnabled := false
			if wfResource != nil {
				wfEnabled = wfResource.Meta.Enabled
			}
			triggerActions = append(triggerActions, TriggerAction{
				Trigger:  trig,
				Playbook: pb,
				Workflow: wf,
				Enabled:  wfEnabled,
			})
		}
		return triggerActions
	}

	btnNodes := a.FindUnique(appNode, NewWalkOpts(Descend, WithMaxDepth(1)), graph.PlaybookButtonNode)
	pbButtons := createTriggersFn(btnNodes)

	evtNodes := a.FindUnique(appNode, NewWalkOpts(Descend, WithMaxDepth(1)), graph.RecordEventNode)
	recordEvents := createTriggersFn(evtNodes)

	return &ApplicationTriggersResult{
		ButtonTriggers:      pbButtons,
		RecordEventTriggers: recordEvents,
	}
}

// ApplicationAccessedBy analyzes which components and playbook-workflow actions access records in this application.
func (a *Analyzer) ApplicationAccessedBy(appNode *graph.Node) []AccessAction {
	recordAccessActions := a.FindAll(appNode, NewWalkOpts(Descend, WithFollowEdgeTypes(graph.AccessedByEdge)), graph.RecordCreateActionNode, graph.RecordUpsertActionNode, graph.RecordUpdateActionNode, graph.RecordDeleteActionNode, graph.RecordSearchActionNode, graph.RecordExportActionNode)

	accessActions := make([]AccessAction, 0, 10)

	for _, actionNode := range recordAccessActions {
		wfNode := a.GetWorkflowForAction(actionNode)
		if wfNode == nil {
			continue
		}

		wfResource := a.GetWorkflowResource(wfNode)
		wfEnabled := false
		if wfResource != nil {
			wfEnabled = wfResource.Meta.Enabled
		}

		// Action is part of a component
		compNode := a.GetComponentForWorkflow(wfNode)
		if compNode != nil {
			accessActions = append(accessActions, AccessAction{
				Action:    actionNode,
				Component: compNode,
				Workflow:  wfNode,
			})
		}

		// Action is part of a playbook-workflow
		pbNode := a.GetPlaybookForWorkflow(wfNode)
		if pbNode != nil {
			accessActions = append(accessActions, AccessAction{
				Action:   actionNode,
				Playbook: pbNode,
				Workflow: wfNode,
				Enabled:  wfEnabled,
			})
		}

	}

	return accessActions
}

// ApplicationFieldModifiedBy analyzes which components and playbook-workflow actions modify the given application field.
func (a *Analyzer) ApplicationFieldModifiedBy(appNode *graph.Node, field *laneclient.ApplicationField) []AccessAction {
	accessActions := a.ApplicationAccessedBy(appNode)

	fieldProperties := map[string]struct{}{
		"appendValues": {},
		"patchValues":  {},
		"removeValues": {},
		"fields":       {},
	}

	touchesFieldFn := func(inputsMap map[string]any) (bool, error) {
		for fieldProp := range fieldProperties {
			fg, ok := inputsMap[fieldProp]
			if !ok {
				continue
			}

			fgMap, ok := fg.(map[string]any)
			if !ok {
				continue
			}

			if ref, isRef := fgMap["$:ref"]; isRef {
				return false, fmt.Errorf("field property %s is a reference (%v), cannot inspect", fieldProp, ref)
			}

			if _, touchesField := fgMap[field.Key]; touchesField {
				return true, nil
			}
		}
		return false, nil
	}

	filteredActions := make([]AccessAction, 0, len(accessActions))
	for _, action := range accessActions {

		switch action.Action.Meta.Type {
		case graph.RecordSearchActionNode, graph.RecordExportActionNode, graph.RecordDeleteActionNode:
			continue
		}

		actionResource := a.GetActionResource(action.Workflow, action.Action)
		if actionResource == nil {
			action.InspectionErr = fmt.Errorf("action resource not found")
			filteredActions = append(filteredActions, action)
			continue
		}

		inputsMap, ok := actionResource.Inputs.(map[string]any)
		if !ok {
			action.InspectionErr = fmt.Errorf("action inputs malformed or missing")
			filteredActions = append(filteredActions, action)
			continue
		}

		fieldTouched, err := touchesFieldFn(inputsMap)
		if err != nil {
			action.InspectionErr = err
			filteredActions = append(filteredActions, action)
			continue
		}

		if fieldTouched {
			filteredActions = append(filteredActions, action)
		}
	}

	return filteredActions
}
