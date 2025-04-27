package models

// Interface represents a network interface on a device
type Interface struct {
	id                 string // ID of the interface
	deviceID           string // ID of the device
	macAddress         string // MAC address of the interface
	ipAddress          string // Optional IP Address
	status             string // Up, Down
	speed              string // e.g., 1Gbps, 10Gbps
	connectedDevice    string // ID of the connected device
	connectedInterface string // ID of the connected interface
}

// SetID sets the ID of the interface
func (i *Interface) SetID(id string) {
	i.id = id
}

// GetID gets the ID of the interface
func (i *Interface) GetID() string {
	return i.id
}

// SetDeviceID sets the DeviceID of the interface
func (i *Interface) SetDeviceID(deviceID string) {
	i.deviceID = deviceID
}

// GetDeviceID gets the DeviceID of the interface
func (i *Interface) GetDeviceID() string {
	return i.deviceID
}

// SetMACAddress sets the MACAddress of the interface
func (i *Interface) SetMACAddress(macAddress string) {
	i.macAddress = macAddress
}

// GetMACAddress gets the MACAddress of the interface
func (i *Interface) GetMACAddress() string {
	return i.macAddress
}

// SetIPAddress sets the IPAddress of the interface
func (i *Interface) SetIPAddress(ipAddress string) {
	i.ipAddress = ipAddress
}

// GetIPAddress gets the IPAddress of the interface
func (i *Interface) GetIPAddress() string {
	return i.ipAddress
}

// SetStatus sets the Status of the interface
func (i *Interface) SetStatus(status string) {
	i.status = status
}

// GetStatus gets the Status of the interface
func (i *Interface) GetStatus() string {
	return i.status
}

// SetSpeed sets the Speed of the interface
func (i *Interface) SetSpeed(speed string) {
	i.speed = speed
}

// GetSpeed gets the Speed of the interface
func (i *Interface) GetSpeed() string {
	return i.speed
}

// SetConnectedDevice sets the ConnectedDevice of the interface
func (i *Interface) SetConnectedDevice(connectedDevice string) {
	i.connectedDevice = connectedDevice
}

// GetConnectedDevice gets the ConnectedDevice of the interface
func (i *Interface) GetConnectedDevice() string {
	return i.connectedDevice
}

// SetConnectedInterface sets the ConnectedInterface of the interface
func (i *Interface) SetConnectedInterface(connectedInterface string) {
	i.connectedInterface = connectedInterface
}

// GetConnectedInterface gets the ConnectedInterface of the interface
func (i *Interface) GetConnectedInterface() string {
	return i.connectedInterface
}