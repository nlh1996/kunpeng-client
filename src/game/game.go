package game

import (
	"client/src/conf"
	"client/src/model"
	"log"
)

// Group .
type Group struct {
	Team    *model.Team
	Players []*model.Player
	Mode    string
	Score   int
	Force   string
}

var (
	myGroup    *Group
	enemyGroup *Group
)

func init() {
	myGroup = &Group{}
	enemyGroup = &Group{}
}

// LegStart .
func LegStart(data *model.LegStart) {
	model.Kmap = &data.Map
	if data.Teams[0].ID == conf.TEAMID {
		myGroup.Team = &data.Teams[0]
		myGroup.Force = data.Teams[0].Force
		enemyGroup.Team = &data.Teams[1]
	} else {
		myGroup.Team = &data.Teams[1]
		myGroup.Force = data.Teams[1].Force
		enemyGroup.Team = &data.Teams[0]
	}
}

// Round .
func Round(data *model.Round) model.Msg {
	myGroup.Mode = data.Mode
	model.RoundID = data.ID
	model.Powers = data.Power
	// 清空数据
	myGroup.Players = append([]*model.Player{})
	enemyGroup.Players = append([]*model.Player{})
	model.Closed = append([]*model.Node{})
	model.UnSafe = append([]*model.Node{})
	// 初始分组
	for i, v := range data.Players {
		if v.Team == myGroup.Team.ID {
			data.Players[i].Mode = data.Mode
			data.Players[i].Force = myGroup.Force
			data.Players[i].Enemys = &model.Enemys{}
			data.Players[i].Teammate = &model.Enemys{}
			data.Players[i].Powers = &model.PowerList{}
			myGroup.Players = append(myGroup.Players, &data.Players[i])
		} else {
			enemyGroup.Players = append(enemyGroup.Players, &data.Players[i])
			data.Players[i].CreateNode()
		}
	}

	// 分配能量
	model.DispatchPower(myGroup.Players)

	// 记录敌方信息
	if len(enemyGroup.Players) != 0 {
		for i, _ := range myGroup.Players {
			myGroup.Players[i].CountEnemys(enemyGroup.Players)
		}
	}

	// 得分情况
	for _, v := range data.Teams {
		if v.ID == myGroup.Team.ID {
			myGroup.Team.Point = v.Point
			myGroup.Team.RemainLife = v.RemainLife
		} else {
			enemyGroup.Team.Point = v.Point
			enemyGroup.Team.RemainLife = v.RemainLife
		}
	}
	return myGroup.action()
}

func (g *Group) action() model.Msg {
	plen := len(g.Players)
	ch := make(chan *model.Move, plen)
	for i := 0; i < plen; i++ {
		go func(p *model.Player) {
			// 记录队友位置
			p.CountTeammate(myGroup.Players)
			// 初始化状态
			p.InitState(ch)
		}(g.Players[i])
	}
	var actions []*model.Move
	manager := &model.Manager{}
	for i := 0; i < plen; i++ {
		action := <-ch
		p := getPlayerByID(action.PlayerID)
		// action拦截
		manager.GetWay(p, action)
	}
	// 管理思考决定返回actions
	actions = manager.Think()
	data := model.Action{ID: model.RoundID, Actions: actions}
	msg := model.Msg{Name: "action", Data: data}

	// 记录当前回合的敌方角色信息，用作下一回合的分析
	model.EnemyPlayers = enemyGroup.Players
	return msg
}

func getPlayerByID(id int) *model.Player {
	for i, _ := range myGroup.Players {
		if id == myGroup.Players[i].ID {
			return myGroup.Players[i]
		}
	}
	return nil
}

// LegEnd .
func LegEnd(data *model.LegEnd) {
	for _, v := range data.Teams {
		if v.ID == myGroup.Team.ID {
			myGroup.Score += v.Point
		} else {
			enemyGroup.Score += v.Point
		}
	}
	log.Println("我方得分：", myGroup.Score, "敌方得分:", enemyGroup.Score)
}
