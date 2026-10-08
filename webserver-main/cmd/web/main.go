package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Application dependencies struct to hold the application-wide dependencies
type application struct {
	logger *slog.Logger
}

type serverConfigStruct struct {
	serverPort   string
	dbSourceName string
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	// Define destinct parameter variable to control the server port number

	var serverConfig serverConfigStruct

	// --- Flag section ---
	flag.StringVar(&serverConfig.serverPort, "addr", ":8080", "HTTP Listen address (IP:Port).")
	flag.StringVar(&serverConfig.dbSourceName, "dsn", "webserver:iamwebserver123@/testdb?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Create logger to write into terminal various messages regarding software actions and events
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	// Print starting log message
	logger.Info("Server is running", "address", serverConfig.serverPort)

	// Initialize a new instance of application struct, containing dependencies for structured logging
	app := &application{
		logger: logger,
	}
	// To keep main() tidy
	// Connection to db pool wil be created by separate openDB() function
	// *dsn will be passed as a command-line flag
	db, err := openDB(serverConfig.dbSourceName)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// Use the http.ListenAndServe() function starts a new web server.
	// We pass inside 2 arguments - TCP address to listen on
	// and servemux created earlier.
	// Based on information from documentation http.ListenAndServe always return
	// non-nil error so if we something goes wrong log.Fatal()
	// will print it as a fatal error
	err = http.ListenAndServe(serverConfig.serverPort, app.Route())
	logger.Error(err.Error())
}
