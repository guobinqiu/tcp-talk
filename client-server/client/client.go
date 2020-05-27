package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	host string
	port int
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.StringVar(&host, "h", "127.0.0.1", "server host")
	flag.IntVar(&port, "p", 9090, "server port")
}

func main() {
	flag.Parse()

	addr := host + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("client connected to %s\n", addr)

	defer conn.Close()

	bs := make([]byte, 1024)
	for {
		//read from console
		n, _ := os.Stdin.Read(bs)

		//write to server
		_, err := conn.Write(bs[:n])

		//occurred when server was closed
		//before: server: FIN_WAIT_2, client: CLOSE_WAIT
		//do client close
		//after: server: TIME_WAIT
		if err != nil {
			log.Println(err)
			break
		}
	}
}
