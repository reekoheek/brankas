package main

import (
	"fmt"
	"net/http"
	"os"

	dbsqlx "github.com/jmoiron/sqlx"
	"github.com/reekoheek/brankas/internal/drivers/sqlx"
	"github.com/reekoheek/brankas/pkg/app/sync"
	"github.com/reekoheek/brankas/web/api"
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

	fmt.Printf("Listening at :%s\n", port)

	handler := http.NewServeMux()
	handler.Handle("/api/", http.StripPrefix("/api", api.Routes()))

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		return err
	}

	return nil
}

func migrate(db *dbsqlx.DB) error {
	// db.MustExec(vault.Schema)
	// db.MustExec(auth.Schema)

	// u := &struct{ Username string }{}
	// if err := db.Get(u, `SELECT username FROM USER WHERE username = "admin"`); err != sql.ErrNoRows {
	// 	return nil
	// }

	// pwd, _ := bcrypt.GenerateFromPassword([]byte("password"), 0)
	// db.MustExec(`INSERT INTO user(username, password, creator) VALUES("admin", ?, "")`, pwd)

	return nil
}
