package explorer

import (
	"fmt"
	"os"
	"path/filepath"

	domain "github.com/rm-ryou/gotor/internal/core/domain/explorer"
)

type UseCase struct {
	fs   domain.FSReader
	tree *domain.Tree

	OnFileSelected func(path string)
}

func New(fs domain.FSReader, rootPath string) (*UseCase, error) {
	absRoot, err := resolveRoot(rootPath)
	if err != nil {
		return nil, fmt.Errorf("explorer usecase: resolve root: %w", err)
	}

	root := &domain.Node{
		Name:     filepath.Base(absRoot),
		Path:     absRoot,
		IsDir:    true,
		Depth:    0,
		Expanded: true,
	}

	children, err := fs.ReadDir(absRoot, 1)
	if err != nil {
		return nil, fmt.Errorf("explorer usecase: initial load %s: %w", absRoot, err)
	}
	root.Children = children

	return &UseCase{
		fs:   fs,
		tree: domain.New(root),
	}, nil
}

func (uc *UseCase) Tree() *domain.Tree {
	return uc.tree
}

func (uc *UseCase) ToggleNode(node *domain.Node) error {
	if !node.IsDir {
		return nil
	}

	var children []*domain.Node

	if !node.Expanded && len(node.Children) == 0 {
		loaded, err := uc.fs.ReadDir(node.Path, node.Depth+1)
		if err != nil {
			return fmt.Errorf("explorer usecase: read dir %s: %w", node.Path, err)
		}
		children = loaded
	}

	uc.tree.Toggle(node, children)
	return nil
}

func (uc *UseCase) SelectFile(node *domain.Node) {
	if node.IsDir {
		return
	}
	uc.tree.Select(node)
	if uc.OnFileSelected != nil {
		uc.OnFileSelected(node.Path)
	}
}

func (uc *UseCase) ClearSelection() {
	uc.tree.ClearSelection()
}

func (uc *UseCase) ChangeRoot(path string) error {
	absPath, err := resolveRoot(path)
	if err != nil {
		return fmt.Errorf("explorer usecase: resolve root: %w", err)
	}

	root := &domain.Node{
		Name:     filepath.Base(absPath),
		Path:     absPath,
		IsDir:    true,
		Depth:    0,
		Expanded: true,
	}
	children, err := uc.fs.ReadDir(absPath, 1)
	if err != nil {
		return fmt.Errorf("explorer usecase: read dir %s: %w", absPath, err)
	}
	root.Children = children

	uc.tree.Reset(root)
	return nil
}

func resolveRoot(path string) (string, error) {
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get working directory: %w", err)
		}
		return cwd, nil
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", abs)
	}

	return abs, nil
}
