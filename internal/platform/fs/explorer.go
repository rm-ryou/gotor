package fs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	domain "github.com/rm-ryou/gotor/internal/core/domain/explorer"
)

type ExplorerReader struct {
	showHidden bool
}

func New(showHidden bool) *ExplorerReader {
	return &ExplorerReader{
		showHidden: showHidden,
	}
}

func (er *ExplorerReader) ReadDir(path string, depth int) ([]*domain.Node, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	nodes := make([]*domain.Node, 0, len(entries))
	for _, entry := range entries {
		if !er.showHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		nodes = append(nodes, &domain.Node{
			Name:  entry.Name(),
			Path:  filepath.Join(path, entry.Name()),
			IsDir: entry.IsDir(),
			Depth: depth,
		})
	}

	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].IsDir != nodes[j].IsDir {
			return nodes[i].IsDir
		}
		return strings.ToLower(nodes[i].Name) < strings.ToLower(nodes[j].Name)
	})

	return nodes, nil
}
