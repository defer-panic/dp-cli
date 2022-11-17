package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/defer-panic/dp-cli/internal/config"
	"github.com/defer-panic/dp-cli/internal/url"
	"github.com/jedib0t/go-pretty/v6/text"
	. "github.com/samber/mo"
	"github.com/spf13/cobra"
)

type urlController struct {
	svc      *url.Service
	cfg      *config.Config
	urlCmd   *cobra.Command
	statsCmd *cobra.Command
}

func URL(svc *url.Service) ControllerConstructor {
	return func(cfg *config.Config) (Controller, error) {
		return &urlController{
			svc: svc,
			cfg: cfg,
		}, nil
	}
}

func (c *urlController) Register(root *cobra.Command) {
	c.urlCmd = &cobra.Command{
		Use:   "url",
		Short: "Shorten given URL",
		Args:  cobra.ExactArgs(1),
		RunE:  c.ShortenLink,
	}

	c.urlCmd.Flags().StringP("identifier", "i", "", "custom identifier to use for the short URL")

	c.statsCmd = &cobra.Command{
		Use:   "stats",
		Short: "Show stats for given short URL",
		Args:  cobra.ExactArgs(1),
		RunE:  c.GetStats,
	}

	c.statsCmd.Flags().Bool("json", false, "output stats in JSON format")

	c.urlCmd.AddCommand(c.statsCmd)
	root.AddCommand(c.urlCmd)
}

func (c *urlController) ShortenLink(_ *cobra.Command, args []string) error {
	identifier := None[string]()
	inputIdentifier, err := c.urlCmd.Flags().GetString("identifier")
	if err != nil {
		return err
	}

	if inputIdentifier != "" {
		identifier = Some(inputIdentifier)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	shortURL, err := c.svc.Shorten(ctx, args[0], identifier)
	if err != nil {
		return err
	}

	fmt.Println(shortURL)
	return nil
}

func (c *urlController) GetStats(_ *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	urlStatsJSONOutput, err := c.statsCmd.Flags().GetBool("json")
	if err != nil {
		return err
	}

	stats, err := c.svc.Stats(ctx, args[0])
	if err != nil {
		return err
	}

	if urlStatsJSONOutput {
		return json.NewEncoder(os.Stdout).Encode(stats)
	}

	fmt.Printf("Identifier:   %s\n", text.AlignDefault.Apply(stats.Identifier, 30))
	fmt.Printf("Original URL: %s\n", text.AlignDefault.Apply(stats.OriginalURL, 30))
	fmt.Printf("Created by:   %s\n", text.AlignDefault.Apply(stats.CreatedBy, 30))
	fmt.Printf("Visits:       %s\n", text.AlignDefault.Apply(strconv.Itoa(int(stats.Visits)), 30))
	return nil
}
