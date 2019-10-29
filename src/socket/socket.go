package socket

import (
	"client/src/conf"
	"client/src/game"
	"client/src/model"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// Client .
type Client struct {
	Conn         net.Conn
	legStartChan chan *model.LegStart
	legEndChan   chan *model.LegEnd
	roundChan    chan *model.Round
	errChan      chan error
	terminalChan chan int
}

var client *Client

func init() {
	client = &Client{}
	client.legStartChan = make(chan *model.LegStart)
	client.legEndChan = make(chan *model.LegEnd)
	client.roundChan = make(chan *model.Round)
	client.errChan = make(chan error, 10)
	client.terminalChan = make(chan int)
}

// NewClient .
func NewClient(address string) (*Client, error) {
	// c.conn, err = net.DialTCP("tcp4", nil, tcpAddr)
	var err error
	client.Conn, err = net.DialTimeout("tcp", address, time.Second*1)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Start .
func (c *Client) Start() {
	if err := registrate(); err != nil {
		log.Println("游戏注册失败！", err)
	}
	go c.Read()

gameLoop:
	for {
		select {
		case start := <-c.legStartChan:
			game.LegStart(start)

		case end := <-c.legEndChan:
			game.LegEnd(end)

		case round := <-c.roundChan:
			msg := game.Round(round)
			_, err := c.Write(msg)
			if err != nil {
				c.errChan <- err
				break
			}

		case <-time.After(time.Millisecond * 750):
			log.Printf("Time OUT!")
		case err := <-c.errChan:
			fmt.Printf("ERROR: %v", err)
		case <-c.terminalChan:
			fmt.Println("GAMEOVER!!!")
			break gameLoop
		}
	}
}

func registrate() error {
	msg := model.Msg{Name: "registration", Data: model.Registration{
		TeamID:   conf.TEAMID,
		TeamName: "are_you_ok",
	}}
	_, err := client.Write(msg)
	return err
}

// GetClient .
func GetClient() *Client {
	return client
}

// Write .
func (c *Client) Write(msg interface{}) (int, error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}
	msgLen := len(msgBytes)
	msgLenStr := strconv.FormatInt(int64(msgLen), 10)
	sizeBytes := make([]byte, 0, msgLen+5)
	for i := 0; i < 5-len(msgLenStr); i++ {
		sizeBytes = append(sizeBytes, '0')
	}
	sizeBytes = append(sizeBytes, []byte(msgLenStr)...)

	sendBytes := append(sizeBytes, msgBytes...)

	log.Println("SEND->", string(msgBytes))

	n, err := c.Conn.Write(sendBytes)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (c *Client) Read() {
	for {
		var buf = make([]byte, 4000)
		n, err := c.Conn.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			continue
		}
		var msgData json.RawMessage
		msg := &model.Msg{Data: &msgData}
		if err := json.Unmarshal(buf[5:n], msg); err != nil {
			log.Printf("Unmarshal err: %v", err)
			continue
		}
		log.Println("RECV->", msg.Name)

		switch msg.Name {
		case "leg_start":
			legStart := new(model.LegStart)
			err := json.Unmarshal(msgData, legStart)
			if err != nil {
				c.errChan <- err
				continue
			}
			c.legStartChan <- legStart
		case "round":
			round := new(model.Round)
			err := json.Unmarshal(msgData, round)
			if err != nil {
				c.errChan <- err
				continue
			}
			c.roundChan <- round
		case "leg_end":
			legEnd := new(model.LegEnd)
			err := json.Unmarshal(msgData, legEnd)
			if err != nil {
				c.errChan <- err
				continue
			}
			c.legEndChan <- legEnd
		case "game_over":
			c.terminalChan <- 1
		default:
			c.errChan <- fmt.Errorf("Unknown msg Name: %v", msg.Name)
		}
	}
}
