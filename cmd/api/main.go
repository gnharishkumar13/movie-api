package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"github.com/gnharishkumar13/movie-api/internal/mailer"
	"os"
	"runtime"
	"time"

	"github.com/gnharishkumar13/movie-api/internal/data"
	"github.com/gnharishkumar13/movie-api/internal/jsonlog"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
}

type config struct {
	port     int
	env      string
	database struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 3000, " port to run the application")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.database.dsn, "db-dsn", os.Getenv("MOVIEDB_DB_DSN"), "PostgreSQL DSN")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "", "SMTP sender")


	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	//database setup
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	//migrations

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", migrationDriver)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("database migrations applied", nil)

	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))

	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.New(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	// Call app.serve() to start the server.
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.database.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
