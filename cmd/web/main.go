package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"tarala/snippetbox/pkg/models/mysql"
)

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	//STEP 1 CONFIGURATION
	// 	WAY 1  - arguments params
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web:dupa123@/snippetbox?parseTime=true", "MySQL Connection")
	// 	WAY 2 - os environment vars
	//os.Getenv("SNIPPETBOX_ADDR")
	//	WAY3  - create own cfg struct and use StringVar method
	//var cfg = &Config{}
	//flag.StringVar(&cfg.Addr, "addr", ":4000", "Http Network Address")
	//flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	// this is the way to access environment variables
	// it works with -flags as well, but doesn't let you set default values
	// moreover it's always string type

	//LOGGERS for info
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	//for error
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//DATABASE connection
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	//Cache of templates loading
	cache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	//DESIGN PATTERN - instead of keeping loggers as global variables
	//Initialize application object which holds "global" loggers
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: cache,
	}

	//configure ROUTES here
	mux := app.routes()

	//Create own server which uses the same logger to print logs
	srv := &http.Server{
		ErrorLog: errorLog,
		Addr:     *addr,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %v", *addr)
	err = srv.ListenAndServe()
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
