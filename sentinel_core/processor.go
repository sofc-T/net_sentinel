package sentinel

import (
	"encoding/json"
	"log"
)

// NetworkEvent represents incoming data
type NetworkEvent struct {
	DeviceID   string  `json:"device_id"`
	Metric     string  `json:"metric"`
	Value      float64 `json:"value"`
	Timestamp  int64   `json:"timestamp"`
}

// ProcessEvent handles network event processing
func ProcessEvent(data []byte) {
	var event NetworkEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Printf("Error parsing event: %v", err)
		return
	}

	log.Printf("Processing event: %+v", event)

	// Apply anomaly detection rules
	if isAnomalous(event) {
		log.Printf("Anomaly detected in device %s: %s = %f", event.DeviceID, event.Metric, event.Value)
		// Send to AlertManager
	}
}

// isAnomalous checks for basic anomalies
func isAnomalous(event NetworkEvent) bool {
	// Example rule: High CPU usage (>90%)
	if event.Metric == "cpu_usage" && event.Value > 90 {
		return true
	}
	return false
}
