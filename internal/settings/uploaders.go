package settings

import "github.com/sboon-gg/svctl/internal/game"

type Credentials struct {
	Key      string `yaml:"key"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
}

type ArtifactConfig struct {
	Subdir string `yaml:"subdir"`
}

type methodConfig struct {
	Artifacts map[game.ArtifactType]ArtifactConfig `yaml:"artifacts"`
}

type SCPConfig struct {
	methodConfig `yaml:",inline"`
	Destination  string      `yaml:"dest"`
	Credentials  Credentials `yaml:"credentials"`
}

type UploaderConfig struct {
	// SFTP          *SFTPConfig         `yaml:"sftp,omitempty"`
	SCP *SCPConfig `yaml:"scp,omitempty"`
	// HTTP          *HTTPConfig         `yaml:"http,omitempty"`
	// S3            *S3Config           `yaml:"s3,omitempty"`
}
