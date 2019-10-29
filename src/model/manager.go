package model

// Manager .
type Manager struct {
	Nodes []*NextNode
}

// NextNode .
type NextNode struct {
	ID      int
	Players []*Player
	Actions []*Move
}

// GetWay .
func (m *Manager) GetWay(p *Player, action *Move) {
	node := &NextNode{}
	node.Players = append(node.Players, p)
	node.Actions = append(node.Actions, action)
	switch action.Move[0] {
	case "up":
		id := p.X + (p.Y-1)*Kmap.Width
		node.ID = id
		break
	case "down":
		id := p.X + (p.Y+1)*Kmap.Width
		node.ID = id
		break
	case "left":
		id := p.X + p.Y*Kmap.Width - 1
		node.ID = id
		break
	case "right":
		id := p.X + (p.Y-1)*Kmap.Width + 1
		node.ID = id
		break
	default:
		id := p.X + p.Y*Kmap.Width
		node.ID = id
		break
	}
	if len(m.Nodes) == 0 {
		m.Nodes = append(m.Nodes, node)
		return
	}
	for i, v := range m.Nodes {
		if node.ID == v.ID {
			m.Nodes[i].Players = append(m.Nodes[i].Players, p)
			m.Nodes[i].Actions = append(m.Nodes[i].Actions, action)
			return
		}
	}
	m.Nodes = append(m.Nodes, node)
}

// Think .
func (m *Manager) Think() []*Move {
	var actions []*Move
	for _, v := range m.Nodes {
		if len(v.Players) == 2 {
			action := v.Players[0].ReThink(v.Actions[0].Move[0])
			actions = append(actions, action)
			actions = append(actions, v.Actions[1])
			continue
		}
		if len(v.Players) == 3 {
			action := v.Players[0].ReThink(v.Actions[0].Move[0])
			actions = append(actions, action)
			action = v.Players[1].ReThink(v.Actions[1].Move[0])
			actions = append(actions, action)
			actions = append(actions, v.Actions[2])
			continue
		}
		actions = append(actions, v.Actions[0])
	}
	// log.Println(len(m.Nodes), actions)
	return actions
}
