package cmd

import (
	"github.com/asyauqi15/payslip-system/internal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type registry struct {
	db *gorm.DB
}

func initRegistry(cfg internal.Config) *registry {
	db, err := gorm.Open(postgres.Open(cfg.Database.Source), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	return &registry{
		db: db,
	}
}
