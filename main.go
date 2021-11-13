package main

import (
	"bufio"
	"context"
	"fmt"
	"homework/tcpserver"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	tcpServer := &tcpserver.Server{
		Addr: ":7171",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go tcpServer.Run(ctx)

	serverShutdown := make(chan os.Signal, 1)
	signal.Notify(serverShutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	time.Sleep(time.Second)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := net.Dial("tcp", "127.0.0.1"+tcpServer.Addr)
				if err != nil {
					log.Fatalln(err)
				}

				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Write a number: ")

				text, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println(err)
				}

				_, err = fmt.Fprintf(conn, text)
				if err != nil {
					fmt.Println(err)
					break
				}

				message, err := bufio.NewReader(conn).ReadString('\n')

				fmt.Println("Response:", message)
			}
		}
	}(ctx)

	<-serverShutdown
	fmt.Println()
	log.Println("Server Stopped")
}
