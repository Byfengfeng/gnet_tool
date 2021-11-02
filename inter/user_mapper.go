package inter

type IUserMapper interface {
	Response(address string,data []byte)
	IsLogin(where interface{}) bool
	UserKickOut(addr string,cid int64)
	AddUserByAddr(netWork INetwork)
	AddUserByCid(addr string,cid int64)
	GetUserByAddr(addr string) INetwork
}