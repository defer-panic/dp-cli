package cmd

import (
	_ "embed"
	"html/template"
	"os"
	"strconv"

	"github.com/erikgeiser/promptkit/selection"
	"github.com/erikgeiser/promptkit/textinput"
	"github.com/spf13/cobra"
)

var articleGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate article",
	Args:  cobra.ExactArgs(1),
	RunE:  executeArticleGenerate,
}

func executeArticleGenerate(_ *cobra.Command, args []string) error {
	data, err := readDataFromInput()
	if err != nil {
		return err
	}

	if err := generate(data, args[0]); err != nil {
		return err
	}

	return nil
}

func readDataFromInput() (*Data, error) {
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
	mainFontInput.InitialValue = "Times New Roman"

	mainFont, err := mainFontInput.RunPrompt()
	if err != nil {
		return nil, err
	}

	monoFontInput := textinput.New("What monospace font should be used?")
	monoFontInput.InitialValue = "Fira Code"

	monoFont, err := monoFontInput.RunPrompt()
	if err != nil {
		return nil, err
	}

	return &Data{
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

//go:embed resources/article.tpl.md
var tpl string

func generate(data *Data, outputFileName string) error {
	tpl, err := template.New("article").Parse(tpl)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return tpl.Execute(outFile, data)
}

type Data struct {
	Title         string
	Subtitle      string
	Author        string
	Language      string
	DocumentClass string
	PaperSize     string
	LineStretch   float64
	MainFont      string
	MonoFont      string
}
