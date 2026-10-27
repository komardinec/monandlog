package main

import (
	"html/template"
	"path/filepath"
	"time"

	models "debutwebapp.komardinec/internal/models"
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

	// filepath.Glob() function to get a slice of all filepath that
	// match pattern "./ui/html/pages/*.tmpl.html".
	// This is essentially gives us a slice of all the filepath for
	// pages templates.
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one
	for _, page := range pages {
		// Extract the filename and assign to variable
		name := filepath.Base(page)

		// Parse the base template file into a template set
		// The template.FuncMap must be registered with the template set before we call
		// the ParseFile() method
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Parse the files into a template set.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() on the tempalate set to add the page template
		ts, err = ts.ParseFiles(page)
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
