package main

import (
	"fmt"
	"html/template" // library to handle HTML in a separate file
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct { // web wiki consist of interconnected web pages. This structure describes how page data stored in memory
	Title string // every web page have a title. can be used for naming
	Body  []byte // and page content (body). Body have type []byte instead of string, because io libraries expecting information in byte format
}

func getTemplateFileName(filename string) string {
	templateFolder := "tmpl/"
	filename = templateFolder + filename + ".html"
	return filename
}

func getDataFileName(filename string) string {
	dataFolder := "data/"
	filename = dataFolder + filename + ".txt"
	return filename
}

// Every time we use template.ParseFiles program parsing files again and again. To prevent excessive disk reading, we can implement caching.
// var templates will call ParseFiles once at programm initialization and will parse them into a single *Template
var templates = template.Must(template.ParseFiles(getTemplateFileName("edit"), getTemplateFileName("view"))) // template.Must is a wrapper that panics once it receive non-nil error code
// otherwise returns the *Template unaltered.
// A panic is appropriate here; if the template can't be loadd the only sensible thing to do is exit the program

var validPath = regexp.MustCompile("/(edit|save|view)/([a-zA-Z0-9]+)$") // regular expression path, used to secure url paths
// regexp.MustCompile will compile every time user enter a path and panic when path of user not applicable to regular expresion string

// For persistent storage of web pages, we can create a save method for Page.
func (p *Page) save() error { // This method save() take a reference to p Page as a receiver, returning a value of error.
	// error return required by os.WriteFile, to let the application handle if anythind happen to the file.
	// if everything is ok Page.save() will return nil (0)
	//filename := p.Title + ".txt"
	filename := getDataFileName(p.Title)
	return os.WriteFile(filename, p.Body, 0600)
}

//func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
//	m := validPath.FindStringSubmatch(r.URL.Path) // FindStringSubmatch will return a slice of substrings from r.URL.Path if they follow expression from ValidPath
//	if m == nil {                                 // if slice of substrings are empty (nil)
//		http.NotFound(w, r)                         // return 404
//		return "", errors.New("invalid Page Title") // return some error and empty title
//	}
//	return m[2], nil // The title is the second subexpression. If title is correct it will be returned along with nil error
//}

// Web server obviously should handle reading of web pages
// function loadPage accept title string, and returning link of Page structure (title, body) by reading body from file
func loadPage(title string) (*Page, error) {
	// by title name. File should be in the same directory as web server.
	//filename := title + ".txt"         // conventional name of file
	filename := getDataFileName(title)
	body, err := os.ReadFile(filename) // os.ReadFile return []byte body of web page file and error code of reading process
	if err != nil {                    // returning empty page with error code if error code not nil
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil // function creating reference to Page structure and returning it with nil error code if page file read was succefull
}

//func viewHandler(w http.ResponseWriter, r *http.Request) {
//title := r.URL.Path[len("/view/"):] // extract title from URL to find web page file, by ignoring "/view/" part of link path
//	title, err := getTitle(w, r) // get title by function
//	if err != nil {
//		return
//	}
//	p, err := loadPage(title) // load file by extracted title
//fmt.Fprintf(w, "<h1>Viewing %s</h1><div>%s</div>", p.Title, p.Body) // hardcoded render of file content
//t, _ := template.ParseFiles("view.html") // view will also be using template files
//t.Execute(w, p)
//	if err != nil {

//func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
//	m := validPath.FindStringSubmatch(r.URL.Path) // FindStringSubmatch will return a slice of substrings from r.URL.Path if they follow expression from ValidPath
//	if m == nil {                                 // if slice of substrings are empty (nil)
//		http.NotFound(w, r)                         // return 404
//		return "", errors.New("invalid Page Title") // return some error and empty title
//	}
//	return m[2], nil // The title is the second subexpression. If title is correct it will be returned along with nil error
//}

//func viewHandler(w http.ResponseWriter, r *http.Request) {
//title := r.URL.Path[len("/view/"):] // extract title from URL to find web page file, by ignoring "/view/" part of link path
//	title, err := getTitle(w, r) // get title by function
//	if err != nil {
//		return
//	}
//	p, err := loadPage(title) // load file by extracted title
//fmt.Fprintf(w, "<h1>Viewing %s</h1><div>%s</div>", p.Title, p.Body) // hardcoded render of file content
//t, _ := template.ParseFiles("view.html") // view will also be using template files
//t.Execute(w, p)
//	if err != nil {
//		http.Redirect(w, r, "/edit/"+title, http.StatusFound) // if page in a view path was not found, server will redirect client to
// edit Page form (code 302) with empty contents, so the contents may be created
// http.Redirect will add code 302 to the header of the HTTP response
//		return
//	}
//	renderTemplate(w, "view", p) // use a dedicated function to render page template
//}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	frontPage := "/view/FrontPage"
	http.Redirect(w, r, frontPage, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc { // wrapper function that takes one of the handler function,
	// and returns a function of type of http.HandlerFunc (suitable to be passed to the function http.HandleFunc)
	// // returned function is called a closure because it encloses values defined outside of it.
	// In this case, the variable fn (the single argument to makeHandler) is enclosed by the closure. The variable fn will be one of our handlers themself.
	return func(w http.ResponseWriter, r *http.Request) { // this closure return by makeHandler (return value of http.HandlerFunc)
		m := validPath.FindStringSubmatch(r.URL.Path) // closure extract the path and validating it
		if m == nil {                                 // if title is invalid ResponseWriter will get an error
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2]) // If everything with path is fine, closure will continue execution inside one of handler function using arguments from closure block
	}
}

//func editHandler(w http.ResponseWriter, r *http.Request) {
//	//title := r.URL.Path[len("/edit/"):]
//	title, err := getTitle(w, r)
//	if err != nil {
//		return
//	}
//	p, err := loadPage(title)
//	if err != nil {
//		p = &Page{Title: title} // if loadPage() return not nil error code, edit handler will create new page
//	}
//fmt.Fprintf(w, "<h1>Editing %s</h1>"+
//	"<form action=\"/save/%s\" method=\"POST\">"+
//	"<textarea name=\"body\">%s</textarea><br>"+
//	"<input type=\"submit\" value=\"Save\">"+
//	"</form>",
//	p.Title, p.Title, p.Body) // hardcoded HTML form
//t, _ := template.ParseFiles("edit.html") // now we parsing .html template file for our web page, returning *template.Template
//t.Execute(w, p) // method will apply *template.Template for web page. t.Execute will edit template by substituition of .Title and .Body in edit.html
// for p.Title and p.Body
// html/template library will guarantee that only safe and correct-looking HTML is generated by template actions.
// For instance, it automatically escapes any greateer than sign (>), replacing it with &gt; to make sure that user data does not corrupt the form HTML
//	renderTemplate(w, "edit", p)
//}

//func saveHandler(w http.ResponseWriter, r *http.Request) { // saveHandler will handle the submission of forms on the edit pages
//	//title := r.URL.Path[len("/len/"):]
//	title, err := getTitle(w, r)
//	if err != nil {
//		return
//	}
//	body := r.FormValue("body")                  // r.FormValue will return string text from textarea named "body" in edit.html
//	p := &Page{Title: title, Body: []byte(body)} // new page is created. And because body is string, we need to convert it into Byte
//	err = p.save()                              // save page file
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError) // print to user the Internal Server Error error screen
//		return
//	}
//	http.Redirect(w, r, "/view/"+title, http.StatusFound) // redirect to the /view/ of the new page
//}

// since we are using the same construction to render templates
// why do not make a function for it?
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	tmplBaseName := tmpl + ".html"
	err := templates.ExecuteTemplate(w, tmplBaseName, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	//http.HandleFunc("/view/", viewHandler) // http.HandleFunc tells server to handle view action on "/view/" path by using viewhandler
	//http.HandleFunc("/edit/", editHandler) // http.HandleFunc tells server to handle edit action on "/edit/" path by using editHandler
	//http.HandleFunc("/save/", saveHandler) // http.HandleFunc tells server to handle edit action on "/save/" path by using saveHandler
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	fmt.Println("Server is running. Ctrl + C to stop it...")
	server := http.ListenAndServe("127.0.0.1:8080", nil)
	if server != nil {
		log.Fatal(server)
	}

	// Page creation require making reference to Page structure in order to make a page
	//test_page_write := &Page{"TestPage", []byte("Any fool can write code that a computer can understand. Good programmers write code that humans can understand.")}
	//test_page_write.save() // method that will save page locally

	//test_page_read, err := loadPage("TestPage") // loadpage() returning Page structure reference, containing title and body
	//if err != nil {                             // stop server if page was not found
	//	fmt.Println("Error happened on page reading")
	//	return
	//}

	//fmt.Printf("%s\n", test_page_read.Body) // print contents of web-page
}
