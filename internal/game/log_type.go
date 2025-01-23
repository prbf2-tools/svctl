package game

import "fmt"

//go:generate go run golang.org/x/tools/cmd/stringer -type=LogType -linecomment -output=log_type_strings.go

type LogType int

const (
	LogTypeChat               LogType = 1 << iota // chatlog
	LogTypeCoincidentIPs                          // coincidentips
	LogTypeAdmin                                  // adminlog
	LogTypeBan                                    // banlog
	LogTypeTickets                                // tickets
	LogTypeJoin                                   // joinlog
	LogTypePlayerProfiles                         // playerprofiles
	LogTypePlayerDataErrors                       // playerdataerrors
	LogTypePythonErrors                           // pythonerrors
	LogTypePythonLaunchErrors                     // pythonlauncherrors
)

func (a LogType) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *LogType) UnmarshalText(text []byte) error {
	v, err := ParseLogType(string(text))
	if err != nil {
		return err
	}
	*a = v
	return nil
}

func ParseLogType(s string) (LogType, error) {
	var v LogType
	switch s {
	case LogTypeChat.String():
		v = LogTypeChat
	case LogTypeCoincidentIPs.String():
		v = LogTypeCoincidentIPs
	case LogTypeAdmin.String():
		v = LogTypeAdmin
	case LogTypeBan.String():
		v = LogTypeBan
	case LogTypeTickets.String():
		v = LogTypeTickets
	case LogTypeJoin.String():
		v = LogTypeJoin
	case LogTypePlayerProfiles.String():
		v = LogTypePlayerProfiles
	case LogTypePlayerDataErrors.String():
		v = LogTypePlayerDataErrors
	case LogTypePythonErrors.String():
		v = LogTypePythonErrors
	case LogTypePythonLaunchErrors.String():
		v = LogTypePythonLaunchErrors
	default:
		return 0, fmt.Errorf("unknown value %s", s)
	}
	return v, nil
}
