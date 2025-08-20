package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/x14n/evgateway/internal/gateway"
	"github.com/x14n/evgateway/internal/protocol"
)

// HandleStatusReport 处理状态上报
func HandleStatusReport(gw *gateway.Gateway, session *gateway.Session, frame protocol.Frame) error {
	var st map[string]any
	if err := json.Unmarshal(frame.Payload, &st); err != nil {
		return fmt.Errorf("bad status payload: %w", err)
	}
	fmt.Printf("[handler] status from %s: %+v\n", session.ID, st)
	return nil
}
