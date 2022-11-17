package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/defer-panic/dp-cli/internal/auth"
	"github.com/defer-panic/dp-cli/internal/config"
	"github.com/spf13/cobra"
)

type authController struct {
	svc *auth.Service
	cfg *config.Config
	cmd *cobra.Command
}

func Auth(svc *auth.Service) ControllerConstructor {
	return func(cfg *config.Config) (Controller, error) {
		return &authController{
			cfg: cfg,
		}, nil
	}
}

func (c *authController) Register(root *cobra.Command) {
	c.cmd = &cobra.Command{
		Use:   "login",
		Short: "Login to your account on https://dfrp.cc (default) or other instance of dfrp-like infrastructure",
		Args:  cobra.MaximumNArgs(1),
		RunE:  c.GetAuthLink,
	}

	c.cmd.PersistentFlags().StringP("server", "s", "https://dfrp.cc", "server to use for authentication")

	saveTokenCmd := &cobra.Command{
		Use:   "save-token",
		Short: "Save token to config file",
		Args:  cobra.ExactArgs(1),
		RunE:  c.SaveToken,
	}

	c.cmd.AddCommand(saveTokenCmd)
	root.AddCommand(c.cmd)
}

func (c *authController) GetAuthLink(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	link, err := c.svc.GetAuthLink(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Please open the following URL in your browser and follow the instructions:")
	fmt.Println(link)
	return nil
}

func (c *authController) SaveToken(_ *cobra.Command, args []string) error {
	loginServer, err := c.cmd.Flags().GetString("server")
	if err != nil {
		return err
	}

	if err := c.svc.SaveToken(args[0], loginServer); err != nil {
		return err
	}

	return nil
}
