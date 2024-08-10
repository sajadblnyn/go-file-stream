package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type FileServer struct {
}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3001")
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go fs.readLoop(conn)

	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buff := new(bytes.Buffer)
	for {
		var size int64

		binary.Read(conn, binary.LittleEndian, &size)

		n, err := io.CopyN(buff, conn, size)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println(buff.Bytes())
		fmt.Printf("received %d bytes over the network \n", n)

	}
}

func sendFile(size int) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3001")
	if err != nil {
		return err
	}
	binary.Write(conn, binary.LittleEndian, int64(size))

	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("received %d bytes over the network \n", n)
	return nil
}

func main() {
	go func() {
		sendFile(1000)
	}()
	fileServer := &FileServer{}
	fileServer.start()
}
