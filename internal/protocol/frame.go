package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

// header(2) + version(1) + cmd(1) + length(4) +Payload(N)+ crc(2) + tail(2)

type Frame struct {
	Version byte
	Cmd     byte
	Payload []byte
}

func NewFrame(version byte, cmd byte, payload []byte) *Frame {
	return &Frame{
		Version: version,
		Cmd:     cmd,
		Payload: payload,
	}
}

// pack to make Frame into io.Write (BigEndian)
func (f *Frame) Packe(w io.Writer) error {
	//check if the payload size is within limits
	payloadLength := uint32(len(f.Payload))
	if payloadLength > DefaultMaxPayloadSize {
		return PayloadTooLargeError
	}

	//header
	if err := write(w, FrameHeader); err != nil {
		return err
	}

	// write version and cmd
	if err := write(w, f.Version); err != nil {
		return err
	}

	// Length of payload
	if err := write(w, payloadLength); err != nil {
		return err
	}

	//payload
	if payloadLength > 0 {
		if err := write(w, f.Payload); err != nil {
			return err
		}
	}

	// CRC calc on (version|cmd|len|payload)
	crcData := make([]byte, 0, 1+1+4+len(f.Payload))
	crcData = append(crcData, f.Version, f.Cmd)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, payloadLength)
	crcData = append(crcData, lenBuf...)
	crcData = append(crcData, f.Payload...)

	crc := CRC16CCITT(crcData)
	if err := binary.Write(w, binary.BigEndian, crc); err != nil {
		return err
	}

	//tail
	if err := write(w, FrameTail); err != nil {
		return err
	}
	return nil

}

//	if err := binary.Write(w, binary.BigEndian, FrameHeader); err != nil {
//
// return err
// }
func write(w io.Writer, data any) error {
	switch v := data.(type) {
	case []byte:
		_, err := w.Write(v)
		return err
	default:
		if err := binary.Write(w, binary.BigEndian, data); err != nil {
			return err
		}
	}

	return errors.New("unsupported type for write")
}

func CRC16CCITT(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if (crc & 0x8000) != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return crc & 0xFFFF
}
