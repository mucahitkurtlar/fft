package fft

import (
	"fmt"
	"log"
	"os"

	"github.com/mucahitkurtlar/fft/pkg/model"
	"github.com/playwright-community/playwright-go"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fft <url> [flags]",
	Short: "Find Font Types",
	Example: `	fft https://www.example.com
	fft https://www.example.com -m 20`,
	Long: `Find Font Types is a tool to find font types used on a website.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteUrl := args[0]
		maxPageCount, err := cmd.Flags().GetUint16("max-pages")
		if err != nil {
			log.Fatalf("Error getting max-pages flag: %s", err)
		}
		goRoutineCount, err := cmd.Flags().GetUint8("go-routines")
		if err != nil {
			log.Fatalf("Error getting go-routines flag: %s", err)
		}
		goToTimeout, err := cmd.Flags().GetFloat64("go-to-timeout")
		if err != nil {
			log.Fatalf("Error getting go-to-timeout flag: %s", err)
		}
		netIdleTimeout, err := cmd.Flags().GetFloat64("net-idle-timeout")
		if err != nil {
			log.Fatalf("Error getting net-idle-timeout flag: %s", err)
		}

		crawlerOpts := &model.CrawlerOpts{
			MaxPageCount:   maxPageCount,
			GoRoutineCount: goRoutineCount,
			GoToTimeout:    goToTimeout,
			NetIdleTimeout: netIdleTimeout,
		}

		log.Printf("Website URL: %s", siteUrl)
		log.Printf("Max URLs: %d", maxPageCount)
		log.Printf("Go Routines: %d", goRoutineCount)
		log.Printf("Go To Timeout: %f", goToTimeout)
		log.Printf("Net Idle Timeout: %f", netIdleTimeout)
		Root(siteUrl, crawlerOpts)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Error executing root command: %s", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Uint16P("max-pages", "m", 1, "maximum number of pages to crawl")
	rootCmd.Flags().Uint8P("go-routines", "g", 3, "number of go routines to use")
	rootCmd.Flags().Float64P("go-to-timeout", "t", 30000, "timeout for page navigation (ms)")
	rootCmd.Flags().Float64P("net-idle-timeout", "n", 30000, "timeout for network idle (ms)")
}

func Root(siteUrl string, crawlerOpts *model.CrawlerOpts) {
	err := playwright.Install()
	if err != nil {
		log.Fatal(err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer pw.Stop()

	err = model.StartCrawler(pw, siteUrl, crawlerOpts)
	if err != nil {
		log.Fatal(err)
	}
}
