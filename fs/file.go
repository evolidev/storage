package fs

type File struct {
	storage Disk
	path    string
	name    string
}

func (f *File) Write(p []byte) (n int, err error) {
	return len(p), f.Put(p, PUBLIC)
}

func NewFile(store Disk, path string, name string) *File {
	return &File{
		storage: store,
		path:    path,
		name:    name,
	}
}

func (f *File) Put(content []byte, visibility Visibility) error {
	return f.storage.Put(f.fullName(), content, visibility)
}

func (f *File) Get() ([]byte, error) {
	return f.storage.Get(f.fullName())
}

func (f *File) Delete() error {
	return f.storage.Delete(f.fullName())
}

func (f *File) Size() int64 {
	return f.storage.Size(f.fullName())
}

func (f *File) LastModified() int64 {
	return f.storage.LastModified(f.fullName())
}

func (f *File) Path() string {
	return f.storage.Path(f.fullName())
}

func (f *File) Prepend(content []byte) error {
	return f.storage.Prepend(f.fullName(), content)
}

func (f *File) Append(content []byte) error {
	return f.storage.Append(f.fullName(), content)
}

func (f *File) Copy(target string) error {
	return f.storage.Copy(f.fullName(), target)
}

func (f *File) Move(target string) error {
	return f.storage.Move(f.fullName(), target)
}

func (f *File) Name() string {
	return f.name
}

func (f *File) fullName() string {
	if f.path == "" {
		return f.Name()
	}

	return f.path + "/" + f.Name()
}
