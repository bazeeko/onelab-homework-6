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
	// ln, err := net.Listen("tcp", s.Addr)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer ln.Close()

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", s.Addr)
	if err != nil {
		fmt.Println(err)
	}
	defer listener.Close()

	s.done = make(chan struct{})

	// serverShutdown := make(chan os.Signal, 1)
	// signal.Notify(serverShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

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

	// <-serverShutdown
	log.Println("Server Stopped")
}

func handlerRequest(conn net.Conn) {
	defer conn.Close()

	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	str := strings.Replace(request, "\n", "", -1)

	number, err := strconv.Atoi(str)
	if err != nil {
		fmt.Println("Error: NaN", err)
	}

	conn.Write([]byte(strconv.Itoa(number * number)))
}
