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

type PP2PLink struct {
	Ind   chan PP2PLink_Ind_Message
	Req   chan PP2PLink_Req_Message
	Run   bool
	Cache map[string]net.Conn
}

// Inicia
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
			// fmt.Println(err)
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
						//fmt.Println("!!!!!!!!"+string(content))
						// fmt.Println("????????"+actual)
						msg := PP2PLink_Ind_Message{
							From:    conn.RemoteAddr().String(),
							Message: s[1],
							IpCorreto: s[0]}

						module.Ind <- msg
						// fmt.Println(msg)
					}
				}
			}()
		}
	}()

	go func() {
		for {
			message := <-module.Req
			module.Send(message)
		}
	}()

}

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

	fmt.Fprintf(conn, message.IpCorreto + "," +message.Message) // Escreve na conexão o ip recebido, mensagem recebida
	
}
