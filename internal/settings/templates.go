package settings

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/goccy/go-yaml"
	"github.com/prbf2-tools/svctl/pkg/templates"
)

func (s *Settings) TemplateData() (templates.Values, *GameConfig, error) {
	var allValuesSources []templates.Values

	config, err := s.Config()
	if err != nil {
		return nil, nil, err
	}

	for _, source := range config.Values {
		if source.File != "" {
			sourceFile := source.File
			if !filepath.IsAbs(sourceFile) {
				sourceFile = filepath.Join(s.Path, sourceFile)
			}

			content, err := os.ReadFile(sourceFile)
			if err != nil {
				return nil, nil, err
			}

			var values templates.Values
			err = yaml.Unmarshal(content, &values)
			if err != nil {
				return nil, nil, err
			}

			allValuesSources = append(allValuesSources, values)
		} else if source.Values != nil {
			allValuesSources = append(allValuesSources, source.Values)
		}
	}

	allValues, err := templates.MergeValues(allValuesSources...)
	if err != nil {
		return nil, nil, err
	}

	return allValues, &config.Game, nil
}

func cloneTemplates(path, repoURL, token string) error {
	var auth transport.AuthMethod
	if token != "" {
		auth = &http.BasicAuth{
			Username: "git",
			Password: token,
		}
	}

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:             repoURL,
		Auth:            auth,
		InsecureSkipTLS: true,
	})
	return err
}

func writeValues(path string, content []byte) error {
	commented := commentOutWholeYamlFile(string(content))

	return os.WriteFile(filepath.Join(path, defaultValuesFile), []byte(commented), 0644)
}

func commentOutWholeYamlFile(content string) string {
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)

	builder := strings.Builder{}

	for scanner.Scan() {
		text := scanner.Text()

		if !strings.HasPrefix(strings.TrimSpace(text), "#") {
			if len(strings.TrimSpace(text)) > 0 {
				text = fmt.Sprintf("# %s", text)
			}
		}

		builder.WriteString(text)
		builder.WriteString("\n")
	}

	return builder.String()
}
