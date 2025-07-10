package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var file string
var model string

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a new task",
	RunE: func(cmd *cobra.Command, args []string) error {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {

			}
		}(f)

		fw, err := writer.CreateFormFile("file", filepath.Base(file))
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, f)
		if err != nil {
			return err
		}

		_ = writer.WriteField("model", model)

		err = writer.Close()
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", "http://localhost:8080/tasks", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
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

		fmt.Println("Task submitted. ID:", result["taskId"])
		return nil
	},
}

func init() {
	submitCmd.Flags().StringVar(&file, "file", "", "Path to input file")
	submitCmd.Flags().StringVar(&model, "model", "default", "Model to use")
	submitCmd.MarkFlagRequired("file")
}
