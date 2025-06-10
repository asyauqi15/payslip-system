package cmd

import (
	"fmt"
	"github.com/asyauqi15/payslip-system/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	migrateCmd.Flags().BoolVarP(&MigrateRollback, "rollback", "r", false, "to rollback the latest version of sql migration")
	migrateCmd.PersistentFlags().StringVarP(&MigrateDir, "dir", "d", "db/migrations", "sql migrations directory")

	seedCmd.PersistentFlags().StringVarP(&SeedDir, "dir", "d", "db/seeds", "sql seeds directory")

	rootCmd.AddCommand(httpServerCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig(path string) (internal.Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("env")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return internal.Config{}, err
	}

	var cfg internal.Config
	if err = viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
