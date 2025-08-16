package graph

import (
	"fmt"
	"strings"
	"swimpeek/internal/lanedump"
	"swimpeek/pkg/laneclient"
)

type Warnings struct {
	Warns []error
}

func (w *Warnings) Add(err error) {
	if err != nil {
		w.Warns = append(w.Warns, err)
	}
}

// newWarnings creates a new Warnings instance.
func newWarnings() *Warnings {
	return &Warnings{
		Warns: make([]error, 0),
	}
}

// linkGraph expands the graph by linking nodes based on the relationships inferred from LaneState.
func linkGraph(graph *Graph, laneState *lanedump.LaneState) ([]error, error) {
	warns := newWarnings()

	// Link workflows to playbooks and components.
	wfNodes, err := linkWorkflows(warns, graph, laneState)
	if err != nil {
		return nil, err
	}

	// Find triggers and link them to workflows and applications
	trNodes, err := linkTriggers(warns, graph, laneState, wfNodes)
	if err != nil {
		return nil, err
	}
	graph.Resources.TriggersById = trNodes

	// Traverse the chain of actions in each workflows and link the resources they reference.
	for wfId, wf := range laneState.WorkflowsById {
		wfNode, exists := wfNodes[wfId]
		if !exists {
			// Skip orphan workflows.
			continue
		}
		err := linkWorkflowActions(warns, graph, wf.Playbook, wfNode)
		if err != nil {
			return nil, err
		}
	}

	return warns.Warns, nil
}

// linkWorkflows links workflows to their playbooks and components.
func linkWorkflows(warns *Warnings, graph *Graph, laneState *lanedump.LaneState) (map[string]*Node, error) {
	wfNodes := make(map[string]*Node, len(laneState.WorkflowsById))

	// Create workflow nodes for playbooks and link them.
	for pbId, pb := range laneState.PlaybooksById {
		// Find the playbook node.
		pbNode, exists := graph.Resources.PlaybooksById[pbId]
		if !exists {
			warns.Add(fmt.Errorf("playbook node %s not found", pbId))
			continue
		}

		for idx, wfId := range pb.PlaybookIds {
			wf, exists := laneState.WorkflowsById[wfId]
			if !exists {
				warns.Add(fmt.Errorf("playbook %s references unknown workflow %s", pbId, wfId))
				continue
			}

			// If a workflows title is equal to its internal name, it's not customized and the label becomes Flow <idx>
			label := wf.Playbook.Title
			if wf.Playbook.Title == wf.Playbook.Name {
				label = fmt.Sprintf("Flow %d", idx+1)
			}

			wfNode := newNode(newMeta(wfId, WorkflowNode, label, wf.Playbook.Description))
			wfNodes[wfId] = wfNode

			// Link it to the playbook node.
			newEdge(pbNode, wfNode, WorkflowEdge, nil)
		}

	}

	// Create workflow nodes for components and link them.
	for compId, comp := range laneState.ComponentsById {
		// Find the component node.
		compNode, exists := graph.Resources.ComponentsById[compId]
		if !exists {
			warns.Add(fmt.Errorf("component node %s not found", compId))
			continue
		}
		wfId := comp.PlaybookId
		wf, exists := laneState.WorkflowsById[wfId]
		if !exists {
			warns.Add(fmt.Errorf("component %s references unknown workflow %s", compId, wfId))
			continue
		}
		wfNode := newNode(newMeta(wfId, WorkflowNode, wf.Playbook.Title, wf.Playbook.Description))
		wfNodes[wfId] = wfNode

		// Link it to the component node.
		newEdge(compNode, wfNode, WorkflowEdge, nil)
	}

	// Orphan workflows that are not linked to any playbook or component.
	for wfId, wf := range laneState.WorkflowsById {
		if _, exists := wfNodes[wfId]; exists {
			continue
		}
		warns.Add(fmt.Errorf("orphan workflow %s found with title %s (solution: %s)", wfId, wf.Playbook.Title, wf.Meta.SolutionId))
	}

	return wfNodes, nil
}

// linkTriggers links triggers to workflows and applications.
func linkTriggers(warns *Warnings, graph *Graph, laneState *lanedump.LaneState, wfNodes map[string]*Node) (map[string]*Node, error) {
	trNodes := make(map[string]*Node)

	// Sensor-type triggers (webhook, flow event)
	for _, sensor := range laneState.SensorsById {
		switch sensor.Sensor.Type {
		case "webhook":
			trNode := newNode(newMeta(sensor.Meta.Name, WebhookNode, sensor.Meta.Title, ""))
			trNodes[sensor.Meta.Name] = trNode
		case "flow":
			trNode := newNode(newMeta(sensor.Meta.Name, FlowEventNode, sensor.Meta.Title, ""))
			trNodes[sensor.Meta.Name] = trNode
		default:
			warns.Add(fmt.Errorf("unknown sensor type %s for sensor %s", sensor.Sensor.Type, sensor.Meta.Name))
			continue
		}
	}

	// Orchestration task triggers (record event, playbook button)
	for _, task := range laneState.OrchestrationTasks {
		app, exists := graph.Resources.AppsById[task.ApplicationId]
		if !exists {
			warns.Add(fmt.Errorf("orchestration task %s references unknown application %s", task.Id, task.ApplicationId))
			continue
		}

		wfNode, exists := wfNodes[task.PlaybookId]
		if !exists {
			warns.Add(fmt.Errorf("orchestration task %s references unknown workflow %s", task.Id, task.PlaybookId))
			continue
		}

		// Tasks without triggers are PlaybookButtons
		if len(task.Triggers) == 0 {
			trNode := newNode(newMeta(task.Id, PlaybookButtonNode, task.Name, ""))
			trNodes[task.Id] = trNode
			newEdge(app, trNode, HasActionEdge, nil)
			newEdge(trNode, wfNode, TriggersWorkflowEdge, nil)
			continue
		}

		// Tasks with triggers are RecordEvents
		recordTriggers := make([]string, 0)
		for _, trigger := range task.Triggers {
			if trigger.OnRecordCreate {
				recordTriggers = append(recordTriggers, "on_create")
			}
			if trigger.OnRecordUpdate {
				recordTriggers = append(recordTriggers, "on_update")
			}
			if trigger.OnCorrelationActionComplete {
				recordTriggers = append(recordTriggers, "on_correlated")
			}
		}
		trNode := newNode(newMeta(task.Id, RecordEventNode, strings.Join(recordTriggers, ", "), ""))
		trNodes[task.Id] = trNode
		newEdge(app, trNode, HasEventEdge, nil)
		newEdge(trNode, wfNode, TriggersWorkflowEdge, nil)
	}

	// Enumerate the workflows to find cron triggers and link the sensor/flow triggers.
	for wfId, wf := range laneState.WorkflowsById {
		wfNode, exists := wfNodes[wfId]
		if !exists {
			// Skip orphan workflows.
			continue
		}
		for trigType, trigConf := range wf.Playbook.Triggers {
			switch trigType {
			case "schedules":
				trId := fmt.Sprintf("%s_cron", wfId)
				schedule, err := reflectCronTrigger(trigConf)
				if err != nil {
					warns.Add(fmt.Errorf("failed to reflect cron trigger for workflow %s: %w", wfId, err))
					continue
				}
				label := fmt.Sprintf("Scheduled (%s)", schedule)
				tfNode := newNode(newMeta(trId, CronEventNode, label, ""))
				newEdge(tfNode, wfNode, TriggersWorkflowEdge, nil)
				trNodes[trId] = tfNode
			case "sensors", "flows":
				sensor, err := reflectSensorTrigger(trigConf)
				if err != nil {
					warns.Add(fmt.Errorf("failed to reflect sensor trigger for workflow %s: %w", wfId, err))
					continue
				}
				sensNode, exists := trNodes[sensor]
				if !exists {
					warns.Add(fmt.Errorf("sensor trigger %s not found for workflow %s", sensor, wfId))
					continue
				}
				newEdge(sensNode, wfNode, TriggersWorkflowEdge, nil)
			}
		}
	}

	return trNodes, nil
}

// linkWorkflowActions links the action chain in the workflow.
func linkWorkflowActions(warns *Warnings, graph *Graph, wfPlaybook laneclient.Playbook, wfNode *Node) error {

	// Chain the actions, starting from the entrypoint actions.
	if err := chainActions(warns, graph, wfNode, wfPlaybook.Actions, wfPlaybook.Entrypoints...); err != nil {
		return fmt.Errorf("failed to chain actions for workflow %s: %w", wfNode.Meta.Id, err)
	}

	return nil
}
