package files

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
)

type SearchFiles struct {
}

func NewSearchFiles() SearchFiles {
	return SearchFiles{}
}

func (h *SearchFiles) Find(ctx context.Context, searchPath string) ([]string, error) {
	pattern, err := regexp.Compile(`\.proto$`)
	if err != nil {
		return nil, err
	}

	protoList, err := readDir(ctx, searchPath, pattern)
	if err != nil {
		return nil, err
	}

	return protoList, nil
}

func readDir(ctx context.Context, path string, regexpExpr *regexp.Regexp) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []string

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, nil
		default:
		}

		if entry.IsDir() {
			files, err := readDir(ctx, filepath.Join(path, entry.Name()), regexpExpr)
			if err != nil {
				return nil, err
			}

			result = append(result, files...)
		} else {
			if regexpExpr.MatchString(entry.Name()) {
				result = append(result, filepath.Join(path, entry.Name()))
			}
		}
	}

	return result, nil
}
