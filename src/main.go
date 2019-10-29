package main

import (
	"client/src/conf"
	"client/src/socket"
	"flag"
	"log"
	"time"
)

func init() {
	flag.IntVar(&conf.TEAMID, "team", 10000, "TeamID")
	flag.StringVar(&conf.IP, "ip", "127.0.0.1", "ServerIP")
	flag.StringVar(&conf.PORT, "port", "6001", "Server Port")
}

func main() {
	flag.Parse()
	address := conf.IP + ":" + conf.PORT
	var (
		client *socket.Client
		err    error
	)
	for i := 0; i < 30; i++ {
		if client, err = socket.NewClient(address); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Println("连接失败！")
		return
	}
	log.Printf("%v join KunPengBattle (%v: %v): \n", conf.TEAMID, conf.IP, conf.PORT)
	client.Start()
}
