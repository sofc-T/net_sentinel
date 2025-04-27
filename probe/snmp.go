package probe

import (
	"fmt"
	// "log"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/sofc-t/sentinel/domain/models"
)

// SNMPConfig holds the settings for SNMP communication
type SNMPConfig struct {
	Target    string
	Port      uint16
	Version   gosnmp.SnmpVersion
	Community string  // For SNMP v2c
	Timeout   time.Duration
	Retries   int
}

// NewSNMPClient initializes an SNMP client
func NewSNMPClient(config SNMPConfig) *gosnmp.GoSNMP {
	return &gosnmp.GoSNMP{
		Target:    config.Target,
		Port:      config.Port,
		Community: config.Community,
		Version:   config.Version,
		Timeout:   config.Timeout,
		Retries:   config.Retries,
		MaxOids:   gosnmp.MaxOids, // Ensures batch querying
	}
}

// FetchMetrics queries SNMP for device metrics
func FetchMetrics(config SNMPConfig, oids []string) (*models.SNMPResult, error) {
	client := NewSNMPClient(config)
	err := client.Connect()
	if err != nil {
		return nil, fmt.Errorf("SNMP connection error: %v", err)
	}
	defer client.Conn.Close()

	// Perform SNMP GET request
	result, err := client.Get(oids)
	if err != nil {
		return nil, fmt.Errorf("SNMP GET failed: %v", err)
	}

	// Process SNMP response
	metrics := make(map[string]string)
	for _, v := range result.Variables {
		metrics[v.Name] = fmt.Sprintf("%v", v.Value)
	}

	return &models.SNMPResult{
		DeviceID:  config.Target,
		Metrics: &models.SNMPMetric{Values: metrics}, 
	}, nil
}

// Example Usage
// func main() {
// 	config := SNMPConfig{
// 		Target:    "192.168.1.1",
// 		Port:      161,
// 		Version:   gosnmp.Version2c,
// 		Community: "public",
// 		Timeout:   2 * time.Second,
// 		Retries:   3,
// 	}

// 	oids := []string{
// 		"1.3.6.1.2.1.1.3.0",  // System uptime
// 		"1.3.6.1.2.1.2.2.1.10.2", // Interface in-octets
// 		"1.3.6.1.2.1.2.2.1.16.2", // Interface out-octets
// 	}

// 	result, err := FetchMetrics(config, oids)
// 	if err != nil {
// 		log.Fatalf("SNMP Error: %v", err)
// 	}

// 	fmt.Printf("SNMP Metrics from %s: %+v\n", config.Target, result)
// }
