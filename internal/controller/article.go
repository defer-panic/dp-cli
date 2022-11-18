package controller

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/defer-panic/dp-cli/internal/article"
	"github.com/defer-panic/dp-cli/internal/config"
	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/spf13/cobra"
)

type articleController struct {
	svc         *article.Service
	articleCmd  *cobra.Command
	generateCmd *cobra.Command
	exportCmd   *cobra.Command
}

func Article(svc *article.Service) ControllerConstructor {
	return func(_ *config.Config) (Controller, error) {
		return &articleController{svc: svc}, nil
	}
}

func (c *articleController) Register(root *cobra.Command) {
	c.articleCmd = &cobra.Command{
		Use:   "article",
		Short: "Manage articles",
	}
	c.generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate article",
		Args:  cobra.ExactArgs(1),
		RunE:  c.Generate,
	}
	c.exportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export article to PDF (pandoc and xelatex are required)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  c.Export,
	}

	c.exportCmd.Flags().Bool("toc", false, "Generate table of contents")
	c.exportCmd.Flags().String("pdf-engine", "xelatex", "PDF engine to use")

	c.articleCmd.AddCommand(c.generateCmd, c.exportCmd)
	root.AddCommand(c.articleCmd)
}

func (c *articleController) Generate(_ *cobra.Command, args []string) error {
	data, err := c.readDataForTemplate()
	if err != nil {
		return err
	}

	data.OutputFilename = args[0]

	return c.svc.Generate(*data)
}

func (c *articleController) Export(_ *cobra.Command, args []string) error {
	var outputFilename string

	if len(args) > 1 {
		outputFilename = args[1]
	} else {
		outputFilename = strings.TrimRight(args[0], filepath.Ext(args[0])) + ".pdf"
	}

	toc, err := c.exportCmd.Flags().GetBool("toc")
	if err != nil {
		return err
	}

	pdfEngine, err := c.exportCmd.Flags().GetString("pdf-engine")
	if err != nil {
		return err
	}

	input := article.ExportInput{
		InputFilename:  args[0],
		OutputFilename: outputFilename,
		GenerateTOC:    toc,
		PDFEngine:      pdfEngine,
	}

	return c.svc.Export(input)
}

func (c *articleController) readDataForTemplate() (*article.GenerateInput, error) {
	titleInput := textinput.New("Choose the best article title:")
	titleInput.Placeholder = "Think hard!"

	title, err := titleInput.RunPrompt()
	if err != nil {
		return nil, err
	}

	subtitleInput := textinput.New("And well supporting subtitle:")
	subtitleInput.Placeholder = "Think harder!"

	subtitle, err := subtitleInput.RunPrompt()
	if err != nil {
		return nil, err
	}

	authorInput := textinput.New("Who is the author of this masterpiece?")
	authorInput.Placeholder = "Look at yourself in the mirror"
	authorInput.InitialValue = "Ильдар Карымов <hi@ildarkarymov.ru>, Алексей Ким <me@ameyuuno.io>"

	author, err := authorInput.RunPrompt()
	if err != nil {
		return nil, err
	}

	languageSelection := selection.New(
		"What language is this article written in?",
		[]*selection.Choice{
			selection.NewChoice("ru"),
			selection.NewChoice("en"),
		},
	)
	languageSelection.Filter = nil

	language, err := languageSelection.RunPrompt()
	if err != nil {
		return nil, err
	}

	documentClassSelection := selection.New(
		"What document class should be used?",
		[]*selection.Choice{
			selection.NewChoice("report"),
			selection.NewChoice("article"),
			selection.NewChoice("book"),
		},
	)
	documentClassSelection.Filter = nil

	documentClass, err := documentClassSelection.RunPrompt()
	if err != nil {
		return nil, err
	}

	paperSizeSelection := selection.New(
		"What paper size should be used?",
		[]*selection.Choice{
			selection.NewChoice("a4"),
			selection.NewChoice("letter"),
		},
	)
	paperSizeSelection.Filter = nil

	paperSize, err := paperSizeSelection.RunPrompt()
	if err != nil {
		return nil, err
	}

	lineStretchInput := textinput.New("What line stretch should be used?")
	lineStretchInput.InitialValue = "1.5"

	lineStretchStr, err := lineStretchInput.RunPrompt()
	if err != nil {
		return nil, err
	}

	lineStretch, err := strconv.ParseFloat(lineStretchStr, 64)
	if err != nil {
		return nil, err
	}

	mainFontInput := textinput.New("What main font should be used?")
	mainFontInput.InitialValue = "CMUSerif-Roman"

	mainFont, err := mainFontInput.RunPrompt()
	if err != nil {
		return nil, err
	}

	monoFontInput := textinput.New("What monospace font should be used?")
	monoFontInput.InitialValue = "CMUTypewriter-Regular"

	monoFont, err := monoFontInput.RunPrompt()
	if err != nil {
		return nil, err
	}

	return &article.GenerateInput{
		Title:         title,
		Subtitle:      subtitle,
		Author:        author,
		Language:      language.String,
		DocumentClass: documentClass.String,
		PaperSize:     paperSize.String,
		LineStretch:   lineStretch,
		MainFont:      mainFont,
		MonoFont:      monoFont,
	}, nil
}
