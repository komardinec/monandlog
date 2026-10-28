package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	models "debutwebapp.komardinec/internal/models"
	ui "debutwebapp.komardinec/ui"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass inside our web server
type templateData struct {
	CurrentYear     int
	File            models.File
	Files           []models.File
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// niceDate will return pretty formatted string
// representation of a time.Time object
func NiceDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable.
// This is essentially a string-keyed map which acts as a lookup between the names
// of our custom template functions and the functions themselves.
var functions = template.FuncMap{
	"NiceDate": NiceDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// The fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filesystem which match the pattern "html/pages/*.tmpl.html". This
	// essentially gives us a slice of all thee "page" templates for the application, just
	// like before
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the filename and assign to variable
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we
		// want to parse
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		// Use ParseFS() to parse the templates from the ui.Files embedded filesystem
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page
		// as the key
		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
