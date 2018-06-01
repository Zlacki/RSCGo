package main

type packet struct {
	Opcode int
	Length int
	Data   []byte
}

func NewPacket(opcode int, length int, data []byte) packet {
	p := packet{opcode, length, data}
	return p
}
