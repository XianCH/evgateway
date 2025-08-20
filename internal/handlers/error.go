package handlers

import (
	"fmt"

	"github.com/x14n/evgateway/internal/gateway"
	"github.com/x14n/evgateway/internal/protocol"
)

// HandleErrorResponse 处理错误响应
func HandleErrorResponse(gw *gateway.Gateway, session *gateway.Session, frame protocol.Frame) error {
	fmt.Printf("[handler] error from %s: %s\n", session.ID, string(frame.Payload))
	return nil
}
