package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reload-config/model"
	"sync"
	"time"
)

// StartHttpServer starts an HTTP Server on port 9090, handle for path / and /reset and /shutdown
func StartHttpServer(port int, wg *sync.WaitGroup, stop, reset, shutdown chan bool) {
	defer wg.Done()

	// handle /get secret
	http.HandleFunc("/v1/secret/data/aestest1", func(rw http.ResponseWriter, r *http.Request) {
		var response3 secretResponse
		response3.Data.Metadata.DeletionTime = "never"
		response3.Data.Metadata.Destroyed = false
		response3.Data.Metadata.Version = 1
		response3.Data.Metadata.CreatedTime = time.Now()
		response3.Data.Data = map[string]string{"ctrkey": "yshWyt25sf30Uhdj", "bitkey": "128"}

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	http.HandleFunc("/v1/secret/metadata/pilot/certs", func(rw http.ResponseWriter, r *http.Request) {
		var response3 listResponse
		response3.RequestId = "request_id_1"
		response3.LeaseId = "lease_id_1"
		response3.Renewable = true
		response3.LeaseDuration = 10
		response3.Data.Keys = append(response3.Data.Keys, "empty")

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	http.HandleFunc("/v1/secret/metadata/pilot", func(rw http.ResponseWriter, r *http.Request) {
		var response3 listResponse
		response3.RequestId = "request_id_1"
		response3.LeaseId = "lease_id_1"
		response3.Renewable = true
		response3.LeaseDuration = 10
		response3.Data.Keys = append(response3.Data.Keys, "certs/")
		response3.Data.Keys = append(response3.Data.Keys, "elstreamer")
		response3.Data.Keys = append(response3.Data.Keys, "fc")
		response3.Data.Keys = append(response3.Data.Keys, "tdextract2")

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle list
	http.HandleFunc("/v1/secret/metadata/", func(rw http.ResponseWriter, r *http.Request) {
		var response3 listResponse
		response3.RequestId = "request_id_1"
		response3.LeaseId = "lease_id_1"
		response3.Renewable = true
		response3.LeaseDuration = 10
		response3.Data.Keys = append(response3.Data.Keys, "pilot/")

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle /get secret
	http.HandleFunc("/v1/secret/data/pilot/elstreamer", func(rw http.ResponseWriter, r *http.Request) {
		var response3 secretResponse
		response3.Data.Metadata.DeletionTime = "never"
		response3.Data.Metadata.Destroyed = false
		response3.Data.Metadata.Version = 1
		response3.Data.Metadata.CreatedTime = time.Now()
		response3.Data.Data = map[string]string{"elstreamer_data_1": "dummy_data_elstreamer1", "elstreamer_data_2": "dummy_data_elstreamer2"}

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle /get secret
	http.HandleFunc("/v1/secret/data/pilot/fc", func(rw http.ResponseWriter, r *http.Request) {
		var response3 secretResponse
		response3.Data.Metadata.DeletionTime = "never"
		response3.Data.Metadata.Destroyed = false
		response3.Data.Metadata.Version = 1
		response3.Data.Metadata.CreatedTime = time.Now()
		response3.Data.Data = map[string]string{"fc1": "dummy_data_fc1", "fc2": "dummy_data_fc2"}

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle /get secret
	http.HandleFunc("/v1/secret/data/pilot/tdextract2", func(rw http.ResponseWriter, r *http.Request) {
		var response3 secretResponse
		response3.Data.Metadata.DeletionTime = "never"
		response3.Data.Metadata.Destroyed = false
		response3.Data.Metadata.Version = 1
		response3.Data.Metadata.CreatedTime = time.Now()
		response3.Data.Data = map[string]string{"mongo_uri": " mongodb://XXXXXXXXXX@127.0.0.1:27017/FunnelConnectPILOT", "postgresql_uri": "postgres://X@postgres.pilot.teavaro.systems/pilot?sslmode=disable"}

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle /get secret
	http.HandleFunc("/v1/secret/data/testkey", func(rw http.ResponseWriter, r *http.Request) {
		var response3 secretResponse
		response3.Data.Metadata.DeletionTime = "never"
		response3.Data.Metadata.Destroyed = false
		response3.Data.Metadata.Version = 1
		response3.Data.Metadata.CreatedTime = time.Now()
		response3.Data.Data = map[string]string{"publickey": model.TestPublicKey, "privatekey": model.TestPrivateKey}

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle /get secret
	http.HandleFunc("/v1/secret/data/teavaro", func(rw http.ResponseWriter, r *http.Request) {
		var response3 secretResponse
		response3.Data.Metadata.DeletionTime = "never"
		response3.Data.Metadata.Destroyed = false
		response3.Data.Metadata.Version = 1
		response3.Data.Metadata.CreatedTime = time.Now()
		response3.Data.Data = map[string]string{"publickey": model.SpecificPublicKey1}

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle /get secret
	http.HandleFunc("/v1/secret/data/string1", func(rw http.ResponseWriter, r *http.Request) {
		var response3 secretResponse
		response3.Data.Metadata.DeletionTime = "never"
		response3.Data.Metadata.Destroyed = false
		response3.Data.Metadata.Version = 1
		response3.Data.Metadata.CreatedTime = time.Now()
		response3.Data.Data = map[string]string{"publickey": model.PublicKey1, "privatekey": model.PrivateKey1}

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle /get secret
	http.HandleFunc("/v1/secret/data/string3", func(rw http.ResponseWriter, r *http.Request) {
		var response3 secretResponse
		response3.Data.Metadata.DeletionTime = "never"
		response3.Data.Metadata.Destroyed = false
		response3.Data.Metadata.Version = 1
		response3.Data.Metadata.CreatedTime = time.Now()
		response3.Data.Data = map[string]string{"publickey": model.PublicKey3, "privatekey": model.PrivateKey3}

		returnbyte, err := json.Marshal(response3)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

	// handle /login
	http.HandleFunc("/v1/auth/userpass/login/sanjaya@teavaro.com", func(rw http.ResponseWriter, r *http.Request) {

		var requestBody authRequest
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			log.Println("can not decode body requset:", err.Error())
			_, err := rw.Write([]byte("<html>empty body request</html>"))
			if err != nil {
				log.Fatal(err)
			}
			return
		} else {
			log.Println("got http request with", requestBody.Password)
		}
		log.Println("checking http request done")

		var response2 authResponse
		response2.LeaseID = "LeaseID"
		response2.Renewable = true
		response2.LeaseDuration = 10
		response2.Auth.ClientToken = "clienttoken1"
		response2.Auth.Policies = []string{"string1", "string2"}
		response2.Auth.Metadata.Username = "username1"
		response2.Auth.LeaseDuration = 2
		response2.Auth.Renewable = true

		returnbyte, err := json.Marshal(response2)
		if err != nil {
			log.Fatal(err)
			returnbyte = []byte(`{"status":"error"}`)
		}

		_, err = rw.Write(returnbyte)
		if err != nil {
			log.Fatal(err)
		}

	})

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
	log.Println("HttpServer is starting now")
	err := httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	log.Println("HttpServer is shutdown")
}

type authRequest struct {
	Password string `json:"password"`
}

type authResponse struct {
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Auth          struct {
		ClientToken string   `json:"client_token"`
		Policies    []string `json:"policies"`
		Metadata    struct {
			Username string `json:"username"`
		} `json:"metadata"`
		LeaseDuration int  `json:"lease_duration"`
		Renewable     bool `json:"renewable"`
	} `json:"auth"`
}

type secretResponse struct {
	Data struct {
		Metadata struct {
			Version      int       `json:"version"`
			Destroyed    bool      `json:"destroyed"`
			DeletionTime string    `json:"deletion_time"`
			CreatedTime  time.Time `json:"created_time"`
		} `json:"metadata"`
		Data map[string]string `json:"data"`
	} `json:"data"`
}

type listResponse struct {
	RequestId     string `json:"request_id"`
	LeaseId       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Keys []string `json:"keys"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}
