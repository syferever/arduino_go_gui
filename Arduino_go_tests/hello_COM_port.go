package main

import (
	// "fmt"
	"bufio"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{
		Name: "COM10",
		Baud: 9600,
	}
	c.ReadTimeout = time.Millisecond * 100
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(s)

	// time.Sleep(5 * time.Second)
	for {
		// buf := make([]byte, 128)
		// bytes_read, err := s.Read(buf)
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			log.Println(err)
			continue
		}
		// line = line[:len(line)-1]
		val, err := strconv.ParseFloat(line, 64)
		if err != nil {
			log.Println(err)
		}
		log.Println("Received from serial:", val, "bytes:", []rune(line))
		// time.Sleep(time.Second)
	}
}
