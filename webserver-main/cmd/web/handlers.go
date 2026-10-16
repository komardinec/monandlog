package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	models "debutwebapp.komardinec/internal/models"
)

type createFileForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
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

	// Use the render helper
	app.Render(w, r, http.StatusOK, "view.tmpl.html", data)
}

// Handler for creating the files (test function, meant to be depricated later)
func (app *application) CreateFile(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = createFileForm{
		Expires: 365,
	}

	app.Render(w, r, http.StatusOK, "create.tmpl.html", data)
}

// Handler to work with POST method while creating a new file (test function, meant to be depricated later)
func (app *application) CreateFilePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I AM INSIDE CREATEFILEPOST. DO YOU SEE ME??????????????")
	// Parse form content and populate PostForm() map from POST request body
	// Also works for PUT and PATCH
	err := r.ParseForm()
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Because r.PostForm.Get() always return string type value,
	// we have to convert it to the integer type, so day count can be
	// represented correctly
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	//app.logger.Info("debug value", "title", title)
	fmt.Println("title: ", title)
	// Get part of form values straight from PostForm, others will be
	// added here after some preparation
	form := createFileForm{
		//Title:       r.PostForm.Get("title"),
		Title:       title,
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	// Check if title field is not empty or bigger than 100 characters
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 { // here we will be checking by counting UNICODE symbols,
		// to only count real character
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check if content field is not empty
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// Check if expires are matching are required day period
	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// If error map are have some errors, dump them as plain HTTP response
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.Render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.dbWebApp.Insert(form.Title, form.Content, form.Expires)
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

	data := app.newTemplateData(r)
	data.Files = files

	// Use the new render helper
	app.Render(w, r, http.StatusOK, "latest.tmpl.html", data)
}
