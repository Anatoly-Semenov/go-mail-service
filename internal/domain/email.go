package domain

type EmailType string

const (
	RegistrationEmail    EmailType = "registration"
	PaymentReminderEmail EmailType = "payment_reminder"
	PaymentSuccessEmail  EmailType = "payment_success"
	PasswordChangeEmail  EmailType = "password_change"
)

type EmailMessage struct {
	UserID    string    `json:"user_id"`
	EmailType EmailType `json:"email_type"`
}

type EmailTemplate struct {
	Subject string
	Body    string
}
