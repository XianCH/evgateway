package handlers

import (
	"fmt"

	"github.com/x14n/evgateway/internal/gateway"
	"github.com/x14n/evgateway/internal/protocol"
)

// HandleRegister 处理注册命令
func HandleRegister(gw *gateway.Gateway, session *gateway.Session, frame protocol.Frame) error {
	chargerID := string(frame.Payload)
	session.ID = chargerID
	gw.AddSession(session)
	fmt.Println("[handler] register:", chargerID, "from", session.Addr)
	return nil
}
