package PP2PLink

import "fmt"
import "net"
import "strings"

type PP2PLink_Req_Message struct {
	To          string
	Message 	string
	IpCorreto	string
}

type PP2PLink_Ind_Message struct {
	From    string
	Message string
	IpCorreto string
}

type PP2PLink_Recebe_Usuario struct {
	From    string
	IpCorreto string
	Tag string
}

type PP2PLink_Novo_Usuario struct {
	Adress    string
	IpCorreto string
	Tag string
}

type PP2PLink_Recebe_Grupo struct {
	From string
	Adresses    []string
	Historico []string
	IpCorreto string
}

type PP2PLink_Novo_Grupo struct {
	To string
	Addresses    []string
	IpCorreto string
	Historico []string
}

type PP2PLink struct {
	Ind   chan PP2PLink_Ind_Message
	Req   chan PP2PLink_Req_Message
	NovoUsuario   chan PP2PLink_Novo_Usuario
	RecebeUsuario   chan PP2PLink_Recebe_Usuario
	NovoGrupo   chan PP2PLink_Novo_Grupo
	RecebeGrupo   chan PP2PLink_Recebe_Grupo
	Run   bool
	Cache map[string]net.Conn
}

func (module PP2PLink) Init(address string) {
	if module.Run { return 
	}
	module.Cache = make(map[string]net.Conn)
	module.Run = true
	module.Start(address)
}

func (module PP2PLink) Start(address string) {
	go func() {
		listen, _ := net.Listen("tcp4", address)
		for {
			// aceita repetidamente tentativas novas de conexao
			conn, err := listen.Accept()
			if err != nil {
				continue
			}
			go func() {
				// quando aceita, repetidamente recebe mensagens na conexao TCP (sem fechar)
				// e passa para cima
				for {
					
					buf := make([]byte, 1014)
					Len, err := conn.Read(buf) // Leio os dados da conexão
					if err != nil {
						continue
					}

					content := make([]byte, Len)
					
					copy(content, buf)

					if !strings.Contains(string(content), "@$@"){
						//fmt.Println("WHY")
					}
					
					for _,actual := range strings.Split(string(content), "@$@"){
						if len(actual) == 0 {
							continue
						}
						s := strings.Split(actual, ",")

						//Verifica se é pra enviar mensagem ou se e pra enviar um novo usuario ou grupo
						
						if (s[0] == "M"){ // Envia nova mensagem
							msg := PP2PLink_Ind_Message{
							From:    conn.RemoteAddr().String(),
							Message: s[2],
							IpCorreto: s[1]}
							module.Ind <- msg

						}
						if (s[0] == "U") { // Novo usuario
							msg1 := PP2PLink_Recebe_Usuario{
							From:    conn.RemoteAddr().String(),
							IpCorreto: s[1],
							Tag: s[3]}
							module.RecebeUsuario <- msg1
						}
						if (s[0] == "G") { // Novo grupo
							u := strings.Split(s[1], "/")
							h := strings.Split(s[2], "/")
				
							msg2 := PP2PLink_Recebe_Grupo{
							From:    conn.RemoteAddr().String(),
							Adresses:    u,
							IpCorreto: "!",
							Historico: h}
							module.RecebeGrupo <- msg2
						}
					}
				}
			}()
		}
	}()

	// Rotina que envia mensagem
	go func() {
		for {
			mensagem := <-module.Req
			module.Send(mensagem)
		}
	}()

	// Rotina que envia usuario
	go func() {
		for {
			usuario := <-module.NovoUsuario
			module.SendUser(usuario)
		}
	}()

	// Rotina que envia dados de um chat
	go func() {
		for {
			grupo := <-module.NovoGrupo
			module.SendGrupo(grupo)
		}
	}()

}

// Nas funcoes Send todos os dados sao passados por strings

func (module PP2PLink) Send(message PP2PLink_Req_Message) {

	var conn net.Conn
	var ok bool
	var err error
	

	// ja existe uma conexao aberta para aquele destinatario?
	if conn, ok = module.Cache[message.To]; ok {
		//fmt.Printf("Reusing connection %v to %v\n", conn.LocalAddr(), message.To)
	} else { // se nao tiver, abre e guarda na cache
		conn, err = net.Dial("tcp", message.To)
		if err != nil {
			// fmt.Println(err)
			return
		}
		module.Cache[message.To] = conn
	}

	fmt.Fprintf(conn, "M," + message.IpCorreto + "," +message.Message) // Escreve na conexão do ip a mensagem 
	
	
}

func (module PP2PLink) SendUser(message PP2PLink_Novo_Usuario) {

	var conn net.Conn
	var ok bool
	var err error


	if conn, ok = module.Cache[message.Adress]; ok {
	} else {
		conn, err = net.Dial("tcp", message.Adress)
		if err != nil {
			return
		}
		module.Cache[message.Adress] = conn
	}

	fmt.Fprintf(conn, "U," + message.IpCorreto + "," + message.Adress + "," + message.Tag) // Escreve na conexão do ip o usuario
	
}
func (module PP2PLink) SendGrupo(message PP2PLink_Novo_Grupo) {

	var conn net.Conn
	var ok bool
	var err error


	if conn, ok = module.Cache[message.IpCorreto]; ok {
	} else {
		conn, err = net.Dial("tcp", message.IpCorreto)
		if err != nil {
			return
		}
		module.Cache[message.IpCorreto] = conn
	}

	var ad string
	var ad2 string

	for i := 0; i < len(message.Addresses); i++ {
		ad += message.Addresses[i] + "/"
	}
	ad = ad[:len(ad)-1]

	for i := 0; i < len(message.Historico); i++ {
		ad2 += message.Historico[i] + "/"
	}
	ad2 = ad2[:len(ad2)-1]

	fmt.Fprintf(conn, "G," + ad + "," + ad2) // Escreve na conexão do ip o grupo
	
}