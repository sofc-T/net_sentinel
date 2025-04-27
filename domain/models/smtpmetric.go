package models

// SNMPMetric stores SNMP-based device metrics
type SNMPMetric struct {
	deviceID    string  
	metricType  string  
	value       float64 
	timestamp   int64
	Values      map[string]string   
}

type SNMPResult struct {
	DeviceID    string             `json:"device_id"`      // ID of the device queried
	IPAddress   string             `json:"ip_address"`     // IP address of the device
	Success     bool               `json:"success"`        // Whether SNMP query was successful
	Message     string             `json:"message"`        // Additional info (error, success msg)
	Timestamp   int64              `json:"timestamp"`      // When the SNMP request was made
	Metrics     *SNMPMetric         `json:"metrics"`        // Actual collected metrics if success
}

// GetDeviceID returns the DeviceID
func (s *SNMPMetric) GetDeviceID() string {
	return s.deviceID
}

// SetDeviceID sets the DeviceID
func (s *SNMPMetric) SetDeviceID(deviceID string) {
	s.deviceID = deviceID
}

// GetMetricType returns the MetricType
func (s *SNMPMetric) GetMetricType() string {
	return s.metricType
}

// SetMetricType sets the MetricType
func (s *SNMPMetric) SetMetricType(metricType string) {
	s.metricType = metricType
}

// GetValue returns the Value
func (s *SNMPMetric) GetValue() float64 {
	return s.value
}

// SetValue sets the Value
func (s *SNMPMetric) SetValue(value float64) {
	s.value = value
}

// GetTimestamp returns the Timestamp
func (s *SNMPMetric) GetTimestamp() int64 {
	return s.timestamp
}

// SetTimestamp sets the Timestamp
func (s *SNMPMetric) SetTimestamp(timestamp int64) {
	s.timestamp = timestamp
}
