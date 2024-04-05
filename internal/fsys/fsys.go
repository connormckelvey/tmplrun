package fsys

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
)

type LimitedFS struct {
	paths       map[string]bool
	sortedPaths []string
	f           fs.FS
}

type LimitedFSOption func(*LimitedFS) error

func WithPaths(paths []string) LimitedFSOption {
	return func(lf *LimitedFS) error {
		lf.paths = make(map[string]bool)
		for _, path := range paths {
			lf.paths[filepath.Clean(path)] = true
		}
		lf.sortedPaths = make([]string, 0)
		for path := range lf.paths {
			lf.sortedPaths = append(lf.sortedPaths, path)
		}
		slices.Sort(lf.sortedPaths)
		return nil
	}
}

func WithGlobs(globs []string) LimitedFSOption {
	return func(lf *LimitedFS) error {
		var expanded []string
		for _, glob := range globs {
			matches, err := fs.Glob(lf.f, glob)
			if err != nil {
				return err
			}
			expanded = append(expanded, matches...)
		}
		return WithPaths(expanded)(lf)
	}
}

func NewLimitedFS(fsys fs.FS, opts ...LimitedFSOption) (*LimitedFS, error) {
	lf := &LimitedFS{
		paths: map[string]bool{},
		f:     fsys,
	}
	for _, apply := range opts {
		if err := apply(lf); err != nil {
			return nil, err
		}
	}
	return lf, nil
}

func (lf *LimitedFS) validate(name string) error {
	if _, ok := lf.paths[filepath.Clean(name)]; !ok {
		return &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fs.ErrNotExist,
		}
	}
	return nil
}

func (lf *LimitedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if err := lf.validate(name); err != nil {
		return nil, err
	}
	origEntries, err := fs.ReadDir(lf.f, name)
	if err != nil {
		return nil, err
	}
	entries := make([]fs.DirEntry, 0)
	for _, entry := range origEntries {
		if err := lf.validate(filepath.Join(name, entry.Name())); err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (lf *LimitedFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fs.ErrInvalid,
		}
	}
	if err := lf.validate(name); err != nil {
		return nil, fmt.Errorf("invalid filename: %w", err)
	}
	return lf.f.Open(name)
}
