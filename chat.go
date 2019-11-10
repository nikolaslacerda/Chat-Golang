package main

import (
	"bufio"
	"fmt"
	"os"

	BEB "./BEB"
)

type Channel struct {
	mensagem string
}

var channels []Channel

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Especifique pelo menos um endereço: porta!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println("Enderecos:", addresses)

	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}

	beb.Init(addresses[0])

	fmt.Println("-------------COMANDOS--------------")
	fmt.Println("1)Enviar mensagem")
	fmt.Println("2)Visualizar histórico de mensagens")
	fmt.Println("-----------------------------------")

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string
		var op string

		for {
			if scanner.Scan() {
				op = scanner.Text()
				switch op {
				case "1":
					fmt.Print("Enviar msg: ")
					if scanner.Scan() {
						msg = scanner.Text()
						channels = append(channels, Channel{mensagem: msg})
					}
					req := BEB.BestEffortBroadcast_Req_Message{
						Addresses: addresses[1:],
						Message:   msg}
					beb.Req <- req
				case "2":
					fmt.Printf("%+v\n", channels)
				}
			}

		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind
			fmt.Printf("Message from %v: %v\n", in.From, in.Message)
		}
	}()

	blq := make(chan int)
	<-blq
}

/*
go run chat.go 127.0.0.1:5001  127.0.0.1:6001
go run chat.go 127.0.0.1:6001  127.0.0.1:5001
*/
