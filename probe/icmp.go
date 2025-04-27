package probe

import (
	"fmt"
	"net"
	// "net/netip"
	"time"

	"github.com/digineo/go-ping"
	"github.com/sofc-t/sentinel/domain/models"
)

// PingDevice sends an ICMP echo request using digineo/go-ping
func PingDevice(deviceID, ipAddress string, timeout time.Duration) models.PingResult {
	pinger, err := ping.New(ipAddress, "")
	if err != nil {
		fmt.Printf("Error creating pinger for %s: %v\n", ipAddress, err)
		return *models.NewPingResult(deviceID, ipAddress, false, -1, time.Now().Unix())
	}
	defer pinger.Close()


	ip, err := net.ResolveIPAddr("ip", ipAddress)
	if err != nil {
		fmt.Printf("Error resolving IP address for %s: %v\n", ipAddress, err)
		return *models.NewPingResult(deviceID, ipAddress, false, -1, time.Now().Unix())
	}

	rtt, err := pinger.Ping(ip, timeout)
	if err != nil {
		fmt.Printf("Ping failed for %s: %v\n", ipAddress, err)
		return *models.NewPingResult(deviceID, ipAddress, false, -1, time.Now().Unix())
	}

	latency := rtt.Milliseconds()
	return *models.NewPingResult(deviceID, ipAddress, true, int(latency), time.Now().Unix())
}

// PingNetwork pings multiple devices concurrently
func PingNetwork(devices []map[string]string, timeout time.Duration) []models.PingResult {
	results := make([]models.PingResult, len(devices))
	resultChan := make(chan models.PingResult, len(devices))

	for _, device := range devices {
		go func(d map[string]string) {
			resultChan <- PingDevice(d["id"], d["ip"], timeout)
		}(device)
	}

	for i := range devices {
		results[i] = <-resultChan
	}

	return results
}

// Main function for testing
// func main() {
// 	devices := []map[string]string{
// 		{"id": "device-1", "ip": "8.8.8.8"},
// 		{"id": "device-2", "ip": "1.1.1.1"},
// 	}

// 	results := PingNetwork(devices, 3*time.Second)

// 	for _, result := range results {
// 		fmt.Printf("Device %s (%s) - Success: %v, Latency: %dms\n",
// 			result.GetDeviceID(), result.GetIPAddress(), result.GetSuccess(), result.GetLatencyMs())
// 	}
// }
