package server

import (
	"chatServer/protocol"
	"chatServer/server/chatpb"

	l4g "github.com/alecthomas/log4go"
	"google.golang.org/protobuf/proto"
)

const (
	KeyUser     = "USER"
	KeyConfig   = "Config:Normal"
	FiledCurUID = "CurUserID"
)

var (
	ProtoMap = make(map[uint32]IProto)
)

type IProto interface {
	Execute(pMsg *UserMsg)
}

func ProtoRegister() {
	ProtoMap[protocol.Chatpb_C2SLogin] = &C2SLogin{}
	ProtoMap[protocol.Chatpb_C2SRegister] = &C2SRegister{}
	ProtoMap[protocol.Chatpb_C2SOnlineUsers] = &C2SOnlineUsers{}
	ProtoMap[protocol.Chatpb_C2SChat] = &C2SChat{}
	ProtoMap[protocol.Chatpb_C2SChatCntTopUsers] = &C2SChatCntTopUsers{}
	ProtoMap[protocol.Chatpb_C2SChatCnt] = &C2SChatCnt{}
	ProtoMap[protocol.Chatpb_C2SHeartBeat] = &C2SHeartBeat{}
}

type C2SLogin chatpb.C2SLogin

func (this *C2SLogin) Execute(pMsg *UserMsg) {
	loginMap := &chatpb.C2SLogin{}
	if err := proto.Unmarshal(pMsg.msg, loginMap); err != nil {
		l4g.Error("C2SLogin proto unmarshal error:%s", err)
	}
	l4g.Info(loginMap)
}

type C2SRegister chatpb.C2SRegister

func (this *C2SRegister) Execute(pMsg *UserMsg) {
	regMap := &chatpb.C2SRegister{}
	if err := proto.Unmarshal(pMsg.msg, regMap); err != nil {
		l4g.Error("C2SRegister proto unmarshal error:%s", err)
	}
	l4g.Info(regMap)
}

type C2SOnlineUsers chatpb.C2SOnlineUsers

func (this *C2SOnlineUsers) Execute(pMsg *UserMsg) {
	onlineMap := &chatpb.C2SOnlineUsers{}
	if err := proto.Unmarshal(pMsg.msg, onlineMap); err != nil {
		l4g.Error("C2SOnlineUsers proto unmarshal error:%s", err)
	}
	l4g.Info(onlineMap)
}

type C2SChat chatpb.C2SChat

func (this *C2SChat) Execute(pMsg *UserMsg) {
	chatMap := &chatpb.C2SChat{}
	if err := proto.Unmarshal(pMsg.msg, chatMap); err != nil {
		l4g.Error("C2SChat proto unmarshal error:%s", err)
	}
	l4g.Info(chatMap)
}

type C2SChatCntTopUsers chatpb.C2SChatCntTopUsers

func (this *C2SChatCntTopUsers) Execute(pMsg *UserMsg) {
	chatUserMap := &chatpb.C2SChatCntTopUsers{}
	if err := proto.Unmarshal(pMsg.msg, chatUserMap); err != nil {
		l4g.Error("C2SChatCntTopUsers proto unmarshal error:%s", err)
	}
	l4g.Info(chatUserMap)
}

type C2SChatCnt chatpb.C2SChatCnt

func (this *C2SChatCnt) Execute(pMsg *UserMsg) {
	chatCntMap := &chatpb.C2SChatCnt{}
	if err := proto.Unmarshal(pMsg.msg, chatCntMap); err != nil {
		l4g.Error("C2SChatCnt proto unmarshal error:%s", err)
	}
	l4g.Info(chatCntMap)
}

type C2SHeartBeat chatpb.C2SHeartBeat

func (this *C2SHeartBeat) Execute(pMsg *UserMsg) {
	l4g.Info("[heartbeat]connid:%d, userid:%d, nickname:%s", pMsg.user.connid, pMsg.user.id, pMsg.user.nickname)
}

func ProtoDispatcher(pMsg *UserMsg) {
	if cb, has := ProtoMap[pMsg.pid]; has {
		cb.Execute(pMsg)
	} else {
		l4g.Error("ProtoDispatcher: not found proto;%s", pMsg.pid)
	}
}
