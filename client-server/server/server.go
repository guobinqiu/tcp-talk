package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

var port int

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.IntVar(&port, "p", 9090, "server port")
}

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("server listening on port: %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	bs := make([]byte, 1024)
	for {
		//read from a client
		n, err := conn.Read(bs)

		//occurred when a client was closed
		//before: server: CLOSE_WAIT, client: FIN_WAIT_2
		//do server close
		//after: client: TIME_WAIT
		//take care too many TIME_WAIT
		if err != nil {
			log.Println(err)
			break
		}

		//write to console, omit error handle
		os.Stdout.Write(bs[:n])
	}
}
