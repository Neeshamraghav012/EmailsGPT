package googleimap

type Mail struct {
	Subject  string
	From     string
	BodyText string
	Link     string
}
type Mails []Mail
