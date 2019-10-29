package model

// Action .
type Action struct {
	ID      int    `json:"round_id"`
	Actions []*Move `json:"actions"`
}

// Move .
type Move struct {
	Team     int      `json:"team"`
	PlayerID int      `json:"player_id"`
	Move     []string `json:"move"` //每步移动方向只能为up, down, right, left，每回合最多一步。不动为空[]
}
