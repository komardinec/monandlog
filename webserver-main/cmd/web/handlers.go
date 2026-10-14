package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	models "debutwebapp.komardinec/internal/models"
)

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "GO v1.27.x")

	data := app.newTemplateData(r)

	app.Render(w, r, http.StatusOK, "home.tmpl.html", data)
}

// Handler for viewing the files (test function, meant to be depricated later)
func (app *application) ViewFile(w http.ResponseWriter, r *http.Request) {
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

	data := app.newTemplateData(r)
	data.File = file

	//data := templateData{
	//	File: file,
	//}

	// Use the render helper
	app.Render(w, r, http.StatusOK, "view.tmpl.html", data)
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

	//data := templateData{
	//	Files: files,
	//}
	data := app.newTemplateData(r)
	data.Files = files

	// Use the new render helper
	app.Render(w, r, http.StatusOK, "latest.tmpl.html", data)
}
