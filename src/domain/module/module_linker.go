package module

import (
	"path"
	"strings"
)

type PackageNameGenerator struct {
}

func NewPackageNameGenerator() PackageNameGenerator {
	return PackageNameGenerator{}
}

func (h *PackageNameGenerator) Generate(
	files []string,
	filesRootPath string,
	moduleName string,
) []string {
	names := make([]string, len(files), len(files))

	for i, file := range files {
		dir := path.Dir(file)
		dir = strings.Replace(dir, filesRootPath, "", 1)
		dir = strings.Trim(dir, "/")
		names[i] = moduleName + "/" + dir
	}

	return names
}
