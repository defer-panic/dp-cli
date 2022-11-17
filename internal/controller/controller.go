package controller

import (
	"github.com/defer-panic/dp-cli/internal/config"
	"github.com/spf13/cobra"
)

type Controller interface {
	Register(root *cobra.Command)
}

type ControllerConstructor func(cfg *config.Config) (Controller, error)

