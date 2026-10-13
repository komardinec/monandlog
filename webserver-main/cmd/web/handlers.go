package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"debutwebapp.komardinec/internal/models"
)

// It is not recommended to use global variables at general,
// but if we want to use templates for all handlers
// global []string are very good implementation
//var templateFiles = []string{
//	"./ui/html/base.tmpl.html",
//	"./ui/html/pages/home.tmpl.html",
//	"./ui/html/pages/view.tmpl.html",
//	"./ui/html/partials/nav.tmpl.html",
//}

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "GO v1.27.x")

	var templateFiles = []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	// template.ParseFiles will read the template file and turn its contents
	// into a template set. If there is an error occured we will get into terminal the error message
	// and user will get nice little plain text "Internal Server Error" response
	// and that case handler will return and no subsequent code will be executed.
	ts, err := template.ParseFiles(templateFiles...) // we unpacking string slice of templateFiles paths
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	// ExecuteTemplate() method will write contents of the "base"
	// template as the response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.ServerError(w, r, err)
	}
}

// Handler for viewing the files (test function, meant to be depricated later)
func (app *application) ViewFile(w http.ResponseWriter, r *http.Request) {
	var templateFiles = []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if id < 1 || err != nil {
		http.NotFound(w, r)
		return
	}

	// Use the SnippetModel's Get() method to retreive the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response
	file, err := app.dbWebApp.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.ServerError(w, r, err)
		}
		return
	}

	ts, err := template.ParseFiles(templateFiles...)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	data := templateData{
		File: file,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.ServerError(w, r, err)
	}

	// Plain-text response body
	//fmt.Fprintf(w, "%+v", file)
}

// Handler for creating the files (test function, meant to be depricated later)
func (app *application) CreateFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("File creation form"))
}

// Handler to work with POST method while creating a new file (test function, meant to be depricated later)
func (app *application) CreateFilePost(w http.ResponseWriter, r *http.Request) {
	//w.Header().Add("Server", "Test Go Backend Server")
	//w.Header().Add("Author", "komardinec")

	title := "POST test DB insert"
	content := "I am thou, thou art I"
	expires := 365

	id, err := app.dbWebApp.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/file/view/%d", id), http.StatusSeeOther)
}

func (app *application) Latest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Handler-Type", "Latest")

	files, err := app.dbWebApp.Latest()
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	//for _, file := range files {
	//	fmt.Fprintf(w, "%+v\n", file)
	//}

	data := templateData{
		Files: files,
	}

	var templateFiles = []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/latest.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	ts, err := template.ParseFiles(templateFiles...)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.ServerError(w, r, err)
	}

}
