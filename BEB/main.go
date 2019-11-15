package BEB

import "fmt"

import PP2PLink "../Link"

// Canal que envia mensagem aos usuarios do chat
type Envia_Mensagem struct {
	Addresses []string
	IpCorreto string
	Message   string
}

// Canal que recebe a mensagem de algum usuario
type Recebe_Mensagem struct {
	From    string
	Message string
	IpCorreto string
}

// Um controlador
type Modulo struct {
	Ind      chan Recebe_Mensagem
	Req      chan Envia_Mensagem
	Pp2plink PP2PLink.PP2PLink
}

// Inicia o controlador
func (module Modulo) Init(address string) {

	module.Pp2plink = PP2PLink.PP2PLink{
		Req: make(chan PP2PLink.PP2PLink_Req_Message),
		Ind: make(chan PP2PLink.PP2PLink_Ind_Message)}
	module.Pp2plink.Init(address)
	module.Start()

}

// Deixa o controlador rodando
func (module Modulo) Start() {

	go func() {
		for {
			select {
			case y := <-module.Req: // Caso em que é solicitado ao controlador enviar uma mensagem, vem lá do chat
				module.FazEnvioDaMensagem(y)
			case y := <-module.Pp2plink.Ind: // Caso em que é solicitado ao controlador receber uma mensagem
				module.fazRecebimentoDaMensagem(PP2PLink2BEB(y))
			}
		}
	}()

}

// Funcao responsavel por enviar uma mensagem
func (module Modulo) FazEnvioDaMensagem(message Envia_Mensagem) { 

	for i := 0; i < len(message.Addresses); i++ {
		msg := criaMensagem(message) // Cria uma mensagem do tipo PP2 para poder enviar pro PP2
		msg.To = message.Addresses[i]
		module.Pp2plink.Req <- msg //Envia a mensagem pro PP2 que é quem envia a mensagem de Fato
		fmt.Println("Enviado para " + message.Addresses[i])
	}

}

//Funcao responsavel por receber uma mensagem
func (module Modulo) fazRecebimentoDaMensagem(message Recebe_Mensagem) {

	module.Ind <- message

}

func criaMensagem(message Envia_Mensagem) PP2PLink.PP2PLink_Req_Message {

	return PP2PLink.PP2PLink_Req_Message{
		To:      message.Addresses[0],
		IpCorreto:		 message.IpCorreto,
		Message: message.Message}

}

func PP2PLink2BEB(message PP2PLink.PP2PLink_Ind_Message) Recebe_Mensagem {

	return Recebe_Mensagem{
		From:    message.From,
		Message: message.Message,
		IpCorreto: message.IpCorreto}

}
