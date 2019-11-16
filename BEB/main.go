package BEB

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

// Canal que envia um novo usuario a algum usuario
type Envia_Novo_Usuario struct {
	Address string
	IpCorreto string
	Tag string
}

//Canal que recebe um novo usuario
type Recebe_Usuario struct {
	From string
	IpCorreto string
	Tag string
}

type Envia_Novo_Grupo struct {
	Addresses []string
	Historico []string
	IpCorreto string
}

type Recebe_Grupo struct {
	Addresses []string
	Historico []string
	From string
}


// Um controlador
type Modulo struct {
	Ind      	chan Recebe_Mensagem
	Req      	chan Envia_Mensagem
	NovoUsuario chan Envia_Novo_Usuario
	RecebeUsuario chan Recebe_Usuario
	NovoGrupo chan Envia_Novo_Grupo
	RecebeGrupo chan Recebe_Grupo
	Pp2plink PP2PLink.PP2PLink
}

// Inicia o controlador
func (module Modulo) Init(address string) {

	module.Pp2plink = PP2PLink.PP2PLink{
		Req: make(chan PP2PLink.PP2PLink_Req_Message),
		Ind: make(chan PP2PLink.PP2PLink_Ind_Message),
		NovoUsuario: make(chan PP2PLink.PP2PLink_Novo_Usuario),
		RecebeUsuario: make(chan PP2PLink.PP2PLink_Recebe_Usuario),
		NovoGrupo: make(chan PP2PLink.PP2PLink_Novo_Grupo),
		RecebeGrupo: make(chan PP2PLink.PP2PLink_Recebe_Grupo)}
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
			case y := <-module.NovoUsuario: // Caso em que é solicitado ao controlador enviar um usuario
				module.InsereNovoUsuario(y)
			case y := <-module.Pp2plink.RecebeUsuario: // Caso em que é solicitado ao controlador receber um usuario
				module.fazRecebimentoDoUsuario(PP2PLink2BEB2(y))
			case y := <-module.NovoGrupo: // Caso em que é solicitado ao controlador enviar um grupo
				module.InsereNovoGrupo(y)
			case y := <-module.Pp2plink.RecebeGrupo: // Caso em que é solicitado ao controlador receber um grupo
				module.fazRecebimentoDoGrupo(PP2PLink2BEB3(y))
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
		//fmt.Println("Enviado para " + message.Addresses[i])
	}

}

// Funcao responsavel por enviar uma mensagem
func (module Modulo) InsereNovoUsuario(message Envia_Novo_Usuario) { 
	msg := criaNovoUsuario(message)
	module.Pp2plink.NovoUsuario <- msg
	//fmt.Println("Você enviou" + message.IpCorreto + "para" + message.Address)
}

// Funcao responsavel por enviar uma mensagem
func (module Modulo) InsereNovoGrupo(message Envia_Novo_Grupo) { 
	msg := criaNovoGrupo(message)
	module.Pp2plink.NovoGrupo <- msg
	//fmt.Println("Você enviou o grupo")
	//fmt.Print(message.Addresses)
	//fmt.Print(" para " + message.IpCorreto)
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

func (module Modulo) fazRecebimentoDoUsuario(message Recebe_Usuario) {
	module.RecebeUsuario <- message
}

func (module Modulo) fazRecebimentoDoGrupo(message Recebe_Grupo) {
	module.RecebeGrupo <- message
}

func criaNovoUsuario(message Envia_Novo_Usuario) PP2PLink.PP2PLink_Novo_Usuario {
	return PP2PLink.PP2PLink_Novo_Usuario{
		Adress:      message.Address,
		Tag: message.Tag,
		IpCorreto:		 message.IpCorreto}
}

func criaNovoGrupo(message Envia_Novo_Grupo) PP2PLink.PP2PLink_Novo_Grupo {
	return PP2PLink.PP2PLink_Novo_Grupo{
		Addresses:      message.Addresses,
		Historico: 		message.Historico,
		IpCorreto:		message.IpCorreto}
}

func PP2PLink2BEB(message PP2PLink.PP2PLink_Ind_Message) Recebe_Mensagem {

	return Recebe_Mensagem{
		From:    message.From,
		Message: message.Message,
		IpCorreto: message.IpCorreto}

}

func PP2PLink2BEB2(message PP2PLink.PP2PLink_Recebe_Usuario) Recebe_Usuario {

	return Recebe_Usuario{
		From:    message.From,
		Tag: message.Tag,
		IpCorreto: message.IpCorreto}

}

func PP2PLink2BEB3(message PP2PLink.PP2PLink_Recebe_Grupo) Recebe_Grupo {

	return Recebe_Grupo{
		Addresses:    message.Adresses,
		Historico: 	  message.Historico,
		From: message.From}

}