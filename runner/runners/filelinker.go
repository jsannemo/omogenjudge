package runners

import (
	"github.com/google/logger"

	"github.com/jsannemo/omogenjudge/util/go/files"
	"github.com/jsannemo/omogenjudge/util/go/users"
)

// FileLinker can be used to easily link in any files into two directories based on its writability.
//
// The main use case it for easily swapping out readable and writable files that should have the
// same name within a container executing a program multiple times, such as input and output files.
// It keeps all links in a small number of directories, which makes it easy to clear up the file system
// environment between runs.
type FileLinker struct {
	readBase  *files.FileBase
	writeBase *files.FileBase
}

// NewFileLinker returns a new file linker, rooted at the given path.
func NewFileLinker(dir string) (*FileLinker, error) {
	base := files.NewFileBase(dir)
	base.Gid = users.OmogenClientsID()
	if err := base.Mkdir("."); err != nil {
		return nil, err
	}
	reader, err := base.SubBase("read")
	if err != nil {
		return nil, err
	}
	writer, err := base.SubBase("write")
	if err != nil {
		return nil, err
	}
	linker := &FileLinker{
		readBase:  &reader,
		writeBase: &writer,
	}
	linker.writeBase.GroupWritable = true
	if err := linker.writeBase.Mkdir("."); err != nil {
		return nil, err
	}
	if err := linker.readBase.Mkdir("."); err != nil {
		return nil, err
	}
	return linker, nil
}

func (fl *FileLinker) base(writeable bool) *files.FileBase {
	if writeable {
		return fl.writeBase
	} else {
		return fl.readBase
	}
}

// PathFor returns the path that a file will get inside the linker.
func (fl *FileLinker) PathFor(inName string, writeable bool) string {
	str, err := fl.base(writeable).FullPath(inName)
	if err != nil {
		logger.Fatalf("Tried to use an env with relative path: %v", err)
	}
	return str
}

// LinkFile hard links the file path into the inside root.
func (fl *FileLinker) LinkFile(path, inName string, writeable bool) error {
	return fl.base(writeable).LinkInto(path, inName)
}

// Clear resets the environment for a new execution.
func (fl *FileLinker) Clear() error {
	rerr := fl.readBase.RemoveContents(".")
	werr := fl.writeBase.RemoveContents(".")
	if rerr != nil {
		return rerr
	}
	if werr != nil {
		return werr
	}
	return nil
}
