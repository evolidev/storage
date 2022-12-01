package storage

import (
	"github.com/evolidev/blitza/disk"
	"github.com/evolidev/blitza/fs"
	"strings"
	"testing"
)

func TestDisk(t *testing.T) {
	t.Parallel()
	d := getStorage().Disk("local")

	if _, ok := d.(*disk.Local); ok == false {
		t.Errorf("Wrong disk returned")
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/put should crate file", func(t *testing.T) {
			file := "test_create.txt"
			createFile(t, storage, file)
			defer clear(storage, file)

			if !storage.Exists(file) {
				t.Errorf("File not created %s", file)
			}
		})

		t.Run(name+"/put with parents directories should also create all directories", func(t *testing.T) {
			file := "tmp/test/test_create.txt"
			createFile(t, storage, file)
			defer clear(storage, file)

			if !storage.Exists(file) {
				t.Errorf("File not created %s", file)
			}
		})
	}
}

func TestFile(t *testing.T) {
	t.Parallel()
	storage := getStorage()
	file := storage.File("test_file.txt")

	if file.Name() != "test_file.txt" {
		t.Errorf("file not created")
	}
}

func TestGet(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/get should get file content", func(t *testing.T) {
			base := "get"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			content, err := storage.Get(file)

			check(t, err, "Failed to read file %s", file)
			if string(content) != "test" {
				t.Errorf("Content %s does not match %s", content, "test")
			}
		})
	}
}

func TestExists(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/exists should return true if file exists", func(t *testing.T) {
			base := "existing"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			result := storage.Exists(file)

			if !result {
				t.Errorf("File does not exists %s", file)
			}
		})

		t.Run(name+"/exists should return true file exists even if it is in sub directories", func(t *testing.T) {
			file := "existing/deep/test_existing.txt"
			createFile(t, storage, file)
			defer clear(storage, file)

			result := storage.Exists(file)

			if !result {
				t.Errorf("File does not exists %s", file)
			}
		})

		t.Run(name+"/exists should return false file does not exists even if it is in sub directories", func(t *testing.T) {
			fileExisting := "existing/deep/test_existing.txt"
			createFile(t, storage, fileExisting)
			defer clear(storage, fileExisting)
			file := "existing/deep/not_existing.txt"

			result := storage.Exists(file)

			if result {
				t.Errorf("File does not exists but function returns true. %s", file)
			}
		})
	}
}

func TestMissing(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/get should get file content", func(t *testing.T) {
			file := "not-existing.txt"

			result := storage.Missing(file)

			if !result {
				t.Errorf("File does not exists %s", file)
			}
		})
	}
}

func TestMetadata(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/get size information", func(t *testing.T) {
			base := "size"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			s := storage.Size(file)

			exp := int64(len([]byte("test")))
			if s != exp {
				t.Errorf("%d does not match expected %d", s, exp)
			}
		})
		t.Run(name+"/get attributes", func(t *testing.T) {
			base := "attributes"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			a := storage.Attributes(file)

			exp := int64(len([]byte("test")))
			if a.Size != exp {
				t.Errorf("%d does not match expected %d", a.Size, exp)
			}
		})
		t.Run(name+"/get modified information", func(t *testing.T) {
			base := "modified"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			s := storage.LastModified(file)

			if s == 0 {
				t.Errorf("Last modified not set")
			}
		})

		t.Run(name+"/get modified information should return an error if file not exists", func(t *testing.T) {
			base := "modified"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/file_not_exists.txt"

			s := storage.LastModified(file)

			if s != 0 {
				t.Errorf("Last modified is set")
			}
		})
	}
}

func TestPath(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)

		t.Run(name+"/path should have the right paht", func(t *testing.T) {
			base := "prepend"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			p := storage.Path(file)

			if !strings.HasSuffix(p, file) {
				t.Errorf("%s path is not valid", p)
			}
		})

	}
}

func TestPrepend(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/prepend should prepend content", func(t *testing.T) {
			base := "prepend"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			err := storage.Prepend(file, []byte("prepend "))

			check(t, err, "Failed to write data %s", err)
			fileContent, err := storage.Get(file)
			check(t, err, "Failed to read file %s with error %s", file, err)
			if string(fileContent) != "prepend test" {
				t.Errorf("Content '%s' does not match expected '%s'", fileContent, "prepend test")
			}
		})

		t.Run(name+"/prepend on none existing file should return error", func(t *testing.T) {
			base := "prepend"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/file_not_exists.txt"

			err := storage.Prepend(file, []byte("prepend "))

			if err == nil {
				t.Errorf("Writing to none existing file")
			}
		})
	}
}

func TestAppend(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/append should append content", func(t *testing.T) {
			base := "append"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			err := storage.Append(file, []byte(" append"))

			check(t, err, "Failed to write data %s", err)
			fileContent, err := storage.Get(file)
			check(t, err, "Failed to read file %s with error %s", file, err)
			if string(fileContent) != "test append" {
				t.Errorf("Content '%s' does not match expected '%s'", fileContent, "test append")
			}
		})

		t.Run(name+"/append on none existing file should return error", func(t *testing.T) {
			base := "append"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/file_not_exists.txt"

			err := storage.Append(file, []byte("append "))

			if err == nil {
				t.Errorf("Writing to none existing file")
			}
		})
	}
}

func TestCopy(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/copy should copy file", func(t *testing.T) {
			base := "copy"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"
			targetFile := base + "/sub/test_copy.txt"

			err := storage.Copy(file, targetFile)

			check(t, err, "Failed to copy file %s, with error: ", file)
			if !storage.Exists(targetFile) {
				t.Errorf("File got not copied")
			}
			content, err := storage.Get(file)
			check(t, err, "Could not get file %s", file)
			copiedContent, err := storage.Get(targetFile)
			check(t, err, "Could not get file %s", targetFile)
			if string(content) != string(copiedContent) {
				t.Errorf("Contents does not match")
			}
		})

		t.Run(name+"/copy on none existing file should return error", func(t *testing.T) {
			base := "copy"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/file_not_exists.txt"

			err := storage.Copy(file, "tmp")

			if err == nil {
				t.Errorf("Writing to none existing file")
			}
		})
	}
}

func TestMove(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/move should move file", func(t *testing.T) {
			base := "move"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"
			targetFile := base + "/sub/test_move_to.txt"

			err := storage.Move(file, targetFile)

			check(t, err, "Failed to move file %s", file)
			if !storage.Exists(targetFile) {
				t.Errorf("File got not copied")
			}
			if storage.Exists(file) {
				t.Errorf("File still exists %s", file)
			}
		})

		t.Run(name+"/move on none existing file should return error", func(t *testing.T) {
			base := "move"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/file_not_exists.txt"

			err := storage.Move(file, "tmp_move")

			if err == nil {
				t.Errorf("Writing to none existing file")
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/delete single file", func(t *testing.T) {
			base := "delete_single"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"

			err := storage.Delete(file)

			check(t, err, "Failed to delete file %s", file)
			if storage.Exists(file) {
				t.Errorf("File still exists %s", file)
			}
		})
		t.Run(name+"/delete multiple files", func(t *testing.T) {
			base := "delete_multiple"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/files.txt"
			file2 := base + "/sub2/files.txt"

			err := storage.Delete(file, file2)

			check(t, err, "Failed to delete file %s", file)
			if storage.Exists(file) {
				t.Errorf("File still exists %s", file)
			}
			if storage.Exists(file2) {
				t.Errorf("File still exists %s", file)
			}
		})
		t.Run(name+"/delete none existing file", func(t *testing.T) {
			base := "delete_single_not_existing"
			d := setup(t, storage, base)
			defer d()
			file := base + "/sub/file_not_exists.txt"

			err := storage.Delete(file)

			if err == nil {

			}
		})
	}
}

func TestMakeDirectory(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/make single directory", func(t *testing.T) {
			dir := "test_make"
			err := storage.MakeDirectory(dir, fs.PUBLIC)
			defer clearDir(storage, dir)

			check(t, err, "Failed to create dir %s", dir)

			if !storage.Exists(dir) {
				t.Errorf("Dir does not exists %s", dir)
			}
		})
	}
}

func TestDeleteDirectory(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/delete empty directory should work", func(t *testing.T) {
			dir := "test_delete"
			err := storage.MakeDirectory(dir, fs.PUBLIC)
			defer clearDir(storage, dir)
			check(t, err, "Failed to create dir %s", dir)
			if !storage.Exists(dir) {
				t.Errorf("Dir does not exists %s", dir)
			}

			err = storage.DeleteDirectory(dir)

			check(t, err, "Failed to delete dir %s", dir)
			if storage.Exists(dir) {
				t.Errorf("Dir still exists %s", dir)
			}
		})

		t.Run(name+"/delete should delete all directories and sub directories", func(t *testing.T) {
			dir := "test_delete_multiple/test/tmp"
			err := storage.MakeDirectory(dir, fs.PUBLIC)
			check(t, err, "Could not create dir %s", dir)
			if !storage.Exists(dir) {
				t.Errorf("Could not create dir %s", dir)
			}

			err = storage.DeleteDirectory("test_delete_multiple")

			check(t, err, "could not delete dir %s", "test_delete_multiple")
			if storage.Exists("test_delete_multiple/test/tmp") {

				t.Errorf("Directory still exists %s", "test_delete_multiple/test/tmp")
			}

			if storage.Exists("test_delete_multiple/test") {
				t.Errorf("Directory still exists %s", "test_delete_multiple/test")
			}

			if storage.Exists("test_delete_multiple") {
				t.Errorf("Directory still exists %s", "test_delete_multiple")
			}
		})
	}
}

func TestFiles(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/return all files of directory", func(t *testing.T) {
			d := setup(t, storage, "list")
			defer d()

			files := storage.Files("list")

			if len(files) != 1 {
				t.Errorf("Expected count of %d does not match current %d", 1, len(files))
			}
		})

		t.Run(name+"/return all files of directory but not of sub directory", func(t *testing.T) {
			d := setup(t, storage, "list_all_with_sub_without_sub_sub")
			defer d()

			files := storage.Files("list_all_with_sub_without_sub_sub")

			if len(files) != 1 {
				t.Errorf("Expected count of %d does not match current %d", 1, len(files))
			}
		})

		t.Run(name+"/return all files of directory and sub directories", func(t *testing.T) {
			base := "list_all_with_sub"
			d := setup(t, storage, base)
			defer d()

			files := storage.AllFiles(base)

			if len(files) != 4 {
				t.Errorf("Expected count of %d does not match current %d", 4, len(files))
			}
		})

		t.Run(name+"/return files of not existing directory should return 0", func(t *testing.T) {
			base := "list_all_with_sub"
			d := setup(t, storage, base)
			defer d()

			files := storage.AllFiles("list_all_with_sub_not_existing")

			if len(files) != 0 {
				t.Errorf("Expected count of %d does not match current %d", 4, len(files))
			}
		})

		t.Run(name+"/return files on a file should return 0", func(t *testing.T) {
			base := "list_all_with_sub_as_file"
			d := setup(t, storage, base)
			defer d()

			files := storage.AllFiles("list_all_with_sub_as_file/files.txt")

			if len(files) != 0 {
				t.Errorf("Expected count of %d does not match current %d", 4, len(files))
			}
		})
	}
}

func TestDirectories(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/directories should return only the direct child directories", func(t *testing.T) {
			d := setup(t, storage, "list_dir")
			defer d()

			files := storage.Directories("list_dir")

			if len(files) != 3 {
				t.Errorf("Expected count of %d does not match current %d", 3, len(files))
			}
		})

		t.Run(name+"/AllDirectories should return all sub directories", func(t *testing.T) {
			d := setup(t, storage, "list_dir_with_sub")
			defer d()

			files := storage.AllDirectories("list_dir_with_sub")

			if len(files) != 4 {
				t.Errorf("Expected count of %d does not match current %d", 4, len(files))
			}
		})
	}
}

func TestPrefix(t *testing.T) {
	t.Parallel()
	testStorage := getStorage()
	for name, _ := range testStorage.disks {
		storage := testStorage
		storage.Default(name)
		t.Run(name+"/directories should return only the direct child directories", func(t *testing.T) {
			d := setup(t, storage, "prefix")
			defer d()

			prefixed := storage.Prefix("prefix")
			if !prefixed.Exists("files.txt") {
				t.Errorf("failed to get file over prefix")
			}

			if prefixed.Cwd() != "prefix" {
				t.Errorf("wrong current working dir")
			}

			if storage.Cwd() != "" {
				t.Errorf("original config got adjusted")
			}
		})
	}
}

func createFile(t *testing.T, storage fs.Disk, path string) {
	t.Helper()

	err := storage.Put(path, []byte("test"), fs.PUBLIC)

	check(t, err, "Could not write content to file %s", path)
}

func check(t *testing.T, err error, message string, args ...any) {
	t.Helper()

	args = append(args, err)
	if err != nil {
		t.Errorf(message+", with error: ", args...)
	}
}

func clear(storage fs.Disk, file string) {
	directories := strings.Split(file, "/")
	directories = directories[:len(directories)-1]
	dirPath := strings.Join(directories, "/")
	if dirPath != "" {
		storage.DeleteDirectory(directories[0])
	}
	storage.Delete(file)
}

func clearDir(storage fs.Disk, path string) {
	storage.DeleteDirectory(path)
}

func setup(t *testing.T, storage fs.Disk, base string) func() {
	t.Helper()

	file := base + "/files.txt"
	subFile := base + "/sub/files.txt"
	sub2File := base + "/sub2/files.txt"
	subSubFile := base + "/sub/sub/files.txt"
	createFile(t, storage, file)
	createFile(t, storage, subFile)
	createFile(t, storage, sub2File)
	createFile(t, storage, subSubFile)
	storage.MakeDirectory(base+"/empty", fs.PUBLIC)

	return func() {
		clearDir(storage, base+"/empty")
		clear(storage, file)
		clear(storage, subFile)
		clear(storage, subSubFile)
	}
}

func getStorage() *Storage {
	s3Config := disk.S3Config{Client: disk.NewMemoryClient()}
	localConfig := disk.LocalConfig{}
	memoryConfig := disk.MemoryConfig{}

	adapters := map[string]fs.Disk{
		"memory": disk.NewMemory(memoryConfig),
		"local":  disk.NewLocal(localConfig),
		"s3":     disk.NewS3(s3Config),
	}

	return New(adapters)
}
