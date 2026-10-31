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

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) HomePage(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)

	app.Render(w, r, http.StatusOK, "home.tmpl.html", data)
}

// Handler for viewing the files (test function, meant to be depricated later)
func (app *application) ViewFile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if id < 1 || err != nil {
		http.NotFound(w, r)
		return
	}

	// Use the dbWebApp's Get() method to retreive the data for a
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

	data := app.NewTemplateData(r)
	data.File = file

	// Use the render helper
	app.Render(w, r, http.StatusOK, "view.tmpl.html", data)
}

// Handler for creating the files (test function, meant to be depricated later)
func (app *application) CreateFile(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)

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
	err := app.DecodePostForm(r, &form)
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
		data := app.NewTemplateData(r)
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

	data := app.NewTemplateData(r)
	data.Files = files

	// Use the new render helper
	app.Render(w, r, http.StatusOK, "latest.tmpl.html", data)
}

// Authentification and authorization section

// UserSignup() handler displays
func (app *application) UserSignup(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = userSignupForm{}
	app.Render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	// Declare an zero-valued instance of our userSignupForm stryct
	var form userSignupForm

	// Parse the form data into the userSignupForadm
	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper function
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters")

	// If there are any errors, redisplay the signup form along with a 422
	// status code
	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.ServerError(w, r, err)
		}

		return
	}

	// Otherwise add a confirmation flash message to the session confirming that
	// their signup worked
	app.sessionManager.Put(r.Context(), "flash", "Signup was successful! Now you can log in.")

	// And redirect the user to the login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) UserLogin(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = userLoginForm{}
	app.Render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	// Decode form data inside login form
	var form userLoginForm

	err := app.DecodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Do some validation checks on the form. We check that both email and password
	// are provided, and also check the format of the email address as
	// a UX-nicety (in case the user makes a typo)
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must contain a email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password are incorrect")

			data := app.NewTemplateData(r)
			data.Form = form
			app.Render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.ServerError(w, r, err)
		}
		return
	}

	// Use the RenewToken() method on the current session to change the session ID.
	// It is a good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g login and logout operations)
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	// Add the id of the current user to the session,
	// now we are counted as logged in
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	app.sessionManager.Put(r.Context(), "flash", "Succesfully logged in!")

	// Redirect to Home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Renew session ID, because we changed the session status
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.ServerError(w, r, err)
		return
	}

	// Remove the authenticatedUserID from the session data so the user
	// is logged out
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// Add a flash message, about out logout
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Ping handler, that are used to check if our server are alive
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}
