package fsutil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// OpenZip reads the whole zip archive at path into memory.
func OpenZip(path string) (*zip.Reader, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
}

// CreateZip creates a zip archive at path.
func CreateZip(path string) (*os.File, *zip.Writer, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}

	return file, zip.NewWriter(file), nil
}

// OpenPack opens a resource pack stored either as a zip archive (.zip, .mcpack, ...) or as a folder.
func OpenPack(path string) (fs.FS, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return os.DirFS(path), nil
	}
	return OpenZip(path)
}

// PackName derives a pack name from its path: a folder's name, or an archive's file name without its
// extension.
func PackName(path string) string {
	name := filepath.Base(filepath.Clean(path))
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return name
	}
	if trimmed := strings.TrimSuffix(name, filepath.Ext(name)); trimmed != "" {
		return trimmed
	}
	return name
}

// FindPackRoot returns the shallowest directory of f that contains marker, such as "pack.mcmeta" or
// "manifest.json". Packs zipped with their enclosing folder keep their files one level down.
func FindPackRoot(f fs.FS, marker string) (root string, found bool) {
	depth := -1
	_ = fs.WalkDir(f, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || d.Name() != marker {
			return nil
		}
		dir := filepath.ToSlash(filepath.Dir(p))
		if n := strings.Count(dir, "/"); depth == -1 || n < depth || (n == depth && dir == ".") {
			root, depth, found = dir, n, true
		}
		return nil
	})
	return root, found
}

// Join joins slash-separated fs.FS path elements, treating "." as the root.
func Join(elem ...string) string {
	var parts []string
	for _, e := range elem {
		if e != "" && e != "." {
			parts = append(parts, e)
		}
	}
	if len(parts) == 0 {
		return "."
	}
	return strings.Join(parts, "/")
}

// IsDir reports whether p is a directory in f.
func IsDir(f fs.FS, p string) bool {
	info, err := fs.Stat(f, p)
	return err == nil && info.IsDir()
}

// FindDirectory finds the shallowest directory called name in a fs.FS.
func FindDirectory(f fs.FS, name string) (dirPath string, found bool) {
	depth := -1
	walkDirFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.Name() == name && d.IsDir() {
			if n := strings.Count(path, "/"); depth == -1 || n < depth {
				dirPath, depth, found = path, n, true
			}
		}
		return nil
	}

	_ = fs.WalkDir(f, ".", walkDirFunc)
	return
}

// Unzip extracts the zip archive at src into dest.
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
