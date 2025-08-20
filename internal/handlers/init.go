package handlers

import (
	"github.com/x14n/evgateway/internal/gateway"
	"github.com/x14n/evgateway/internal/protocol"
)

// RegisterAllHandlers 把所有命令处理器注册到 Dispatcher
func RegisterAllHandlers(d *gateway.Dispatcher) {
	d.RegisterHandler(protocol.CmdRegister, HandleRegister)
	d.RegisterHandler(protocol.CmdHeartbeat, HandleHeartbeat)
	d.RegisterHandler(protocol.CmdStatus, HandleStatusReport)
	d.RegisterHandler(protocol.CmdError, HandleErrorResponse)
}
