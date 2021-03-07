package server

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/fzzy/radix/redis"
	"time"
)

var (
	RedisCli *redis.Client
	svr *Server
	ticker *time.Timer
)

func main()  {
	l4g.LoadConfiguration("../config/log.xml")
	if cli, err := redis.DialTimeout("tcp", "127.0.0.1:7777", time.Duration(10)*time.Second); err == nil {
		l4g.Info("redis:%s connected", "127.0.0.1:7777")
		RedisCli = cli
	} else {
		l4g.Info("redis err:%s", err)
	}
	defer RedisCli.Close()
	ticker = time.NewTimer(30 * time.Second)

}

func Run()  {
	for {
		select {
		case msg := <-svr.inChannel:
			
		}
	}
}
