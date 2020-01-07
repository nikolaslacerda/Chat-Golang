# Chat com gravação da conversação

Implementação de um chat via rede capaz de fazer a troca de mensagens entre usuários, adicionar novos membros a conversa e manter o histórico da conversação. Programa criado utilizando a linguagem Golang  
Trabalho final da disciplina **Modelos para Computação Concorrente**

## Tecnologias Usadas:
- Go

## Como usar:

Abra o terminal na pasta *projeto* e digite o seguinte comando:

```
go run chat.go ip_1  ip_2
````

Onde os ip's devem ser preenchidos, por exemplo:

```
go run chat.go 127.0.0.1:5001  127.0.0.1:5002
````

Uma janela de chat conectada do ip_1 para o ip_2 será aberta, o mesmo deve ser feito do ip_2 para o ip_1

```
go run chat.go 127.0.0.1:5002  127.0.0.1:5001
````

Deste modo, o chat já estará executando
