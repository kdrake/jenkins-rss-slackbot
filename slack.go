package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// Message to slack
type Message struct {
	Text        string        `json:"text,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

// Attachment pin to Message
type Attachment struct {
	Fallback   string `json:"fallback,omitempty"`
	Color      string `json:"color,omitempty"`
	Title      string `json:"title,omitempty"`
	Text       string `json:"text,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`
	Timestamp  int64  `json:"ts,omitempty"`
}

type Slack struct {
	WebhookURL string
}

// Post pull builds info to Slack
func (s *Slack) Post(item *gofeed.Item) error {
	changeSet, err := getChangeSet(item.Link)
	if err != nil {
		return err
	}

	message, err := getMessage(item, changeSet)
	if err != nil {
		return err
	}

	buf, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}

func getMessage(item *gofeed.Item, changeSet *changeSet) (*Message, error) {
	t, err := time.Parse(time.RFC3339, item.Published)
	if err != nil {
		return nil, err
	}

	var color string
	if strings.Contains(item.Title, "broken") {
		color = "danger"
	} else {
		color = "#dddddd"
	}

	message := &Message{Text: fmt.Sprintf("*%s*\n<%s | Click here> for details.", item.Title, item.Link)}
	var text bytes.Buffer
	for _, i := range changeSet.Items {
		text.WriteString(fmt.Sprintf("%s: %s\n", i.Author.FullName, i.Msg))
	}
	attachment := &Attachment{
		item.Title,
		color,
		"",
		text.String(),
		"",
		"",
		t.Unix(),
	}

	message.Attachments = append(message.Attachments, attachment)

	return message, nil
}

func getChangeSet(url string) (*changeSet, error) {
	r, err := http.Get(fmt.Sprintf("%sapi/json", url))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var details Details
	if err := json.Unmarshal(body, &details); err != nil {
		return nil, err
	}

	return &details.ChangeSet, nil
}

type author struct {
	FullName    string
	AbsoluteURL string
}

type items struct {
	Id            string
	Msg           string
	Comment       string
	Author        author
	Timestamp     int64
	AffectedPaths []string
}

type changeSet struct {
	Items []items
}

type Details struct {
	ChangeSet changeSet
}
