package main

import (
	"github.com/unsafe9/studies/go_unix_net/netutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	netutil.SetFDLimit(61500)

	l, err := net.Listen("tcp", ":3030")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer l.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-sigCh
		log.Printf("got signal, attempting graceful shutdown")
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("accept: %v", err)
		}

		log.Printf("Accepted connection from: %s", conn.RemoteAddr())

		fd := netutil.GetTCPSocketFD(conn)
		log.Printf("File Descriptor: %d", fd)

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("read: %v", err)
		}

		log.Printf("Received: %s", buf[:n])
		conn.Close()
	}
}
