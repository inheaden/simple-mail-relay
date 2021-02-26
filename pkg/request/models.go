package request

type MailRequest struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type NonceResponse struct {
	Nonce string `json:"nonce"`
}
