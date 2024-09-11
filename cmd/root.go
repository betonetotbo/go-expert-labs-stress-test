package cmd

import (
	"betonetotbo/go-expert-labs-stress-test/internal/stresstest"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "Stress test tool",
	Long:  `Stress test for make HTTP requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			return fmt.Errorf("no url provided")
		}

		requests, err := cmd.Flags().GetInt("requests")
		if err != nil {
			return err
		}

		concurrency, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			return err
		}

		method, err := cmd.Flags().GetString("method")
		if err != nil {
			return err
		}

		headers, err := cmd.Flags().GetStringSlice("header")
		if err != nil {
			return err
		}
		if headers != nil && !stresstest.ValidateHeaders(headers) {
			return fmt.Errorf("invalid headers format")
		}

		timeout, err := cmd.Flags().GetDuration("timeout")
		if err != nil {
			return err
		}

		body, err := cmd.Flags().GetString("body")
		if err != nil {
			return err
		}
		var data []byte
		if body != "" {
			data, err = base64.StdEncoding.DecodeString(body)
			if err != nil {
				return fmt.Errorf("invalid body base64: %w", err)
			}
		}

		report := stresstest.StressTest(url, method, headers, data, timeout, requests, concurrency)

		fmt.Printf("%+v\n", report)

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("url", "u", "", "Target URL")
	rootCmd.Flags().IntP("requests", "r", 10, "Number of requests")
	rootCmd.Flags().IntP("concurrency", "c", 100, "Number of concurrent requests")
	rootCmd.Flags().StringP("method", "m", "GET", "Request method")
	rootCmd.Flags().StringSliceP("header", "H", []string{}, "Header in format of NAME:VALUE")
	rootCmd.Flags().DurationP("timeout", "t", time.Second*5, "Each request timeout, example: 5s")
	rootCmd.Flags().StringP("body", "b", "", "Request body in base64 format")
}
