package models

type EmailResponseEntry struct {
	ID      int    `json:"id"`
	Subject string `json:"subject"`
	From    string `json:"from"`
	Summary string `json:"summary"`
	Link    string `json:"link"`
}

type SummarizerResponse struct {
	Emails []EmailResponseEntry `json:"critical_emails"`
}
