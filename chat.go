package main

import (
	"bufio"
	"fmt"
	"os"

	BEB "./BEB"
)

var historico []string // Historico do chat

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Especifique pelo menos um endereço: porta!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println("Enderecos:", addresses)

	beb := BEB.Modulo{
		EnviaMensagem:  make(chan BEB.Envia_Mensagem),		// Canal que envia mensagem
		RecebeMensagem: make(chan BEB.Recebe_Mensagem),		// Canal que recebe mensagem
		NovoUsuario:    make(chan BEB.Envia_Novo_Usuario),	// Canal que envia um novo usuario
		RecebeUsuario:  make(chan BEB.Recebe_Usuario),		// Canal que recebe um novo usuario
		NovoGrupo:      make(chan BEB.Envia_Novo_Grupo),	// Canal que recebe dados de um grupo (usuarios e historico)
		RecebeGrupo:    make(chan BEB.Recebe_Grupo)}		// Canal que recebe dados de um grupo (usuarios e historico)

	beb.Init(addresses[0])

	fmt.Println("-------------COMANDOS---------------")
	fmt.Println("1) Enviar mensagem")
	fmt.Println("2) Visualizar histórico de mensagens")
	fmt.Println("3) Entrar em um chat")
	fmt.Println("4) Visualizar membros do chat")
	fmt.Println("------------------------------------")

	// Enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string
		var op string
		var ip string

		for {
			if scanner.Scan() {
				op = scanner.Text()
				switch op {
				case "1":
					fmt.Print("Enviar mensagem: ")
					if scanner.Scan() {
						msg = scanner.Text()
						historico = append(historico, addresses[0] + ":" + msg)
					}
					req := BEB.Envia_Mensagem{
						Addresses: addresses[1:],
						IpCorreto: addresses[0],
						Message:   msg}
					beb.EnviaMensagem <- req
				case "2":
					fmt.Println("----HISTORICO----")
					for _, i := range historico {
						fmt.Println(i)
					}
					fmt.Println("-----------------")
				case "3":
					fmt.Print("Pedir p/ Ip: ") // Ip de algum usuario do chat que deseja entrar
					if scanner.Scan() {
						ip = scanner.Text()
					}
					req := BEB.Envia_Novo_Usuario{ // Envia o novo usuario para esse ip
						Address:   ip,
						IpCorreto: addresses[0],
						Tag:       "0"}
					beb.NovoUsuario <- req
					fmt.Println("Você entrou no chat!")
				case "4":
					fmt.Println("-----MEMBROS-----")
					for _, i := range addresses {
						fmt.Println(i)
					}
					fmt.Println("-----------------")
				}
			}
		}
	}()

	// Rotina responsavel por receber novas mensagens
	go func() {
		for {
			mensagemRecebida := <-beb.RecebeMensagem
			fmt.Printf("Mensagem de %v: %v\n", mensagemRecebida.IpCorreto, mensagemRecebida.Message)
			historico = append(historico, mensagemRecebida.IpCorreto + ":" + mensagemRecebida.Message)
		}
	}()

	// Rotina responsavel por receber dados de um grupo
	go func() {
		for {
			grupoRecebido := <-beb.RecebeGrupo
			addresses = append(addresses, grupoRecebido.Addresses...)
			historico = append(historico, grupoRecebido.Historico...)

		}
	}()

	// Rotina responsavel por receber usuarios novos
	go func() {
		for {
			usuarioRecebido := <-beb.RecebeUsuario

			// Tag responsavel por indicar se o usuario deve espalhar o novo usuario
			// 1 - Usuario novo solicita para um usuario do chat que quer entrar no chat
			// 2 - Usuario do chat manda o ip desse usuario novo para todos os outros usuarios do chat, para que todos consigam o adicionar e conversar
			// 3 - Porem os usuarios que receberam nao devem espalhar o ip de novo, pois isso ja esta sendo feito
			// A tag controla isso

			// Se a tag for 0 o ip deve ser espalhado
			if usuarioRecebido.Tag == "0" {
				// Manda o novo usuario para todos os outros usuarios
				for i := 1; i < len(addresses); i++ {
					req := BEB.Envia_Novo_Usuario{
						Address:   addresses[i],
						IpCorreto: usuarioRecebido.IpCorreto,
						Tag:       "1"} // Adiciona tag 1 para avisar que nao deve espalhar mais
					beb.NovoUsuario <- req
				}

				// Manda todos os usuarios e o historico de conversa para o novo membro
				req := BEB.Envia_Novo_Grupo{
					Addresses: addresses,
					IpCorreto: usuarioRecebido.IpCorreto,
					Historico: historico}
				beb.NovoGrupo <- req

			}
			fmt.Println(usuarioRecebido.IpCorreto + " entrou no chat!")
			// Adiciona o novo usuario a lista de participantes
			addresses = append(addresses, usuarioRecebido.IpCorreto)
		}
	}()

	blq := make(chan int)
	<-blq
}
