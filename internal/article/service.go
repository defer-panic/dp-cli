package article

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
)

//go:embed resources/article.tpl.md
var tpl string

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Generate(input GenerateInput) error {
	tpl, err := template.New("article").Parse(tpl)
	if err != nil {
		return err
	}

	// ensure output directory exists
	if err := os.MkdirAll(path.Dir(input.OutputFilename), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.Create(input.OutputFilename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return tpl.Execute(outFile, input)
}

func (s *Service) Export(input ExportInput) error {
	cmdArgs := []string{
		"-o", input.OutputFilename,
		input.InputFilename,
		fmt.Sprintf("--pdf-engine=%s", input.PDFEngine),
	}

	if input.GenerateTOC {
		cmdArgs = append(cmdArgs, "--toc")
	}

	out, err := exec.Command("pandoc", cmdArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to export article: %s", out)
	}

	fmt.Println(string(out))

	return nil
}

type GenerateInput struct {
	OutputFilename string
	Title          string
	Subtitle       string
	Author         string
	Language       string
	DocumentClass  string
	PaperSize      string
	LineStretch    float64
	MainFont       string
	MonoFont       string
}

type ExportInput struct {
	InputFilename  string
	OutputFilename string
	PDFEngine      string
	GenerateTOC    bool
}
