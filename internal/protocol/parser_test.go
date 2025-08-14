package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestTryParse_ValidFrame(t *testing.T) {
	payload := []byte("hello")
	length := uint32(len(payload))

	// Create a byte slice that simulates a valid frame
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, FrameHeader) // Header
	buf.WriteByte(1)                                  // Version
	buf.WriteByte(2)                                  // Cmd
	binary.Write(&buf, binary.BigEndian, length)      // Length
	buf.Write(payload)                                // Payload

	// CRC
	crcData := make([]byte, 0, 1+1+4+len(payload))
	crcData = append(crcData, 1, 2)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, length)
	crcData = append(crcData, lenBuf...)
	crcData = append(crcData, payload...)
	crc := CRC16CCITT(crcData)
	binary.Write(&buf, binary.BigEndian, crc)

	binary.Write(&buf, binary.BigEndian, FrameTail)

	p := NewParser(nil)
	p.buf = buf.Bytes()

	frame, consumed, err := p.tryParse()
	fmt.Printf("buf: %v\n", buf.Bytes())
	fmt.Println("frame:", frame)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if consumed != buf.Len() {
		t.Errorf("expected consumed=%d, got %d", buf.Len(), consumed)
	}
	if frame.Version != 1 || frame.Cmd != 2 {
		t.Errorf("unexpected frame header: %+v", frame)
	}
	if string(frame.Payload) != "hello" {
		t.Errorf("unexpected payload: %s", string(frame.Payload))
	}
}

//
// func TestTryParse(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		buf         []byte
// 		expected    *Frame
// 		expectedLen int
// 		expectErr   error
// 	}{
// 		{
// 			name:        "Valid frame",
// 			buf:         []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
// 			expected:    &Frame{Version: 1, Cmd: 2, Payload: []byte{0x03, 0x04, 0x05}},
// 			expectedLen: 6,
// 			expectErr:   nil,
// 		},
// 		{
// 			name:        "Incomplete frame",
// 			buf:         []byte{0x00, 0x01},
// 			expected:    nil,
// 			expectedLen: 2,
// 			expectErr:   ErrNeedMoreData,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &Parser{buf: tt.buf}
// 			frame, length, err := p.tryParse()
// 			if !errors.Is(err, tt.expectErr) {
// 				t.Errorf("expected error %v, got %v", tt.expectErr, err)
// 			}
// 			if frame != nil && *frame != *tt.expected {
// 				t.Errorf("expected frame %v, got %v", tt.expected, frame)
// 			}
// 			if length != tt.expectedLen {
// 				t.Errorf("expected length %d, got %d", tt.expectedLen, length)
// 			}
// 		})
// 	}
// }
