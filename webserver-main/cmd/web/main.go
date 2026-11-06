package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	// Import module package of interanl db module
	models "debutwebapp.komardinec/internal/models"
	"github.com/alexedwards/scs/mysqlstore"
	scs "github.com/alexedwards/scs/v2"
	form "github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql" // underscore is here because we don't need anything directly
	// from mysql, except init() function of the driver. This is standard practice for most GO's SQL drivers.
)

// Application dependencies struct to hold the application-wide dependencies
type application struct {
	logger         *slog.Logger
	dbWebApp       models.WebAppModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	serverFlags    *serverConfigStruct
	metricsServer  *Metrics
}

type serverConfigStruct struct {
	serverPort   string
	dbSourceName string
	debugMode    bool
}

func openDB(dsn string) (*sql.DB, error) {
	// The sql.Open() function initializes a new sql.DB object, which is essentially a
	// pool of database connections.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// db.Ping() will test if connection pool to database are actually works
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	// Our servers will be in personal goroutines, so waitgroups will be handy
	var wg sync.WaitGroup
	var serverConfig serverConfigStruct

	// --- Flag section ---
	flag.StringVar(&serverConfig.serverPort, "addr", ":8080", "HTTP Listen address (IP:Port).")
	flag.StringVar(&serverConfig.dbSourceName, "dsn", "webserver:iamwebserver123@/debut?parseTime=true", "MySQL data source name")
	flag.BoolVar(&serverConfig.debugMode, "debug", false, "Flag that enables the debug mode of application. Absence of that flag implying that debug mode is off, if flag is used, that means that debug mode is on.")

	flag.Parse()

	// Create logger to write into terminal various messages regarding software actions and events
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	// To keep the main() function openDB() function,
	// with code for creating a DB connection pool.
	// We will get a DSN from a cmd flag
	db, err := openDB(serverConfig.dbSourceName)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() exits
	defer db.Close()

	// Initialize a new template cache
	templateCache, err := NewTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager.
	// Then we configure it to use our MySQL database as the session store
	// and set a lifetime of 1 hour (so that sessions automatically expire
	// 1 hour after first being created).
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 1 * time.Hour
	sessionManager.Cookie.Secure = true

	// Initialize metrics server with custom register from metrics model
	metricsServer := NewCustomMetrics()
	// Optionally we can set the listen adress of the metrics server
	// metricsServer.Address = ":2112"

	// Initialize a new instance of application struct, containing dependencies for structured logging
	app := &application{
		logger:         logger,
		dbWebApp:       &models.WebAppModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		serverFlags:    &serverConfig,
		metricsServer:  metricsServer,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings
	// we want server to use.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialize the new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before
	srv := &http.Server{
		Addr:         serverConfig.serverPort,
		Handler:      app.Route(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Use the http.ListenAndServe() function starts a new web server.
	// We pass inside 2 arguments - TCP address to listen on
	// and servemux created earlier.
	// Based on information from documentation http.ListenAndServe always return
	// non-nil error so if we something goes wrong logger.Error()
	// will print it as a fatal error.

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Print starting log message
		logger.Info("Server is running. Ctrl + C to stop it...", "address", serverConfig.serverPort)
		err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		logger.Error(err.Error())
		os.Exit(1)
	}()
	// Additional servers can fail as much as they want, while main server is still serving
	// If main process fails, all additional subroutines will be closed too
	metricsServer.MetricsListenAndServe()
	wg.Wait()
}
