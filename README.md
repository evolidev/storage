![Build status](https://github.com/evolidev/blitza/actions/workflows/main.yml/badge.svg)
[![codecov](https://codecov.io/github/evolidev/storage/branch/main/graph/badge.svg?token=QU5LPLJSJ4)](https://codecov.io/github/evolidev/storage)

# evolidev/storage

A simple storage wrapper with some basic operations. 
If you have some basic operations like put, get, delete, copy, move etc. then this package may be something for you. 
If you need more special operations in a desired storage type then this package may not be for you.

## How to use

### Setup

```go
s3Config := disk.S3Config{Client: disk.NewMemoryClient()}
localConfig := disk.LocalConfig{}
memoryConfig := disk.MemoryConfig{}

adapters := map[string]fs.Disk{
    "memory": disk.NewMemory(memoryConfig),
    "local":  disk.NewLocal(localConfig),
    "s3":     disk.NewS3(s3Config),
}

// first added adapter is default
storage := storage.New(adapters)

// change default
storage.Default("local")
```

### Writing

To create file you can simple call the storage `put` method or create a file struct and call its `put` method. 
If the file exists the content will be overridden. 
All necessary directories will be created for you automatically.  
```go
// use storage
err := storage.Put("path/to/file.txt", []byte("hello world"), fs.PUBLIC)

// use file
file := fs.NewFile(Storage::Disk("local"), "path/to/file", "test.txt")
err := file.Put([]byte("test"), fs.PUBLIC)
```

It is possible to append or prepend content. 

```go
// use storage
err := storage.Append("path/to/file.txt", []byte("appended content"))
err := storage.Prepend("path/to/file.txt", []byte("prepended content"))

// use file
file := fs.NewFile(Storage::Disk("local"), "path/to/file", "test.txt")
err := file.Append([]byte("appended content"))
err := file.Prepend([]byte("prepended content"))
```

The copy/move commands can be used to copy/move files to a new location on the same disk
```go
// use storage
err := storage.Move("source/file.txt", "destination/file.txt")
err := storage.Copy("source/file.txt", "destination/file.txt")

// use file
file := fs.NewFile(Storage::Disk("local"), "source/path", "file.txt")
err := file.Move("destination/path")
err := file.Copy("destination/path")
```

To create a directory use `MakeDirectory`. 
Since S3 does not have directories it will be an empty object. 

```go
err := storage.MakeDirectory("directory")
```

### Retrieving

`Directories` will return a slice of fs.Disk interface. It is always the adapter of calling storage.
To get all directories including subdirectories use `AllDirectories`.
```go
// files are a slice of fs.Disk
// the type of the default disk
dirs := storage.Directories("path/to/directroy")
// given you have following structure
// | path
// |- to
// |-- directory
// |--- subdir
// |---- subsubdir
// then you will get a slice only with subdir

dirs := storage.AllDirectories("path/to/directroy")
// based on the example above you will get a slice with subdir and subsubdir

// now dirs will be a slice of *disk.Local  
dirs := storage.Disk("local").Directories("path/to/directroy")
```

For getting files in a directory use `Files`. 
To get all files including files of subdirectories use `AllFiles`
```go
// files are a slice of *fs.File
files := storage.Files("path/to/directroy")
files := storage.AllFiles("path/to/directroy")
// the result will be as in the example of directories
```

`Cwd` (current working directory) will return the path to directory
```go
dirs := storage.AllDirectories("path/to/directroy")
dir := dirs[0]
// will print "path/to/directory/subdir"
fmt.Println(dir.Cwd())
```

To get the content of file use the `Get` method
```go
// use storage
content, err := storage.Get("path/to/file.txt")

// if you have a file struct from Files() or AllFiles()
content, err := file.Get()
```

Use `Prefix` to get a sub storage of calling storage. 
The resulting storages of `Directories` will prefix the storages.
```go
// given you have following structure
// | path
// |- to
// |-- directory
// |--- subdir
// |---- subsubdir
// you can now traverse it like following
s := storage.Prefx("path")
s = s.Prefix("to")
s = s.Prefix("directory")
// will hold a slice with "subdir" as storage
dirs := s.Directories()
s = dirs[0].Prefix("subsubdir")
```

### Attributes

Checking for existence or missing files can be done with `Exists` and `Missing`
```go
if(storage.Exists("path/to/file.txt")) {
	fmt.Println("file exists")
}

if(storage.Missing("path/to/file.txt")) {
    fmt.Println("file does not exists")
}
```

There are `Size` and `LastModified` to get desired information. 
Also `Attributes` could be used to get the information. 
```go
// use helper functions
size := storage.Size("path/to/file.txt")
lastModified := storage.LastModified("path/to/file.txt")

// use attributes
size := storage.Attributes("path/to/file.txt").Size
lastModified := storage.Attributes("path/to/file.txt").LastModified

// if you have a file struct in your hand
// use helper functions
size := file.Size()
lastModified := file.LastModified()

// use attributes
size := storage.Attributes().Size
lastModified := storage.Attributes().LastModified
```

### Deleting

`Delete` will delete a single file. To delete a directory use `DeleteDirectory`. 
`DeleteDirectory` will remove everything inside! 
```go
err := storage.Delete("path/to/file.txt")
// the file also got a delete method
err := file.Delete()

err := storage.DeleteDirectory("path/to/dir")
```

## TODO

* Visibility of files in S3 adapter
* Embed adapter
* Streaming