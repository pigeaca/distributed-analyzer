package commands

import (
	"bufio"
	"context"
	"fmt"
	"github.com/DistributedMarketplace/internal/cli/config"
	"os"
	"strings"
)

func StartConsole(ctx context.Context, config config.Config) error {
	fmt.Println("Welcome to AI Task CLI")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		lineCh := make(chan string, 1)
		errCh := make(chan error, 1)

		go func() {
			line, err := reader.ReadString('\n')
			if err != nil {
				errCh <- err
			} else {
				lineCh <- line
			}
		}()

		select {
		case <-ctx.Done():
			fmt.Println("\nInterrupt received, exiting...")
			return nil
		case err := <-errCh:
			return err
		case line := <-lineCh:
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if line == "exit" || line == "quit" {
				fmt.Println("Exit!")
				return nil
			}

			args := strings.Fields(line)
			root := NewRootCommand()
			root.SetArgs(args)
			if err := root.Execute(); err != nil {
				fmt.Println("error:", err)
			}
		}
	}
}
