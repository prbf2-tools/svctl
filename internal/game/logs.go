package game

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	raConfigFilePattern = `(?m)^\s*%s\s*=\s*"(.*?)"`
)

var realityConfigAdminSource = map[LogType]sourceConfig{
	LogTypeChat: {
		Path: source{VariableName: "log_chat_path"},
		File: source{VariableName: "log_chat_file"},
	},
	LogTypeCoincidentIPs: {
		Path: source{VariableName: "log_IP_coincidence_path"},
		File: source{VariableName: "log_IP_coincidence_file"},
	},
	LogTypeAdmin: {
		Path: source{VariableName: "log_admins_path"},
		File: source{VariableName: "log_admins_file"},
	},
	LogTypeBan: {
		Path: source{VariableName: "log_bans_path"},
		File: source{VariableName: "log_bans_file"},
	},
	LogTypeTickets: {
		Path: source{VariableName: "log_tickets_path"},
		File: source{VariableName: "log_tickets_file"},
	},
	LogTypeJoin: {
		Path: source{Value: "admin/logs"},
		File: source{Value: "joinlog.log"},
	},
	LogTypePlayerProfiles: {
		Path: source{Value: "admin/logs"},
		File: source{Value: "playerprofiles.log"},
	},
	LogTypePlayerDataErrors: {
		Path: source{Value: "admin/logs"},
		File: source{Value: "playerdataerrors.log"},
	},
	LogTypePythonErrors: {
		Path: source{Value: "[MOD]/settings"},
		File: source{Value: "python_errors_v([0-9]+.[0-9]+.[0-9]+.[0-9]+).log"},
	},
	LogTypePythonLaunchErrors: {
		Path: source{Value: "[MOD]/settings"},
		File: source{Value: "python_launch_error.log"},
	},
}

func (s *Server) LogsConfig() (map[LogType]ArtifactConfig, error) {
	configMap := make(map[LogType]ArtifactConfig)

	content, err := os.ReadFile(filepath.Join(s.Path, raConfigFile))
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	for typ, sourceConfig := range realityConfigAdminSource {
		config := ArtifactConfig{}
		if sourceConfig.Path.Value != "" {
			config.Path = sourceConfig.Path.Value
		} else {
			config.Path = extractWithPattern(contentStr, fmt.Sprintf(raConfigFilePattern, sourceConfig.Path.VariableName))
		}

		if sourceConfig.File.Value != "" {
			config.File = sourceConfig.File.Value
		} else {
			config.File = extractWithPattern(contentStr, fmt.Sprintf(raConfigFilePattern, sourceConfig.File.VariableName))
		}

		config.Path = strings.ReplaceAll(config.Path, "[MOD]", modPath)

		configMap[typ] = config
	}

	return configMap, nil
}
