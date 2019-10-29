package model

import (
	"client/src/common"
)

const (
	DEFAULT = 0 // 默认状态
	WARNING = 1 // 警告状态
	DANGER  = 2 // 危险状态
	ATTACK  = 3 // 进攻状态
	SEDUCE  = 4 // 引诱状态
	WORKING = 5 // 得分状态
	LUCK    = 6 // 潜伏状态
)

// Player 包含自己，同时也包含视野范围内的敌人.
type Player struct {
	ID       int `json:"id"`
	Score    int `json:"score"`
	Sleep    int `json:"sleep"`
	Team     int `json:"team"`
	X        int `json:"x"`
	Y        int `json:"y"`
	State    int
	Mode     string
	Force    string
	ch       chan *Move
	Enemys   *Enemys
	Teammate *Enemys
	Powers   *PowerList
}

// PowerList .
type PowerList struct {
	List []Power
}

// Enemys .
type Enemys struct {
	List []Enemy
}

// Enemy .
type Enemy struct {
	Player   *Player
	Distance int
}

// vector 向量坐标.
type vector struct {
	x   int
	y   int
	dir string
}

// CountTeammate .
func (p *Player) CountTeammate(team []*Player) {
	for i, _ := range team {
		if p.ID != team[i].ID {
			e := Enemy{Player: team[i]}
			e.Distance = common.ComputeDistance(p.X, p.Y, team[i].X, team[i].Y)
			p.Teammate.List = append(p.Teammate.List, e)
		}
	}
}

// CountEnemys .
func (p *Player) CountEnemys(enemys []*Player) {
	for i := 0; i < len(enemys); i++ {
		e := Enemy{Player: enemys[i]}
		if common.ComputeDistance(p.X, p.Y, enemys[i].X, enemys[i].Y) <= 8 {
			p.Enemys.List = append(p.Enemys.List, e)
		}
	}
	if p.Force != p.Mode {
		for i, _ := range Closed {
			Closed[i].next()
		}
		for i, v := range Closed {
			if p.X+p.Y*Kmap.Width == v.ID {
				e := Enemy{Player: Closed[i].Player}
				p.Enemys.List = append(p.Enemys.List, e)
			}
		}
		for i, v := range UnSafe {
			if p.X+p.Y*Kmap.Width == v.ID {
				e := Enemy{Player: UnSafe[i].Player}
				p.Enemys.List = append(p.Enemys.List, e)
			}
		}
	}
}

// InitState 初始化状态
func (p *Player) InitState(ch chan *Move) {
	p.ch = ch
	if p.Force != p.Mode && len(p.Enemys.List) >= 1 {
		p.setState(WARNING)
		return
	}
	if len(p.Powers.List) >= 1 {
		p.setState(WORKING)
		return
	}
	if p.Force == p.Mode && p.ID%4 != 0 {
		if RoundID < 140 || RoundID >= 150 {
			p.setState(ATTACK)
			return
		}
	}
	if p.Force != p.Mode {
		if RoundID < 120 {
			p.setState(DEFAULT)
			return
		}
	}
	p.setState(DEFAULT)
}

func (p *Player) setState(state int) {
	p.State = state
	p.exec()
}

func (p *Player) next() []vector {
	vec := vector{}
	var list []vector
	vec.x = p.X
	vec.y = p.Y - 1
	vec.dir = "up"
	list = append(list, vec)

	vec.x = p.X
	vec.y = p.Y + 1
	vec.dir = "down"
	list = append(list, vec)

	vec.x = p.X - 1
	vec.y = p.Y
	vec.dir = "left"
	list = append(list, vec)

	vec.x = p.X + 1
	vec.y = p.Y
	vec.dir = "right"
	list = append(list, vec)
	return list
}

// move 移动
func (p *Player) move(str string) {
	var dir []string
	dir = append(dir, str)
	move := Move{
		Team:     p.Team,
		PlayerID: p.ID,
		Move:     dir,
	}
	p.ch <- &move
}

func (p *Player) exec() {
	switch p.State {
	case WARNING:
		p.avoid()
		break
	case DANGER:
		p.escape()
		break
	case WORKING:
		p.work()
		break
	case ATTACK:
		p.attack()
		break
	case DEFAULT:
		p.explore()
		break
	}
}

// 回避
func (p *Player) avoid() {
	list := p.next()
	var nodes []*Node
	var safe []*Node
	for _, v := range list {
		node := NewNode(v.x, v.y, nil, v.dir)
		if node != nil {
			if node.Type == 上隧道 && node.Dir == "down" {
				continue
			}
			if node.Type == 下隧道 && node.Dir == "up" {
				continue
			}
			if node.Type == 右隧道 && node.Dir == "left" {
				continue
			}
			if node.Type == 左隧道 && node.Dir == "right" {
				continue
			}
			nodes = append(nodes, node)
		}
	}
	for i, v := range nodes {
		for j := 0; j < len(Closed); j++ {
			if v.ID == Closed[j].ID {
				nodes[i].Type = 陨石
			}
		}
		if nodes[i].Type != 陨石 {
			safe = append(safe, nodes[i])
		}
	}
	for _, v := range safe {
		if v.Type == 能量 || v.Type == 虫洞 {
			p.move(v.Dir)
			return
		}
	}
	if len(safe) == 0 {
		p.setState(DANGER)
		return
	}
	if len(safe) == 1 {
		p.move(safe[0].Dir)
		return
	}
	if len(safe) == 2 {
		if safe[0].Type == 空块 && safe[1].Type != 空块 {
			p.move(safe[0].Dir)
			return
		}
		if safe[0].Type != 空块 && safe[1].Type == 空块 {
			p.move(safe[1].Dir)
			return
		}
	}
	for _, v := range safe {
		for _, v2 := range p.Teammate.List {
			if v2.Distance <= 18 {
				if common.ComputeDistance(v.X, v.Y, v2.Player.X, v2.Player.Y) > v2.Distance && v.Type == 空块 {
					p.move(v.Dir)
					return
				}
			}
		}
	}
	if RoundID%2 == 1 {
		var dir string
		if 2*p.X <= Kmap.Width {
			dir = "right"
		} else {
			dir = "left"
		}
		for _, v := range safe {
			if v.Dir == dir && v.Type == 空块 {
				p.move(dir)
				return
			}
		}
	}
	if RoundID%2 == 0 {
		var dir string
		if 2*p.Y <= Kmap.Height {
			dir = "down"
		} else {
			dir = "up"
		}
		for _, v := range safe {
			if v.Dir == dir && v.Type == 空块 {
				p.move(dir)
				return
			}
		}
	}
	if len(safe) == 3 {
		for i, v := range safe {
			if v.Type == 空块 {
				p.move(safe[i].Dir)
				return
			}
		}
	}
	if len(safe) == 4 {
		for i, v := range safe {
			if v.Type == 空块 {
				p.move(safe[i].Dir)
				return
			}
		}
	}
	p.move(safe[0].Dir)
}

// 极限逃生
func (p *Player) escape() {
	list := p.next()
	var nodes []*Node
	var safe []*Node
	for _, v := range list {
		node := NewNode(v.x, v.y, nil, v.dir)
		if node != nil {
			nodes = append(nodes, node)
		}
	}
	for i, v := range nodes {
		for _, v2 := range p.Enemys.List {
			if v.ID == v2.Player.X+Kmap.Width*v2.Player.Y {
				nodes[i].Type = 陨石
			}
		}
	}
	for i, v := range nodes {
		if v.Type != 陨石 {
			safe = append(safe, nodes[i])
		}
	}
	if len(safe) == len(nodes) {
		p.move("")
		return
	}
	var unsafe []*Node
	for i, v := range p.Enemys.List {
		if v.Distance > 1 {
			node := p.Enemys.List[i].Player.maybeMove()
			if node != nil {
				unsafe = append(unsafe, node)
			}
		}
	}
	for i, v := range safe {
		for _, v2 := range unsafe {
			if v.ID == v2.ID {
				safe[i].Type = 陨石
			}
		}
	}
	// 一线生机
	for _, v := range safe {
		if v.Type != 陨石 {
			p.move(v.Dir)
			return
		}
	}
	// 看运气逃生
	if len(safe) > 1 {
		p.move(safe[RoundID%2].Dir)
		return
	}
	// 执行到这里基本必死
	p.move("")
}

// DispatchPower 分配能量
func DispatchPower(players []*Player) {
	for i, v := range Powers {
		var min int
		var p *Player
		for j, v2 := range players {
			dis := common.ComputeDistance(v2.X, v2.Y, v.X, v.Y)
			if min == 0 {
				min = dis
				p = players[j]
			}
			if dis < min {
				min = dis
				p = players[j]
			}
		}
		p.Powers.List = append(p.Powers.List, Powers[i])
	}
}

// 吃豆
func (p *Player) work() {
	var min int
	var tag *Node
	for i := 0; i < len(p.Powers.List); i++ {
		node := AStar(vector{x: p.X, y: p.Y}, vector{x: p.Powers.List[i].X, y: p.Powers.List[i].Y})
		if node != nil {
			if min == 0 {
				min = node.G
				tag = node
			}
			if node.G < min {
				min = node.G
				tag = node
			}
		}
	}
	if tag != nil {
		for {
			if tag.G <= 1 {
				break
			}
			tag = tag.Parent
		}
		p.move(tag.Dir)
	} else {
		p.setState(DEFAULT)
	}
}

// 探索
func (p *Player) explore() {
	index := p.ID % 4
	var node *Node
	// 0，2 钻洞
	if index == 0 || index == 2 {
		if p.X+p.Y*Kmap.Width == Kmap.Wormhole[index].X+Kmap.Wormhole[index].Y*Kmap.Width {
			index++
		}
		if p.Force == p.Mode {
			node = AStar(vector{x: p.X, y: p.Y}, vector{x: Kmap.Wormhole[index].X, y: Kmap.Wormhole[index].Y})
		} else {
			if lastPoint0 == 1 && index == 0 {
				index += lastPoint0
				node = AStar(vector{x: p.X, y: p.Y}, vector{x: Kmap.Wormhole[index].X, y: Kmap.Wormhole[index].Y})
				if node = node.Parent; node != nil {
					if node.G <= 1 {
						lastPoint0 = 0
					}
				}
			}
			if lastPoint0 == 0 && index == 0 {
				node = AStar(vector{x: p.X, y: p.Y}, vector{x: Kmap.Wormhole[index].X, y: Kmap.Wormhole[index].Y})
				if node = node.Parent; node != nil {
					if node.G <= 1 {
						lastPoint0 = 1
					}
				}
			}
			if lastPoint2 == 1 && index == 2 {
				index += lastPoint2
				node = AStar(vector{x: p.X, y: p.Y}, vector{x: Kmap.Wormhole[index].X, y: Kmap.Wormhole[index].Y})
				if node = node.Parent; node != nil {
					if node.G <= 1 {
						lastPoint2 = 0
					}
				}
			}
			if lastPoint2 == 0 && index == 2 {
				node = AStar(vector{x: p.X, y: p.Y}, vector{x: Kmap.Wormhole[index].X, y: Kmap.Wormhole[index].Y})
				if node = node.Parent; node != nil {
					if node.G <= 1 {
						lastPoint2 = 1
					}
				}
			}
		}
		if node != nil {
			for {
				if node.G == 1 {
					break
				}
				if node.Parent != nil {
					node = node.Parent
				} else {
					break
				}
			}
		}
		if node != nil {
			p.move(node.Dir)
			return
		}
	}
	// 1，3 巡逻
	if index == 1 {
		if p.X+p.Y*Kmap.Width == Kmap.Width-5+2*Kmap.Width {
			lastPoint1 = 1
		}
		if p.X+p.Y*Kmap.Width == 2+6*Kmap.Width {
			lastPoint1 = 0
		}
		if lastPoint1 == 0 {
			node = AStar(vector{x: p.X, y: p.Y}, vector{x: Kmap.Width - 5, y: 2})
		} else {
			node = AStar(vector{x: p.X, y: p.Y}, vector{x: 2, y: 6})
		}
		if node != nil {
			for {
				if node.G == 1 {
					break
				}
				if node.Parent != nil {
					node = node.Parent
				} else {
					break
				}
			}
		}
		if node != nil {
			p.move(node.Dir)
			return
		}
	}
	if index == 3 {
		if p.X+p.Y*Kmap.Width == Kmap.Width-3+(Kmap.Height-3)*Kmap.Width {
			lastPoint3 = 1
		}
		if p.X+p.Y*Kmap.Width == 4+(Kmap.Height-4)*Kmap.Width {
			lastPoint3 = 0
		}
		if lastPoint3 == 0 {
			node = AStar(vector{x: p.X, y: p.Y}, vector{x: Kmap.Width - 3, y: Kmap.Height - 3})
		} else {
			node = AStar(vector{x: p.X, y: p.Y}, vector{x: 4, y: Kmap.Height - 4})
		}
		if node != nil {
			for {
				if node.G == 1 {
					break
				}
				if node.Parent != nil {
					node = node.Parent
				} else {
					break
				}
			}
		}
		if node != nil {
			p.move(node.Dir)
			return
		}
	}

	// 异常情况移动
	if RoundID%3 == 0 {
		p.move("left")
		return
	}
	if RoundID%3 == 1 {
		p.move("up")
		return
	}
	if RoundID%3 == 2 {
		p.move("right")
		return
	}
}

// 进攻
func (p *Player) attack() {
	if p.ID%4 == 0 {
		p.setState(DEFAULT)
		return
	}
	if len(p.Enemys.List) != 0 {
		p.runAndHit(p.Enemys.List[0].Player.X, p.Enemys.List[0].Player.Y)
	} else {
		p.setState(DEFAULT)
	}
}

func (p *Player) runAndHit(x int, y int) {
	node := AStar(vector{x: p.X, y: p.Y}, vector{x: x, y: y})
	if node != nil {
		for {
			if node.G == 1 {
				break
			}
			if node.Parent != nil {
				node = node.Parent
			} else {
				break
			}
		}
	}
	if node != nil {
		p.move(node.Dir)
		return
	}
	p.setState(WORKING)
}

// CreateNode .
func (p *Player) CreateNode() {
	list := p.next()
	for _, v := range list {
		node := NewNode(v.x, v.y, nil, v.dir)
		if node != nil {
			for node.Type == 上隧道 {
				node = NewNode(node.X, node.Y-1, nil, node.Dir)
			}
			for node.Type == 下隧道 {
				node = NewNode(node.X, node.Y+1, nil, node.Dir)
			}
			for node.Type == 左隧道 {
				node = NewNode(node.X-1, node.Y, nil, node.Dir)
			}
			for node.Type == 右隧道 {
				node = NewNode(node.X+1, node.Y, nil, node.Dir)
			}
			node.Player = p
			Closed = append(Closed, node)
		}
	}
	// 当前位置
	node := NewNode(p.X, p.Y, nil, "")
	node.Player = p
	Closed = append(Closed, node)
}

// ReThink .
func (p *Player) ReThink(str string) *Move {
	var dir []string
	switch p.State {
	case WARNING:
		list := p.next()
		var nodes []*Node
		var safe []*Node
		for _, v := range list {
			node := NewNode(v.x, v.y, nil, v.dir)
			if node != nil {
				nodes = append(nodes, node)
			}
		}
		for i, v := range nodes {
			for j := 0; j < len(Closed); j++ {
				if v.ID == Closed[j].ID {
					nodes[i].Type = 陨石
				}
			}
			if nodes[i].Type != 陨石 {
				safe = append(safe, nodes[i])
			}
		}
		for _, v := range safe {
			if v.Dir != str {
				dir = append(dir, v.Dir)
			}
		}
		if len(safe) == 1 {
			dir = append(dir, str)
		}
		break
	case LUCK:
		dir = append(dir, str)
		break
	case ATTACK:
		list := p.next()
		var nodes []*Node
		for _, v := range list {
			node := NewNode(v.x, v.y, nil, v.dir)
			if node != nil {
				nodes = append(nodes, node)
			}
		}
		for _, v := range nodes {
			if v.Dir != str {
				dir = append(dir, v.Dir)
			}
		}
		break
	}
	return &Move{
		Team:     p.Team,
		PlayerID: p.ID,
		Move:     dir,
	}
}

// 敌方角色下一回合可能移动的位置
func (p *Player) maybeMove() *Node {
	for _, v := range EnemyPlayers {
		if v.ID == p.ID {
			if p.X > v.X {
				return NewNode(p.X+1, p.Y, nil, "")
			}
			if p.X < v.X {
				return NewNode(p.X-1, p.Y, nil, "")
			}
			if p.Y > v.Y {
				return NewNode(p.X, p.Y+1, nil, "")
			}
			if p.Y < v.Y {
				return NewNode(p.X, p.Y-1, nil, "")
			}
		}
	}
	return nil
}
