package cmd

import (
	"fmt"

	"github.com/defer-panic/dp-cli/internal/article"
	"github.com/defer-panic/dp-cli/internal/auth"
	"github.com/defer-panic/dp-cli/internal/config"
	"github.com/defer-panic/dp-cli/internal/controller"
	"github.com/defer-panic/dp-cli/internal/url"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dp-cli",
	Short: "Toolkit for managing Defer Panic articles and other stuff",
}

func init() {
	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		fmt.Printf("Error loading config: %v", err)
		return
	}

	controllerConstructors := []controller.ControllerConstructor{
		controller.Article(article.NewService()),
		controller.URL(url.NewService(cfg.Server, cfg.JWT)),
		controller.Auth(auth.NewService(cfg.Server)),
	}

	for _, c := range controllerConstructors {
		ctrl, err := c(cfg)
		if err != nil {
			fmt.Printf("Error building controller: %v", err)
			return
		}

		ctrl.Register(rootCmd)
	}
}

func Execute() error {
	return rootCmd.Execute()
}
