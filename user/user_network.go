package user

import (
	"github.com/Byfengfeng/gnet_tool/inter"
	"github.com/Byfengfeng/gnet_tool/log"
	"go.uber.org/zap"
	"sync"
)

type userMapperService struct {
	addressMapperINetwork map[string]inter.INetwork
	cidMapperAddress map[int64]string
	delLock sync.Mutex
}

func NewUserMapperService() inter.IUserMapper {
	return &userMapperService{
		make(map[string]inter.INetwork,0),
		make(map[int64]string,0),
		sync.Mutex{}}
}

func (u *userMapperService) Response(address string,data []byte)  {
	iNetwork := u.GetUserByAddr(address)
	if iNetwork != nil{
		iNetwork.Action(func() {
			iNetwork.WriteWriteChan(data)
		})
	}
}

func (u *userMapperService) IsLogin(where interface{}) bool {
	switch ifData := where.(type) {
	case string:
		_,ok := u.addressMapperINetwork[ifData]
		return ok
	case int64:
		addr,ok := u.cidMapperAddress[ifData]
		if ok {
			_,ok1 := u.addressMapperINetwork[addr]
			return ok1
		}
	}
	return false
}

func (u *userMapperService) AddUserByAddr(netWork inter.INetwork)  {
	u.delLock.Lock()
	defer u.delLock.Unlock()
	u.addressMapperINetwork[netWork.GetAddr()] = netWork
}

func (u *userMapperService) AddUserByCid(addr string,cid int64)  {
	u.delLock.Lock()
	defer u.delLock.Unlock()
	user,ok := u.addressMapperINetwork[addr]
	if ok {
		user.SetCid(cid)
		u.cidMapperAddress[cid] = addr
	}else{
		log.Logger.Error("找不到netWork",zap.String("addr:",addr))
	}
}

func (u *userMapperService) AddUserByUid(addr string,uid int64)  {
	u.delLock.Lock()
	defer u.delLock.Unlock()
	user,ok := u.addressMapperINetwork[addr]
	if ok {
		user.SetUid(uid)
		u.cidMapperAddress[uid] = addr
	}else{
		log.Logger.Error("找不到netWork",zap.String("addr:",addr))
	}
}

func (u *userMapperService) GetUserByAddr(addr string) inter.INetwork {
	u.delLock.Lock()
	defer u.delLock.Unlock()
	n,ok := u.addressMapperINetwork[addr]
	if ok {
		return n
	}
	return nil
}

func (u *userMapperService) GetUserByCid(cid int64) string {
	u.delLock.Lock()
	defer u.delLock.Unlock()
	a,ok := u.cidMapperAddress[cid]
	if ok {
		return a
	}
	return ""
}

func (u *userMapperService) UserKickOut(addr string,cid int64,isExit bool) {
	if !isExit {
		//... 被踢下线、登出
	}
	log.Logger.Info("user close ",zap.Int64("cid:",cid),zap.String("addr:",addr))
	var user inter.INetwork
	user = u.GetUserByAddr(addr)
	if user != nil {
		u.delLock.Lock()
		delete(u.addressMapperINetwork,addr)
		delete(u.cidMapperAddress,cid)
		u.delLock.Unlock()
		user.SetIsClose()
	}
	return
}