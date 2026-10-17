package pkg

// Dep represents a package dependency.
type Dep struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Type    string `json:"type,omitempty"` // runtime, build, opt
}

// Package represents metadata for a native or foreign package.
type Package struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Provider    string   `json:"provider"` // waynux, fedora, debian, arch, aur, local
	Description string   `json:"description,omitempty"`
	License     string   `json:"license,omitempty"`
	URL         string   `json:"url,omitempty"`
	Checksum    string   `json:"checksum,omitempty"`
	Depends     []Dep    `json:"depends,omitempty"`
	BuildDate   int64    `json:"build_date,omitempty"`
	InstallDate int64    `json:"install_date,omitempty"`
	Files       []string `json:"files,omitempty"`
}
