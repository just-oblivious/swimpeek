package lanedump

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/just-oblivious/swimpeek/internal/config"
	"github.com/just-oblivious/swimpeek/pkg/laneclient"

	"github.com/charmbracelet/log"

	"golang.org/x/sync/errgroup"
)

var logger *log.Logger = config.GetLogger("lanedump")

// LoadFromTenant loads the orchestration state from a tenant.
func LoadFromTenant(ctx context.Context, laneClient *laneclient.TenantClient) (*LaneState, error) {
	laneState := LaneState{
		TimeStamp: time.Now(),
		Tenant:    laneClient.Tenant,
	}

	eg, _ctx := errgroup.WithContext(ctx)

	logger.Info("Enumerating tenant...")

	// Playbooks
	eg.Go(func() error {
		playbooks, err := laneClient.GetPlaybooks(_ctx)
		if err != nil {
			return fmt.Errorf("failed to get playbooks: %w", err)
		}
		laneState.PlaybooksById = make(map[string]laneclient.OrchestrationSolution, len(playbooks))
		for _, solution := range playbooks {
			laneState.PlaybooksById[solution.Id] = solution
		}
		return nil
	})

	// Components
	eg.Go(func() error {
		components, err := laneClient.GetComponents(_ctx)
		if err != nil {
			return fmt.Errorf("failed to get components: %w", err)
		}
		laneState.ComponentsById = make(map[string]laneclient.OrchestrationSolution, len(components))
		for _, component := range components {
			laneState.ComponentsById[component.Id] = component
		}
		return nil
	})

	// Playbook workflows
	eg.Go(func() error {
		workflows, err := laneClient.GetPlaybookWorkflows(_ctx)
		if err != nil {
			return fmt.Errorf("failed to get workflows: %w", err)
		}
		laneState.WorkflowsById = make(map[string]laneclient.Workflow, len(workflows))
		for _, workflow := range workflows {
			laneState.WorkflowsById[workflow.Id] = workflow
		}
		return nil
	})

	// Applications
	eg.Go(func() error {
		applications, err := laneClient.GetApplications(_ctx)
		if err != nil {
			return fmt.Errorf("failed to get applications: %w", err)
		}
		laneState.ApplicationsById = make(map[string]laneclient.Application, len(applications))
		for _, app := range applications {
			laneState.ApplicationsById[app.Id] = app
		}
		return nil
	})

	// Connectors
	eg.Go(func() error {
		connectors, err := laneClient.GetConnectors(_ctx)
		if err != nil {
			return fmt.Errorf("failed to get connectors: %w", err)
		}
		laneState.ConnectorsById = make(map[string]laneclient.Connector, len(connectors))
		for _, connector := range connectors {
			laneState.ConnectorsById[connector.Id] = connector
		}
		return nil
	})

	// Sensors
	eg.Go(func() error {
		sensors, err := laneClient.GetSensors(_ctx)
		if err != nil {
			return fmt.Errorf("failed to get sensors: %w", err)
		}
		laneState.SensorsById = make(map[string]laneclient.Sensor, len(sensors))
		for _, sensor := range sensors {
			laneState.SensorsById[sensor.Id] = sensor
		}
		return nil
	})

	// Orchestration tasks
	eg.Go(func() error {
		otasks, err := laneClient.GetOrchestrationTasks(_ctx)
		if err != nil {
			return fmt.Errorf("failed to get orchestration tasks: %w", err)
		}
		laneState.OrchestrationTasks = otasks
		return nil
	})

	if err := eg.Wait(); err != nil {
		return &laneState, err
	}

	logger.Info("Done!", "playbooks", len(laneState.PlaybooksById),
		"components", len(laneState.ComponentsById),
		"workflows", len(laneState.WorkflowsById),
		"applications", len(laneState.ApplicationsById),
		"connectors", len(laneState.ConnectorsById),
		"sensors", len(laneState.SensorsById),
		"orchestrationTasks", len(laneState.OrchestrationTasks))

	return &laneState, nil
}

// LoadFromDisk loads an orchestration state from a JSON file on disk.
func LoadFromDisk(path string) (*LaneState, error) {
	laneState := LaneState{}

	data, err := os.ReadFile(path)
	if err != nil {
		return &laneState, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, &laneState); err != nil {
		return &laneState, fmt.Errorf("failed to unmarshal JSON from %s: %w", path, err)
	}

	return &laneState, nil
}

// WriteToDisk writes an orchestration state to a JSON file on disk.
func WriteToDisk(laneState *LaneState, path string) error {
	json, err := json.MarshalIndent(laneState, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lane state to JSON: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = file.Write(json)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file %s: %w", path, err)
	}
	return nil
}
