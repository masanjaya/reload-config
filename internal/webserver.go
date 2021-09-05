package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// StartHttpServer starts an HTTP Server on port 9090, handle for path / and /reset and /shutdown
func StartHttpServer(port int, wg *sync.WaitGroup, stop, reset, shutdown chan bool) {
	defer wg.Done()

	// handle /
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>click <a href='/reset'>reset</a> to reset Publisher</html>"))
		if err != nil {
			log.Fatal(err)
		}
	})
	// handle /reset
	http.HandleFunc("/reset", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>Publisher has been reset<br/><a href='/'>back to main page</a></html>"))
		if err != nil {
			log.Fatal(err)
		}
		// send reset signal
		log.Println("got http request to reset publisher")
		reset <- true
	})
	// handle /shutdown
	http.HandleFunc("/shutdown", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>app is about to shutdown</html>"))
		if err != nil {
			log.Fatal(err)
		}
		// send shutdown signal
		log.Println("got http request to shutdown the app")
		shutdown <- true
	})

	httpServer := http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	// monitor stop signal
	go func() {
		<-stop
		log.Println("HttpServer is shutting down")
		httpServer.Shutdown(context.Background())
	}()

	// start serving
	log.Println("HttpServer is starting")
	err := httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("HttpServer is shutdown")
}
