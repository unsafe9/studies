package main

import (
	"errors"
	"github.com/unsafe9/studies/go_unix_net/netutil"
	"golang.org/x/sys/unix"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
)

const (
	maxEvents  = 128
	maxTasks   = 10000
	maxWorkers = 10
)

type epoller struct {
	fd         int
	eventQueue chan *epollConn
	pool       sync.Map
}

type epollConn struct {
	fd     int
	conn   net.Conn
	closed atomic.Bool
}

func main() {
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		log.Fatalf("epoll_create1: %v", err)
	}

	defer unix.Close(epfd)

	ep := &epoller{
		fd:         epfd,
		eventQueue: make(chan *epollConn, maxTasks),
	}

	for i := 0; i < maxWorkers; i++ {
		go startWorker(ep)
	}

	go startListenServer(ep)

	startReadEvents(ep)
}

func startListenServer(ep *epoller) {
	ln, err := net.Listen("tcp", ":3030")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				log.Printf("accept timeout: %v", err)
				continue
			}
			log.Printf("accept error: %v", err)
			return
		}

		fd := netutil.GetTCPSocketFD(conn)
		err = unix.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{
			Events: unix.POLLIN | unix.POLLHUP,
			Fd:     int32(fd),
		})
		if err != nil {
			log.Printf("epollctl add failed: %v", err)
			conn.Close()
			continue
		}

		ep.pool.Store(fd, &epollConn{
			fd:   fd,
			conn: conn,
		})
	}
}

func startWorker(ep *epoller) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("worker panic: %v", err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		c := <-ep.eventQueue
		if c == nil {
			continue
		}

		log.Printf("worker: %d, %s", c.fd, c.conn.RemoteAddr())

		n, err := c.conn.Read(buf)
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				log.Printf("read timeout: %v", err)
				continue
			}

			if !c.closed.Load() {
				c.closed.Store(true)
				for {
					err := unix.EpollCtl(ep.fd, syscall.EPOLL_CTL_DEL, c.fd, nil)
					if err != nil {
						if errors.Is(err, unix.EINTR) {
							continue
						}
						log.Printf("epollctl del failed: %v", err)
					}
					break
				}
				ep.pool.Delete(c.fd)
			}

			if errors.Is(err, io.EOF) {
				log.Print("eof")
			} else {
				log.Printf("read error: %v", err)
			}
			continue
		}

		log.Printf("Received: %s", buf[:n])
	}
}

func epolleventsread(events uint32) bool {
	return events&(unix.EPOLLIN|unix.EPOLLRDHUP|unix.EPOLLHUP|unix.EPOLLERR) != 0
}
func epolleventswrite(events uint32) bool {
	return events&(unix.EPOLLOUT|unix.EPOLLHUP|unix.EPOLLERR) != 0
}

func startReadEvents(ep *epoller) {
	var (
		n      int
		err    error
		events = make([]unix.EpollEvent, maxEvents)
	)
	for {
		for {
			n, err = unix.EpollWait(ep.fd, events, maxEvents)
			if err != nil {
				if errors.Is(err, unix.EINTR) {
					continue
				}
				log.Panicf("startReadEvents: %v", err)
			}
			break
		}
		if n == 0 {
			continue
		}

		// FIXME : The sessions that received EOF are polled here until the worker closes them.

		for i := 0; i < n; i++ {
			event, ok := ep.pool.Load(int(events[i].Fd))
			if ok && epolleventsread(events[i].Events) {
				conn := event.(*epollConn)
				if !conn.closed.Load() {
					ep.eventQueue <- event.(*epollConn)
				}
			}
		}
	}
}
