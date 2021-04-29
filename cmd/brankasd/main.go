package main

import (
	"fmt"
	"net/http"
	"os"

	dbsqlx "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/reekoheek/brankas/internal/drivers/sqlx"
	"github.com/reekoheek/brankas/pkg/app/sync"
	"github.com/reekoheek/brankas/web/api"
	"github.com/reekoheek/brankas/web/auth"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbfile := os.Getenv("DBFILE")
	if dbfile == "" {
		dbfile = "brankas.db"
	}

	db, err := dbsqlx.Connect("sqlite3", dbfile)
	if err != nil {
		return err
	}

	defer db.Close()

	if err := migrate(db); err != nil {
		return err
	}

	// authService := auth.NewService(db, []byte(secret), EXPIRES)
	vaultRepo := sqlx.NewVaultRepository(db)
	syncService := sync.New(vaultRepo)

	api := api.New(syncService)
	auth := auth.New()

	fmt.Printf("Listening at :%s\n", port)

	handler := http.NewServeMux()

	handler.Handle("/auth/", http.StripPrefix("/auth", auth.Routes()))
	handler.Handle("/api/", http.StripPrefix("/api", api.Routes()))
	handler.Handle("/", http.FileServer(&spaFileSystem{http.Dir("./web/ui/www")}))

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		return err
	}

	return nil
}

func migrate(db *dbsqlx.DB) error {
	db.MustExec(sqlx.VaultSchema)
	// db.MustExec(auth.Schema)

	// u := &struct{ Username string }{}
	// if err := db.Get(u, `SELECT username FROM USER WHERE username = "admin"`); err != sql.ErrNoRows {
	// 	return nil
	// }

	// pwd, _ := bcrypt.GenerateFromPassword([]byte("password"), 0)
	// db.MustExec(`INSERT INTO user(username, password, creator) VALUES("admin", ?, "")`, pwd)

	return nil
}

type spaFileSystem struct {
	root http.FileSystem
}

func (fs *spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		return fs.root.Open("index.html")
	}
	return f, err
}
