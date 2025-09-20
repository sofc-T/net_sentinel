package main

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/sofc-t/sentinel/domain/models"
	"github.com/sofc-t/sentinel/probe"
	sentinel "github.com/sofc-t/sentinel/sentinel_core"
)

func lookupVendorFromMAC(mac string) string {
	// Placeholder for MAC vendor lookup – integrate with a local OUI DB if available
	if mac == "" {
		return ""
	}
	// Simple heuristic: return prefix as "Vendor-<prefix>"
	parts := strings.Split(mac, ":")
	if len(parts) > 2 {
		return "Vendor-" + strings.ToUpper(parts[0]) + strings.ToUpper(parts[1])
	}
	return "UnknownVendor"
}

func main() {

	allDevices := []sentinel.DeviceRecord{}
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
		lldpDevices, err = probe.CaptureLLDP(interfaceName, 50*time.Second)
		if err != nil {
			log.Println("LLDP/CDP Capture Error:", err)
		}
		close(done)
	}()

	// format and append lldpDevices to allDevices
	for _, d := range lldpDevices {
		allDevices = append(allDevices, sentinel.DeviceRecord{
			DeviceID:  d.GetID(),
			Hostname:  d.GetHostname(),
			IP:        d.GetIPAddress(),
			MAC:       d.GetMACAddress(),
			Status:    d.GetStatus(),
			Type:      d.GetDeviceType(),
			Vendor:    d.GetVendor(),
			Protocols: strings.Join(d.GetMonitoringProtocols(), ","),
		})
	}

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

	// format and append ipDevices to allDevices
	for _, d := range ipDevices {

		allDevices = append(allDevices, sentinel.DeviceRecord{
			DeviceID:  d.GetID(),
			Hostname:  d.GetHostname(),
			IP:        d.GetIPAddress(),
			Status:    d.GetStatus(),
			Type:      d.GetDeviceType(),
			Vendor:    d.GetVendor(),
			Protocols: strings.Join(d.GetMonitoringProtocols(), ","),
		})
	}

	// format and append arpDevices to allDevices
	for _, d := range arpDevices {
		allDevices = append(allDevices, sentinel.DeviceRecord{
			DeviceID:  d.GetID(),
			Hostname:  d.GetHostname(),
			IP:        d.GetIPAddress(),
			MAC:       d.GetMACAddress(),
			Status:    d.GetStatus(),
			Type:      d.GetDeviceType(),
			Vendor:    d.GetVendor(),
			Protocols: strings.Join(d.GetMonitoringProtocols(), ","),
		})
	}

	log.Println("Discovered Devices:")
	for _, dev := range allDevices {
		log.Printf("- IP: %s, Status: %s, Protocols: %v\n", dev.IP, dev.Status, dev.Protocols)
	}

	// ping for reachability and latency
	var devicesToPing []map[string]string
	for _, dev := range allDevices {
		if dev.IP != "" {
			devicesToPing = append(devicesToPing, map[string]string{
				"id": dev.DeviceID,
				"ip": dev.IP,
			})
		}
	}

	pingResults := probe.PingNetwork(devicesToPing, 20*time.Second)
	log.Println("Ping Results:")
	for _, r := range pingResults {
		log.Printf("Device %s (%s) - Success: %v, Latency: %dms", r.GetDeviceID(), r.GetIPAddress(), r.GetSuccess(), r.GetLatencyMs())
	}

	// Fetch SNMP metrics for discovered devices
	for i, dev := range allDevices {
		if dev.IP != "" {
			// Reverse DNS lookup for hostname
			if dev.Hostname == "" {
				if names, err := net.LookupAddr(dev.IP); err == nil && len(names) > 0 {
					allDevices[i].Hostname = strings.TrimSuffix(names[0], ".")
				}
			}

			// MAC vendor lookup
			if dev.Vendor == "" && dev.MAC != "" {
				allDevices[i].Vendor = lookupVendorFromMAC(dev.MAC)
			}

			// Optional Nmap fingerprinting (if implemented in probe)
			if descr, protos := probe.NmapFingerprint(dev.IP); descr != "" {
				allDevices[i].Descr = descr
				if allDevices[i].Protocols != "" {
					allDevices[i].Protocols += "," + protos
				} else {
					allDevices[i].Protocols = protos
				}
			}

			snmpConfig := probe.SNMPConfig{
				Target:    dev.IP,
				Port:      161,
				Version:   gosnmp.Version2c,
				Community: "public",
				Timeout:   20 * time.Second,
				Retries:   1,
			}
			oids := []string{
				"1.3.6.1.2.1.1.3.0",      // sysUpTime
				"1.3.6.1.2.1.1.5.0",      // sysName
				"1.3.6.1.2.1.1.1.0",      // sysDescr
				"1.3.6.1.2.1.2.2.1.10.1", // ifInOctets (interface 1)
				"1.3.6.1.2.1.2.2.1.16.1", // ifOutOctets (interface 1)
				"1.3.6.1.2.1.2.2.1.14.1", // ifInErrors
				"1.3.6.1.2.1.2.2.1.20.1", // ifOutErrors
				"1.3.6.1.2.1.25.1.1.0",   // hrSystemUptime
				"1.3.6.1.2.1.25.2.2.0",   // hrMemorySize (total RAM)
				"1.3.6.1.4.1.9.2.1.57.0", // Cisco CPU utilization (example vendor-specific)
			}

			metrics, err := probe.FetchMetrics(snmpConfig, oids)
			if err != nil {
				log.Printf("SNMP fetch failed for %s: %v", dev.IP, err)
				continue
			}

			log.Printf("SNMP Metrics for %s: %+v", dev.IP, metrics.Metrics.Values)
		}
	}

	sentinel.DisplayTable(allDevices)

	// 😉 default username and password
	// for _, dev := range allDevices {
	// 	if dev.GetIPAddress() != "" {
	// 		sshConfig := probe.SSHConfig{
	// 			Host:     dev.GetIPAddress(),
	// 			Port:     "22",
	// 			Username: "admin",
	// 			Password: "password",
	// 			Timeout:  3 * time.Second,
	// 		}
	// 		output, err := probe.RunSSHCommand(sshConfig, "show interfaces")
	// 		if err == nil {
	// 			log.Printf("SSH Output for %s:\n%s", dev.GetIPAddress(), output)
	// 		}

	// 		// Example Telnet
	// 		output, err = probe.RunTelnetCommand(dev.GetIPAddress()+":23", "show ip route")
	// 		if err == nil {
	// 			log.Printf("Telnet Output for %s:\n%s", dev.GetIPAddress(), output)
	// 		}
	// 	}
	// }

}
