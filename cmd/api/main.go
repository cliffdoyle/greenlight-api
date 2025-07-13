package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"
)

// Declare a string containing the application version number. Later in the book we'll
// generate this automatically at build time, but for now we'll just store the version
// number as a hard-coded global constant.
const version = "1.0.0"

type olifo string

// Define a config struct to hold all the configuration settings for our application.
// For now, the only configuration settings will be the network port that we want the
// server to listen on, and the name of the current operating environment for the
// application (development, staging, production, etc.). We will read in these
// configuration settings from command-line flags when the application starts.
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	// logger := slog.New(slog.NewTextHandler(os.Stdout,nil))
	// logger.Info("This is informational message","user_id",45,"ip","127.0.0.1")

	// Declare an instance of the config struct
	var cfg config

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	//Read the DSN value from the db-dsn command-line flag into the config struct.
	//We default to using our development DSN if no flag is provided
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	// Initialize a new structured logger which writes log entries to the standard out
	// stream
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Declare an instance of the application struct, containing the config
	// struct and the logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare a new servemux and add a /v1/healthcheck route which dispatches requests
	// to the healthcheckerhandler method
	// mux := app.routes()

	// Declare a HTTP server which listens on the port provided in the config struct,
	// uses the servemux we created above as the handler, has some sensible timeout
	// settings and writes any log messages to the structured logger at Error level.

	fmt.Println("num of goroutines",runtime.NumGoroutine())
	fmt.Println("logical CPUs available",runtime.NumCPU())
	//Call the openDB() to create the connection pool passing in
	//the config struct.
	db, err := openDB(cfg)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	//Defer a call to db.Close() so that the connection pool is closed
	defer db.Close()

	logger.Info("database connection pool established")
	//before the main() function exits
	//Use the httprouter instance returned by app.routes() as the server handler

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the HTTP server
	logger.Info("Starting server", "addr", server.Addr, "env", cfg.env)

	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
