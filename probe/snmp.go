package probe

import (
	"fmt"
	"log"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/sofc-t/sentinel/domain/models"
)

// SNMPConfig holds settings for SNMP communication.
type SNMPConfig struct {
	Target    string
	Port      uint16
	Version   gosnmp.SnmpVersion
	Community string
	Timeout   time.Duration
	Retries   int
}

// NewSNMPClient initializes an SNMP client.
func NewSNMPClient(cfg SNMPConfig) *gosnmp.GoSNMP {
	return &gosnmp.GoSNMP{
		Target:    cfg.Target,
		Port:      cfg.Port,
		Community: cfg.Community,
		Version:   cfg.Version,
		Timeout:   cfg.Timeout,
		Retries:   cfg.Retries,
		MaxOids:   gosnmp.MaxOids,
	}
}

// FetchMetrics queries SNMP for a list of OIDs and returns results as a map.
func FetchMetrics(cfg SNMPConfig, oids []string) (*models.SNMPResult, error) {
	client := NewSNMPClient(cfg)

	// Establish connection
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("[SNMP] connection failed for %s: %v", cfg.Target, err)
	}
	defer client.Conn.Close()

	// Perform GET request
	pdu, err := client.Get(oids)
	if err != nil {
		return nil, fmt.Errorf("[SNMP] GET failed for %s: %v", cfg.Target, err)
	}

	metrics := make(map[string]string)
	for _, variable := range pdu.Variables {
		metrics[variable.Name] = fmt.Sprintf("%v", variable.Value)
	}

	return &models.SNMPResult{
		DeviceID: cfg.Target,
		Metrics:  &models.SNMPMetric{Values: metrics},
	}, nil
}

// FetchCommonDeviceMetrics retrieves uptime, CPU, and memory utilization if available.
func FetchCommonDeviceMetrics(cfg SNMPConfig) (uptime string, cpu, mem float64) {
	commonOIDs := map[string]string{
		"sysUpTime": ".1.3.6.1.2.1.1.3.0",  // Uptime
		"cpuLoad":   ".1.3.6.1.4.1.2021.10.1.3.1", // CPU (example for UCD-SNMP)
		"memAvail":  ".1.3.6.1.4.1.2021.4.6.0",   // Memory (example for UCD-SNMP)
	}

	res, err := FetchMetrics(cfg, []string{
		commonOIDs["sysUpTime"], commonOIDs["cpuLoad"], commonOIDs["memAvail"],
	})
	if err != nil {
		log.Printf("[SNMP] Failed to fetch common metrics from %s: %v", cfg.Target, err)
		return "", 0, 0
	}

	if val, ok := res.Metrics.Values[commonOIDs["sysUpTime"]]; ok {
		uptime = val
	}
	if val, ok := res.Metrics.Values[commonOIDs["cpuLoad"]]; ok {
		fmt.Sscanf(val, "%f", &cpu)
	}
	if val, ok := res.Metrics.Values[commonOIDs["memAvail"]]; ok {
		fmt.Sscanf(val, "%f", &mem)
	}

	return
}

// BulkWalkMetrics performs a BULK WALK for a base OID, useful for interfaces or routing tables.
func BulkWalkMetrics(cfg SNMPConfig, baseOID string) (map[string]string, error) {
	client := NewSNMPClient(cfg)
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("[SNMP] BulkWalk connection failed for %s: %v", cfg.Target, err)
	}
	defer client.Conn.Close()

	metrics := make(map[string]string)
	err := client.BulkWalk(baseOID, func(pdu gosnmp.SnmpPDU) error {
		metrics[pdu.Name] = fmt.Sprintf("%v", pdu.Value)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("[SNMP] BulkWalk failed: %v", err)
	}

	return metrics, nil
}
