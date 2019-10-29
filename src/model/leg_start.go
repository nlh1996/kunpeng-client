package model

// LegStart .
type LegStart struct {
	Map   Map    `json:"map"`
	Teams []Team `json:"teams"`
}

// Map 地图信息
type Map struct {
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Version  int        `json:"version"`
	Meteor   []Meteor   `json:"meteor"`
	Tunnel   []Tunnel   `json:"tunnel"`
	Wormhole []Wormhole `json:"wormhole"`
}

// Meteor 陨石坐标
type Meteor struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Tunnel 时空隧道的坐标和方向
type Tunnel struct {
	Direction string `json:"direction"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
}

// Wormhole 虫洞的坐标和名称
type Wormhole struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

// Team 各队及Tank ID
type Team struct {
	ID         int    `json:"id"`
	Players    []int  `json:"players"`
	Force      string `json:"force"`
	Point      int    `json:"point"`
	RemainLife int    `json:"remain_life"`
}
