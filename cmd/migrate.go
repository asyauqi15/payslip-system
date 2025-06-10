package cmd

import (
	"context"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"log"
)

var (
	MigrateDir      string
	MigrateRollback bool
)

var migrateCmd = &cobra.Command{
	Run:   runMigrate,
	Use:   "migrate",
	Short: "Run database migrations",
}

func runMigrate(_ *cobra.Command, _ []string) {
	ctx := context.TODO()

	cfg, err := loadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	db, err := goose.OpenDBWithDriver("pgx", cfg.Database.Source)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}
	goose.SetTableName("schema_migrations")

	if MigrateRollback {
		if err := goose.RunContext(ctx, "down", db, MigrateDir); err != nil {
			log.Fatalf("goose down: %v", err)
		}
	}

	if err := goose.RunWithOptionsContext(ctx, "up", db, MigrateDir, []string{}, goose.WithAllowMissing()); err != nil {
		log.Fatalf("goose up: %v", err)
	}
}
