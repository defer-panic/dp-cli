package article_test

import (
	"os"
	"testing"

	"github.com/defer-panic/dp-cli/internal/article"
	"github.com/stretchr/testify/require"
)

func TestService_Generate(t *testing.T) {
	t.Run("generates article boilerplate file", func(t *testing.T) {
		t.Cleanup(func() {
			os.Remove("testdata/article.md")
		})

		var (
			input = article.GenerateInput{
				OutputFilename: "testdata/article.md",
				Title:          "Test Article",
				Subtitle:       "Test Subtitle",
				Author:         "Test Author",
				Language:       "en",
				DocumentClass:  "article",
				PaperSize:      "a4",
				LineStretch:    1.5,
				MainFont:       "Times New Roman",
				MonoFont:       "Courier New",
			}
			svc = article.NewService()
		)

		require.NoError(t, svc.Generate(input))
		require.FileExists(t, input.OutputFilename)

		outputFileContent, err := os.ReadFile(input.OutputFilename)
		require.NoError(t, err)

		expectedFileContent, err := os.ReadFile("testdata/article.example.md")
		require.NoError(t, err)

		require.Equal(t, string(expectedFileContent), string(outputFileContent))
	})
}

func TestService_Export(t *testing.T) {
	t.Run("generates PDF from article using pandoc", func(t *testing.T) {
		t.Cleanup(func() {
			os.Remove("testdata/article.pdf")
		})

		var (
			input = article.ExportInput{
				InputFilename:  "testdata/article.example.md",
				OutputFilename: "testdata/article.pdf",
				PDFEngine:      "xelatex",
				GenerateTOC:    true,
			}
			svc = article.NewService()
		)

		require.NoError(t, svc.Export(input))
		require.FileExists(t, input.OutputFilename)
	})
}
