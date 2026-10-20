package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	models "debutwebapp.komardinec/internal/models"
	validator "debutwebapp.komardinec/internal/validator"
)

type createFileForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
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

	// Initialize a new createFileForm, this will generate fields with default nil
	// values for type fields, or we can set our own default values
	data.Form = createFileForm{
		Expires: 365,
	}

	app.Render(w, r, http.StatusOK, "create.tmpl.html", data)
}

// Handler to work with POST method while creating a new file (test function, meant to be depricated later)
func (app *application) CreateFilePost(w http.ResponseWriter, r *http.Request) {
	// Declare a new empty instance if the createFileForm struct
	var form createFileForm

	// Parse form content and populate PostForm() map from POST request body
	// Also works for PUT and PATCH
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Call the Decode() method of the form decoder, passing in the current
	// request and a POINTER to our createFileForm struct.
	// This will essentially fill the struct with the relevant values from the HTML form.
	// If there is a problem, we return 400 Bad Request response to the client
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Now validator struct are included inside form variable, and we can use
	// form.CheckField() to check and add new errors to error map
	// validator functions with functionality for checking the specific
	// type of error. And also, because it is all happens inside one structure,
	// any values or functions inside it, are shared between each other.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be equal 1, 7, or 365")

	// If form validation function will fail, server will render the error for user
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.Render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	// Data source for db model will be from form structure entity
	id, err := app.dbWebApp.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	// Use the Put() method to add a string value and the corresponding value
	// to the session data. Used to say that file has been created.
	app.sessionManager.Put(r.Context(), "flash", "File sucessfully created!")

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

// Authentification and authorization section

func (app *application) UserSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "signup form for users")
}

func (app *application) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST method for user creation")
}

func (app *application) UserLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "login form for user")
}

func (app *application) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "authenticate and login user, verify against our information")
}

func (app *application) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "logout user, close session")
}
