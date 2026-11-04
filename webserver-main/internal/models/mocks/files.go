package mocks

import (
	"time"

	"debutwebapp.komardinec/internal/models"
)

// That mockFile instance should have simple string without
// any unecessary UTF-8 special symbols, otherwise
// http/template will cook it
var mockFile = models.File{
	ID:      1,
	Title:   "This is a new reality",
	Content: "Do you hear me?",
	Created: time.Now(),
	Expires: time.Now(),
}

type WebAppModel struct{}

func (m *WebAppModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *WebAppModel) Get(id int) (models.File, error) {
	switch id {
	case 1:
		return mockFile, nil
	default:
		return models.File{}, models.ErrNoRecord
	}
}

func (m *WebAppModel) Latest() ([]models.File, error) {
	return []models.File{mockFile}, nil
}
