package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/go-co-op/gocron"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"inheaden.io/services/simple-mail-api/pkg/config"
	"inheaden.io/services/simple-mail-api/pkg/db"
	"inheaden.io/services/simple-mail-api/pkg/mail"
)

var database = db.NewDB()
var generalConfig config.Config

func sendMail(w http.ResponseWriter, r *http.Request) {
	// check content type
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// decode request
	var request MailRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate request
	err = validator.New().Struct(request)
	if err != nil {
		writeError(w, err)
		return
	}

	// ratelimit request
	ok, err := RatelimitIP(r)
	if err != nil {
		writeError(w, err)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	// verify hash
	if generalConfig.RequireNonce {
		ok, err = VerifyNonce(r, request.Nonce, request.Hash)
		if err != nil {
			writeError(w, err)
			return
		}

		if !ok {
			log.Infof("Wrong hash")
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	// verify email is on allow list
	allowList := generalConfig.AllowList
	if !contains(allowList, request.To) {
		log.Infof("%s is not on the allow list", request.To)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// send mail
	err = mail.Sendmail(request.To, request.Subject, request.Body)
	if err != nil {
		writeError(w, err)
		return
	}

	// done
	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, err error) {
	log.Error(err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

func getNonce(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	nonce, err := GetNonce(r)
	if err != nil {
		writeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(NonceResponse{Nonce: nonce})
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// HandleRequests starts the web server
func HandleRequests() {
	generalConfig = config.GetConfig()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/send", sendMail).Methods("POST")
	router.HandleFunc("/nonce", getNonce).Methods("GET")
	router.HandleFunc("/health", health).Methods("GET")
	router.Use(mux.CORSMethodMiddleware(router))

	s := gocron.NewScheduler(time.UTC)
	s.Every(60).Seconds().Do(func() {
		log.Debug("Running cleanups")
		err1 := database.CleanIPs(generalConfig.MaxAgeSeconds)
		err2 := database.CleanNonces(generalConfig.MaxAgeSeconds)
		if err1 != nil || err2 != nil {
			log.Error(err1.Error())
			log.Error(err2.Error())
		}
	})
	s.StartAsync()

	log.Infof("Listening on port %s", generalConfig.Port)
	log.WithError(http.ListenAndServe(fmt.Sprintf(":%s", generalConfig.Port), router)).Fatal("Failed")
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
