package code_tool

type IRequestCtx struct {
	cid int64
	addr string
}

//处理用户请求
type FuncUserResponse func(ctx IRequestCtx, pkt interface{}, resCh chan<- interface{})

