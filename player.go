package main

import "fmt"
import "io"
import "time"

type Player struct {
	session  *Session
	username string
	password string
	/* TODO: A lot of player data goes here. */
}

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

func (self Player) process() {
	ph := NewPacketHandler(self)
	for {
		if !self.isConnected() {
			break
		}
		buffer := make([]byte, 5000)
		bytesRead, err := self.session.conn.Read(buffer)
		if checkErr(err) {
			fmt.Println("Error preceding from attempting to read packet frame.")
			break
		}
		if bytesRead < 2 {
			continue
		}
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
		 * I imagine Jagex made this SmartUint16 type to
		 * save bandwidth, as this is a circa-2001 MMORPG.
		 * This type is unsigned.
		 */
		length := int(buffer[caret] & 0xFF)
		caret++

		/**
		 * Here we check if it's encoded as a smartUint16.
		 * Any length of 160 or over should be encoded to
		 * a smartUint16.
		 */
		if length >= 160 {
			length = (length-160)*256 + int(buffer[caret]&0xFF)
			caret++
		}
		if length >= bytesRead-2 {
			/**
			 * If the packet frame is a smartUint8, the second byte
			 * written to the buffer is the last byte of the packet
			 * frame.  I can't be sure why they did this, yet.
			 */
			var lastByte uint8 = 0
			if length < 160 && length > 1 {
				lastByte = uint8(buffer[caret])
				caret++
				length--
			}

			/**
			 * This reads the packet frame's opcode.
			 * It's like an identifier code for what operation
			 * this packet is supposed to execute on the server.
			 *
			 * Maximum is 256, so max user operations is 256.
			 * Potential to expand, but probably never will
			 * need to.
			 *
			 * Potentially call it frameID?
			 */
			opcode := uint8(buffer[caret] & 0xFF)
			caret++
			length--
			payload := buffer[caret:]
			if length < 160 {
				payload[length] = lastByte
			}
			p := NewPacket(opcode, length, payload)
			/* TODO: Maybe goroutine for handling packets, to read next one instantly? */
			ph.HandlePacket(&p)
		} else {
			fmt.Println("Error reading packet frame.  Expected frame length less than available bytes in buffer.")
		}
	}
}

func NewPlayer(session *Session, username string, password string) Player {
	p := Player{session, username, password}
	return p
}
