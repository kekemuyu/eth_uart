package main

import (
	"fmt"
	"log"
	"net"
	"test/eth_uart/com"
	"time"
)

func main() {
	com, err := com.New("/dev/ttyUSB2", 115200, "/proc/eth_uart/uio")
	if err != nil {
		log.Fatal(err)
	}
	go com.Run()

	laddr := &net.TCPAddr{
		IP:   net.ParseIP("192.168.4.2"),
		Port: 9000,
	}

	raddr := &net.TCPAddr{
		IP:   net.ParseIP("192.168.4.1"),
		Port: 9000,
	}

	conn, cerr := net.DialTCP("tcp", laddr, raddr)
	fmt.Println(conn, cerr)
	for {

		conn, cerr = net.DialTCP("tcp", laddr, raddr)

		fmt.Println("dial tcp", conn, err)
		if conn != nil {
			break
		}
	}

	for {
		conn.Write([]byte("hello"))
		time.Sleep(time.Second)
	}

}
