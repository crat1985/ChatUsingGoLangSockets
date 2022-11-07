package main

import (
	"log"
	"net"
)

type user struct {
	socket net.Conn
	pseudo string
}

var users []user
var pseudoString string

func main() {
	server, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go processClient(conn)
	}
}

func processClient(conn net.Conn) {
	for {
		pseudoOk, pseudoInfos := verifyPseudo(conn)
		if !pseudoOk {
			conn.Write([]byte(pseudoInfos))
			continue
		}
		_, err := conn.Write([]byte("next"))
		if err != nil {
			return
		}
		break
	}
	userTemp := user{socket: conn, pseudo: pseudoString}
	users = append(users, userTemp)
	go listenForMsg(userTemp)
}

func verifyPseudo(conn net.Conn) (bool, string) {
	pseudoSlice := make([]byte, 1024)
	n, err := conn.Read(pseudoSlice)
	if err != nil {
		return false, "Erreur lors de la lecture du pseudo !"
	}
	pseudoString = string(pseudoSlice[:n])
	if len(pseudoString) < 5 {
		return false, "Pseudo trop court : le pseudo doit au moins contenir 5 caractères !"
	}
	return true, ""
}

func listenForMsg(u user) {
	slice := make([]byte, 1024)
	for {
		n, err := u.socket.Read(slice)
		if err != nil {
			log.Printf("%s disconnected\n", u.pseudo)
			removeElement(u)
			break
		}
		log.Printf("New message from %s : %s", u.pseudo, string(slice[:n]))
		for _, value := range users {
			_, err := value.socket.Write([]byte(u.pseudo + "> " + string(slice[:n])))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func removeElement(element user) {
	var index int = -1
	for i, value := range users {
		if value == element {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	users = append(users[:index], users[index+1:]...)
}
