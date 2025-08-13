package protocol

const (
	FrameHeader uint16 = 0xAA55 // 16-bit header for frame start
	FrameTail   uint16 = 0x55AA // 16-bit tail for frame end

	// header(2) + version(1) + cmd(1) + length(4) + crc(2) + tail(2)
	//MinFrameSize
	DefaultMinFrameSize   = 2 + 1 + 1 + 4 + 2 + 2
	DefaultMaxPayloadSize = 64 * 1024 //64kb
)

// cmd commands
const (
	CmdRegister  byte = 1 // Register a new client
	CmdHeartbeat byte = 2 // Heartbeat signal
	CmdStatus    byte = 3 // Status update
	CmdError     byte = 4 // Error message
)
