package lanedump

import (
	"time"

	"github.com/just-oblivious/swimpeek/pkg/laneclient"
)

// LaneState holds the state of the SwimLane tenant.
type LaneState struct {
	TimeStamp          time.Time
	Tenant             laneclient.Tenant                           // The tenant this state is for.
	PlaybooksById      map[string]laneclient.OrchestrationSolution // Playbooks are orchestration solutions referencing one or more workflows.
	ComponentsById     map[string]laneclient.OrchestrationSolution // Components are orchestration solutions referencing exactly one workflow.
	WorkflowsById      map[string]laneclient.Workflow              // Workflows describe the chain of actions to be performed (the "playbook").
	ApplicationsById   map[string]laneclient.Application           // Applications define the shape of the records that can be referenced.
	ConnectorsById     map[string]laneclient.Connector             // Connectors are called by workflows to perform a variety of actions.
	SensorsById        map[string]laneclient.Sensor                // Sensors are event listeners like webhooks or flow events.
	OrchestrationTasks []laneclient.OrchestrationTask              // Orchestration tasks are references between workflows and applications (e.g. recordAction, playbookButton).
}
