package model

// Node A*算法移动的节点.
type Node struct {
	ID     int
	X      int
	Y      int
	G      int    // 移动量
	H      int    // 估值
	F      int    // 和值
	Dir    string // 移动方向
	Parent *Node
	Type   int
	Player *Player
}

const (
	空块  = 0
	能量  = 1
	陨石  = 2
	虫洞  = 3
	上隧道 = 4
	下隧道 = 5
	左隧道 = 6
	右隧道 = 7
)

// AStar A*寻路算法。
func AStar(start vector, end vector) *Node {
	var (
		open   []*Node
		closed []*Node
	)
	// 开始点
	node := NewNode(start.x, start.y, nil, "")
	if node == nil {
		return nil
	}
	open = append(open, node)
	node.Next(&open, &closed)
	// 开启循环，直到找到终点。
	for i := 0; i < 2000; i++ {
		var temp []*Node
		// 计算和值
		for i, v := range open {
			open[i].H = computeH(v.X, v.Y, end.x, end.y)
			open[i].F = v.G + v.H
			if v.Type == 上隧道 || v.Type == 下隧道 || v.Type == 左隧道 || v.Type == 右隧道 {
				open[i].F += 20
			}
			if v.Type == 虫洞 {
				open[i].F += 2
			}
		}
		// 找到最小和
		var min int
		for _, v := range open {
			if min == 0 {
				min = v.F
			}
			if v.F < min {
				min = v.F
			}
		}
		// 将最小和节点加入closed表，如果没有到达终点继续寻找
		for i, v := range open {
			if v.F == min {
				closed = append(closed, open[i])
				if open[i].ID == (end.x + end.y*Kmap.Width) {
					return open[i]
				}
				open[i].Next(&temp, &closed)
			} else {
				temp = append(temp, open[i])
			}
		}
		open = temp
	}
	return nil
}

// NewNode .
func NewNode(x int, y int, parent *Node, dir string) *Node {
	if x < 0 || x >= Kmap.Width || y < 0 || y >= Kmap.Height {
		return nil
	}
	id := x + y*Kmap.Width
	t := 空块
	for _, v := range Kmap.Meteor {
		if id == v.X+v.Y*Kmap.Width {
			t = 陨石
			return nil
		}
	}
	for _, v := range Kmap.Wormhole {
		if id == v.X+v.Y*Kmap.Width {
			t = 虫洞
		}
	}
	for _, v := range Powers {
		if id == v.X+v.Y*Kmap.Width {
			t = 能量
		}
	}
	for _, v := range Kmap.Tunnel {
		if id == v.X+v.Y*Kmap.Width {
			switch v.Direction {
			case "up":
				t = 上隧道
				break
			case "down":
				t = 下隧道
				break
			case "left":
				t = 左隧道
				break
			case "right":
				t = 右隧道
				break
			}
		}
	}
	var steps int
	if parent != nil {
		steps = parent.G + 1
	}
	return &Node{ID: id, Type: t, X: x, Y: y, Parent: parent, G: steps, Dir: dir}
}

// Next .
func (n *Node) Next(open *[]*Node, closed *[]*Node) {
	var (
		x   int
		y   int
		dir string
		tag bool
	)
	x = n.X
	y = n.Y - 1
	dir = "up"
	tag = true
	for _, v := range *closed {
		if x+y*Kmap.Width == v.ID {
			tag = false
			break
		}
	}
	for _, v := range *open {
		if x+y*Kmap.Width == v.ID {
			tag = false
			break
		}
	}
	if tag {
		if newNode := NewNode(x, y, n, dir); newNode != nil {
			newNode.isOpen(open, closed)
		}
	}
	x = n.X
	y = n.Y + 1
	dir = "down"
	tag = true
	for _, v := range *closed {
		if x+y*Kmap.Width == v.ID {
			tag = false
			break
		}
	}
	for _, v := range *open {
		if x+y*Kmap.Width == v.ID {
			tag = false
			break
		}
	}
	if tag {
		if newNode := NewNode(x, y, n, dir); newNode != nil {
			newNode.isOpen(open, closed)
		}
	}

	x = n.X - 1
	y = n.Y
	dir = "left"
	tag = true
	for _, v := range *closed {
		if x+y*Kmap.Width == v.ID {
			tag = false
			break
		}
	}
	for _, v := range *open {
		if x+y*Kmap.Width == v.ID {
			tag = false
			break
		}
	}
	if tag {
		if newNode := NewNode(x, y, n, dir); newNode != nil {
			newNode.isOpen(open, closed)
		}
	}

	x = n.X + 1
	y = n.Y
	dir = "right"
	tag = true
	for _, v := range *closed {
		if x+y*Kmap.Width == v.ID {
			tag = false
			break
		}
	}
	for _, v := range *open {
		if x+y*Kmap.Width == v.ID {
			tag = false
			break
		}
	}
	if tag {
		if newNode := NewNode(x, y, n, dir); newNode != nil {
			newNode.isOpen(open, closed)
		}
	}
}

func (n *Node) isOpen(open *[]*Node, closed *[]*Node) {
	tag := true
	switch n.Type {
	case 陨石:
		tag = false
		break
	case 上隧道:
		if n.Dir == "down" {
			tag = false
		}
		break
	case 下隧道:
		if n.Dir == "up" {
			tag = false
		}
		break
	case 左隧道:
		if n.Dir == "right" {
			tag = false
		}
		break
	case 右隧道:
		if n.Dir == "left" {
			tag = false
		}
		break
	default:
		tag = true
		break
	}
	if tag {
		*open = append(*open, n)
	}
}

func computeH(x1 int, y1 int, x2 int, y2 int) int {
	var h int
	if x1 >= x2 {
		h += x1 - x2
	} else {
		h += x2 - x1
	}
	if y1 >= y2 {
		h += y1 - y2
	} else {
		h += y2 - y1
	}
	return h
}

func (n *Node) next() {
	if n.Dir != "up" {
		if node := NewNode(n.X, n.Y+1, nil, "down"); node != nil {
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
			node.Player = n.Player
			UnSafe = append(UnSafe, node)
		}
	}
	if n.Dir != "down" {
		if node := NewNode(n.X, n.Y-1, nil, "up"); node != nil {
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
			node.Player = n.Player
			UnSafe = append(UnSafe, node)
		}	
	}
	if n.Dir != "left" {
		if node := NewNode(n.X+1, n.Y, nil, "right"); node != nil {
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
			node.Player = n.Player
			UnSafe = append(UnSafe, node)
		}
	}
	if n.Dir != "right" {
		if node := NewNode(n.X-1, n.Y, nil, "left"); node != nil {
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
			node.Player = n.Player
			UnSafe = append(UnSafe, node)
		}
	}
}
