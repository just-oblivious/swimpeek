package graph

import "fmt"

// reflectCronTrigger extracts the cron schedule from a cron trigger.
func reflectCronTrigger(conf any) (string, error) {
	if scheduleMaps, ok := conf.([]any); ok {
		for _, scheduleMap := range scheduleMaps {
			if schedule, ok := scheduleMap.(map[string]any); ok {
				for _, cron := range schedule {
					if cronStr, ok := cron.(string); ok {
						return cronStr, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("invalid cron trigger configuration: %v", conf)
}

// reflectSensorTrigger extracts the sensor name from a sensor trigger.
func reflectSensorTrigger(conf any) (string, error) {
	if sensorList, ok := conf.([]any); ok {
		if sensorMap, ok := sensorList[0].(map[string]any); ok {
			for name := range sensorMap {
				return name, nil // Return the first sensor name found.
			}
		}
	}
	return "", fmt.Errorf("invalid sensor trigger configuration: %v", conf)
}

// reflectEmitAction extracts the sensor name from the input data of an emit action.
func reflectEmitAction(inputs any) (string, error) {
	if emitInputs, ok := inputs.(map[string]any); ok {
		if sensorName, ok := emitInputs["sensorName"].(string); ok {
			return sensorName, nil
		}
	}
	return "", fmt.Errorf("invalid emit action inputs: %v", inputs)
}

// reflectRecordActionAppId extracts the application ID from the input data of a record action.
func reflectRecordActionAppId(inputs any) (string, error) {
	cfg, ok := inputs.(map[string]any)
	if ok {

		// The Application id may be a string or a map describing a dynamic reference, we're unable to resolve the latter.
		if appId, ok := cfg["applicationId"].(string); ok {
			return appId, nil
		}
		if appIdMap, ok := cfg["applicationId"].(map[string]any); ok {
			return "", fmt.Errorf("dynamic application ID reference not supported: %v", appIdMap)
		}
	}
	return "", fmt.Errorf("invalid record action inputs: %v", inputs)
}

// reflectRecordActionType determines the type of record action based on its inputs.
func reflectRecordActionType(inputs any) (NodeType, error) {
	cfg, ok := inputs.(map[string]any)
	if ok {
		if _, hasPatchValues := cfg["patchValues"]; hasPatchValues {
			return RecordUpdateActionNode, nil
		}
		if _, hasFields := cfg["fields"]; hasFields {
			return RecordCreateActionNode, nil
		}
		if _, hasFilter := cfg["filters"]; hasFilter {
			return RecordSearchActionNode, nil
		}
	}
	return "", fmt.Errorf("invalid record action inputs: %v", inputs)
}
