package game

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"time"

	strftime "github.com/ncruces/go-strftime"
)

type ArtifactState int

const (
	ArtifactStateUnknown ArtifactState = iota
	ArtifactStateAppending
	ArtifactStateStopped
	ArtifactStateFinished
)

type ArtifactType int

const (
	ArtifactTypeChatLog ArtifactType = 1 << iota
	ArtifactTypeCoincidentIPsLog
	ArtifactTypeAdminLog
	ArtifactTypeBanLog
	ArtifactTypeTicketsLog
	ArtifactTypeJoinLog
	ArtifactTypePlayerProfilesLog
	ArtifactPlayerDataErrorsLog
	ArtifactTypePythonErrorLog
	ArtifactTypePythonLaunchErrorLog
	ArtifactTypeBF2Demo
	ArtifactTypePRDemo
	ArtifactTypePRDemoPrivate
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
	raConfigFile = "mods/pr/python/game/realityconfig_admin.py"
)

type source struct {
	VariableName string
	Value        string
}

type artifactSourceConfig struct {
	Path source
	File source
}

var (
	realityConfigAdminSource = map[ArtifactType]artifactSourceConfig{
		ArtifactTypeChatLog: {
			Path: source{
				VariableName: "log_chat_path",
			},
			File: source{
				VariableName: "log_chat_file",
			},
		},
		ArtifactTypeCoincidentIPsLog: {
			Path: source{
				VariableName: "log_IP_coincidence_path",
			},
			File: source{
				VariableName: "log_IP_coincidence_file",
			},
		},
		ArtifactTypeAdminLog: {
			Path: source{
				VariableName: "log_admins_path",
			},
			File: source{
				VariableName: "log_admins_file",
			},
		},
		ArtifactTypeBanLog: {
			Path: source{
				VariableName: "log_bans_path",
			},
			File: source{
				VariableName: "log_bans_file",
			},
		},
		ArtifactTypeTicketsLog: {
			Path: source{
				VariableName: "log_tickets_path",
			},
			File: source{
				VariableName: "log_tickets_file",
			},
		},
		ArtifactTypeJoinLog: {
			Path: source{
				Value: "admin/logs",
			},
			File: source{
				Value: "joinlog.log",
			},
		},
		ArtifactTypePlayerProfilesLog: {
			Path: source{
				Value: "admin/logs",
			},
			File: source{
				Value: "playerprofiles.log",
			},
		},
		ArtifactPlayerDataErrorsLog: {
			Path: source{
				Value: "admin/logs",
			},
			File: source{
				Value: "playerdataerrors.log",
			},
		},
		ArtifactTypePythonErrorLog: {
			Path: source{
				Value: "[MOD]/settings",
			},
			File: source{
				Value: "python_errors_v([0-9]+.[0-9]+.[0-9]+.[0-9]+).log",
			},
		},
		ArtifactTypePythonLaunchErrorLog: {
			Path: source{
				Value: "[MOD]/settings",
			},
			File: source{
				Value: "python_launch_error.log",
			},
		},
	}

	realityConfigTrackerSource = map[ArtifactType]artifactSourceConfig{
		ArtifactTypeBF2Demo: {
			Path: source{
				Value: "[MOD]/demos",
			},
			File: source{
				Value: "demo_%Y_%m_%d_%H_%M_%S_/map_/mode_/layer",
			},
		},
		ArtifactTypePRDemo: {
			Path: source{
				VariableName: "PUBLIC_FOLDER",
			},
			File: source{
				VariableName: "FILE_NAME",
			},
		},
		ArtifactTypePRDemoPrivate: {
			Path: source{
				VariableName: "PRIVATE_FOLDER",
			},
			File: source{
				VariableName: "FILE_NAME",
			},
		},
	}
)

type artifactConfig struct {
	Path string
	File string
}

const (
	raConfigFilePattern = `^\s*%s\s*=\s*"(.*?)"`
	rtConfigFilePattern = `^\s*C\['%s'\]\s*=\s*'(.*?)'`
)

func generateConfig(path string) (map[ArtifactType]artifactConfig, error) {
	config := make(map[ArtifactType]artifactConfig)

	content, err := os.ReadFile(filepath.Join(path, raConfigFile))
	if err != nil {
		return nil, err
	}

	for typ, sourceConfig := range realityConfigAdminSource {

	}

	return config, nil
}

func getArtifactsConfig(path string) (*artifactsConfig, error) {
	config := make(map[ArtifactType]artifactConfig)

	content, err := os.ReadFile(filepath.Join(path, raConfigFile))
	if err != nil {
		return nil, err
	}

	fileContent := string(content)

	// Match variables
	for _, variable := range variables {
		pattern := fmt.Sprintf(regexPattern, variable)
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(fileContent)
		if len(matches) > 2 {
			fmt.Printf("%s = %s\n", matches[1], matches[2])
		} else {
			fmt.Printf("%s not found\n", variable)
		}
	}

}

type Artifact struct {
	Name    string
	State   ArtifactState
	Type    ArtifactType
	Created time.Time
}

func (s *Server) listChatLogs() ([]Artifact, error) {
	dir := filepath.Join(s.Path, "admin/logs")

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filenameFormat := "chatlog_%Y-%m-%d_%H%Ms.txt"

	layout, err := strftime.Layout(filenameFormat)
	if err != nil {
		return nil, err
	}

	logs := make([]Artifact, 0, len(entries))
	for _, entry := range entries {
		t, err := time.Parse(layout, entry.Name())
		if err != nil {
			continue
		}

		logs = append(logs, Artifact{
			Name:    entry.Name(),
			State:   ArtifactStateFinished,
			Type:    ArtifactTypeChatLog,
			Created: t,
		})
	}

	slices.SortFunc(logs, func(a, b Artifact) int {
		return int(a.Created.Sub(b.Created))
	})

	// TODO: on running server mark last as appending

	return logs, nil
}
