package main

import (
	"bufio"
	"fmt"
	"os"

	BEB "./BEB"
)

type Channel struct {
	EnviadoPor string
	Mensagem   string
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
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message),
		NewUser: make(chan string), 
		RefreshList: make(chan string)}

	beb.Init(addresses[0])

	fmt.Println("-------------COMANDOS--------------")
	fmt.Println("1)Enviar mensagem")
	fmt.Println("2)Visualizar histórico de mensagens")
	fmt.Println("3)Pedir para entrar em um chat")
	fmt.Println("4)Mostrar participantes")
	fmt.Println("-----------------------------------")

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string
		var ip string
		var ipL []string
		var op string

		for {
			if scanner.Scan() {
				op = scanner.Text()
				switch op {
				case "1":
					fmt.Print("Enviar msg: ")
					if scanner.Scan() {
						msg = scanner.Text()
						channels = append(channels, Channel{EnviadoPor: addresses[0], Mensagem: msg})
					}
					req := BEB.BestEffortBroadcast_Req_Message{
						Addresses: addresses[1:],
						Message:   msg}
					beb.Req <- req
				case "2":
					fmt.Printf("%+v\n", channels)
				case "3":
					fmt.Printf("IP: ")
					if scanner.Scan() {
						ip = scanner.Text()
						ipL = append(ipL, ip)
						// msg = addresses[0] + " Gostaria de participar do chat"
						msg = "Entrou no chat!"
					}
					req := BEB.BestEffortBroadcast_Req_Message{
						Addresses: ipL,
						Message:   msg}
					beb.Req <- req
				case "4":
					fmt.Println("Participantes:", addresses)
				}
			}

		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind
			fmt.Printf("Mensagem de %v: %v\n", in.From, in.Message)
			channels = append(channels, Channel{EnviadoPor: in.From, Mensagem: in.Message})
		}
	}()

	go func() {
		for {
			newU := <-beb.NewUser
			req := BEB.BestEffortBroadcast_Req_Message{
				Addresses: addresses[1:],
				Message:   "Atualizem ai!",
				Ip: newU}

			beb.Req <- req

			addresses = append(addresses, newU)

		}
	}()

	go func() {
		for {
			newU := <-beb.RefreshList
			addresses = append(addresses, newU)
		}
	}()

	blq := make(chan int)
	<-blq
}

/*
go run chat.go 127.0.0.1:5001  127.0.0.1:5002
go run chat.go 127.0.0.1:5002  127.0.0.1:5001
*/
