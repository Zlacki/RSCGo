package main

import "fmt"

type PacketHandler struct {
	client *Session
	packet *packet
}

var handledPackets = []int{32}

func (handler *PacketHandler) sessionIDRequest() {
	fmt.Println("Session ID Request handled")
}

func (handler *PacketHandler) HandlePacket() {
	switch handler.packet.Opcode {
	case 32:
		handler.sessionIDRequest()
		break
	default:
		break
	}
	/*
		for opcode := range handledPackets {
			switch opcode {
			case 32:
				handler.sessionIDRequest()
				break
			default:
				break
			}
		}
	*/
}

func NewPacketHandler(client *Session, packet *packet) PacketHandler {
	ph := PacketHandler{client, packet}
	return ph
}
