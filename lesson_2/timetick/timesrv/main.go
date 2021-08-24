package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}
	l, err := cfg.Listen(ctx, "tcp", "localhost:9001")
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	log.Println("im started!")

	go func() {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		} else {
			wg.Add(1)
			go handleConn(ctx, conn, wg)
			go handleText(ctx, conn, wg)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("done")
			l.Close()
			wg.Wait()
			log.Println("exit")
			return
		}
	}
}

func handleConn(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	// каждую 1 секунду отправлять клиентам текущее время сервера
	tck := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-tck.C:
			fmt.Fprintf(conn, "now: %s\n", t)
		}
	}
}

// stdin сервера перенаправляется на соеденение через io.Copy (в отдельной горутине)
func handleText(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	io.Copy(conn, os.Stdin)
}
