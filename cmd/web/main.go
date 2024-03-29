package main

import (
	// "crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"jeisaRaja.git/snippetbox/pkg/models/mysql"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	session       *sessions.Session
	users         *mysql.UserModel
}

func main() {

	// tlsConfig := &tls.Config{
	// 	PreferServerCipherSuites: true,
	// 	CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	// }

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	MYSQLUSER := os.Getenv("MYSQLUSER")
	MYSQLPASSWORD := os.Getenv("MYSQLPASSWORD")
	MYSQLDATABASE := os.Getenv("MYSQLDATABASE")
	MYSQLURL := os.Getenv("MYSQLURL")
	productionDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", MYSQLUSER, MYSQLPASSWORD, MYSQLURL, MYSQLDATABASE)

	fmt.Println(MYSQLDATABASE)
	fmt.Println(productionDSN)
	// defaultDSN := "web:mysqlCirebon01@/snippetbox?parseTime=true"
	// addr := flag.String("addr", "4000", "HTTP network address")
	port := os.Getenv("PORT")
	dsn := flag.String("dsn", productionDSN, "MYSQL Database Pool")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key for session")
	flag.Parse()

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.SameSite = http.SameSiteStrictMode
	db, err := openDB(*dsn)
	if err != nil {
		fmt.Println("DB Connection Error:", err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errlog.Fatal()
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errlog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
		users:         &mysql.UserModel{DB: db},
	}
	addr := ":" + port

	fmt.Println(addr)
	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errlog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", addr)
	fmt.Println("running now")
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
