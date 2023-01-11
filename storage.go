package storage

import "github.com/evolidev/storage/fs"

type Storage struct {
	disks    map[string]fs.Disk
	fallback string
}

func New(disks map[string]fs.Disk) *Storage {
	s := &Storage{disks: make(map[string]fs.Disk)}
	for name, disk := range disks {
		s.AddDisk(name, disk)
	}

	return s
}

func (s *Storage) AddDisk(name string, disk fs.Disk) {
	if len(s.disks) == 0 {
		s.fallback = name
	}
	s.disks[name] = disk
}

func (s *Storage) Disk(name string) fs.Disk {
	return s.disks[name]
}

func (s *Storage) Attributes(file string) fs.Attributes {
	return s.disk().Attributes(file)
}

func (s *Storage) Put(file string, content []byte, visibility fs.Visibility) error {
	return s.disk().File(file).Put(content, visibility)
}

func (s *Storage) Get(file string) ([]byte, error) {
	return s.disk().File(file).Get()
}

func (s *Storage) Exists(file string) bool {
	return s.disk().Exists(file)
}

func (s *Storage) Missing(file string) bool {
	return s.disk().Missing(file)
}

func (s *Storage) Size(file string) int64 {
	return s.disk().File(file).Size()
}

func (s *Storage) LastModified(file string) int64 {
	return s.disk().File(file).LastModified()
}

func (s *Storage) Path(file string) string {
	return s.disk().File(file).Path()
}

func (s *Storage) Prepend(file string, content []byte) error {
	return s.disk().File(file).Prepend(content)
}

func (s *Storage) Append(file string, content []byte) error {
	return s.disk().File(file).Append(content)
}

func (s *Storage) Copy(source string, destination string) error {
	return s.disk().File(source).Copy(destination)
}

func (s *Storage) Move(source string, destination string) error {
	return s.disk().File(source).Move(destination)
}

func (s *Storage) Delete(files ...string) error {
	return s.disk().Delete(files...)
}

func (s *Storage) MakeDirectory(dir string, visibility fs.Visibility) error {
	return s.disk().MakeDirectory(dir, visibility)
}

func (s *Storage) DeleteDirectory(dir string) error {
	return s.disk().DeleteDirectory(dir)
}

func (s *Storage) File(file string) *fs.File {
	return s.disk().File(file)
}

func (s *Storage) Files(dir string) []*fs.File {
	return s.disk().Files(dir)
}

func (s *Storage) AllFiles(dir string) []*fs.File {
	return s.disk().AllFiles(dir)
}

func (s *Storage) Directories(dir string) []fs.Disk {
	return s.disk().Directories(dir)
}

func (s *Storage) AllDirectories(dir string) []fs.Disk {
	return s.disk().AllDirectories(dir)
}

func (s *Storage) Prefix(prefix string) fs.Disk {
	return s.disk().Prefix(prefix)
}

func (s *Storage) Cwd() string {
	return s.disk().Cwd()
}

func (s *Storage) Default(name string) {
	s.fallback = name
}

func (s *Storage) disk() fs.Disk {
	return s.disks[s.fallback]
}
