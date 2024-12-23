package main

import (
	"fmt"
	"log"
	"net"
)
type Message struct {
	from string 
	payload []byte
}
type Server struct {
	servAddr string
	quitch chan struct{}
	ln net.Listener
	msgch chan Message
}
func NewServer(servAddr string) *Server {
	return &Server{
		servAddr: servAddr,
		quitch: make(chan struct{}),
		msgch: make(chan Message , 10),
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.servAddr)
	if err!= nil {
        return err
    }
	fmt.Println("Starting")
	defer ln.Close()
	s.ln = ln
	go s.AcceptLoop()
	<-s.quitch
	close(s.msgch)
	return nil
}
func (s *Server) AcceptLoop()  {
	// Implement the accept loop here
	for{
		conn, err := s.ln.Accept()
        if err!= nil {
            fmt.Println("Accept error", err)
			continue
        }
		go s.readLoop(conn)
	}
}
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for{
		n, err := conn.Read(buf)
        if err!= nil {
            fmt.Println("Read error", err)
            continue
        }
        // Process the received data here
        s.msgch <- Message{
			from: conn.RemoteAddr().String(),
            payload: buf[:n],
		}
		conn.Write([]byte("Thanks for reading"))
	}
}
func (s *Server) Stop() {
    close(s.quitch)
    s.ln.Close()
}
func main() {
	serv := NewServer(":8080")
    log.Fatal( serv.Start())
    // if err!= nil {
    //     fmt.Println("Error starting server:", err)
    //     return
    // }
    // defer serv.Stop()
    // serv.AcceptLoop()
	go func(){
		for msg := range serv.msgch {
            fmt.Println("Broadcasting:", string(msg.from) ,"payload" , msg.payload)
        }
	}()
}