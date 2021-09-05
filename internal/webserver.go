package internal

import (
	"context"
	"log"
	"net/http"
	"sync"
)

// startHttpServer starts an HTTP Server on port 9090, handle for path / and /reset
func StartHttpServer(wg *sync.WaitGroup, stop, reset, shutdown chan bool) {
	defer wg.Done()

	// handle path /
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>click <a href='/reset'>reset</a> to reset Publisher</html>"))
		if err != nil {
			log.Fatal(err)
		}
	})
	// handle path /reset
	http.HandleFunc("/reset", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("<html>Publisher has been reset<br/><a href='/'>back to main page</a></html>"))
		if err != nil {
			log.Fatal(err)
		}
		// send reset signal
		log.Println("got http request to reset publisher")
		reset <- true
	})
	// handle path /stop
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
		Addr: ":9090",
	}

	// monitor stop signal
	go func() {
		<-stop
		log.Println("HttpServer is shutting down")
		httpServer.Shutdown(context.Background())
	}()

	log.Println("HttpServer is starting")
	err := httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("HttpServer is shutdown")
}
