package reservation

// Describes a Consumable
type Consumable struct {
	ID        string
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	Stock     string `json:"Stock"`
	Available string `json:"Available"`
}

type User struct {
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Email     string `json:"Email"`
	Admin     string `json:"Admin"`
}

type Equipment struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}

type Reservations struct {
	From               string `json:"From"`
	To                 string `json:"To"`
	Rooms_and_desks_id string `json:"Rooms_and_desks_id"`
}

type RoomsAndDesks struct {
	Wing   string `json:"Wing"`
	Floor  string `json:"Floor"`
	Number string `json:"Number"`
	Places string `json:"Places"`
	Type   string `json:"Type"`
}
