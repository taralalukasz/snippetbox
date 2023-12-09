package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"

	"tarala/snippetbox/pkg/models"
	"tarala/snippetbox/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// this is how you use interfaces to pass various implementaitons to the context.
	// both sql.SnippetModel and mock.SnippetModel have exactly the same method signatures, so can be used if passed that way
	snippets interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}

	users interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	templateCache map[string]*template.Template
	session       *sessions.Session
}

func main() {
	// STEP 1 CONFIGURATION
	// 	WAY 1  - arguments params
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web:dupa123@/snippetbox?parseTime=true", "MySQL Connection")
	// 	WAY 2 - os environment vars
	// os.Getenv("SNIPPETBOX_ADDR")
	//	WAY3  - create own cfg struct and use StringVar method
	// var cfg = &Config{}
	// flag.StringVar(&cfg.Addr, "addr", ":4000", "Http Network Address")
	// flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	// this is the way to access environment variables
	// it works with -flags as well, but doesn't let you set default values
	// moreover it's always string type

	// we need a secret to encrypt session cookies
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret to encrypt session cookies")

	flag.Parse()

	// LOGGERS for info
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// for error
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// DATABASE connection
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Cache of templates loading
	cache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// SESSION CREATION
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	// DESIGN PATTERN - instead of keeping loggers as global variables
	// Initialize application object which holds "global" loggers
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: cache,
		session:       session,
		users:         &mysql.UserModel{DB: db},
	}

	// configure ROUTES here
	mux := app.routes()

	// initialize  custom TLS configuration struct
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		//CipherSuites: []uint16{
		//	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		//	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		//	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		//	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		//},
		//MinVersion: tls.VersionTLS12,
		//MaxVersion: tls.VersionTLS12,
	}

	// Create own server which uses the same logger to print logs
	srv := &http.Server{
		TLSConfig: tlsConfig,
		ErrorLog:  errorLog,
		Addr:      *addr,
		Handler:   mux,
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %v", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}
