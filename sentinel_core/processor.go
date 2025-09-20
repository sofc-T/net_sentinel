package sentinel

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// DeviceRecord holds all collected data for a single device.
type DeviceRecord struct {
	DeviceID   string
	Hostname   string
	IP         string
	MAC        string
	Status     string
	PingMs     int64
	LLDP       string
	CPU        float64
	Mem        float64
	IntIn      int64
	IntOut     int64
	InErrors   int64
	OutErrors  int64
	Uptime     string
	Descr      string
	Type       string
	Vendor     string
	Protocols  string
	LastSeen   time.Time
	SysName    string
}

// Processor stores device data and handles display.
type Processor struct {
	devices map[string]DeviceRecord
}

// NewProcessor creates a new Processor instance.
func NewProcessor() *Processor {
	return &Processor{
		devices: make(map[string]DeviceRecord),
	}
}

// UpdateDevice updates or inserts a device record.
func (p *Processor) UpdateDevice(d DeviceRecord) {
	if existing, ok := p.devices[d.IP]; ok {
		d.LastSeen = time.Now()
		// Preserve missing fields from the previous record.
		if d.DeviceID == "" {
			d.DeviceID = existing.DeviceID
		}
		if d.Hostname == "" {
			d.Hostname = existing.Hostname
		}
	}
	if d.LastSeen.IsZero() {
		d.LastSeen = time.Now()
	}
	p.devices[d.IP] = d
}

// DisplayTable prints all stored device info in a table.
func (p *Processor) DisplayTable() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{
		"DeviceID", "Hostname", "IP", "MAC", "Status", "Ping(ms)", "LLDP", "CPU%", "Mem%",
		"InOctets", "OutOctets", "InErr", "OutErr", "Uptime", "Descr", "Type", "Vendor",
		"Protocols", "SysName", "LastSeen",
	})

	// Sort by IP for consistency
	var ips []string
	for ip := range p.devices {
		ips = append(ips, ip)
	}
	sort.Strings(ips)

	for _, ip := range ips {
		d := p.devices[ip]
		t.AppendRow(table.Row{
			d.DeviceID, d.Hostname, d.IP, d.MAC, d.Status, d.PingMs, d.LLDP, d.CPU, d.Mem,
			d.IntIn, d.IntOut, d.InErrors, d.OutErrors, d.Uptime, d.Descr, d.Type, d.Vendor,
			d.Protocols, d.SysName, d.LastSeen.Format("15:04:05"),
		})
	}

	if len(ips) == 0 {
		fmt.Println("No devices discovered.")
		return
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Descr", WidthMax: 20, Align: text.AlignLeft},
	})
	t.Render()
}

func DisplayTable(devices []DeviceRecord) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = false

	t.AppendHeader(table.Row{
		"DeviceID", "Hostname", "IP", "MAC", "Status", "Ping(ms)", "LLDP", "CPU%", "Mem%",
		"InOctets", "OutOctets", "InErr", "OutErr", "Uptime", "Descr", "Type", "Vendor",
		"Protocols", "SysName", "LastSeen",
	})

	for _, d := range devices {
		t.AppendRow(table.Row{
			d.DeviceID, d.Hostname, d.IP, d.MAC, d.Status, d.PingMs, d.LLDP, d.CPU, d.Mem,
			d.IntIn, d.IntOut, d.InErrors, d.OutErrors, d.Uptime, d.Descr, d.Type, d.Vendor,
			d.Protocols, d.SysName, d.LastSeen.Format("15:04:05"),
		})
	}

	if len(devices) == 0 {
		fmt.Println("No devices discovered.")
		return
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Descr", WidthMax: 20, Align: text.AlignLeft},
	})
	t.Render()
}