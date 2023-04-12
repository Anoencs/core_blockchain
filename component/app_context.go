package component

import "projectx/network"

type AppContext interface {
	GetTrHandler() *network.LocalTransport
}

type AppCtx struct {
	handler *network.LocalTransport
}

func NewAppCtx(tr *network.LocalTransport) *AppCtx {
	return &AppCtx{
		handler: tr,
	}
}

func (ctx *AppCtx) GetTrHandler() *network.LocalTransport {
	return ctx.handler
}
