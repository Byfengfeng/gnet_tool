package user

import (
	"github.com/Byfengfeng/gnet_tool/inter"
	"github.com/Byfengfeng/gnet_tool/log"
	"go.uber.org/zap"
)

type userMapperService struct {
	addressMapperINetwork map[string]inter.INetwork
	cidMapperAddress map[int64]string
}

func NewUserMapperService() inter.IUserMapper {
	return &userMapperService{make(map[string]inter.INetwork,0),make(map[int64]string,0)}
}

func (u *userMapperService) Response(address string,data []byte)  {
	iNetwork := u.addressMapperINetwork[address]
	if iNetwork != nil && !iNetwork.GetClose(){
		if !iNetwork.GetClose() {
			iNetwork.WriteWriteChan(data)
		}
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
	u.addressMapperINetwork[netWork.GetAddr()] = netWork
}

func (u *userMapperService) AddUserByCid(addr string,cid int64)  {
	_,ok := u.addressMapperINetwork[addr]
	if ok {
		u.cidMapperAddress[cid] = addr
	}else{
		log.Logger.Error("找不到netWork",zap.String("addr:",addr))
	}
}

func (u *userMapperService) GetUserByAddr(addr string) inter.INetwork {
	n,ok := u.addressMapperINetwork[addr]
	if ok {
		return n
	}
	return nil
}

func (u *userMapperService) UserKickOut(where interface{}) {
	var user inter.INetwork
	switch ifData := where.(type) {
	case string:
		newUser,ok := u.addressMapperINetwork[ifData]
		if ok {
			user = newUser
		}
	case int64:
		addr,ok := u.cidMapperAddress[ifData]
		if ok {
			newUser,ok1 := u.addressMapperINetwork[addr]
			if ok1 {
				user = newUser
			}else{
				delete(u.addressMapperINetwork,addr)
			}
		}
	}
	if user != nil {
		user.SetIsClose()
	}
	return
}