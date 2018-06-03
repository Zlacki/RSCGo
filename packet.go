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
	l := uint64(p.data[p.currentOffset]<<56 | p.data[p.currentOffset+1]<<48 | p.data[p.currentOffset+2]<<40 | p.data[p.currentOffset+3]<<32 | p.data[p.currentOffset+4]<<24 | p.data[p.currentOffset+5]<<16 | p.data[p.currentOffset+6]<<8 | p.data[p.currentOffset+7])
	p.currentOffset += 8
	return l
}

func (p *packet) readInt() uint32 {
	i := uint32(p.data[p.currentOffset]<<24 | p.data[p.currentOffset+1]<<16 | p.data[p.currentOffset+2]<<8 | p.data[p.currentOffset+3])
	p.currentOffset += 4
	return i
}

func (p *packet) readShort() uint16 {
	si := uint16(((p.data[p.currentOffset] & 0xFF) << 8) | (p.data[p.currentOffset+1] & 0xFF))
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
