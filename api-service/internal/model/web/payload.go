package web

type EmailPayload struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

type WebhookPayload struct {
	URL     string            `json:"url" validate:"required,url"`
	Method  string            `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
}
