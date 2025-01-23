package game

import "fmt"

//go:generate go run golang.org/x/tools/cmd/stringer -type=ArtifactType -linecomment -output=artifact_type_strings.go

type ArtifactType int

const (
	ArtifactTypeBF2Demo       ArtifactType = 1 << iota // bf2demo
	ArtifactTypePRDemo                                 // prdemo
	ArtifactTypePRDemoPrivate                          // prdemoprivate
	ArtifactTypeJSONSummary                            // jsonsummary
)

func (a ArtifactType) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *ArtifactType) UnmarshalText(text []byte) error {
	v, err := ParseArtifactType(string(text))
	if err != nil {
		return err
	}
	*a = v
	return nil
}

func ParseArtifactType(s string) (ArtifactType, error) {
	var v ArtifactType
	switch s {
	case ArtifactTypeBF2Demo.String():
		v = ArtifactTypeBF2Demo
	case ArtifactTypePRDemo.String():
		v = ArtifactTypePRDemo
	case ArtifactTypePRDemoPrivate.String():
		v = ArtifactTypePRDemoPrivate
	case ArtifactTypeJSONSummary.String():
		v = ArtifactTypeJSONSummary
	default:
		return 0, fmt.Errorf("unknown value %s", s)
	}
	return v, nil
}
