package controller

import (
	"fmt"

	"github.com/defer-panic/dp-cli/internal/build"
	"github.com/defer-panic/dp-cli/internal/config"
	"github.com/spf13/cobra"
)

type versionController struct {}

func Version() ControllerConstructor {
	return func(cfg *config.Config) (Controller, error) {
		return &versionController{}, nil
	}
}

func (c *versionController) Register(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Get dp-cli version",
		Long:  `All software has versions. This is dp-cli's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("dp-cli %s (%s)\n", build.Version, build.CommitHash)
		},
	})
}
