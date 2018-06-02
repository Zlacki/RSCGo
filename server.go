package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var port = flag.Int("port", 43594, "The port for the server to listen on")

func main() {
	fmt.Println("Basic RSC Server in Go")

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

		go processClient(clientSocket)
	}
}

func processClient(conn net.Conn) {
	session := NewSession(conn)
	defer conn.Close()
	remoteAddress := conn.RemoteAddr().String()
	fmt.Println("Incoming client from: ", remoteAddress)
	reader := bufio.NewReader(session.conn)
	for {
		if reader.Size() < 2 {
			fmt.Println("Sleep")
			time.Sleep(600 * time.Millisecond)
			continue
		}
		var lastByte byte = 0
		var size int = 0
		tmpSize, _ := reader.ReadByte()
		size = int(tmpSize)
		if size >= 160 {
			tmpSize2, _ := reader.ReadByte()
			size = (size-160)*256 + int(tmpSize2)
		}
		if size >= reader.Buffered() {
			if size < 160 && size > 1 {
				lastByteTmp, _ := reader.ReadByte()
				lastByte = lastByteTmp
				size--
			}
			payload := make([]byte, size)
			opcode, _ := reader.ReadByte()
			size--
			if size >= 160 {
				i, _ := reader.Read(payload)
				fmt.Println(i, " read to payload")
			} else if size > 0 {
				i, _ := reader.Read(payload)
				fmt.Println(i, " read to payload")
			}
			if size < 160 {
				payload[size] = lastByte
				size++
				/* Increase size cos index 0 is length 1... */
			}
			p := NewPacket(int(opcode), size, payload)
			fmt.Printf("Incoming packet[opcode:%d;length:%d;]\n", p.Opcode, p.Length)
			ph := NewPacketHandler(session, &p)
			ph.HandlePacket()
		} else {
			fmt.Println("Error with incoming packet.")
		}
	}
}
