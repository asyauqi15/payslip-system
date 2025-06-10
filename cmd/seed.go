package cmd

import (
	"context"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"log"
)

var SeedDir string

var seedCmd = &cobra.Command{
	Run:   runSeed,
	Use:   "seed",
	Short: "Run database seed",
}

func runSeed(_ *cobra.Command, _ []string) {
	ctx := context.TODO()

	cfg, err := loadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	db, err := goose.OpenDBWithDriver("pgx", cfg.Database.Source)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	if err := goose.RunWithOptionsContext(ctx, "up", db, SeedDir, []string{}, goose.WithAllowMissing()); err != nil {
		log.Fatalf("goose up: %v", err)
	}
}
