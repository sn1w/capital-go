package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/sn1w/capital-go/config"
	"github.com/sn1w/capital-go/entities/infrastructures/bitflyer"
	"github.com/sn1w/capital-go/entities/interface/cli"
	"github.com/sn1w/capital-go/entities/usecases"
	cerror "github.com/sn1w/capital-go/error"
	"github.com/spf13/cobra"
)

var bitflyerCmd = &cobra.Command{
	Use:   "bitflyer",
	Short: "Invoke actions related to bitflyer",
}

var usecase = usecases.BitFlyerUseCase{
	Client: bitflyer.NewBitFlyer(config.NewConfig()),
}
var bf = cli.BitFlyerCLI{
	UseCase: usecase,
}

var showMarkets = &cobra.Command{
	Use:   "markets",
	Short: "show avaiable markets",
	Run: func(cmd *cobra.Command, args []string) {
		markets, err := bf.GetAvaiableMarkets()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(markets)
	},
}

var showBoards = &cobra.Command{
	Use:   "board [product_code]",
	Short: "show current board",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		daemon, _ := cmd.Flags().GetBool("daemon")

		boards, err := bf.GetBoard(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(boards)

		if daemon {
			for {
				time.Sleep(time.Second * 5)
				c := exec.Command("clear")
				c.Stdout = os.Stdout
				c.Run()
				boards, err := bf.GetBoard(args[0])
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(boards)

			}
		}
	},
}

var getBalance = &cobra.Command{
	Use:   "balance",
	Short: "show current balance (required authorization)",
	Run: func(cmd *cobra.Command, args []string) {
		balance, err := bf.GetBalance()
		if err != nil {
			if errors.Is(err, cerror.ErrUnAuthorized) {
				fmt.Println("authorization key is missing or invalid. please check your configuration.")
				return
			}

			fmt.Println(err.Error())
			return
		}
		fmt.Println(balance)
	},
}

func init() {
	showBoards.Flags().BoolP("daemon", "d", false, "Using auto reloading")
	commands := []*cobra.Command{
		showMarkets, showBoards, getBalance,
	}
	for _, v := range commands {
		bitflyerCmd.AddCommand(v)
	}
	rootCmd.AddCommand(bitflyerCmd)
}
