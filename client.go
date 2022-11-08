package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var addr string = "localhost:8888"
var protocol string = "tcp"
var input string
var message string = "Pseudo : "
var slice []byte
var stringDatas string
var quit bool = false

func main() {
	askForDetails()
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		sendMessages(conn)
		listenForMsg(conn)
		if stringDatas == "next" {
			message = "Entre ton message (:quit pour quitter) :\n"
			break
		}
	}
	go func() {
		for {
			listenForMsg(conn)
		}
	}()
	for !quit {
		sendMessages(conn)
	}
}

func sendMessages(conn net.Conn) {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		log.Println(err)
		return
	}
	if isPrefix {
		fmt.Println("Message to long !")
	}
	input = string(line)
	if input == ":quit" {
		quit = true
		conn.Close()
		log.Println("Left chat")
		return
	}
	conn.Write([]byte(input))
}

func askForDetails() {
	fmt.Printf("Addresse of server (d for %s) : ", addr)
	var tempAddr string
	fmt.Scan(&tempAddr)
	if strings.ToLower(tempAddr) != "d" {
		addr = tempAddr
	}
	fmt.Printf("Protocol of server (d for %s) : ", protocol)
	var tempProtocol string
	fmt.Scan(&tempProtocol)
	if strings.ToLower(tempProtocol) != "d" {
		protocol = tempProtocol
	}
}

func listenForMsg(conn net.Conn) {
	slice = make([]byte, 1024)
	n, err := conn.Read(slice)
	if err != nil {
		if !quit {
			log.Fatal(err)
		}
	}
	stringDatas = string(slice[:n])
	if stringDatas != "next" {
		fmt.Println(stringDatas)
	}
}
