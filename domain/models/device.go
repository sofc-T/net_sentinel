package models

// Device represents a network device such as a router, switch, or firewall
type Device struct {
	id                 string     
	hostname           string     
	ipAddress          string     
	deviceType         string     
	vendor             string     
	status             string     
	monitoringProtocols []string  
	interfaces         []Interface
	macAddress        string 
}


// DeviceConfig represents the configuration of a device
type DeviceConfig struct {
	Hostname           string
	IPAddress          string
	DeviceType         string
	Vendor             string
	Status             string
	MonitoringProtocols []string
	Interfaces         []Interface
	MACAddress        string
}

// NewDeviceConfig creates a new DeviceConfig instance
func NewDevice(device DeviceConfig) *Device {
	return &Device{
		hostname:           device.Hostname,
		ipAddress:          device.IPAddress,
		deviceType:         device.DeviceType,
		vendor:             device.Vendor,
		status:             device.Status,
		monitoringProtocols: device.MonitoringProtocols,
		interfaces:         device.Interfaces,
		macAddress:        device.MACAddress,
	}
}


// SetID sets the ID of the device
func (d *Device) SetID(id string) {
	d.id = id
}

// GetID gets the ID of the device
func (d *Device) GetID() string {
	return d.id
}

// SetHostname sets the hostname of the device
func (d *Device) SetHostname(hostname string) {
	d.hostname = hostname
}

// GetHostname gets the hostname of the device
func (d *Device) GetHostname() string {
	return d.hostname
}

// SetIPAddress sets the IP address of the device
func (d *Device) SetIPAddress(ipAddress string) {
	d.ipAddress = ipAddress
}

// GetIPAddress gets the IP address of the device
func (d *Device) GetIPAddress() string {
	return d.ipAddress
}

// SetDeviceType sets the device type
func (d *Device) SetDeviceType(deviceType string) {
	d.deviceType = deviceType
}

// GetDeviceType gets the device type
func (d *Device) GetDeviceType() string {
	return d.deviceType
}

// SetVendor sets the vendor of the device
func (d *Device) SetVendor(vendor string) {
	d.vendor = vendor
}

// GetVendor gets the vendor of the device
func (d *Device) GetVendor() string {
	return d.vendor
}

// SetStatus sets the status of the device
func (d *Device) SetStatus(status string) {
	d.status = status
}

// GetStatus gets the status of the device
func (d *Device) GetStatus() string {
	return d.status
}

// SetMonitoringProtocols sets the monitoring protocols of the device
func (d *Device) SetMonitoringProtocols(protocols []string) {
	d.monitoringProtocols = protocols
}

// GetMonitoringProtocols gets the monitoring protocols of the device
func (d *Device) GetMonitoringProtocols() []string {
	return d.monitoringProtocols
}

// SetInterfaces sets the interfaces of the device
func (d *Device) SetInterfaces(interfaces []Interface) {
	d.interfaces = interfaces
}

// GetInterfaces gets the interfaces of the device
func (d *Device) GetInterfaces() []Interface {
	return d.interfaces
}

func (d *Device) GetMACAddress() string {
	return d.macAddress
}