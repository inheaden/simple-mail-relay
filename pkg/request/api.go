package request

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"simple-mail-api/pkg/mail"
)

func send(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	receiverEmail := vars["receiverEmailAddress"]
	emailSubject := vars["emailSubject"]
	emailBody := vars["emailBody"]

	mail.Sendmails(receiverEmail, emailSubject, emailBody)
	mail.SendmailNotification(receiverEmail, emailSubject, emailBody)
}
func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/sendmail/{receiverEmailAddress}/{emailSubject}/{emailBody}", send).Methods("POST")
	log.Fatal(http.ListenAndServe(":8082", myRouter))
}
