package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileBase represents a given directory, and allows operations to be performded on files and
// directories within it. It makes sure all operations are only performed within that directory,
// providing extra protection against malicious input or erronous operations.
type FileBase struct {
	base          string
	Gid           int
	Uid           int
	GroupWritable bool
}

func NewFileBase(path string) FileBase {
	return FileBase{
		base:          filepath.Clean(path),
		Gid:           os.Getgid(),
		Uid:           os.Getuid(),
		GroupWritable: false,
	}
}

func (fb FileBase) Path() string {
	return fb.base
}

func (fb FileBase) SubBase(path string) (FileBase, error) {
	nbase, err := fb.FullPath(path)
	if err != nil {
		return FileBase{}, nil
	}
	fb.base = nbase
	return fb, nil
}

func (fb FileBase) FullPath(subPath string) (string, error) {
	fullPath := filepath.Clean(filepath.Join(fb.base, subPath))
	baseComponents := strings.Split(string(filepath.Separator), fb.base)
	newComponents := strings.Split(string(filepath.Separator), fullPath)
	if len(newComponents) < len(baseComponents) {
		return "", fmt.Errorf("path %s traverses upwards (too few components)", subPath)
	}
	for i := 0; i < len(baseComponents); i++ {
		if baseComponents[i] != newComponents[i] {
			return "", fmt.Errorf("path %s traverses upwards from %s (%s != %s)", fullPath, fb.base, baseComponents[i], newComponents[i])
		}
	}
	return fullPath, nil
}

func (fb FileBase) ReadFile(subPath string) ([]byte, error) {
	npath, err := fb.FullPath(subPath)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(npath)
}

func (fb FileBase) LinkInto(fromPath string, toSubPath string) error {
	npath, err := fb.FullPath(toSubPath)
	if err != nil {
		return err
	}
	return os.Link(fromPath, npath)
}

func (fb FileBase) RemoveContents(subPath string) error {
	npath, err := fb.FullPath(subPath)
	if err != nil {
		return err
	}
	dir, err := ioutil.ReadDir(npath)
	if err != nil {
		return err
	}
	for _, d := range dir {
		if err := os.RemoveAll(filepath.Join(npath, d.Name())); err != nil {
			return err
		}
	}
	return nil
}

func (fb FileBase) Exists(subPath string) (bool, error) {
	npath, err := fb.FullPath(subPath)
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(npath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		panic(err)
	}
	return true, nil
}

func (fb FileBase) Copy(srcFile, dstFile string) error {
	npath, err := fb.FullPath(dstFile)
	if err != nil {
		return err
	}
	out, err := os.Create(npath)
	defer out.Close()
	if err != nil {
		return err
	}

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func (fb FileBase) FixOwners(subPath string) error {
	npath, err := fb.FullPath(subPath)
	if err != nil {
		return err
	}
	return os.Chown(npath, fb.Uid, fb.Gid)
}

func (fb FileBase) setMode(subPath string, mode os.FileMode) error {
	npath, err := fb.FullPath(subPath)
	if err != nil {
		return err
	}
	return os.Chmod(npath, mode)
}

func (fb FileBase) SetModeExec(subPath string) error {
	mode := 0750
	if fb.GroupWritable {
		mode = mode | 0020
	}
	return fb.setMode(subPath, os.FileMode(mode))
}

func (fb FileBase) SetMode(subPath string) error {
	mode := 0640
	if fb.GroupWritable {
		mode = mode | 0020
	}
	return fb.setMode(subPath, os.FileMode(mode))
}

func (fb FileBase) Mkdir(subPath string) error {
	npath, err := fb.FullPath(subPath)
	if err != nil {
		return err
	}
	if err := os.Mkdir(npath, 0750); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	if err := fb.FixOwners(subPath); err != nil {
		return err
	}
	if err := fb.SetModeExec(subPath); err != nil {
		return err
	}
	return nil
}

func (fb FileBase) WriteFile(subPath string, contents []byte) error {
	npath, err := fb.FullPath(subPath)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(npath, contents, 0750); err != nil {
		return err
	}
	if err := fb.FixOwners(subPath); err != nil {
		return err
	}
	if err := fb.SetMode(subPath); err != nil {
		return err
	}
	return nil
}

func (fb FileBase) CopyInto(srcDir string) error {
	entries, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(srcDir, entry.Name())
		if entry.IsDir() {
			if err := fb.Mkdir(entry.Name()); err != nil {
				return err
			}
			subfb, err := fb.SubBase(entry.Name())
			if err != nil {
				return err
			}
			if err := subfb.CopyInto(filepath.Join(srcDir, entry.Name())); err != nil {
				return err
			}
		} else {
			if err := fb.Copy(sourcePath, entry.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}
