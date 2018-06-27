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
	fmt.Println("Checking for updates...")
	equinoxUpdate()
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
		/* TODO: Collections of server entities, incl. Players */
		go player.process()
	}
}
