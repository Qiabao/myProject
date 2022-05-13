package main

import (
	"context"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"log"
	"myProject/lib"
	"sync"
	"time"

	"errors"
	"fmt"
	"github.com/spf13/cobra"
)
type task struct {
	to    crypto.Address
	value int
}

var (
	cli, _   = http.New("http://127.0.0.1:26657", "/websocket")
	logger     *log.Logger
    dbs        []dbm.DB
	channels   []chan *task
	wallet     *lib.Wallet
)

func main() {
	rootCmd := &cobra.Command{
		Use: "cli",
	}

	walletCmd := &cobra.Command{
		Use: "init-wallet",
		Run: func(cmd *cobra.Command, args []string) { initWallet() },
	}

	issueCmd := &cobra.Command{
		Use: "issue-tx",
		Run: func(cmd *cobra.Command, args []string) { issue() },
	}

	transferCmd := &cobra.Command{
		Use: "transfer-tx",
		Run: func(cmd *cobra.Command, args []string) { transfer() },
	}

	queryCmd := &cobra.Command{
		Use: "query",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("query who?")
			}
			label := args[0]
			query(label)
			return nil
		},
	}

	infoCmd:=&cobra.Command{
		Use: "info",
		Run: func(cmd *cobra.Command, args []string)  {
			info()
		},
	}
	rootCmd.AddCommand(walletCmd)
	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(transferCmd)
	rootCmd.AddCommand(queryCmd)

	rootCmd.AddCommand(infoCmd)

	rootCmd.Execute()

	//initWallet()
	//issue()
	//transfer()
	//query("michael")
	//query("britney")
}

func info() {
	ret, err := cli.ABCIInfo(context.Background())
	if err!=nil {
		panic(err)
	}
	fmt.Println(ret)
}

func issue() {
	//wallet := lib.LoadWallet("./wallet")
	tx := lib.NewTx(lib.NewIssuePayload(
		wallet.GetAddress("issuer"),
		wallet.GetAddress("michael"),
		1000))
	tx.Sign(wallet.GetPrivKey("issuer"))
	bz, err := lib.MarshalBinary(tx)
	if err != nil {
		panic(err)
	}
	ret, err := cli.BroadcastTxCommit(context.Background(),bz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("issue ret => %+v\n", ret)
}

func transfer() {

	tx := lib.NewTx(lib.NewTransferPayload(
		wallet.GetAddress("michael"),
		wallet.GetAddress("britney"),
		100))
	tx.Sign(wallet.GetPrivKey("michael"))
	bz, err := lib.MarshalBinary(tx)
	if err != nil {
		panic(err)
	}
	//发送交易
	ret, err := cli.BroadcastTxCommit(context.Background(), bz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("issue ret => %+v\n", ret)
}

func query(label string)() {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(6)
	for i := 0; i < 6; i++ {
		go func() {
			defer wg.Done()
			for i:=0;i<300000;i++{
				cli.BroadcastTxSync(context.Background(), types.Tx(string(i)))
			}
		}()
	}
	wg.Wait()
	end:=time.Now()

	cli.ABCIQuery(context.Background(), "", bytes.HexBytes(label))

	fmt.Println(end.Sub(start))
}

func initWallet() {
	wallet := lib.NewWallet()
	wallet.GenPrivKey("issuer")
	wallet.GenPrivKey("michael")
	wallet.GenPrivKey("britney")
	wallet.Save("./wallet")
}
