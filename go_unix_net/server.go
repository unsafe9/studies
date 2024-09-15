package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

func main() {
	SetFDLimit(61500)

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

		fd := GetTCPSocketFD(conn)
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

func SetFDLimit(max int) uint64 {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		log.Fatalf("getrlimit: %v", err)
	}
	log.Printf("FDSoftLimit:%d, FDHardLimit:%d", limit.Cur, limit.Max)

	if max > 0 {
		if uint64(max) < limit.Max {
			limit.Cur = uint64(max)
		} else {
			limit.Cur = limit.Max
		}
	} else {
		// Set to the maximum value if max is 0
		limit.Cur = limit.Max
	}
	if limit.Cur < limit.Max {
		if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
			log.Fatalf("setrlimit: %v", err)
		}
		log.Printf("Set FDSoftLimit to %d", limit.Cur)
	}
	return limit.Cur
}

func GetTCPSocketFD(conn net.Conn) int {
	//tls := reflect.TypeOf(conn.UnderlyingConn()) == reflect.TypeOf(&tls.Conn{})
	// Extract the file descriptor associated with the connection
	//connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	//if tls {
	//	tcpConn = reflect.Indirect(tcpConn.Elem())
	//}
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
