package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Application dependencies struct to hold the application-wide dependencies
type application struct {
	logger *slog.Logger
}

type serverConfigStruct struct {
	serverPort string
}

func main() {
	// Define destinct parameter variable to control the server port number
	//serverPort := 8080

	var serverConfig serverConfigStruct

	// --- Flag section ---
	flag.StringVar(&serverConfig.serverPort, "addr", ":8080", "HTTP Listen address (IP:Port).")

	flag.Parse()

	// Create logger to write into terminal various messages regarding software actions and events
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	// Print starting log message
	logger.Info("Server is running", "address", serverConfig.serverPort)

	// Initialize a new instance of application struct, containing dependencies for structured logging
	app := &application{
		logger: logger,
	}

	// Use the http.ListenAndServe() function starts a new web server.
	// We pass inside 2 arguments - TCP address to listen on
	// and servemux created earlier.
	// Based on information from documentation http.ListenAndServe always return
	// non-nil error so if we something goes wrong log.Fatal()
	// will print it as a fatal error
	err := http.ListenAndServe(serverConfig.serverPort, app.Route())
	logger.Error(err.Error())
}
