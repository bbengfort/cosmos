package pkg

import "fmt"

// Version components.
const (
	VersionMajor  = 1
	VersionMinor  = 0
	VersionPatch  = 0
	VersionSerial = 0
)

// Version returns the semvar version string.
func Version() string {
	if VersionPatch > 0 {
		return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	}
	return fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
}
