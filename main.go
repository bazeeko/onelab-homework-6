package main

import (
	"bufio"
	"context"
	"fmt"
	"homework/tcpserver"
	"log"
	"net"
	"os"
	"time"
)

// func main() {
// 	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
// 		do(request.Context())
// 		_, _ = writer.Write([]byte("all done"))
// 	})

// 	done := make(chan os.Signal, 1)
// 	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// 	srv := &http.Server{
// 		Addr:    ":8080",
// 		Handler: nil,
// 	}

// 	go func() {
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("listen: %s\n", err)
// 		}
// 	}()
// 	log.Print("Server Started")

// 	<-done
// 	log.Print("Server Stopped")

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer func() {
// 		cancel()
// 	}()

// 	if err := srv.Shutdown(ctx); err != nil {
// 		log.Fatalf("Server Shutdown Failed:%+v", err)
// 	}
// 	log.Print("Server Exited Properly")
// }

// func do(ctx context.Context) {
// 	log.Println("started task")
// 	var wg sync.WaitGroup
// 	sm := make(chan struct{}, 3)

// 	for i := 0; i < 10 && ctx.Err() == nil; i++ {
// 		select {
// 		case <-ctx.Done():
// 			log.Println("context cancelled. Stop working")
// 		case sm <- struct{}{}:
// 			wg.Add(1)
// 			go func() {
// 				defer func() { <-sm }()
// 				defer wg.Done()
// 				heavyOperation()
// 			}()
// 		}
// 	}

// 	wg.Wait()
// 	log.Println("finished task")
// }

// func heavyOperation() {
// 	for i := 0; i < 1e8; i++ {
// 		_ = strconv.Itoa(i)
// 	}
// 	fmt.Println("done")
// }

func main() {
	tcpServer := &tcpserver.Server{
		Addr: ":7171",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go tcpServer.Run(ctx)

	time.Sleep(time.Second)

	for {
		if !tcpServer.IsAlive() {
			break
		}

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
