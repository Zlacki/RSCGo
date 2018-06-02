package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

var port = flag.Int("port", 43594, "The port for the server to listen on")

func init() {
	fmt.Println("Basic RSC Server in Go initializing...")
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	fmt.Print("Attempting to bind to port ", strconv.Itoa(*port), "...")
	socketListener, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		fmt.Println("ERROR: Could not bind to specified port.")
		os.Exit(1)
	}
	fmt.Println("done")
	defer socketListener.Close()

	for {
		clientSocket, err := socketListener.Accept()
		if err != nil {
			fmt.Println("Some issue accepting a client.  I don't think it's fatal; they can try to reconnect.")
			fmt.Println("Error: ", err.Error())
			continue
		}
		s := NewSession(clientSocket)
		player := NewPlayer(&s, "", "")
		go player.process()
		//		go handleIncomingPackets(session)
	}
}

/*func handleIncomingPackets(session Session) {
	reader := bufio.NewReader(session.conn)
	for {
		var lastByte byte = 0
		var size int = 0
		tmpSize, err := reader.ReadByte()
		if err != nil {
			break
		}
		size = int(tmpSize)
		if size >= 160 {
			tmpSize2, err := reader.ReadByte()
			if err != nil {
				break
			}
			size = (size-160)*256 + int(tmpSize2)
		}
		if size >= reader.Buffered() {
			if size < 160 && size > 1 {
				lastByteTmp, err := reader.ReadByte()
				if err != nil {
					break
				}
				lastByte = lastByteTmp
				size--
			}
			payload := make([]byte, size)
			opcode, err := reader.ReadByte()
			if err != nil {
				break
			}
			size--
			if size >= 160 {
				i, err := reader.Read(payload)
				if i < size || err != nil {
					break
				}
			} else if size > 0 {
				i, err := reader.Read(payload)
				if i < size || err != nil {
					break
				}
			}
			if size < 160 {
				payload[size] = lastByte
				size++
			}
			p := NewPacket(int(opcode), size, payload)
			fmt.Printf("Incoming packet[opcode:%d;length:%d;]\n", p.opcode, p.length)
			ph := NewPacketHandler(NewPlayer(session, "", ""), &p)
			ph.HandlePacket()
		} else {
			fmt.Printf("Error with incoming packet.  reader.Size():%d; reader.Buffered():%d\n", reader.Size(), reader.Buffered())
		}
	}

	fmt.Printf("Closing connection from: %s\n", session.ipAddress)
	session.conn.Close()
	if session.conn != nil {
		session.conn = nil
	}
}
*/
