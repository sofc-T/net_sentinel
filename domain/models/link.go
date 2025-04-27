package models

// Link represents a connection between two devices
type Link struct {
	id                   string // ID of the link
	sourceDevice         string // ID of the source device
	destinationDevice    string // ID of the destination device
	sourceInterface      string // Interface ID on source device
	destinationInterface string // Interface ID on destination device
	status               string // Up, Down
}

// GetID returns the ID of the link
func (l *Link) GetID() string {
	return l.id
}

// SetID sets the ID of the link
func (l *Link) SetID(id string) {
	l.id = id
}

// GetSourceDevice returns the source device ID
func (l *Link) GetSourceDevice() string {
	return l.sourceDevice
}

// SetSourceDevice sets the source device ID
func (l *Link) SetSourceDevice(sourceDevice string) {
	l.sourceDevice = sourceDevice
}

// GetDestinationDevice returns the destination device ID
func (l *Link) GetDestinationDevice() string {
	return l.destinationDevice
}

// SetDestinationDevice sets the destination device ID
func (l *Link) SetDestinationDevice(destinationDevice string) {
	l.destinationDevice = destinationDevice
}

// GetSourceInterface returns the source interface ID
func (l *Link) GetSourceInterface() string {
	return l.sourceInterface
}

// SetSourceInterface sets the source interface ID
func (l *Link) SetSourceInterface(sourceInterface string) {
	l.sourceInterface = sourceInterface
}

// GetDestinationInterface returns the destination interface ID
func (l *Link) GetDestinationInterface() string {
	return l.destinationInterface
}

// SetDestinationInterface sets the destination interface ID
func (l *Link) SetDestinationInterface(destinationInterface string) {
	l.destinationInterface = destinationInterface
}

// GetStatus returns the status of the link
func (l *Link) GetStatus() string {
	return l.status
}

// SetStatus sets the status of the link
func (l *Link) SetStatus(status string) {
	l.status = status
}
