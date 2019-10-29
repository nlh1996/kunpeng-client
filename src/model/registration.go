package model

// Msg .
type Msg struct {
	Name string      `json:"msg_name"`
	Data interface{} `json:"msg_data"`
}

// Registration .
type Registration struct {
	TeamID   int    `json:"team_id"`
	TeamName string `json:"team_name"`
}
