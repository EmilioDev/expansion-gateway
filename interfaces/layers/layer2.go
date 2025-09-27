// file: /interfaces/layers/layer2.go
package layers

type Layer2 interface {
	SmartLayer
	GetActiveSessions() int32
	HasSession(int64) bool
}
