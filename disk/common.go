package disk

import (
	"github.com/evolidev/blitza/fs"
	"strings"
)

type Common struct {
	disk fs.Disk
}

func NewCommon(disk fs.Disk) *Common {
	return &Common{
		disk: disk,
	}
}

func (c *Common) Missing(file string) bool {
	return !c.disk.Exists(file)
}

func (c *Common) Prepend(file string, content []byte) error {
	oldContent, err := c.disk.Get(file)

	if err != nil {
		return err
	}

	content = append(content, oldContent...)

	return c.disk.Put(file, content, fs.PUBLIC)
}

func (c *Common) Append(file string, content []byte) error {
	oldContent, err := c.disk.Get(file)

	if err != nil {
		return err
	}

	content = append(oldContent, content...)

	return c.disk.Put(file, content, fs.PUBLIC)
}

func (c *Common) Copy(source string, destination string) error {
	content, err := c.disk.Get(source)

	if err != nil {
		return err
	}

	return c.disk.Put(destination, content, fs.PUBLIC)
}

func (c *Common) Move(source string, destination string) error {
	content, err := c.disk.Get(source)

	if err != nil {
		return err
	}

	err = c.disk.Put(destination, content, fs.PUBLIC)

	if err != nil {
		return err
	}

	return c.disk.Delete(source)
}

func (c *Common) Size(file string) int64 {
	return c.disk.Attributes(file).Size
}

func (c *Common) LastModified(file string) int64 {
	return c.disk.Attributes(file).LastModified
}

func (c *Common) AllDirectories(dir string) []fs.Disk {
	dirs := c.disk.Directories(dir)
	r := dirs

	for _, d := range dirs {
		r = append(r, c.disk.Directories(d.Cwd())...)
	}

	return r
}

func (c *Common) AllFiles(dir string) []*fs.File {
	dirs := c.AllDirectories(dir)
	r := c.disk.Files(dir)

	for _, d := range dirs {
		r = append(r, c.disk.Files(d.Cwd())...)
	}

	return r
}

func (c *Common) File(file string) *fs.File {
	parts := strings.Split(file, "/")
	name := parts[len(parts)-1]
	parts = parts[:len(parts)-1]

	return fs.NewFile(c.disk, strings.Join(parts, "/"), name)
}
