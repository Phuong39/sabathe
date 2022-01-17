package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sabathe/server/register"
	"strings"
)

var RpcxServerAddress string

var rootCmd = &cobra.Command{
	Use:   "sabathe",
	Short: "\n                                                       \n               ,--.             ,--.  ,--.             \n ,---.  ,--,--.|  |-.  ,--,--.,-'  '-.|  ,---.  ,---.  \n(  .-' ' ,-.  || .-. '' ,-.  |'-.  .-'|  .-.  || .-. : \n.-'  `)\\ '-'  || `-' |\\ '-'  |  |  |  |  | |  |\\   --. \n`----'  `--`--' `---'  `--`--'  `--'  `--' `--' `----' \n                                                       \n",
	Run: func(cmd *cobra.Command, args []string) {
		addr := strings.Split(RpcxServerAddress, ":")
		if len(addr) == 2 {
			register.RegisterService(fmt.Sprintf("%v:%v", addr[0], addr[1]))
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&RpcxServerAddress, "rpcx", "0.0.0.0:8443", "set rpcx address")
	rootCmd.Flags().StringVar(&register.AuthPassword, "pass", "kali123", "set client conn server pass")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(fmt.Sprintf("%v", err))
		os.Exit(1)
	}
}
