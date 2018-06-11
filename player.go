package main

import (
	"fmt"
	"io"
	"time"
)

type Player struct {
	session  *Session
	username string
	password string
	/* TODO: A lot of player data goes here. */
}

/**
 * Detect if Player connection has been reset by peer.
 *
 * Side-effect: If the connection isn't closed, set the
 * default read timeout to 1 minute.
 *
 * @return {true} if connection isn't closed, otherwise {false}
 */
func (self Player) isConnected() bool {
	one := []byte{}
	self.session.conn.SetReadDeadline(time.Now())
	if _, err := self.session.conn.Read(one); err == io.EOF {
		fmt.Printf("Connection closed for %s\n", self.session)
		self.session.conn.Close()
		self.session.conn = nil
		return false
	} else {
		self.session.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	}
	return true
}

/**
 * This is the main loop that each player has for processing
 * incoming packets.  It reads from their {Session}'s connection
 * buffer, and attempts to decode packet frames from the data it
 * reads.
 *
 * If a packet frame is detected and successfully decoded, then
 * it is passed to a {PacketHandler} instance that belongs to this
 * {Player}, to perform whatever actions that packet is intended to.
 *
 * If the connection times out, or is closed by other means,
 * the loop will end, and the players {Session}'s connection
 * will be set to nil.
 *
 * TODO: Possibly add more processing to this loop.
 * TODO: Handle errors while decoding frames more gracefully?
 *       That is to say, I don't believe we should end the {Session}
 *       because of a malformed packet frame.  Maybe just log it.
 */
func (self Player) process() {
	ph := NewPacketHandler(self)
	for {
		// End player process routine if the peer closed connection
		if !self.isConnected() {
			break
		}
		buffer := make([]byte, 3)
		bytesRead, err := self.session.conn.Read(buffer)
		// Connection reset by peer
		if err == io.EOF {
			fmt.Printf("Connection closed for %s\n", self.session)
			self.session.conn.Close()
			self.session.conn = nil
			break
		} else if err != nil || bytesRead < 3 {
			// The header for every packet is minimum 2 bytes for raw packets, 3 for frames
			continue
		}
		// position we're reading at in the buffer
		caret := 0
		/**
		 * Start by reading the packet frames length.
		 * The length is held in a Jagex-specific data type,
		 * commonly referred to as a SmartInt.
		 *
		 * If the length is less than 160, the length is held
		 * in a single byte, the first byte of the frame, and
		 * the second byte of the frame is given value of the last
		 * byte to be read in the frame.
		 *
		 * If the length is 160 or more, the length is held
		 * in the first 2 bytes of the frame, forming a
		 * short integer(16 bits).
		 *
		 * I imagine Jagex made this SmartInt type to
		 * save bandwidth, as this is a circa-2001 MMORPG.
		 * This type is unsigned.
		 */
		length := int(buffer[caret] & 0xFF)
		caret++

		var lastByte uint8 = 0
		/**
		 * Here we check if it's encoded as a short integer.
		 * Any length of 160 or over should be encoded to
		 * a short integer.
		 */
		if length >= 160 {
			length = (length-160)*256 + int(buffer[caret]&0xFF)
			caret++
		} else if length > 1 {
			lastByte = uint8(buffer[caret])
			caret++
			length--
		}
		opcode := uint8(buffer[caret] & 0xFF)
		caret++
		//		length--
		buffer = make([]byte, length)
		bytesRead, err = self.session.conn.Read(buffer)
		// Connection reset by peer
		if err == io.EOF {
			fmt.Printf("Connection closed for %s\n", self.session)
			self.session.conn.Close()
			self.session.conn = nil
			break
		} else if err != nil || bytesRead < length {
			fmt.Printf("packet read error:invalid length; expected %d, got %d\n", length, bytesRead)
			continue
		}
		fmt.Printf("%d\n", length)
		caret = 0
		/**
		 * If the packet frame is a byte, the second byte
		 * written to the buffer is the last byte of the packet
		 * frame.  I can't be sure why they did this, yet.
		 */
		//			var lastByte uint8 = 0
		//			if length < 160 && length > 1 {
		//				lastByte = uint8(buffer[caret])
		//				caret++
		//				length--
		//			}
		if length < 160 && length > 0 {
			buffer[length] = lastByte
			p := NewPacket(opcode, length, buffer)
			ph.HandlePacket(&p)
		} else if length == 0 {
			buffer = []byte{lastByte}
			p := NewPacket(opcode, length, buffer)
			ph.HandlePacket(&p)
		} else {
			p := NewPacket(opcode, length, buffer)
			ph.HandlePacket(&p)
		}
		//		p := NewPacket(opcode, length, buffer)
		/* TODO: Maybe goroutine for handling packets, to read next one instantly? */
		//		ph.HandlePacket(&p)
	}
}

func NewPlayer(session *Session, username string, password string) Player {
	p := Player{session, username, password}
	return p
}
