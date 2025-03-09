package loadbalancers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	baseUrl = "http://localhost:800"
)
type LoadBalancer struct {
	RevProxy *httputil.ReverseProxy
}

type Endpoints struct {
	List []*url.URL
}

func (e *Endpoints) Shuffle() {
	temp := e.List[0]
	e.List = e.List[1:]
	e.List = append(e.List, temp)
}

func MakeLoadBalancer(amount int) {
	var lb LoadBalancer
	var ep Endpoints

	router := mux.NewRouter()
	server := http.Server {
		Addr: ":8080",
		Handler: router,
	}

	for i := 0; i < amount; i++ {
		fmt.Println(baseUrl)
		ep.List = append(ep.List, createEndpoint(baseUrl, i))
	}

	router.PathPrefix("/").HandlerFunc(makeRequest(&lb, &ep))
	router.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test page"))
	})

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			log.Println("Registered route:", path)
		}
		return nil
	})

	log.Println("Load balancer running on 8080")
	log.Fatal(server.ListenAndServe())
}

func makeRequest(lb *LoadBalancer, ep *Endpoints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(ep.List) == 0 {
            http.Error(w, "No backend servers available", http.StatusServiceUnavailable)
            return
        }
		backend := ep.List[0]

		for !testServer(backend.String()) {
			ep.Shuffle()
		}
		log.Printf("Forwarding request to backend: %s\n", backend.String())
		lb.RevProxy = httputil.NewSingleHostReverseProxy(backend)
		ep.Shuffle()
		lb.RevProxy.ServeHTTP(w, r)
	}
}

func createEndpoint(endpoint string, idx int) *url.URL {
	link := endpoint + strconv.Itoa(idx)
	parsedURL, err := url.Parse(link)

	if err != nil {
		log.Fatalf("Error passing URL %s: %v", link, err)
	}
	return parsedURL
}

func testServer(endpoint string) bool {
	resp, err := http.Get(endpoint)

	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	return true

}