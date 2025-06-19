package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var taskID string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check task status",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/tasks/%s", taskID))
		if err != nil {
			return err
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		var result map[string]any

		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return err
		}

		fmt.Printf("Task ID: %s\nStatus: %s\n", result["id"], result["status"])
		return nil
	},
}

func init() {
	statusCmd.Flags().StringVar(&taskID, "id", "", "Task ID")
	statusCmd.MarkFlagRequired("id")
}
