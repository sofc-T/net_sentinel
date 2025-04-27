package models

// PingResult stores ICMP ping test results
type PingResult struct {
	deviceID   string
	ipAddress  string
	success    bool
	latencyMs  int
	timestamp  int64
}

// NewPingResult creates a new PingResult instance
func NewPingResult(deviceID, ipAddress string, success bool, latencyMs int, timestamp int64) *PingResult {
	return &PingResult{
		deviceID:  deviceID,
		ipAddress: ipAddress,
		success:   success,
		latencyMs: latencyMs,
		timestamp: timestamp,
	}
}

// ConfigurePingResult updates the fields of an existing PingResult
func (p *PingResult) ConfigurePingResult(deviceID, ipAddress string, success bool, latencyMs int, timestamp int64) {
	p.deviceID = deviceID
	p.ipAddress = ipAddress
	p.success = success
	p.latencyMs = latencyMs
	p.timestamp = timestamp
}

// GetDeviceID returns the DeviceID
func (p *PingResult) GetDeviceID() string {
	return p.deviceID
}

// SetDeviceID sets the DeviceID
func (p *PingResult) SetDeviceID(deviceID string) {
	p.deviceID = deviceID
}

// GetIPAddress returns the IPAddress
func (p *PingResult) GetIPAddress() string {
	return p.ipAddress
}

// SetIPAddress sets the IPAddress
func (p *PingResult) SetIPAddress(ipAddress string) {
	p.ipAddress = ipAddress
}

// GetSuccess returns the Success status
func (p *PingResult) GetSuccess() bool {
	return p.success
}

// SetSuccess sets the Success status
func (p *PingResult) SetSuccess(success bool) {
	p.success = success
}

// GetLatencyMs returns the LatencyMs
func (p *PingResult) GetLatencyMs() int {
	return p.latencyMs
}

// SetLatencyMs sets the LatencyMs
func (p *PingResult) SetLatencyMs(latencyMs int) {
	p.latencyMs = latencyMs
}

// GetTimestamp returns the Timestamp
func (p *PingResult) GetTimestamp() int64 {
	return p.timestamp
}

// SetTimestamp sets the Timestamp
func (p *PingResult) SetTimestamp(timestamp int64) {
	p.timestamp = timestamp
}
