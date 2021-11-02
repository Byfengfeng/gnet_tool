package code_tool

import (
	"github.com/Byfengfeng/gnet_tool/inter"
	"github.com/Byfengfeng/gnet_tool/log"
	"github.com/Byfengfeng/gnet_tool/pb"
	"github.com/Byfengfeng/gnet_tool/user"
	"github.com/Byfengfeng/gnet_tool/utils"
	"go.uber.org/zap"
	"time"
)

type IRequestCtx struct {
	Cid int64
	Addr string
}

var(
	_codeResponse map[uint16]FuncUserResponse
	_codePkt *CodecBase
	_users inter.IUserMapper
)

func init() {
	_codeResponse = make(map[uint16]FuncUserResponse)
	_codePkt = NewCodecProto()
	_users = user.NewUserMapperService()
}
func NewIRequestCtx(cid int64,addr string) *IRequestCtx {
	return &IRequestCtx{cid,addr}
}

//处理用户请求
type FuncUserResponse func(ctx IRequestCtx, pkt interface{}, resCh chan<- interface{})

func Request(address string,netWork inter.INetwork,code uint16,data []byte)  {
	reqFn,ok := _codeResponse[code]
	if !ok {
		log.Logger.Error("err req code not fail ")
		return
	}

	reqPkt,err := _codePkt.DecodeReq(code,data)
	if err != nil {
		log.Logger.Error("err req code not fail ",zap.Error(err))
		return
	}
	//hanDel
	resChan := make(chan interface{},1)
	ctx := netWork.GetCtx().(*IRequestCtx)
	reqFn(*ctx,reqPkt,resChan)
	timer := time.NewTimer(5 * time.Second)
	select {
	case <- timer.C:
		log.Logger.Error("err res time out ")
		_users.UserKickOut(address,ctx.Cid)
	case res := <-resChan:
		switch resData := res.(type) {
		case func(interface{}):
			timer.Stop()
			_codePkt.PutReqPkt(code,reqPkt)
			bytes := _codePkt.EncodeRes(code, resData)
			_users.Response(address,bytes)
		}
	}
}

func OffLine(addr string,cid int64)  {
	_users.UserKickOut(addr,cid)
}

func NewChannel(n inter.INetwork)  {
	_users.AddUserByAddr(n)
}

func GetNetWorkByAddr(addr string) inter.INetwork {
	return _users.GetUserByAddr(addr)
}

func GetCodeResponse() map[uint16]FuncUserResponse {
	return _codeResponse
}

func GetCodePkt() *CodecBase {
	return _codePkt
}

func UserAddCid(addr string,cid int64)  {
	_users.AddUserByCid(addr,cid)
}

func Set() {
	_codePkt.BindPool(1001, func() interface{} {
		return &pb.ReqLogin{}
	}, func() interface{} {
		return &pb.ResLogin{}
	})
	_codeResponse[1001] = func(ctx IRequestCtx, pkt interface{}, resCh chan<- interface{}) {
		req := pkt.(*pb.ReqLogin)
		if len(req.Token) > 0{
			ctx.Cid = utils.GetSnowflakeId()
			_users.AddUserByCid(ctx.Addr,ctx.Cid)
		}
		resCh <- func(pkt interface{}) {
			res := pkt.(*pb.ResLogin)
			res.Cid = ctx.Cid
		}
	}
}