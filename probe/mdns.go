package probe

import (
	"context"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/sofc-t/sentinel/domain/models"
)

// CaptureMDNS discovers devices broadcasting mDNS/Bonjour services
func CaptureMDNS(timeout time.Duration) ([]models.Device, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	devices := []models.Device{}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			device := models.NewDevice(models.DeviceConfig{
				Hostname: entry.HostName,
				IPAddress: "", // Can fill from entry.AddrIPv4 later
				DeviceType: "unknown",
				Vendor: "",
				Status: "active",
				MonitoringProtocols: []string{"mDNS"},
			})
			devices = append(devices, *device)
		}
	}(entries)

	if err := resolver.Browse(ctx, "_services._dns-sd._udp", "local.", entries); err != nil {
		log.Println("mDNS browse error:", err)
	}

	<-ctx.Done()
	close(entries)

	log.Printf("[mDNS] Found %d device(s)\n", len(devices))
	return devices, nil
}
