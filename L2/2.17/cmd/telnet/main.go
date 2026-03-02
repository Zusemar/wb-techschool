package main

import (
	"fmt"
	"os"
	"time"

	"telnet/internal/client"

	"github.com/spf13/cobra"
)

func main() {
	var timeout time.Duration

	rootCmd := &cobra.Command{
		Use:   "telnet <host> <port>",
		Short: "Simple TCP telnet client",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := args[0]
			port := args[1]

			client := client.New(host, port, timeout)
			fmt.Println("host:   ", host)
			fmt.Println("port:   ", port)
			fmt.Println("timeout:", timeout)
			return client.Run()
		},
	}

	rootCmd.Flags().DurationVar(
		&timeout,
		"timeout",
		10*time.Second,
		"connection timeout (e.g. 5s, 500ms)",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
