package version

import "fmt"

var (
	// Name of the service
	Name = "dasoft"
	// APIVersion of the API this service provide
	APIVersion = "v1"
	// Version is the current echo version in SemVer format
	SemVer = "v0.0.1"
	// BuildDate represents the date this service was built
	BuildDate string
	// GitCommit is the hash of the git commit
	GitCommit string
)

type Version struct {
	APIVersion string `json:"api_version"`
	Version    string `json:"version"`
	GitCommit  string `json:"git_commit"`
	BuildDate  string `json:"build_date"`
}

func Get() *Version {
	v := &Version{
		APIVersion: APIVersion,
		Version:    SemVer,
		GitCommit:  GitCommit,
		BuildDate:  BuildDate,
	}
	return v
}

func (v *Version) String() string {
	return fmt.Sprintf("%s (%s) %s", SemVer, GitCommit, BuildDate)
}
