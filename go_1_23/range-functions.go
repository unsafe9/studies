package main

import (
	"iter"
	"log"
	"math"
	"math/rand/v2"
)

type Server struct {
	ID      uint64
	Healthy bool
}

var (
	servers        []*Server
	healthyServers int
)

func init() {
	serverCount := rand.IntN(100)
	servers = make([]*Server, serverCount)

	for i := 0; i < serverCount; i++ {
		servers[i] = &Server{
			ID:      rand.N[uint64](math.MaxInt64), // same with rand.Uint64()
			Healthy: rand.IntN(2) == 0,
		}
		if servers[i].Healthy {
			healthyServers++
		}
	}
}

func HealthyServerIDs() iter.Seq[uint64] {
	return func(yield func(uint64) bool) {
		for _, server := range servers {
			if server.Healthy {
				yield(server.ID)
			}
		}
	}
}

func main() {
	n := 1
	for id := range HealthyServerIDs() {
		log.Printf("%2d: %d\n", n, id)
		n++
	}
	log.Printf("servers: %d, healthy: %d\n", len(servers), healthyServers)
}
