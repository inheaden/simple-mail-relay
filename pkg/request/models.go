package request

// MailRequest defines the request body to send a mail
type MailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
	Nonce   string `json:"nonce"`
	Hash    string `json:"hash"`
	From    string `json:"from"`
}

// ErrorResponse defines the response body if an error occures
type ErrorResponse struct {
	Error string `json:"error"`
}

// NonceResponse defines the request body when a nonce is returned
type NonceResponse struct {
	Nonce string `json:"nonce"`
}
