package main

import (
	"io"
	"log"
	"os"

	"github.com/ThePianist/flowkick/cmd"
	"github.com/ThePianist/flowkick/logger"
	"github.com/ThePianist/flowkick/store"
	tea "github.com/charmbracelet/bubbletea"

	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	store := new(store.Store)
	if err := store.Init(); err != nil {
		log.Fatalf("unable to init store: %v", err)
	}
	db := store.Conn

	// Mute Goose logs by default, enable with GOOSE_LOG=on
	if os.Getenv("GOOSE_LOG") != "on" {
		goose.SetLogger(log.New(io.Discard, "", 0))
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	if err := store.Init(); err != nil {
		log.Fatalf("unable to init store: %v", err)
	}

	logger.Init("flowkick.log")
	p := tea.NewProgram(cmd.InitialAppModel(store))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
