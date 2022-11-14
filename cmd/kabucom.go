package cmd

import (
	"fmt"

	"github.com/sn1w/capital-go/config"
	"github.com/sn1w/capital-go/entities/infrastructures/kabucom"
	"github.com/sn1w/capital-go/entities/interface/cli"
	"github.com/sn1w/capital-go/entities/usecases"
	"github.com/spf13/cobra"
)

var kabucomCmd = &cobra.Command{
	Use:   "kabucom",
	Short: "Actions related to kabucom",
}

var kb = cli.NewKabucomCli(
	usecases.NewKabucomUseCase(
		kabucom.NewKabucomClient(config.NewConfig()),
	),
)

var doAuth = func() *cobra.Command {
	var pwd string
	cmd := &cobra.Command{
		Use:   "authorize",
		Short: "Auth kabucom service",
		Run: func(cmd *cobra.Command, args []string) {
			output, err := kb.Authorization(pwd)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(output)
		},
	}

	cmd.Flags().StringVarP(&pwd, "password", "p", "", "api password (required)")
	if err := cmd.MarkFlagRequired("password"); err != nil {
		panic(err)
	}

	return cmd
}

var getKabucomBalance = func() *cobra.Command {
	var apiKey string

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Get Balance",
		Run: func(cmd *cobra.Command, args []string) {
			output, err := kb.GetBalance(
				cli.Parameters{
					APIKey: apiKey,
				},
			)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(output)
		},
	}

	cmd.Flags().StringVarP(&apiKey, "token", "t", "", "api token (required)")
	if err := cmd.MarkFlagRequired("token"); err != nil {
		panic(err)
	}

	return cmd
}

func init() {
	subCommands := []*cobra.Command{
		doAuth(), getKabucomBalance(),
	}

	for _, v := range subCommands {
		kabucomCmd.AddCommand(v)
	}
	rootCmd.AddCommand(kabucomCmd)
}
