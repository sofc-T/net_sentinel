package probe

import (
	"fmt"
	"net"
	"time"

	"github.com/sofc-t/sentinel/domain/models"
)

// NetBIOSScan tries to resolve NetBIOS names for devices in subnet
func NetBIOSScan(ip string) string {
	timeout := 500 * time.Millisecond
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: 137, // NetBIOS Name Service
	}
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		return ""
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(timeout))

	// minimal query (broadcast), ignoring response parsing for now
	buf := []byte{
		0x00, 0x00, // Transaction ID
		0x00, 0x10, // Flags
		0x00, 0x01, // Questions
		0x00, 0x00, // Answer RRs
		0x00, 0x00, // Authority RRs
		0x00, 0x00, // Additional RRs
		// NetBIOS name query header follows...
	}
	_, err = conn.Write(buf)
	if err != nil {
		return ""
	}

	resp := make([]byte, 1024)
	_, err = conn.Read(resp)
	if err != nil {
		return ""
	}

	// crude: just return "NetBIOS-device" placeholder
	return fmt.Sprintf("NetBIOS-%s", ip)
}

// CaptureNetBIOS scans a list of IPs and returns Device objects
func CaptureNetBIOS(ips []string) []models.Device {
	devices := []models.Device{}
	for _, ip := range ips {
		name := NetBIOSScan(ip)
		if name != "" {
			device := models.NewDevice(models.DeviceConfig{
				Hostname: name,
				IPAddress: ip,
				DeviceType: "unknown",
				Vendor: "",
				Status: "active",
				MonitoringProtocols: []string{"NetBIOS"},
			})
			devices = append(devices, *device)
		}
	}
	return devices
}
