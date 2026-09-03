// adapters/driving/dto/beacon.go (100行以下 - SPEC-PRINCIPLE-001)
package dto

// BeaconStatus represents the status of an external service monitored by the PS1 conductor.
type BeaconStatus string

const (
	BeaconStatusReady   BeaconStatus = "ready"
	BeaconStatusBusy    BeaconStatus = "busy"
	BeaconStatusStopped BeaconStatus = "stopped"
)

// BeaconRequestDTO is the payload sent from the PS1 conductor script to the internal Go API.
type BeaconRequestDTO struct {
	Service string       `json:"service"` // "stash", "thunder", "motrix"
	Status  BeaconStatus `json:"status"`
	Details string       `json:"details,omitempty"`
}
