package handlers

import (
	"fmt"

	"github.com/x14n/evgateway/internal/gateway"
	"github.com/x14n/evgateway/internal/protocol"
)

// HandleHeartbeat 处理心跳命令
func HandleHeartbeat(gw *gateway.Gateway, session *gateway.Session, frame protocol.Frame) error {
	session.UpdateLastSeen()
	fmt.Println("[handler] heartbeat:", session.ID)
	return nil
}
