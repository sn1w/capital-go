package cmd

import (
	"fmt"

	"github.com/sn1w/capital-go/entities/usecases"
	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Say hello to user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(usecases.Hello())
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
}
