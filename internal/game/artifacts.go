package game

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// chatlog, newest appends
// adminlog appending
// joinlog appending
// banlog appending
// ticketslog appending
// coincidentipslog appending
// pythonerrorlog newest appends
// pythonlaunchlog newest appends
// bf2demo finished when newer file created (when running)
// prdemo newest file
// prdemo_private newest file

const (
	modPath           = "mods/pr"
	raConfigFile      = modPath + "/python/game/realityconfig_admin.py"
	trackerConfigFile = modPath + "/python/game/realityconfig_tracker.py"
)

type source struct {
	VariableName string
	Value        string
}

type sourceConfig struct {
	Path source
	File source
}

var (
	realityConfigTrackerSource = map[ArtifactType]sourceConfig{
		ArtifactTypeBF2Demo: {
			Path: source{Value: "[MOD]/demos"},
			File: source{Value: "auto_%Y_%m_%d_%H_%M_%S.bf2demo"},
		},
		ArtifactTypePRDemo: {
			Path: source{VariableName: "PUBLIC_FOLDER"},
			File: source{VariableName: "FILE_NAME"},
		},
		ArtifactTypePRDemoPrivate: {
			Path: source{VariableName: "PRIVATE_FOLDER"},
			File: source{VariableName: "FILE_NAME"},
		},
		ArtifactTypeJSONSummary: {
			Path: source{VariableName: "JSON_FOLDER"},
			File: source{VariableName: "FILE_NAME"},
		},
	}
)

type ArtifactConfig struct {
	Path string
	File string
}

const (
	rtConfigFilePattern = `(?m)^\s*C\['%s'\]\s*=\s*'(.*?)'`
)

func (s *Server) ArtifactsConfig() (map[ArtifactType]ArtifactConfig, error) {
	configMap := make(map[ArtifactType]ArtifactConfig)

	content, err := os.ReadFile(filepath.Join(s.Path, trackerConfigFile))
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	for typ, sourceConfig := range realityConfigTrackerSource {
		config := ArtifactConfig{}
		if sourceConfig.Path.Value != "" {
			config.Path = sourceConfig.Path.Value
		} else {
			config.Path = extractWithPattern(contentStr, fmt.Sprintf(rtConfigFilePattern, sourceConfig.Path.VariableName))
		}

		if sourceConfig.File.Value != "" {
			config.File = sourceConfig.File.Value
		} else {
			config.File = extractWithPattern(contentStr, fmt.Sprintf(rtConfigFilePattern, sourceConfig.File.VariableName))
		}

		config.Path = strings.ReplaceAll(config.Path, "[MOD]", modPath)

		configMap[typ] = config
	}

	postProcessConfigMap(configMap)

	return configMap, nil
}

func postProcessConfigMap(configMap map[ArtifactType]ArtifactConfig) {
	if entry, ok := configMap[ArtifactTypePRDemo]; ok {
		entry.File = entry.File + ".PRdemo"
		configMap[ArtifactTypePRDemo] = entry
	}

	if entry, ok := configMap[ArtifactTypePRDemoPrivate]; ok {
		entry.File = entry.File + ".p.PRdemo"
		configMap[ArtifactTypePRDemoPrivate] = entry
	}

	if entry, ok := configMap[ArtifactTypeJSONSummary]; ok {
		entry.File = entry.File + ".json"
		configMap[ArtifactTypeJSONSummary] = entry
	}
}

func extractWithPattern(content string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
