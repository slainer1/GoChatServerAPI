package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	//Write to server
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//read response
	bs := make([]byte, 1024)
	n, err := conn.Read(bs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", string(bs[:n]))
}
