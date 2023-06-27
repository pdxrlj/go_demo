package lib

import (
	"fmt"
)

type Instruction struct {
	Opcode string
	Args   []string
	cache  string
}

func NewInstruction(opcode string, args ...string) *Instruction {
	return &Instruction{
		Opcode: opcode,
		Args:   args,
	}
}

// 6.select,3.vnc;
// 4.size,1.0,4.1024,3.768;
func (i *Instruction) String() string {
	if len(i.cache) > 0 {
		return i.cache
	}

	i.cache = fmt.Sprintf("%d.%s", len(i.Opcode), i.Opcode)
	for _, value := range i.Args {
		i.cache += fmt.Sprintf(",%d.%s", len(value), value)
	}
	i.cache += ";"

	return i.cache
}

func (i *Instruction) Bytes() []byte {
	return []byte(i.String())
}

func Parse(s *Stream) (*Instruction, error) {
	read, err := s.Read()
	if err != nil {
		return nil, err
	}
	fmt.Printf("read: %s\n", read)
	return nil, nil
}
