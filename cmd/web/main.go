package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"jeisaRaja.git/snippetbox/pkg/models/mysql"
)

type application struct {
	infoLog    *log.Logger
	errorLog   *log.Logger
	snippetbox *mysql.SnippetModel
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	defaultDSN := "web:Cipinang01@/snippetbox?parseTime=true"
	addr := flag.String("addr", "4000", "HTTP network address")
	dsn := flag.String("dsn", defaultDSN, "MYSQL Database Pool")

	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		errlog.Println("DB Connection Error!")
	}
	defer db.Close()

	app := &application{
		infoLog:    infoLog,
		errorLog:   errlog,
		snippetbox: &mysql.SnippetModel{DB: db},
	}
	*addr = "127.0.0.1:" + *addr

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errlog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	srv.ListenAndServe()

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
