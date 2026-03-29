package usecase

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	domain "github.com/rm-ryou/gotor/internal/core/domain/explorer"
)

type Explorer struct {
	fs   domain.FSReader
	fsw  domain.FSWriter
	tree *domain.Tree

	OnFileSelected func(path string) error
}

func NewExplorer(fs domain.FSReader, fsw domain.FSWriter, rootPath string) (*Explorer, error) {
	absRoot, err := resolveRoot(rootPath)
	if err != nil {
		return nil, NewError("Failed to open the workspace.", err)
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
		return nil, NewError("Failed to load the workspace.", err)
	}
	root.Children = children

	return &Explorer{
		fs:   fs,
		fsw:  fsw,
		tree: domain.New(root),
	}, nil
}

func (e *Explorer) Tree() *domain.Tree {
	return e.tree
}

func (e *Explorer) ToggleNode(node *domain.Node) error {
	if !node.IsDir {
		return nil
	}

	var children []*domain.Node

	if !node.Expanded && len(node.Children) == 0 {
		loaded, err := e.fs.ReadDir(node.Path, node.Depth+1)
		if err != nil {
			return NewError("Failed to load the folder.", err)
		}
		children = loaded
	}

	e.tree.Toggle(node, children)
	return nil
}

func (e *Explorer) SelectFile(node *domain.Node) error {
	if node.IsDir {
		return nil
	}
	e.tree.Select(node)
	if e.OnFileSelected != nil {
		if err := e.OnFileSelected(node.Path); err != nil {
			e.tree.ClearSelection()
			return NewError("Failed to open the selected file.", err)
		}
	}

	return nil
}

func (e *Explorer) CreateFile(dirPath, name string) error {
	path := filepath.Join(dirPath, name)
	if err := e.fsw.CreateFile(path); err != nil {
		return NewError("Failed to create the file.", err)
	}

	parent := e.tree.FindNode(dirPath)
	if parent != nil && parent.Expanded {
		children, err := e.fs.ReadDir(dirPath, parent.Depth+1)
		if err != nil {
			return NewError("Failed to reload the folder.", err)
		}
		e.tree.Refresh(parent, children)
	}

	return nil
}

func (e *Explorer) ClearSelection() {
	e.tree.ClearSelection()
}

func (e *Explorer) ChangeRoot(path string) error {
	absPath, err := resolveRoot(path)
	if err != nil {
		return NewError("Failed to change the workspace.", err)
	}

	root := &domain.Node{
		Name:     filepath.Base(absPath),
		Path:     absPath,
		IsDir:    true,
		Depth:    0,
		Expanded: true,
	}
	children, err := e.fs.ReadDir(absPath, 1)
	if err != nil {
		return NewError("Failed to load the workspace.", err)
	}
	root.Children = children

	e.tree.Reset(root)
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
		return "", errors.New("path is not a directory")
	}

	return abs, nil
}
