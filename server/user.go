package server

import (
	"bytes"
	"chatServer/common"
	"fmt"
	"net"

	l4g "github.com/alecthomas/log4go"
	"github.com/golang/protobuf/proto"
)

type User struct {
	connid, id, sid uint32
	nickname        string
	net.Conn
	status                            uint32
	lastLoginTime, regTime, leaveTime uint64
}

type UserManager struct {
	UserPool map[uint32]*User
}

type UserMsg struct {
	pid  uint32
	msg  []byte
	user *User
}

var (
	OnlineUsers = UserManager{
		UserPool: make(map[uint32]*User),
	}
	ConnUsers = UserManager{
		UserPool: make(map[uint32]*User),
	}
)

func (conn *User) Send(msg proto.Message, pid uint32) error {
	if _, has := ProtoMap[pid]; !has {
		l4g.Error("protocol id error:%d", pid)
		err := fmt.Errorf("protocol id error:%d", pid)
		return err
	}
	buf, err := proto.Marshal(msg)
	if err != nil {
		l4g.Error("Send Proto Marahal failed:%s", err)
		return err
	}
	headerBuf := common.CreateHeader(buf, pid, conn.id, conn.sid)
	var sendBuf bytes.Buffer
	sendBuf.Write(headerBuf)
	sendBuf.Write(buf)
	conn.Write(sendBuf.Bytes())
	return nil
}

func (conn *User) SendAll(msg proto.Message, pid uint32) error {
	for k, v := range OnlineUsers.UserPool {
		if k != conn.id {
			if err := v.Send(msg, pid); err != nil {
				return err
			}
		}
	}
	return nil
}

func (conn *User) SendTo(msg proto.Message, pid uint32, targetId uint32) error {
	if pUser, has := OnlineUsers.UserPool[targetId]; has {
		return pUser.Send(msg, pid)
	} else {
		l4g.Error("%d send to %d failed", conn.id, targetId)
	}
	return nil
}

func (pUM *UserManager) Append(pU *User) error {
	if _, has := pUM.UserPool[pU.id]; has {
		return fmt.Errorf("id:%d,addr:%s,already had", pU.id, pU.RemoteAddr().Network())
	}
	pUM.UserPool[pU.id] = pU
	return nil
}

func (pUM *UserManager) Delete(pU *User) error {
	if _, has := pUM.UserPool[pU.id]; !has {
		return fmt.Errorf("delete id:%d, addr: %s not found", pU.id, pU.RemoteAddr().Network())
	}
	delete(pUM.UserPool, pU.id)
	return nil
}

func (pUM *UserManager) getOnlineList() []uint32 {
	var idList []uint32
	for k := range pUM.UserPool {
		idList = append(idList, k)
	}
	return idList
}

func (pUM *UserManager) getUserByName(id uint32) (*User, error) {
	if pUser, has := pUM.UserPool[id]; has {
		return pUser, nil
	}
	return nil, fmt.Errorf("not found user id:%d", id)
}
