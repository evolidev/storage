package disk

import (
	"github.com/evolidev/storage/fs"
	"os"
	"path/filepath"
	"strings"
)

type Local struct {
	*Common
	config LocalConfig
	prefix string
}

func NewLocal(config LocalConfig) *Local {
	disk := &Local{}
	disk.Common = NewCommon(disk)

	if config.PermModeFilePublic == 0 {
		config.PermModeFilePublic = 0644
	}

	if config.PermModeFilePrivate == 0 {
		config.PermModeFilePrivate = 0600
	}

	if config.PermModeDirectoryPublic == 0 {
		config.PermModeDirectoryPublic = 0777
	}

	if config.PermModeDirectoryPrivate == 0 {
		config.PermModeDirectoryPrivate = 0700
	}

	disk.config = config

	return disk
}

func (l *Local) Prefix(prefix string) fs.Disk {
	c := l.config
	c.Prefix = prefix

	return NewLocal(c)
}

func (l *Local) Cwd() string {
	return l.config.Prefix
}

func (l *Local) Put(file string, content []byte, visibility fs.Visibility) error {
	directories := strings.Split(file, "/")
	directories = directories[:len(directories)-1]
	dirPath := strings.Join(directories, "/")
	if dirPath != "" {
		err := l.MakeDirectory(l.getPath(dirPath), visibility)

		if err != nil {
			return err
		}
	}

	return l.write(file, content, visibility)
}

func (l *Local) Get(file string) ([]byte, error) {
	f, err := os.ReadFile(l.getPath(file))

	if err != nil {
		return nil, err
	}

	return f, nil
}

func (l *Local) Attributes(file string) fs.Attributes {
	var s, m int64
	s = 0
	m = 0
	stats, err := os.Stat(l.getPath(file))

	if err == nil {
		s = stats.Size()
		m = stats.ModTime().Unix()
	}

	return fs.Attributes{
		Size:         s,
		LastModified: m,
	}
}

func (l *Local) Exists(file string) bool {
	_, err := os.Stat(l.getPath(file))

	if err == nil {
		return true
	}

	return false
}

func (l *Local) Path(file string) string {
	p, _ := filepath.Abs(l.getPath(file))

	return p
}

func (l *Local) Delete(files ...string) error {
	for _, file := range files {
		err := os.Remove(l.getPath(file))
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Local) MakeDirectory(dir string, visibility fs.Visibility) error {
	mode := l.config.PermModeDirectoryPublic

	if visibility == fs.PRIVATE {
		mode = l.config.PermModeDirectoryPrivate
	}

	return os.MkdirAll(l.getPath(dir), mode)
}

func (l *Local) DeleteDirectory(dir string) error {
	return os.RemoveAll(l.getPath(dir))
}

func (l *Local) Files(dir string) []*fs.File {
	result := make([]*fs.File, 0)

	files := l.getFiles(dir)

	for _, v := range files {
		if !v.IsDir() {
			result = append(result, fs.NewFile(l, l.getPath(dir), v.Name()))
		}
	}

	return result
}

func (l *Local) Directories(dir string) []fs.Disk {
	result := make([]fs.Disk, 0)

	files := l.getFiles(dir)

	for _, v := range files {
		if v.IsDir() {
			//result = append(result, fs.NewDirectory(l, l.getPath(dir), v.Name()))
			result = append(result, l.Prefix(l.getPath(dir)+"/"+v.Name()))
		}
	}

	return result
}

func (l *Local) getFiles(dir string) []os.FileInfo {
	result := make([]os.FileInfo, 0)

	f, err := os.Open(l.getPath(dir))
	defer f.Close()

	if err != nil {
		return result
	}

	files, err := f.Readdir(0)
	if err != nil {

		return result
	}

	return files
}

func (l *Local) write(file string, content []byte, visibility fs.Visibility) error {
	mode := l.config.PermModeFilePublic
	if visibility == fs.PRIVATE {
		mode = l.config.PermModeFilePrivate
	}
	f, err := os.OpenFile(l.getPath(file), os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)

	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)

	return err
}

func (l *Local) getPath(path string) string {
	if l.config.Prefix == "" {
		return path
	}

	return l.config.Prefix + "/" + path
}
