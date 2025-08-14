package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"
)

type Parser struct {
	r           io.Reader      // 数据源，用于读取原始数据
	buf         []byte         // 缓冲区，用于存储未解析的数据
	maxPayload  uint32         // 最大允许的帧负载大小
	framesCh    chan Frame     // 用于传递解析出的帧
	errCh       chan error     // 用于传递解析过程中发生的错误
	quit        chan struct{}  // 用于通知解析器停止工作
	wg          sync.WaitGroup // 用于等待所有 goroutine 完成
	readBufSize int            // 每次从 io.Reader 读取的缓冲区大小
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		r:           r,
		buf:         make([]byte, 0, ReadBufSize),
		maxPayload:  DefaultMaxPayloadSize,
		framesCh:    make(chan Frame, 10), // 帧通道，缓冲区大小为10
		errCh:       make(chan error, 1),  // 错误通道，缓冲区大小为1
		quit:        make(chan struct{}),
		readBufSize: ReadBufSize,
	}
}

func (p *Parser) Frames() <-chan Frame {
	return p.framesCh
}

func (p *Parser) Errors() <-chan error {
	return p.errCh
}

func (p *Parser) handleParseError(err error) {
	// 报告错误（非堵塞）
	select {
	case p.errCh <- err:
	default:
	}

	//重同步策略：在剩余的缓冲区中查找下一个帧头
	next := -1
	for i := 1; i <= len(p.buf)-2; i++ {
		if binary.BigEndian.Uint16(p.buf[i:i+2]) == FrameHeader {
			next = i
			break
		}
	}
	if next == -1 {
		// 保留最后 1 字节，避免在帧头分割时丢失
		if len(p.buf) > 1 {
			p.buf = p.buf[len(p.buf)-1:]
		}
	} else {
		p.buf = p.buf[next:]
	}
}

func (p *Parser) Stop() {
	close(p.quit)
	p.wg.Wait() // 等待解析循环结束
}

func (p *Parser) Start() {
	p.wg.Add(1)
	go p.loop()
}

func (p *Parser) loop() {
	defer p.wg.Done()
	readBuf := make([]byte, p.readBufSize)
	for {
		if p.parseFrames() {
			return
		}
		if p.readMoreDat(readBuf) {
			return
		}
		if p.checkQuit() {
			return
		}
	}
}

func (p *Parser) checkQuit() bool {
	select {
	case <-p.quit:
		close(p.framesCh)
		return true
	default:
		return false
	}
}

// 从reader中曲度数据到缓冲区
func (p *Parser) readMoreDat(readBuf []byte) bool {
	n, err := p.r.Read(readBuf)
	if n > 0 {
		p.buf = append(p.buf, readBuf[:n]...)
	}
	if err != nil {
		if err == io.EOF {
			close(p.framesCh)
			return true
		}

		// 其他错误处理
		select {
		case p.errCh <- err:
		default:
		}
		close(p.framesCh)
		return true
	}
	return false
}

func (p *Parser) parseFrames() bool {
	for {
		frame, consumed, err := p.tryParse()
		if err != nil {
			p.handleParseError(err)
			continue
		}
		if consumed > 0 && frame != nil {
			if p.sendFrame(*frame) {
				return true // 发送失败，退出解析
			}
			p.consumeBufferByte(consumed)
			continue
		}
		break
	}
	return false
}

func (p *Parser) sendFrame(frame Frame) bool {
	select {
	case p.framesCh <- frame:
		return false
	case <-p.quit:
		close(p.framesCh)
		return true
	}
}

// 消费已经解析的字节
func (p *Parser) consumeBufferByte(consumed int) {
	if consumed >= len(p.buf) {
		p.buf = p.buf[:0]
	} else {
		p.buf = p.buf[consumed:]
	}
}

func (p *Parser) tryParse() (*Frame, int, error) {
	if len(p.buf) < 2 {
		return nil, 0, ErrNeedMoreData
	}

	//check header
	headIdx := -1
	for i := 0; i <= len(p.buf)-2; i++ {
		if binary.BigEndian.Uint16(p.buf[i:i+2]) == FrameHeader {
			headIdx = i
			break
		}
	}

	//if not get header,save last 1 byte (maybe a part of another header)
	if headIdx == -1 {
		if len(p.buf) > 1 {
			p.buf = p.buf[len(p.buf)-1:]
		}
		return nil, 0, ErrNeedMoreData
	}

	// 如果帧头不在缓冲开头，告诉调用者丢掉前 headIdx 字节
	if headIdx > 0 {
		return nil, headIdx, nil
	}

	//现在 buf[0:2] 是帧头
	if len(p.buf) < DefaultMinFrameSize {
		return nil, 0, errors.New("frame too short")
	}

	version := p.buf[2]
	cmd := p.buf[3]
	length := binary.BigEndian.Uint32(p.buf[4:8])

	if length > p.maxPayload {
		return nil, 0, errors.New("payload too large")
	}

	totalLen := int(2 + 1 + 1 + 4 + length + 2 + 2)

	if len(p.buf) < totalLen {
		return nil, 0, errors.New("data not enough for full frame")
	}

	// 校验tail
	tailPos := 2 + 1 + 1 + 4 + int(length) + 2
	tail := binary.BigEndian.Uint16(p.buf[tailPos : tailPos+2])
	if tail != FrameTail {
		return nil, 0, errors.New("invalid frame tail")
	}

	payloadStart := 2 + 1 + 1 + 4
	payload := make([]byte, length)
	copy(payload, p.buf[payloadStart:payloadStart+int(length)])

	//crc

	crcPos := payloadStart + int(length)
	readCRC := binary.BigEndian.Uint16(p.buf[crcPos : crcPos+2])

	crcData := make([]byte, 0, 1+1+4+len(payload))
	crcData = append(crcData, version, cmd)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, length)
	crcData = append(crcData, lenBuf...)
	crcData = append(crcData, payload...)

	expCRC := CRC16CCITT(crcData)
	if readCRC != expCRC {
		return nil, 0, ErrCRCMismatch
	}

	frame := &Frame{
		Version: version,
		Cmd:     cmd,
		Payload: payload,
	}
	return frame, totalLen, nil
}
