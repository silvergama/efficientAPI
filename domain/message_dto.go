package domain

import (
	"strings"
	"time"

	"github.com/silvergama/efficientAPI/utils/errorutils"
)

type Message struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *Message) Validate() errorutils.MessageErr {
	m.Title = strings.TrimSpace(m.Title)
	m.Body = strings.TrimSpace(m.Body)
	if m.Title == "" {
		return errorutils.NewUnprocessibleEntityError("Please enter a valid title")
	}
	if m.Body == "" {
		return errorutils.NewUnprocessibleEntityError("Please enter a valid body")
	}
	return nil
}
