package fs

//'file' => [
//'public' => 0644,
//'private' => 0600,
//],
//'dir' => [
//'public' => 0755,
//'private' => 0700,
//],

type Visibility uint32

const (
	PUBLIC  Visibility = 0755
	PRIVATE            = 0600
)

type BaseOperation interface {
	Put(file string, content []byte, visibility Visibility) error
	Get(file string) ([]byte, error)
	Attributes(file string) Attributes
	Exists(file string) bool
	Delete(files ...string) error
	Directories(dir string) []Disk
	Files(dir string) []*File
}

type Disk interface {
	BaseOperation
	Missing(file string) bool
	Size(file string) int64
	LastModified(file string) int64
	Path(file string) string
	Prepend(file string, content []byte) error
	Append(file string, content []byte) error
	Copy(source string, destination string) error
	Move(source string, destination string) error
	MakeDirectory(dir string, visibility Visibility) error
	DeleteDirectory(dir string) error
	AllFiles(dir string) []*File
	AllDirectories(dir string) []Disk
	File(file string) *File
	Prefix(prefix string) Disk
	Cwd() string
}
