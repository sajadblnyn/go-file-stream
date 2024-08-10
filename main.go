package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type FileServer struct {
}

func (fs *FileServer) start() {
	ln, err := net.Listen("tcp", ":3005")
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

		fmt.Println(buff.Bytes())
		if err != nil {
			log.Fatal(err.Error())

		}

		f, err := os.OpenFile("2.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		defer f.Close()

		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = io.Copy(f, buff)
		if err != nil {
			log.Fatal(err.Error())

		}
		// fmt.Println(buff.Bytes())
		fmt.Printf("received %d bytes over the network \n", n)

	}
}

func sendFile() error {
	f, err := os.Open("1.txt")
	defer f.Close()
	if err != nil {
		return err
	}
	file := make([]byte, 10000000000)
	m, _ := io.ReadFull(f, file)

	conn, err := net.Dial("tcp", ":3005")
	if err != nil {
		return err
	}
	binary.Write(conn, binary.LittleEndian, int64(m))

	n, err := io.CopyN(conn, bytes.NewReader(file), int64(m))
	if err != nil {
		return err
	}
	fmt.Printf("received %d bytes over the network \n", n)
	return nil
}

func main() {
	go func() {
		err := sendFile()
		if err != nil {
			panic(err)
		}
	}()
	fileServer := &FileServer{}
	fileServer.start()
}
