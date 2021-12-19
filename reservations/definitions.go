package Reservations

// Describes a Consumable
type Consumable struct {
	ID        string
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	Stock     string `json:"Stock"`
	Available string `json:"False"`
}
