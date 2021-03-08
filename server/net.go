package server

import (
	"chatServer/common"
	"net"
	"os"
	"time"

	l4g "github.com/alecthomas/log4go"
)

type Server struct {
	Addr      string
	CurConnId uint32
	CurUserId uint32
	MaxMsgBuf uint32
	inChannel chan *UserMsg
}

func (this *Server) Listen() {
	listener, err := net.Listen("tcp", this.Addr)
	if err != nil {
		l4g.Error(err)
		os.Exit(1)
	}
	defer listener.Close()
	l4g.Info("listen successed")
	for {
		conn, err := listener.Accept()
		if err != nil {
			l4g.Error(err)
			os.Exit(0)
		}
		this.CurConnId++
		user := User{id: this.CurConnId, Conn: conn, connid: this.CurConnId}
		l4g.Info(">>>>> new connection %s >>>>>", conn.RemoteAddr())
		go this.handleConn(&user)
	}
}

func (this *Server) handleConn(pUser *User) {
	defer func() {
		var tmpUser User = *pUser
		tmpUser.id = tmpUser.connid
		l4g.Info("id: %d, connid: %d", pUser.id, pUser.connid)
		pUser.Close()
	}()
	readBuf := make([]byte, this.MaxMsgBuf)
	for {
		pUser.SetDeadline(time.Now().Add(time.Duration(30) * time.Second))
		n, err := pUser.Read(readBuf)
		if err != nil {
			l4g.Error(err)
			break
		}
		startPos := 0
		for n >= 16 {
			headers, err := common.DealHeader(readBuf[startPos:])
			if err != nil {
				l4g.Error(err)
				return
			}
			pUser.sid = headers[3]
			pkgLen := int(headers[0])
			pb := readBuf[16:pkgLen]
			um := &UserMsg{
				pid:  headers[1],
				msg:  pb,
				user: pUser,
			}
			this.inChannel <- um
			l4g.Info("msg push chanlen:%d %d", len(this.inChannel), cap(this.inChannel))
			n = n - int(pkgLen)
			startPos += int(pkgLen)
		}
	}
}
