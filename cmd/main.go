package main

import (
	"log"
	"time"

	"github.com/sofc-t/sentinel/domain/models"
	"github.com/sofc-t/sentinel/probe"
)

func main() {
	interfaceName, subnet, err := probe.FindDefaultInterfaceAndSubnet()
	if err != nil {
		log.Fatalf("Failed to find default network interface: %v", err)
	}
	log.Printf("[Main] Using Interface: %s, Subnet: %s\n", interfaceName, subnet)

	var lldpDevices []models.Device

	log.Println("[Main] Starting LLDP Capture")
	done := make(chan struct{})
	go func() {
		var err error
		lldpDevices, err = probe.CaptureLLDP(interfaceName, 5*time.Second)
		if err != nil {
			log.Println("LLDP/CDP Capture Error:", err)
		}
		close(done)
	}()

	log.Println("[Main] Starting IP Scan")
	ipDevices, err := probe.ScanIPRange(subnet)
	if err != nil {
		log.Fatalf("IP Scan failed: %v", err)
	}

	log.Println("[Main] Starting ARP Scan")
	arpDevices, err := probe.ARPScan(interfaceName)
	if err != nil {
		log.Fatal(err)
	}

	<-done // wait for LLDP

	var allDevices []models.Device
	allDevices = append(allDevices, arpDevices...)
	allDevices = append(allDevices, ipDevices...)
	allDevices = append(allDevices, lldpDevices...)

	log.Println("Discovered Devices:")
	for _, dev := range allDevices {
		log.Printf("- IP: %s, Status: %s, Protocols: %v\n", dev.GetIPAddress(), dev.GetStatus(), dev.GetMonitoringProtocols())
	}

	var devicesToPing []map[string]string
	for _, dev := range allDevices {
		if dev.GetIPAddress() != "" {
			devicesToPing = append(devicesToPing, map[string]string{
				"id": dev.GetID(),
				"ip": dev.GetIPAddress(),
			})
		}
	}

	pingResults := probe.PingNetwork(devicesToPing, 2*time.Second)

	log.Println("Ping Results:")
	for _, r := range pingResults {
		log.Printf("Device %s (%s) - Success: %v, Latency: %dms",
			r.GetDeviceID(), r.GetIPAddress(), r.GetSuccess(), r.GetLatencyMs())
	}
}
