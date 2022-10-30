package config

import (
	"easycoding/internal/config"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type user struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate config.yaml file",
	RunE: func(cmd *cobra.Command, args []string) error {
		dryrun := getDryrun(cmd)
		config.LoadConfig("config.yaml")
		if dryrun {
			// TODO(qujiabao): print pretty
			fmt.Println(viper.AllSettings())
			return nil
		}
		if err := viper.WriteConfig(); err != nil {
			return err
		}

		return nil
	},
}

func initGenerate() {
	configCmd.AddCommand(generateCmd)
}
