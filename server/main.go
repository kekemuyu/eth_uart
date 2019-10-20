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

	l, err := net.Listen("tcp", "192.168.4.2:9000")
	if err != nil {
		log.Fatal(err)
	}

	bs := make([]byte, 1024)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("get conn", err)
			continue
		}

		go func() {
			for {

				conn.Write([]byte("hello2"))
				time.Sleep(time.Second)
			}
		}()
		go func() {
			for {
				n, err := conn.Read(bs)
				if err != nil {
					fmt.Println(err)
					break
				}
				fmt.Println("get client string data:", n, err, string(bs[:n]))
			}
		}()

	}

}
