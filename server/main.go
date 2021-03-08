package server

import (
	"strconv"
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/fzzy/radix/redis"
)

var (
	RedisCli *redis.Client
	svr      *Server
	ticker   *time.Timer
)

func main() {
	l4g.LoadConfiguration("../config/log.xml")
	if cli, err := redis.DialTimeout("tcp", "127.0.0.1:7777", time.Duration(10)*time.Second); err == nil {
		l4g.Info("redis:%s connected", "127.0.0.1:7777")
		RedisCli = cli
	} else {
		RedisErrHandler(err)
	}
	defer RedisCli.Close()
	ProtoRegister()
	ticker = time.NewTimer(30 * time.Second)
	svr = &Server{
		Addr:      "0.0.0.0:8000",
		CurConnId: 0,
		MaxMsgBuf: 10240,
		inChannel: make(chan *UserMsg, 8192),
	}
	go svr.Listen()
	Run()
}

func RedisErrHandler(err error) bool {
	if err != nil {
		l4g.Error("redis err:%s", err)
		return false
	}
	return true
}

func SvrReadConfig() {
	r, err := RedisCli.Cmd("hgetall", KeyConfig).Hash()
	if RedisErrHandler(err) {

	} else {
		id, _ := r[FiledCurUID]
		nId, _ := strconv.Atoi(id)
		svr.CurUserId = uint32(nId)
	}
}

func SvrConfigInit() {
	l4g.Info("svr config install")
	r, err := RedisCli.Cmd("hset", KeyConfig, FiledCurUID, 10000).Bool()
	if RedisErrHandler(err) {
		if r {
			l4g.Info("install config successed")
		} else {
			l4g.Error("install config failed")
		}
	}
}

func Run() {
	for {
		l4g.Info("loop:%d cap:%d", len(svr.inChannel), cap(svr.inChannel))
		select {
		case msg := <-svr.inChannel:
			ProtoDispatcher(msg)
		case tm := <-ticker.C:
			ticker.Reset(30 * time.Second)
			l4g.Info("[TIMER EVENT][%s] ONLINE:%d, CONN:%d", tm, len(OnlineUsers.UserPool), len(ConnUsers.UserPool))
		}
	}
}
