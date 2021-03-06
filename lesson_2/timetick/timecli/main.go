package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)

	d := net.Dialer{
		Timeout:   time.Second,
		KeepAlive: time.Minute,
	}
	conn, err := d.DialContext(ctx, "tcp", "localhost:9001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to server")
	log.Println(io.Copy(os.Stdout, conn))
}
