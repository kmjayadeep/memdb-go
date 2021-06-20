package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	listener         net.Listener
	db               memoryDB
	quit             chan struct{}
	exited           chan struct{}
	connections      map[int]net.Conn
	connCloseTimeout time.Duration
}

func NewServer() *Server {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("failed to create listener", err.Error())
	}

	srv := Server{
		listener:         l,
		quit:             make(chan struct{}),
		exited:           make(chan struct{}),
		db:               newDB(),
		connections:      map[int]net.Conn{},
		connCloseTimeout: 10 * time.Second,
	}
	go srv.serve()
	return &srv
}

func (s *Server) serve() {
	fmt.Println("listening for clients")
	id := 0
	for {
		select {
		case <-s.quit:
			fmt.Println("shutting down the server")
		default:
			tcpListener := s.listener.(*net.TCPListener)
			err := tcpListener.SetDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				fmt.Println("failed to set listener deadline", err.Error())
			}
			conn, err := tcpListener.Accept()

			if oppErr, ok := err.(*net.OpError); ok && oppErr.Timeout() {
				continue
			}

			if err != nil {
				fmt.Println("failed to accept connection", err.Error())
			}

			write(conn, "Welcome to MemoryDB server")
			s.connections[id] = conn

			go func(connID int) {
				fmt.Println("client with id", connID, "joined")
				s.handleConn(conn)
				delete(s.connections, connID)
			}(id)
			id++
		}
	}
}

func write(conn net.Conn, s string) {
	_, err := fmt.Fprintf(conn, "%s\n-> ", s)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		l := strings.ToLower(strings.TrimSpace(scanner.Text()))
		values := strings.Split(l, " ")

		switch {
		case len(values) == 3 && values[0] == "set":
			s.db.set(values[1], values[2])
			write(conn, "OK")
		case len(values) == 2 && values[0] == "get":
			k := values[1]
			val, found := s.db.get(k)
			if !found {
				write(conn, fmt.Sprintf("key %s not found", k))
			} else {
				write(conn, val)
			}
		case len(values) == 2 && values[0] == "delete":
			s.db.delete(values[1])
			write(conn, "OK")
		case len(values) == 1 && values[0] == "exit":
      if err:= conn.Close(); err != nil {
        fmt.Println("could not close connection", err.Error())
      }
    default:
      write(conn, fmt.Sprintf("UNKNOWN : %s", l))
		}
	}

}

func (s *Server) Stop() {
	fmt.Println("Stopping server")
}
