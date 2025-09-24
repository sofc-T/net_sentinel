package main

import (
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"fmt"

	"github.com/sofc-t/sentinel/probe"
	sentinel "github.com/sofc-t/sentinel/sentinel_core"
	"github.com/gosnmp/gosnmp"
)

func lookupVendorFromMAC(mac string) string {
	if mac == "" {
		return ""
	}
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

	// Channels for discovered devices
	devChan := make(chan sentinel.DeviceRecord, 100)
	var wg sync.WaitGroup

	// LLDP Capture
	wg.Add(1)
	go func() {
		defer wg.Done()
		lldpDevices, err := probe.CaptureLLDP(interfaceName, 10*time.Second)
		if err != nil {
			log.Println("LLDP capture error:", err)
			return
		}
		for _, d := range lldpDevices {
			devChan <- sentinel.DeviceRecord{
				DeviceID:  d.GetID(),
				Hostname:  d.GetHostname(),
				IP:        d.GetIPAddress(),
				MAC:       d.GetMACAddress(),
				Status:    d.GetStatus(),
				Type:      d.GetDeviceType(),
				Vendor:    d.GetVendor(),
				Protocols: strings.Join(d.GetMonitoringProtocols(), ","),
			}
		}
	}()

	// IP Scan
	wg.Add(1)
	go func() {
		defer wg.Done()
		ipDevices, err := probe.ScanIPRange(subnet)
		if err != nil {
			log.Println("IP scan error:", err)
			return
		}
		for _, d := range ipDevices {
			devChan <- sentinel.DeviceRecord{
				DeviceID:  d.GetID(),
				Hostname:  d.GetHostname(),
				IP:        d.GetIPAddress(),
				Status:    d.GetStatus(),
				Type:      d.GetDeviceType(),
				Vendor:    d.GetVendor(),
				Protocols: strings.Join(d.GetMonitoringProtocols(), ","),
			}
		}
	}()

	// ARP Scan
	wg.Add(1)
	go func() {
		defer wg.Done()
		arpDevices, err := probe.ARPScan(interfaceName)
		if err != nil {
			log.Println("ARP scan error:", err)
			return
		}
		for _, d := range arpDevices {
			devChan <- sentinel.DeviceRecord{
				DeviceID:  d.GetID(),
				Hostname:  d.GetHostname(),
				IP:        d.GetIPAddress(),
				MAC:       d.GetMACAddress(),
				Status:    d.GetStatus(),
				Type:      d.GetDeviceType(),
				Vendor:    d.GetVendor(),
				Protocols: strings.Join(d.GetMonitoringProtocols(), ","),
			}
		}
	}()

	// Wait for discovery scans
	go func() {
		wg.Wait()
		close(devChan)
	}()

	// Collect devices
	for d := range devChan {
		allDevices = append(allDevices, d)
	}

	log.Printf("[Main] Found %d devices.\n", len(allDevices))

	// Ping all devices concurrently
	pingChan := make(chan sentinel.DeviceRecord, len(allDevices))
	wgPing := sync.WaitGroup{}
	semPing := make(chan struct{}, 50) // limit concurrency

	for i := range allDevices {
		wgPing.Add(1)
		semPing <- struct{}{}
		go func(dev *sentinel.DeviceRecord) {
			defer wgPing.Done()
			defer func() { <-semPing }()
			if dev.IP != "" {
				results := probe.PingNetwork([]map[string]string{{"id": dev.DeviceID, "ip": dev.IP}}, 2*time.Second)
				if len(results) > 0 {
					dev.PingMs = int64(results[0].GetLatencyMs())
					if results[0].GetSuccess() {
						dev.Status = "active"
					} else {
						dev.Status = "inactive"
					}
				}
			}
			pingChan <- *dev
		}(&allDevices[i])
	}

	go func() {
		wgPing.Wait()
		close(pingChan)
	}()

	allDevices = []sentinel.DeviceRecord{}
	for d := range pingChan {
		allDevices = append(allDevices, d)
	}

	// SNMP + Nmap concurrently
	snmpNmapChan := make(chan sentinel.DeviceRecord, len(allDevices))
	wgSNMP := sync.WaitGroup{}
	semSNMP := make(chan struct{}, 20) // limit concurrency

	for i := range allDevices {
		wgSNMP.Add(1)
		semSNMP <- struct{}{}
		go func(dev *sentinel.DeviceRecord) {
			defer wgSNMP.Done()
			defer func() { <-semSNMP }()

			// Reverse DNS
			if dev.Hostname == "" && dev.IP != "" {
				if names, err := net.LookupAddr(dev.IP); err == nil && len(names) > 0 {
					dev.Hostname = strings.TrimSuffix(names[0], ".")
				}
			}

			// Nmap fingerprint
			if dev.IP != "" {
				if descr, protos := probe.NmapFingerprint(dev.IP); descr != "" {
					dev.Descr = descr
					if dev.Protocols != "" {
						dev.Protocols += "," + protos
					} else {
						dev.Protocols = protos
					}
				}
			}

			if  dev.IP != "" {
				// Reverse DNS
				if dev.Hostname == "" {
					if names, err := net.LookupAddr(dev.IP); err == nil && len(names) > 0 {
						dev.Hostname = strings.TrimSuffix(names[0], ".")
					}
				}

				// Port scan
				openPorts := scanCommonPorts(dev.IP, 500*time.Millisecond)
				if len(openPorts) > 0 {
					dev.Protocols += ",ports"
					dev.Descr += fmt.Sprintf("Open ports: %v ", openPorts)

					// Guess OS
					dev.Type = guessOS(openPorts)
				}
			}


			// SNMP metrics
			if dev.IP != "" {
				config := probe.SNMPConfig{
					Target:    dev.IP,
					Port:      161,
					Version:   gosnmp.Version2c,
					Community: "public",
					Timeout:   2 * time.Second,
					Retries:   1,
				}
				oids := []string{
					"1.3.6.1.2.1.1.3.0",
					"1.3.6.1.2.1.1.5.0",
					"1.3.6.1.2.1.1.1.0",
					"1.3.6.1.2.1.2.2.1.10.1",
					"1.3.6.1.2.1.2.2.1.16.1",
					"1.3.6.1.2.1.2.2.1.14.1",
					"1.3.6.1.2.1.2.2.1.20.1",
				}
				metrics, err := probe.FetchMetrics(config, oids)
				if err == nil && metrics != nil && metrics.Metrics != nil {
					if val, ok := metrics.Metrics.Values["1.3.6.1.2.1.2.2.1.10.1"]; ok {
						dev.IntIn = parseInt64(val)
					}
					if val, ok := metrics.Metrics.Values["1.3.6.1.2.1.2.2.1.16.1"]; ok {
						dev.IntOut = parseInt64(val)
					}
					if val, ok := metrics.Metrics.Values["1.3.6.1.2.1.2.2.1.14.1"]; ok {
						dev.InErrors = parseInt64(val)
					}
					if val, ok := metrics.Metrics.Values["1.3.6.1.2.1.2.2.1.20.1"]; ok {
						dev.OutErrors = parseInt64(val)
					}
				}
			}

			// Vendor from MAC if missing
			if dev.Vendor == "" && dev.MAC != "" {
				dev.Vendor = lookupVendorFromMAC(dev.MAC)
			}

			snmpNmapChan <- *dev
		}(&allDevices[i])
	}

	go func() {
		wgSNMP.Wait()
		close(snmpNmapChan)
	}()

	allDevices = []sentinel.DeviceRecord{}
	for d := range snmpNmapChan {
		allDevices = append(allDevices, d)
	}

	// Display final table
	sentinel.DisplayTable(allDevices)
}

// Helper for converting string metrics to int64
func parseInt64(val string) int64 {
	var i int64
	fmt.Sscan(val, &i)
	return i
}


// scan common ports quickly
func scanCommonPorts(ip string, timeout time.Duration) []int {
    ports := []int{22, 80, 135, 139, 445, 443, 3389}
    openPorts := []int{}
    for _, port := range ports {
        addr := fmt.Sprintf("%s:%d", ip, port)
        conn, err := net.DialTimeout("tcp", addr, timeout)
        if err == nil {
            openPorts = append(openPorts, port)
            conn.Close()
        }
    }
    return openPorts
}


func guessOS(openPorts []int) string {
    portSet := map[int]bool{}
    for _, p := range openPorts {
        portSet[p] = true
    }
    switch {
    case portSet[135] || portSet[139] || portSet[445] || portSet[3389]:
        return "Windows (likely)"
    case portSet[22]:
        return "Linux/Unix (likely)"
    default:
        return "Unknown"
    }
}
