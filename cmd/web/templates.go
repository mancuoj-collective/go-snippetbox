package main

import "github.com/mancuoj/snippetbox/internal/models"

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
