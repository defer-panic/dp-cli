package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	articleExportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export article to PDF (pandoc and xelatex are required)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  executeArticleExport,
	}
	generateTOC bool
	pdfEngine   string
)

func init() {
	articleExportCmd.Flags().BoolVar(&generateTOC, "toc", false, "Generate table of contents")
	articleExportCmd.Flags().StringVar(&pdfEngine, "pdf-engine", "xelatex", "PDF engine to use")
}

func executeArticleExport(_ *cobra.Command, args []string) error {
	var (
		outputFilename = getExportFilename(args)
		cmdArgs        = []string{
			"-o", outputFilename,
			args[0],
			fmt.Sprintf("--pdf-engine=%s", pdfEngine),
		}
	)

	if generateTOC {
		cmdArgs = append(cmdArgs, "--toc")
	}

	fmt.Println(args[0], outputFilename, cmdArgs)

	out, err := exec.Command("pandoc", cmdArgs...).CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}

func getExportFilename(args []string) string {
	var outputFilename string

	if len(args) > 1 {
		outputFilename = args[1]
	} else {
		outputFilename = strings.TrimRight(args[0], filepath.Ext(args[0])) + ".pdf"
	}

	return outputFilename
}
