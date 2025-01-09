package game

//go:generate go run golang.org/x/tools/cmd/stringer -type=ArtifactType -output=artifacts_type_strings.go

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
