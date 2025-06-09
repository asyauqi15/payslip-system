package cmd

import (
	"context"
	"fmt"
	"github.com/asyauqi15/payslip-system/internal/transport"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	defaultWaitShutdownDuration = 10 * time.Second
)

var httpServerCmd = &cobra.Command{
	RunE:  runHTTPServer,
	Use:   "http_server",
	Short: "to run http server",
}

func runHTTPServer(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultWaitShutdownDuration)
	defer cancel()

	cfg, err := loadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	r := initRegistry(cfg)

	server, err := transport.NewRESTServer(cfg, r.db)
	if err != nil {
		log.Fatalf("failed to initiate http server: %s", err)
	}

	errCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		log.Println("http server is running")
		if err := server.Start(); err != nil {
			errCh <- fmt.Errorf("failed to run http server: %w", err)
		}
	}()

	go func() {
		<-signalCh
		signal.Reset(os.Interrupt)
		errCh <- fmt.Errorf("interrupted")
	}()

	<-errCh

	if err := server.Stop(ctx); err != nil {
		log.Println(err)
	}

	return nil
}
