package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name      string
		writer    io.Writer
		input     any
		expected  []byte
		expectErr bool
	}{
		{
			name:      "Write byte slice",
			writer:    &bytes.Buffer{},
			input:     []byte{0x01, 0x02, 0x03},
			expected:  []byte{0x01, 0x02, 0x03},
			expectErr: false,
		},
		{
			name:      "Write int32",
			writer:    &bytes.Buffer{},
			input:     int32(42),
			expected:  []byte{0x00, 0x00, 0x00, 0x2A}, // BigEndian representation of 42
			expectErr: false,
		},
		{
			name:   "Write float64",
			writer: &bytes.Buffer{},
			input:  float64(3.14),
			expected: func() []byte {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.BigEndian, float64(3.14))
				return buf.Bytes()
			}(),
			expectErr: false,
		},
		{
			name:      "Unsupported type",
			writer:    &bytes.Buffer{},
			input:     struct{}{},
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Nil writer",
			writer:    nil,
			input:     []byte{0x01, 0x02},
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "Writer error",
			writer:    &errorWriter{},
			input:     []byte{0x01, 0x02},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := tt.writer
			if writer == nil {
				writer = &buf
			}

			err := write(writer, tt.input)

			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}

			if !tt.expectErr && writer == &buf && !bytes.Equal(buf.Bytes(), tt.expected) {
				t.Errorf("expected output: %v, got: %v", tt.expected, buf.Bytes())
			}
		})
	}
}
