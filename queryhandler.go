package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"time"

	"github.com/openp2p-cn/totp"
)

type queryHandler struct {
}

func (h *queryHandler) handleMessage(ctx *msgContext) error {
	head := openP2PHeader{}
	err := binary.Read(bytes.NewReader(ctx.msg[:openP2PHeaderSize]), binary.LittleEndian, &head)
	if err != nil {
		return err
	}
	wsSess, ok := ctx.sess.(*wssSession)
	if !ok {
		gLog.Println(LvERROR, "interface conversion error")
		return errors.New("interface conversion error")
	}
	if head.SubType == MsgQueryPeerInfoReq {
		// TODO: verify token
		req := QueryPeerInfoReq{}
		err := json.Unmarshal(ctx.msg[openP2PHeaderSize:], &req)
		if err != nil {
			gLog.Printf(LvERROR, "%s wrong MsgQueryPeerInfoReq:%s", string(ctx.msg[openP2PHeaderSize:]), err)
			return err
		}
		rsp := QueryPeerInfoRsp{}
		gWSSessionMgr.allSessionsMtx.Lock()
		toSess, ok := gWSSessionMgr.allSessions[nodeNameToID(req.PeerNode)]
		gWSSessionMgr.allSessionsMtx.Unlock()
		if !ok {
			rsp.Online = 0
		} else {
			t := totp.TOTP{Step: totp.RelayTOTPStep}
			if !t.Verify(req.Token, toSess.token, time.Now().Unix()) {
				gLog.Printf(LvERROR, "%s MsgQueryPeerInfoReq denied", req.PeerNode)
				return errors.New("push denied")
			}
			rsp.Online = 1
			rsp.Version = toSess.version
			rsp.IPv4 = toSess.IPv4
			rsp.HasIPv4 = toSess.hasIPv4
			rsp.HasUPNPorNATPMP = toSess.hasUPNPorNATPMP
			rsp.NatType = toSess.natType
			if req.Token == toSess.token {
				rsp.IPv6 = toSess.IPv6 // ipv6 is sensitive, totp token not set ipv6
			}
		}
		wsSess.write(head.MainType, MsgQueryPeerInfoRsp, rsp)
	}
	return nil
}
