package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type ServerList struct {
	Ports []int
}

func (s *ServerList) Populate(amount int) {
	if amount >= 10 {
		log.Fatal("Amount of ports cannot be greater than 10")
	}
	for x := 0; x < amount; x++ {
		s.Ports = append(s.Ports, x)
	}
}

func (s *ServerList) Pop() int {
	port := s.Ports[0]
	s.Ports = s.Ports[1:]

	return port
}

func RunServers(amount int) {
	var myServerList ServerList

	myServerList.Populate(amount)

	var wg sync.WaitGroup
	wg.Add(amount)
	defer wg.Wait()

	for x := 0; x < amount; x++ {
		go makeServers(&myServerList, wg)
	}

}

func makeServers(sl *ServerList, wg sync.WaitGroup) {
	defer wg.Done()
	
	r := mux.NewRouter()
	port := sl.Pop()
	server := http.Server {
		Addr: fmt.Sprintf(":800%d", port),
		Handler: r,
	}
	
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server %d", port)
	})

	r.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("Server shutdown"))
		server.Shutdown(context.Background())
	})
	
	server.ListenAndServe()
}