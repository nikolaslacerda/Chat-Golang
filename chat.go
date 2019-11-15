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

	beb := BEB.Modulo{
		Req: make(chan BEB.Envia_Mensagem),
		Ind: make(chan BEB.Recebe_Mensagem)}

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
					fmt.Print("Enviar mensagem: ")
					if scanner.Scan() {
						msg = scanner.Text()
						channels = append(channels, Channel{EnviadoPor: addresses[0], Mensagem: msg})
					}
					req := BEB.Envia_Mensagem{
						Addresses: addresses[1:],
						IpCorreto: addresses[0],
						Message:   msg}
					beb.Req <- req
				case "2":
					fmt.Printf("%+v\n", channels)
				}
			}

		}
	}()

	// Recebe mensagem
	go func() {
		for {
			mensagemRecebida := <-beb.Ind
			fmt.Printf("Mensagem de %v: %v\n", mensagemRecebida.IpCorreto, mensagemRecebida.Message)
			channels = append(channels, Channel{EnviadoPor: mensagemRecebida.IpCorreto, Mensagem: mensagemRecebida.Message})
		}
	}()

	blq := make(chan int)
	<-blq
}

/*
go run chat.go 127.0.0.1:5001  127.0.0.1:5002
go run chat.go 127.0.0.1:5002  127.0.0.1:5001
*/
