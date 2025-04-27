package models

// Network represents a collection of network devices and their relationships
type Network struct {
	id          string
	name        string
	description string
	ipRanges    []string
	subnetMask  string
	devices     []Device
}

// GetID returns the ID of the network
func (n *Network) GetID() string {
	return n.id
}

// SetID sets the ID of the network
func (n *Network) SetID(id string) {
	n.id = id
}

// GetName returns the name of the network
func (n *Network) GetName() string {
	return n.name
}

// SetName sets the name of the network
func (n *Network) SetName(name string) {
	n.name = name
}

// GetDescription returns the description of the network
func (n *Network) GetDescription() string {
	return n.description
}

// SetDescription sets the description of the network
func (n *Network) SetDescription(description string) {
	n.description = description
}

// GetIPRanges returns the IP ranges of the network
func (n *Network) GetIPRanges() []string {
	return n.ipRanges
}

// SetIPRanges sets the IP ranges of the network
func (n *Network) SetIPRanges(ipRanges []string) {
	n.ipRanges = ipRanges
}

// GetSubnetMask returns the subnet mask of the network
func (n *Network) GetSubnetMask() string {
	return n.subnetMask
}

// SetSubnetMask sets the subnet mask of the network
func (n *Network) SetSubnetMask(subnetMask string) {
	n.subnetMask = subnetMask
}

// GetDevices returns the devices in the network
func (n *Network) GetDevices() []Device {
	return n.devices
}

// SetDevices sets the devices in the network
func (n *Network) SetDevices(devices []Device) {
	n.devices = devices
}

