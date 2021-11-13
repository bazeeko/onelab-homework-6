package tcpserver

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Server struct
type Server struct {
	Addr     string
	MaxConns int
	done     chan struct{}
	wg       sync.WaitGroup
	sync.RWMutex
}

// NewServer returns the new Server struct with the given parameters
func NewServer(addr string, conns int) *Server {
	return &Server{
		Addr:     addr,
		MaxConns: conns,
		done:     make(chan struct{}),
	}
}

// Run starts the server
func (s *Server) Run(ctx context.Context) {
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", s.Addr)
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()

	s.done = make(chan struct{})

	resquestChannel := make(chan net.Conn, s.MaxConns)

	go func() {
		for {
			select {
			case <-s.done:
				return

			case c, ok := <-resquestChannel:
				if ok {
					go handlerRequest(c)
				}
			}
		}
	}()

LOOP:
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		select {
		case <-ctx.Done():
			s.done <- struct{}{}
			close(resquestChannel)
			break LOOP

		default:
			resquestChannel <- conn
		}
	}

	fmt.Println()
	log.Println("Server Stopped")
}

func handlerRequest(conn net.Conn) {
	defer conn.Close()

	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	str := strings.Replace(request, "\n", "", -1)

	number, err := strconv.Atoi(str)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	conn.Write([]byte(strconv.Itoa(number * number)))
}
