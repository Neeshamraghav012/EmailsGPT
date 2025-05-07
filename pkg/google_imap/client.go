package googleimap

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Client struct {
	Username string
	Password string
	Host     string
	Port     int
	UseSSL   bool
}

func NewClient(username, password string) *Client {
	return &Client{
		Username: username,
		Password: password,
		Host:     "imap.gmail.com",
		Port:     993,
		UseSSL:   true,
	}
}

func (c *Client) FetchMailsSince(sinceDate time.Time) ([]Mail, error) {
	// Simulate fetching mails from the server
	target := fmt.Sprintf("%s:%d", c.Host, c.Port)
	cl, err := client.DialTLS(target, nil)
	if err != nil {
		return nil, err
	}
	defer cl.Logout()

	if err := cl.Login(c.Username, c.Password); err != nil {
		return nil, err
	}

	_, err = cl.Select("[Gmail]/All Mail", true)
	if err != nil {
		return nil, err
	}

	// Search emails since the given date
	criteria := imap.NewSearchCriteria()
	date := time.Now().Add(-24 * time.Hour)
	criteria.Since = date.Truncate(24 * time.Hour)
	//criteria.Before = time.Now()

	ids, err := cl.Search(criteria)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []Mail{}, nil // No emails found
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	messages := make(chan *imap.Message, len(ids))
	const GmailMsgID = imap.FetchItem("X-GM-MSGID")
	err = cl.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, GmailMsgID, imap.FetchRFC822Header}, messages)
	if err != nil {
		return nil, err
	}

	var emails Mails
	for msg := range messages {
		if msg == nil {
			continue
		}

		if msg.Envelope == nil {
			continue
		}

		if msg.Envelope.Date.Before(sinceDate) {
			continue
		}

		from := ""
		if len(msg.Envelope.From) > 0 {
			from = msg.Envelope.From[0].Address()
		}

		if from == c.Username {
			continue // Skip emails sent by the user
		}

		subject := msg.Envelope.Subject
		formattedMessageID := msg.Envelope.MessageId[1 : len(msg.Envelope.MessageId)-1]
		link := fmt.Sprintf("https://mail.google.com/mail/?authuser=%s#search/rfc822msgid:%s", c.Username, formattedMessageID)

		emails = append(emails, Mail{
			From:    from,
			Subject: subject,
			Link:    link,
		})
	}

	return emails, nil

}
