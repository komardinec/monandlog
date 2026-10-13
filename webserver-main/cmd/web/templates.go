package main

import (
	models "debutwebapp.komardinec/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	File  models.File
	Files []models.File
}
