package state

// SD for VM
type SD interface {
	Data(f, t int) []byte
}
