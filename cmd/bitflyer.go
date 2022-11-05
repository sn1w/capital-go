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
	Short: "Actions related to bitflyer",
}

var bf = cli.NewBitFlyerCli(
	usecases.NewBitFlyerUseCase(bitflyer.NewBitFlyer(config.NewConfig())),
)

var showMarkets = func() *cobra.Command {
	return &cobra.Command{
		Use:   "markets",
		Short: "Show avaiable markets",
		Run: func(cmd *cobra.Command, args []string) {
			markets, err := bf.GetAvaiableMarkets()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(markets)
		},
	}
}

var showBoards = func() *cobra.Command {
	cmd := cobra.Command{
		Use:   "board [product_code]",
		Short: "Show current board",
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
	cmd.Flags().BoolP("daemon", "d", false, "Using auto reloading")
	return &cmd
}

var getBalance = func() *cobra.Command {
	return &cobra.Command{
		Use:   "balance",
		Short: "Show current balance (required authorization)",
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
}
var sendOrder = func() *cobra.Command {
	var productCode string
	var price float64
	var size float64

	cmd := cobra.Command{
		Use:   "orders",
		Short: "Actions related to order (required authorization)",
	}

	runner := func(buy bool) func(*cobra.Command, []string) {
		return func(*cobra.Command, []string) {
			res, err := bf.CreateOrder(cli.CreateOrderArgument{
				ProductCode: productCode,
				Price:       price,
				Size:        size,
				Buy:         buy,
			})
			if err != nil {
				if errors.Is(err, cerror.ErrUnAuthorized) {
					fmt.Println("authorization key is missing or invalid. please check your configuration.")
					return
				}

				fmt.Println(err.Error())
				return
			}
			fmt.Println(res)
		}
	}

	buy := cobra.Command{
		Use:   "buy",
		Short: "Send 'buy' order (required authorization)",
		Run:   runner(true),
	}

	sell := cobra.Command{
		Use:   "sell",
		Short: "Send 'sell' order (required authorization)",
		Run:   runner(false),
	}

	cmd.AddCommand(&buy)
	cmd.AddCommand(&sell)

	cmd.PersistentFlags().StringVarP(&productCode, "code", "c", "PRODUCT_CODE", "product code you want to order (required)")
	cmd.PersistentFlags().Float64VarP(&price, "price", "p", 0, "order price (required)")
	cmd.PersistentFlags().Float64VarP(&size, "size", "s", 0, "order size (required)")

	cmd.MarkPersistentFlagRequired("product_code")
	cmd.MarkPersistentFlagRequired("price")
	cmd.MarkPersistentFlagRequired("size")

	cmd.Flags().SortFlags = false

	return &cmd
}

func init() {
	commands := []*cobra.Command{
		showMarkets(), showBoards(), getBalance(), sendOrder(),
	}
	for _, v := range commands {
		bitflyerCmd.AddCommand(v)
	}
	rootCmd.AddCommand(bitflyerCmd)
}
