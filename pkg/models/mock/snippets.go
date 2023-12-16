package mock

import (
	"time"

	"tarala/snippetbox/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (m SnippetModel) Get(id int) (*models.Snippet, error) {
	if id == 1 {
		return mockSnippet, nil
	} 
	return nil, models.ErrNoRecord
	
}

func (m SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}