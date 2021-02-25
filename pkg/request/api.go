package request

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
	"inheaden.io/services/simple-mail-api/pkg/config"
	"inheaden.io/services/simple-mail-api/pkg/mail"
)

type MailRequest struct {
	To      string
	Subject string
	Body    string
}

func send(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var request MailRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	allowList := config.GetConfig().AllowList
	if !contains(allowList, request.To) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = mail.Sendmail(request.To, request.Subject, request.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func health(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("nc", "-v", "78.46.5.205", "80")
	stdout, err := cmd.CombinedOutput()

	cmd = exec.Command("nc", "-v", "78.46.5.205", "465")
	stdout, err = cmd.CombinedOutput()

	log.Printf(string(stdout))
	log.Print(err)

	w.WriteHeader(http.StatusOK)
}

// HandleRequests starts the web server
func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/sendmail", send).Methods("POST")
	router.HandleFunc("/health", health).Methods("GET")

	log.Printf("Listening on port %s", config.GetConfig().Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.GetConfig().Port), router))
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
