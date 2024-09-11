package mock

type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

func (e *EmailSender) SendEmail(subject string, body string, emailTo string) error {
	subject = ""
	body = ""
	emailTo = ""

	return nil
}
