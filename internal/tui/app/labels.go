package app

import "swimpeek/internal/graph"

// nodeIcons provides visual symbols for different node types.
var NodeIcons map[graph.NodeType]string = map[graph.NodeType]string{
	graph.ConditionalActionNode: "◆",
	graph.ParallelActionNode:    "‖",
	graph.ForEachLoopAction:     "↻",
	graph.WhileLoopAction:       "↻",

	graph.FlowEventNode:   "✲",
	graph.ComponentNode:   "Σ",
	graph.ApplicationNode: "⌘",
	graph.ConnectorNode:   "⎋",
	graph.WorkflowNode:    "▶",
	graph.PlaybookNode:    "⎔",
}

// nodeLabels provides human-readable labels for different node types.
var NodeLabels map[graph.NodeType]string = map[graph.NodeType]string{
	graph.ConditionalActionNode:    "condition",
	graph.ParallelActionNode:       "parallel",
	graph.ForEachLoopAction:        "for each",
	graph.WhileLoopAction:          "while",
	graph.ComponentActionNode:      "component",
	graph.ConnectorActionNode:      "connector",
	graph.RecordActionNode:         "record",
	graph.RecordCreateActionNode:   "create record",
	graph.RecordUpdateActionNode:   "update record",
	graph.RecordSearchActionNode:   "search records",
	graph.RecordDeleteActionNode:   "delete record",
	graph.RecordUpsertActionNode:   "upsert record",
	graph.RecordExportActionNode:   "export records",
	graph.EmitEventActionNode:      "emit event",
	graph.TransformationActionNode: "transformation",
	graph.PythonActionNode:         "python",
	graph.CreateVarsActionNode:     "create vars",
	graph.UpdateVarsActionNode:     "update vars",
	graph.HTTPActionNode:           "http",
	graph.PlaybookButtonNode:       "playbook button",
	graph.RecordEventNode:          "record event",
	graph.CronEventNode:            "cron event",
	graph.WebhookNode:              "incoming webhook",
}

// edgeLabels provides human-readable labels for different edge types.
var EdgeLabels map[graph.EdgeType]string = map[graph.EdgeType]string{
	graph.OnSuccessEdge:        "on success",
	graph.OnFailureEdge:        "on failure",
	graph.OnCompleteEdge:       "on complete",
	graph.ElseEdge:             "else",
	graph.IfEdge:               "if",
	graph.HasActionEdge:        "action",
	graph.HasEventEdge:         "event",
	graph.TriggersWorkflowEdge: "triggers",
	graph.UnreachableEdge:      "unreachable",
}
