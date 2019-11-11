// Construido como parte da disciplina de Sistemas Distribuidos
// Semestre 2018/2  -  PUCRS - Escola Politecnica
// Estudantes:  Andre Antonitsch e Rafael Copstein
// Professor: Fernando Dotti  (www.inf.pucrs.br/~fldotti)
// Algoritmo baseado no livro:
// Introduction to Reliable and Secure Distributed Programming
// Christian Cachin, Rachid Gerraoui, Luis Rodrigues

package BEB

import (
	"fmt"
	)

import PP2PLink "../Link"

type BestEffortBroadcast_Req_Message struct {
	Addresses []string
	Message   string
	Ip    string
}

type BestEffortBroadcast_Ind_Message struct {
	From    string
	Message string
	Ip    string
}


type BestEffortBroadcast_Module struct {
	Ind      chan BestEffortBroadcast_Ind_Message
	Req      chan BestEffortBroadcast_Req_Message
	NewUser  chan string
	RefreshList chan string
	Pp2plink PP2PLink.PP2PLink
}

func (module BestEffortBroadcast_Module) Init(address string) {

	fmt.Println("Init BEB!")
	module.Pp2plink = PP2PLink.PP2PLink{
		Req: make(chan PP2PLink.PP2PLink_Req_Message),
		Ind: make(chan PP2PLink.PP2PLink_Ind_Message)}
	module.Pp2plink.Init(address)
	module.Start()

}

func (module BestEffortBroadcast_Module) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req:
				module.Broadcast(y)
			case y := <-module.Pp2plink.Ind:
				module.Deliver(PP2PLink2BEB(y))
			}
		}
	}()

}

func (module BestEffortBroadcast_Module) Broadcast(message BestEffortBroadcast_Req_Message) {

	//Até aqui OK
	fmt.Println(message)

	for i := 0; i < len(message.Addresses); i++ {
		msg := BEB2PP2PLink(message)
		msg.To = message.Addresses[i]
		msg.Ip = message.Ip
		//OK
		module.Pp2plink.Req <- msg
		fmt.Println("Sent to " + message.Addresses[i])
	}

}

func (module BestEffortBroadcast_Module) Deliver(message BestEffortBroadcast_Ind_Message) {

	fmt.Println("Received '" + message.Message + "' from " + message.From + "' Ip" + message.Ip)
	module.Ind <- message
	//fmt.Println("# End BEB Received")
	if (message.Message == "Entrou no chat!"){
		module.NewUser <- message.From
	} else if message.Message == "Atualizem ai!" {
		fmt.Println("Velho" + message.Ip)
		fmt.Println("Velho" + message.From)
		module.RefreshList <- message.Ip
	}
	
}

func BEB2PP2PLink(message BestEffortBroadcast_Req_Message) PP2PLink.PP2PLink_Req_Message {

	return PP2PLink.PP2PLink_Req_Message{
		To:      message.Addresses[0],
		Message: message.Message,
		Ip: message.Ip}

}

func PP2PLink2BEB(message PP2PLink.PP2PLink_Ind_Message) BestEffortBroadcast_Ind_Message {

	return BestEffortBroadcast_Ind_Message{
		From:    message.From,
		Message: message.Message,
		Ip: message.Ip}

}
