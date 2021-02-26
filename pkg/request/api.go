package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"inheaden.io/services/simple-mail-api/pkg/config"
	"inheaden.io/services/simple-mail-api/pkg/mail"
	"inheaden.io/services/simple-mail-api/pkg/utils"
)

func sendMail(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	var request MailRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(request)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	allowList := config.GetConfig().AllowList
	if !contains(allowList, request.To) {
		log.Infof("%s is not on the allow list", request.To)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = mail.Sendmail(request.To, request.Subject, request.Body)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getNonce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NonceResponse{Nonce: utils.GetRandomString(10)})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// HandleRequests starts the web server
func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/send", sendMail).Methods("POST")
	router.HandleFunc("/nonce", getNonce).Methods("GET")
	router.HandleFunc("/health", health).Methods("GET")
	router.Use(mux.CORSMethodMiddleware(router))

	log.Infof("Listening on port %s", config.GetConfig().Port)
	log.WithError(http.ListenAndServe(fmt.Sprintf(":%s", config.GetConfig().Port), router)).Fatal("Failed")
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
