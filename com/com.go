package com

import (
	"fmt"
	"io"
	"os"

	"log"

	"github.com/jacobsa/go-serial/serial"
)

type Com struct {
	Irw       io.ReadWriteCloser
	FileHanle *os.File
}

func New(port string, baund uint, filename string) (com *Com, err error) {
	// Set up options.
	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        baund,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	rw, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
		return com, err
	}
	fmt.Println(filename)

	fhanle, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return com, err
	}
	com = &Com{
		Irw:       rw,
		FileHanle: fhanle,
	}

	return com, nil
}

func (c *Com) Run() {
	// Make sure to close it later.
	defer c.Irw.Close()
	go c.write()
	c.read()
}

func (c *Com) read() {
	bs1 := make([]byte, 1)

	for {
		n, err := c.Irw.Read(bs1)
		if n > 0 {

			dataLen := int(bs1[0])
			curLen := 0
			bs2 := make([]byte, dataLen)
			fmt.Println("dataLen:", dataLen)
			for {
				n, err = c.Irw.Read(bs2[curLen:])
				if n > 0 {
					curLen = curLen + n
					if curLen >= dataLen {
						fmt.Println("com read data :", dataLen, err, bs2)
						_, ferr := c.FileHanle.Write(bs2)

						if ferr != nil {
							fmt.Println(ferr)
						}
						dataLen = 0
						curLen = 0
						break
					}

				}
			}
		}

	}
}

func (c *Com) write() {

	bs := make([]byte, 1024)

	for {
		n, err := c.FileHanle.Read(bs[1:])
		if err != nil {
			fmt.Println("readfile:", n, err)
		}

		if n > 0 {
			bs[0] = byte(n)
			fmt.Println("com send data", err, bs[:n+1])
			_, err := c.Irw.Write(bs[:(n + 1)])
			if err != nil {
				log.Fatalf("port.Write: %v", err)
			}
		}

	}
}
