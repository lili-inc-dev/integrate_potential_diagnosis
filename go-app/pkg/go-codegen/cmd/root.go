/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"os"

	"github.com/80andCo/LiLi-LABO/pkg/go-codegen/generator"
	"github.com/spf13/cobra"
)

var (
	inputFilePath string
	model         string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-protogen",
	Short: "create struct for .proto file from golang struct",
	Long:  `create struct for .proto file from golang struct.please specify the import file and output path.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := generator.GenerateProtoFile(inputFilePath, model)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&inputFilePath, "inputFilePath", "example.go", "input file path")
	rootCmd.PersistentFlags().StringVar(&model, "model", "example", "model name")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
