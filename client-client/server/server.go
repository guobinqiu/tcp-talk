package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

var port int

//map is not thread safe, you need mutex lock,
//to keep the code simple i didn't use lock here.
type connList map[net.Conn]net.Conn

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

	conns := make(connList)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		conns[conn] = conn

		go handleConn(conn, conns)
	}
}

//read from client and write to all clients
func handleConn(conn net.Conn, conns connList) {
	bs := make([]byte, 1024)
	for {
		//read from current client
		n, err := conn.Read(bs)

		//occurred when the client was closed
		//before: server: CLOSE_WAIT, client: FIN_WAIT_2
		//do close to avoid CLOSE_WAIT
		//after: client: TIME_WAIT
		//take care too many TIME_WAIT
		if err != nil {
			log.Println(err)
			conn.Close()
			delete(conns, conn)
			break
		}

		//write to other clients (exclude current client)
		for out := range conns {
			if out != conn {
				_, err := out.Write(bs[:n])

				//occurred when the client was closed
				if err != nil {
					log.Println(err)
					out.Close()
					delete(conns, out)
					continue
				}
			}
		}
	}
}
