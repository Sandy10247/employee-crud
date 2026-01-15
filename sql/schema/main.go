// This is custom goose binary with sqlite3 support only.

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	db "server/init"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	// Parse Flags and Args
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("goose: failed to parse flags: %v", err)
	}
	args := flags.Args()

	// check for at least minimum args are provided
	if len(args) != 1 {
		flags.Usage()
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("goose: failed to get Current Working Directory %v", err)
	}

	// Construct the Root Directory Path
	envPath := filepath.Join(cwd, "..", "..", ".env")

	// Load ".env"
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatal("ERROR loading env :", err)
	}

	// Load Config
	c := db.LoadConfig()

	// Construct postgres Connection String
	dsnString := db.DSN(&c)

	// Command to execute
	command := args[0]

	// Connect to DB
	db, err := goose.OpenDBWithDriver("postgres", dsnString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	// Close the DB connetion
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v", err)
		}
	}()

	// Execute Migration
	if err := goose.RunContext(context.Background(), command, db, *dir); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
