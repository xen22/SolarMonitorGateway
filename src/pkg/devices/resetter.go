package devices

// Resetter interface can be implemented by devices that can be reset.
type Resetter interface {
	Reset() error
}
