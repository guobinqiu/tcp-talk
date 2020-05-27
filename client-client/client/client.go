package main

import (
	"flag"
	"fmt"
	"io"
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

	go sendToConsole(conn, os.Stdout)

	sendToServer(os.Stdin, conn)
}

//read from server and write to console
func sendToConsole(in net.Conn, out io.Writer) {
	bs := make([]byte, 1024)
	for {
		//read from server
		n, _ := in.Read(bs)

		//write to console
		out.Write(bs[:n])
	}
}

//read from console and write to server
func sendToServer(in io.Reader, out net.Conn) {
	bs := make([]byte, 1024)
	for {
		//read from console
		n, _ := in.Read(bs)

		//write to server
		_, err := out.Write(bs[:n])

		//occurred when server was closed
		if err != nil {
			log.Println(err)
			break
		}
	}
}
