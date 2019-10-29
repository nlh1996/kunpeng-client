package model

// Round .
type Round struct {
	ID      int      `json:"round_id"`
	Mode    string   `json:"mode"` // "mode"表示本回合优势的能力
	Power   []Power  `json:"power"`
	Players []Player `json:"players"`
	Teams   []Team   `json:"teams"`
}

// Power 只包含视野范围内的矿点.
type Power struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Point int `json:"point"`
}


