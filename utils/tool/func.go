package tool

import (
	"io"
	"os"
	"strings"
)

// HasPrefix returns true if name has any string in given slice as prefix.
func HasPrefix(name string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

// IsEntry returns true if name equals to any string in given slice.
func IsEntry(name string, entries []string) bool {
	for _, e := range entries {
		if e == name {
			return true
		}
	}
	return false
}

// IsFilter returns true if given name matches any of global filter rule.
func IsFilter(name string) bool {
	if strings.Contains(name, ".DS_Store") {
		return true
	}
	return false
}

// IsExist returns true if given path is a file or directory.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// Copy copies file from source to target path.
func Copy(dest, src string) error {
	// Gather file information to set back later.
	si, err := os.Lstat(src)
	if err != nil {
		return err
	}

	// Handle symbolic link.
	if si.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		// NOTE: os.Chmod and os.Chtimes don't recoganize symbolic link,
		// which will lead "no such file or directory" error.
		return os.Symlink(target, dest)
	}

	sr, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sr.Close()

	dw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer dw.Close()

	if _, err = io.Copy(dw, sr); err != nil {
		return err
	}

	// Set back file information.
	if err = os.Chtimes(dest, si.ModTime(), si.ModTime()); err != nil {
		return err
	}
	return os.Chmod(dest, si.Mode())
}
