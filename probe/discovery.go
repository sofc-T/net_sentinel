package probe

import (
	"context"
	"fmt"
	"log"
	"net/netip"

	// "log"
	"net"
	"time"
	"os/exec"
	"bytes"
	"regexp"
	"strings"

	"github.com/Ullaakut/nmap/v2"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/mdlayher/arp"
	"github.com/sofc-t/sentinel/domain/models"
)


func CaptureLLDP(interfaceName string, captureTimeout time.Duration) ([]models.Device, error) {
	var devices []models.Device

	handle, err := pcap.OpenLive(interfaceName, 1600, true, 100*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error opening interface %s: %v", interfaceName, err)
	}
	defer handle.Close()

	if err := handle.SetBPFFilter("ether proto 0x88cc"); err != nil {
		return nil, fmt.Errorf("error setting BPF filter: %v", err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()

	timeout := time.After(captureTimeout)

	LOOP:
		for {
			select {
			case packet, ok := <-packets:
				if !ok {
					log.Println("[LLDP] Packet channel closed.")
					break LOOP
				}
				if packet == nil {
					continue
				}

				fmt.Println("[LLDP] Captured packet:", packet)

				config := models.DeviceConfig{
					Hostname: "lldp-device", // TODO: Parse real hostname from packet
					DeviceType: "switch",
					Status: "active",
					MonitoringProtocols: []string{"LLDP"},
				}
				device := models.NewDevice(config)
				devices = append(devices, *device)

			case <-timeout:
				fmt.Println("[LLDP] Timeout reached, finishing capture.")
				break LOOP
			}
		}

	log.Printf("[LLDP] Capture finished. Found %d device(s).\n", len(devices))
	return devices, nil
}



func ScanIPRange(subnet string) ([]models.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(subnet),
		nmap.WithPingScan(),
		nmap.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating scanner: %v", err)
	}
	result, _, err := scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("scan failed: %v", err)
	}

	var devices []models.Device
	for _, host := range result.Hosts {
		if len(host.Addresses) > 0 {
			// Create a DeviceConfig first
			config := models.DeviceConfig{
				Hostname: "", // Nmap Ping scan may not give hostname directly
				IPAddress: host.Addresses[0].String(),
				DeviceType: "unknown", // deeper scan later
				Vendor: "",
				Status: "active",
				MonitoringProtocols: []string{"Nmap"},
				Interfaces: nil,
			}
			
			// Create Device from DeviceConfig
			device := models.NewDevice(config)
			devices = append(devices, *device)
		}
	}

	log.Println("finished scan ip")
	return devices, nil
}

func ARPScan(interfaceName string) ([]models.Device, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, fmt.Errorf("error finding interface %s: %v", interfaceName, err)
	}

	client, err := arp.Dial(iface)
	if err != nil {
		return nil, fmt.Errorf("error creating ARP client: %v", err)
	}
	defer client.Close()

	var devices []models.Device

	// Define subnet
	subnet := "192.168.107.0/24"
	ipNet, err := netip.ParsePrefix(subnet)
	if err != nil {
		return nil, fmt.Errorf("invalid subnet %s: %v", subnet, err)
	}

	// Channel for IPs
	ipChan := make(chan netip.Addr, 100)
	// Channel for found devices
	deviceChan := make(chan models.Device, 100)

	const (
		numWorkers    = 50
		timeoutPerIP  = 300 * time.Millisecond
	)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for ip := range ipChan {
				if ip.IsMulticast() || ip.IsLinkLocalUnicast() {
					continue
				}
				mac, err := resolveWithTimeout(client, ip, timeoutPerIP)
				if err == nil {
					fmt.Printf("[ARP] Found device: IP=%s, MAC=%s\n", ip, mac)
					config := models.DeviceConfig{
						Hostname: "",
						IPAddress: ip.String(),
						DeviceType: "unknown",
						Vendor: "",
						Status: "active",
						MonitoringProtocols: []string{"ARP"},
						Interfaces: nil,
					}
					device := models.NewDevice(config)
					deviceChan <- *device
				} else {
					fmt.Printf("[ARP] %s not found: %v\n", ip, err)
				}
			}
		}()
	}

	// Feed IPs
	go func() {
		for ip := ipNet.Addr(); ipNet.Contains(ip); ip = ip.Next() {
			ipChan <- ip
		}
		close(ipChan)
	}()

	// Collect devices
	go func() {
		for dev := range deviceChan {
			devices = append(devices, dev)
		}
	}()

	// Wait for all IPs to be scanned
	time.Sleep(time.Duration(len(devices)/numWorkers+1) * timeoutPerIP)

	close(deviceChan)

	return devices, nil
}





func resolveWithTimeout(client *arp.Client, ip netip.Addr, timeout time.Duration) (net.HardwareAddr, error) {
    type result struct {
        mac net.HardwareAddr
        err error
    }

    resultChan := make(chan result, 1)

    go func() {
        mac, err := client.Resolve(ip)
        resultChan <- result{mac: mac, err: err}
    }()

    select {
    case res := <-resultChan:
        return res.mac, res.err
    case <-time.After(timeout):
        return nil, fmt.Errorf("timeout after %v", timeout)
    }
}


func FindDefaultInterfaceAndSubnet() (string, string, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return "", "", fmt.Errorf("failed to list interfaces: %v", err)
    }

    for _, iface := range interfaces {
        if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
            continue // skip down or loopback interfaces
        }

        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }

        for _, addr := range addrs {
            var ipNet *net.IPNet
            switch v := addr.(type) {
            case *net.IPNet:
                ipNet = v
            case *net.IPAddr:
                ipNet = &net.IPNet{IP: v.IP, Mask: v.IP.DefaultMask()}
            }

            if ipNet == nil || ipNet.IP.IsLoopback() {
                continue
            }

            ip4 := ipNet.IP.To4()
            if ip4 == nil {
                continue // not an IPv4 address
            }

            // Found a valid interface
            subnet := fmt.Sprintf("%s/%d", ip4.String(), maskToPrefix(ipNet.Mask))
            return iface.Name, subnet, nil
        }
    }
    return "", "", fmt.Errorf("no active network interface found")
}

func maskToPrefix(mask net.IPMask) int {
    ones, _ := mask.Size()
    return ones
}


func NmapFingerprint(ip string) (string, string) {
	cmd := exec.Command("nmap", "-sV", "-T4", "--open", ip)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		// If nmap fails or is not installed, return empty values
		return "", ""
	}

	output := out.String()

	// Extract open ports and services using regex
	re := regexp.MustCompile(`(?m)^(\d+)/tcp\s+open\s+([\w\-\?\!]+)(?:\s+(.*))?$`)
	matches := re.FindAllStringSubmatch(output, -1)

	if len(matches) == 0 {
		return "", ""
	}

	var protocols []string
	var descParts []string

	for _, m := range matches {
		port := m[1]
		service := m[2]
		info := strings.TrimSpace(m[3])
		if info != "" {
			descParts = append(descParts, port+"/"+service+"("+info+")")
		} else {
			descParts = append(descParts, port+"/"+service)
		}
		protocols = append(protocols, service)
	}

	description := strings.Join(descParts, "; ")
	protoList := strings.Join(protocols, ",")

	return description, protoList
}

