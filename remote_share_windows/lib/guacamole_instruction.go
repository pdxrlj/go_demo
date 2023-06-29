package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
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

func Parse(s []byte) (*Instruction, error) {
	ioReader := bytes.NewReader(s)
	reader := bufio.NewReader(ioReader)
	instruction := &Instruction{}
	for {
		index, err := reader.ReadBytes('.')
		if err != nil {
			return nil, err
		}

		arg, err := reader.ReadSlice(',')
		if err != nil {
			if err == io.EOF {
				indexInt, err := strconv.Atoi(string(index[:len(index)-1]))
				if err != nil {
					return nil, err
				}
				end := bytes.TrimSuffix(s[len(s)-indexInt-1:], []byte{';'})

				instruction.Args = append(instruction.Args, string(end))
				return instruction, nil
			} else {
				return nil, err
			}
		}
		arg = bytes.TrimSuffix(arg, []byte{','})

		if instruction.Opcode == "" {
			instruction.Opcode = string(arg)
			continue
		} else {
			instruction.Args = append(instruction.Args, string(arg))
		}
	}
}
