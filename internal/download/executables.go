package download

import (
	"path/filepath"
	"strings"

	"github.com/elsbrock/go-putio"
)

var executableExtensions = map[string]struct{}{
	".bat": {},
	".cmd": {},
	".exe": {},
	".jar": {},
	".lnk": {},
	".msi": {},
	".pif": {},
	".ps1": {},
	".scr": {},
	".vbs": {},
}

// firstExecutable returns the name of the first file with an executable
// extension, or "" when there is none.
func firstExecutable(files []*putio.File) string {
	for _, file := range files {
		if file == nil {
			continue
		}
		if _, ok := executableExtensions[strings.ToLower(filepath.Ext(file.Name))]; ok {
			return file.Name
		}
	}
	return ""
}
