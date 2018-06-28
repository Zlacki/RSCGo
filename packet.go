package main

type packet struct {
	opcode        uint8
	length        int
	data          []byte
	currentOffset int
}

func (p *packet) incOffset() {
	p.currentOffset++
}

func (p *packet) readLong() uint64 {
	l := uint64(uint64(p.data[p.currentOffset])<<56 | uint64(p.data[p.currentOffset+1])<<48 | uint64(p.data[p.currentOffset+2])<<40 | uint64(p.data[p.currentOffset+3])<<32 | uint64(p.data[p.currentOffset+4])<<24 | uint64(p.data[p.currentOffset+5])<<16 | uint64(p.data[p.currentOffset+6])<<8 | uint64(p.data[p.currentOffset+7]))
	p.currentOffset += 8
	return l
}

func (p *packet) readInt() uint32 {
	i := uint32(uint32(p.data[p.currentOffset])<<24 | uint32(p.data[p.currentOffset+1])<<16 | uint32(p.data[p.currentOffset+2])<<8 | uint32(p.data[p.currentOffset+3]))
	p.currentOffset += 4
	return i
}

func (p *packet) readShort() uint16 {
	si := uint16((uint16(p.data[p.currentOffset] & 0xFF) << 8) | uint16(p.data[p.currentOffset+1] & 0xFF))
	p.currentOffset += 2
	return si
}

func (p *packet) readByte() byte {
	defer p.incOffset()
	return p.data[p.currentOffset]
}

func NewPacket(opcode uint8, length int, data []byte) packet {
	return packet{opcode, length, data, 0}
}
