package lib

import (
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
)

type Stream struct {
	Conn         net.Conn
	parseStart   int
	ConnectionID string
	buffer       []byte
	reset        []byte
	timeout      time.Duration
}

func NewStream(conn net.Conn, timeout time.Duration) *Stream {
	buffer := make([]byte, 0, MaxGuacamoleMessage*3)
	return &Stream{
		Conn:    conn,
		timeout: timeout,
		buffer:  buffer,
		reset:   buffer[:cap(buffer)],
	}
}

func (s *Stream) Handshake(config *GuacamoleConfig) (stream *Stream, err error) {
	selectArg := config.ConnectionID
	if len(selectArg) == 0 {
		selectArg = config.Protocol
	}
	_, err = s.Write(NewInstruction("select", selectArg).Bytes())
	if err != nil {
		return stream, err
	}

	instruction, err := s.AssertOpcode("args")
	if err != nil {
		return stream, errors.WithStack(err)
	}
	argNameS := instruction.Args
	argValueS := make([]string, 0, len(argNameS))
	for _, argName := range argNameS {

		// Retrieve argument name

		// Get defined value for name
		value := config.Parameters[argName]

		// If value defined, set that value
		if len(value) == 0 {
			value = ""
		}
		argValueS = append(argValueS, value)
	}

	// Send size
	_, err = s.Write(NewInstruction("size",
		fmt.Sprintf("%v", config.OptimalScreenWidth),
		fmt.Sprintf("%v", config.OptimalScreenHeight),
		fmt.Sprintf("%v", config.OptimalResolution)).Bytes(),
	)
	if err != nil {
		return stream, err
	}

	// Send supported audio formats
	_, err = s.Write(NewInstruction("audio", config.AudioMimetypes...).Bytes())
	if err != nil {
		return stream, err
	}

	// Send supported video formats
	_, err = s.Write(NewInstruction("video", config.VideoMimetypes...).Bytes())
	if err != nil {
		return stream, err
	}

	// Send supported image formats
	_, err = s.Write(NewInstruction("image", config.ImageMimetypes...).Bytes())
	if err != nil {
		return stream, err
	}

	// Send Args
	_, err = s.Write(NewInstruction("connect", argValueS...).Bytes())
	if err != nil {
		return stream, err
	}

	// Wait for ready, store ID
	ready, err := s.AssertOpcode("ready")
	if err != nil {
		return stream, err
	}
	readyArgs := ready.Args
	if len(readyArgs) == 0 {
		return stream, errors.WithStack(errors.New("ready instruction has no arguments"))
	}
	s.ConnectionID = readyArgs[0]

	return s, nil
}

func (s *Stream) Write(data []byte) (int, error) {
	err := s.Conn.SetWriteDeadline(time.Now().Add(s.timeout))
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return s.Conn.Write(data)
}

func (s *Stream) Read() ([]byte, error) {
	err := s.Conn.SetReadDeadline(time.Now().Add(s.timeout))
	if err != nil {
		return nil, err
	}

	var n int
	// While we're blocking, or input is available
	for {
		// Length of element
		var elementLength int

		// Resume where we left off
		i := s.parseStart

	parseLoop:
		// Parse instruction in buffer
		for i < len(s.buffer) {
			// ReadSome character
			readChar := s.buffer[i]
			i++

			switch readChar {
			// If digit, update length
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				elementLength = elementLength*10 + int(readChar-'0')

			// If not digit, check for end-of-length character
			case '.':
				if i+elementLength >= len(s.buffer) {
					// break for i < s.usedLength { ... }
					// Otherwise, read more data
					break parseLoop
				}
				// Check if element present in buffer
				terminator := s.buffer[i+elementLength]
				// Move to character after terminator
				i += elementLength + 1

				// Reset length
				elementLength = 0

				// Continue here if necessary
				s.parseStart = i

				// If terminator is semicolon, we have a full
				// instruction.
				switch terminator {
				case ';':
					instruction := s.buffer[0:i]
					s.parseStart = 0
					s.buffer = s.buffer[i:]
					return instruction, nil
				case ',':
					// keep going
				default:
					err = errors.New("Element terminator of instruction was not ';' nor ','")
					return nil, err
				}
			default:
				// Otherwise, parse error
				err = errors.New("Non-numeric character in element length:" + string(readChar))
				return nil, err
			}
		}

		if cap(s.buffer) < MaxGuacamoleMessage {
			s.Flush()
		}

		n, err = s.Conn.Read(s.buffer[len(s.buffer):cap(s.buffer)])
		if err != nil && n == 0 {
			switch err.(type) {
			case net.Error:
				ex := err.(net.Error)
				if ex.Timeout() {
					err = errors.New("Connection to guacd timed out." + err.Error())
				} else {
					err = errors.New("Connection to guacd is closed." + err.Error())
				}
			default:
				err = errors.New(err.Error())
			}
			return nil, err
		}
		if n == 0 {
			err = errors.New("read 0 bytes")
		}
		// must reslice so len is changed
		s.buffer = s.buffer[:len(s.buffer)+n]
	}
}

func (s *Stream) AssertOpcode(opcode string) (*Instruction, error) {
	read, err := s.Read()
	if err != nil {
		return nil, err
	}
	instruction, err := Parse(read)
	if err != nil {
		return nil, err
	}
	if instruction.Opcode != opcode {
		return nil, errors.Errorf("expected opcode %s, got %s", opcode, instruction.Opcode)
	}

	return instruction, nil
}

func (s *Stream) Close() error {
	return s.Conn.Close()
}

func (s *Stream) Flush() {
	copy(s.reset, s.buffer)
	s.buffer = s.reset[:len(s.buffer)]
}
