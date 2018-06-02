package main

type packet struct {
	opcode int
	length int
	data   []byte
	currentOffset int
}

func (p *packet) incOffset() {
	p.currentOffset++
}

func (p *packet) readByte() byte {
	defer p.incOffset()
	return p.data[p.currentOffset]
}

func NewPacket(opcode int, length int, data []byte) packet {
	return packet{opcode, length, data, 0}
}
