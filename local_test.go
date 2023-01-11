package storage

import (
	"github.com/evolidev/storage/disk"
	"github.com/evolidev/storage/fs"
	"os"
	"testing"
)

func TestPut(t *testing.T) {
	t.Parallel()
	t.Run("put should return error if os can not write it", func(t *testing.T) {
		//s3Config := disk.S3Config{Client: disk.NewMemoryClient()}
		localConfig := disk.LocalConfig{PermModeFilePublic: 0400, PermModeDirectoryPublic: 0400}
		storage := disk.NewLocal(localConfig)
		result := storage.Put("local/tmp/error.txt", []byte("test"), fs.PUBLIC)
		defer storage.Delete("local/tmp/error.txt")
		defer storage.DeleteDirectory("local")

		if result == nil {
			t.Errorf("File created")
		}
	})

	t.Run("put private file should add private file", func(t *testing.T) {
		localConfig := disk.LocalConfig{PermModeFilePrivate: 0700, PermModeDirectoryPrivate: 0700}
		storage := disk.NewLocal(localConfig)
		result := storage.Put("local_private_put/tmp.txt", []byte("tmp"), fs.PRIVATE)
		defer storage.DeleteDirectory("local_private_put")

		if result != nil {
			t.Errorf("could not created dir")
		}

		f, err := os.Open("local_private_put/tmp.txt")
		defer f.Close()

		if err != nil {
			t.Errorf("failed to get file")
		}

		stat, err := f.Stat()

		if err != nil {
			t.Errorf("failed to stats of file")
		}

		if stat.Mode().Perm() != 0700 {
			t.Errorf("wrong file mode set")
		}
	})

	t.Run("move to none writeable destination should return error", func(t *testing.T) {
		localConfig := disk.LocalConfig{PermModeFilePrivate: 0700, PermModeDirectoryPrivate: 0700}
		storage := disk.NewLocal(localConfig)
		base := "move_not_writeable"
		d := setup(t, storage, base)
		defer d()
		file := base + "/sub/file_exists.txt"
		err := storage.Put(file, []byte("test"), fs.PUBLIC)

		if err != nil {
			t.Errorf("Could not write file, %s", err)
		}

		os.Mkdir(base+"/restricted", 0000)

		err = storage.Move(file, base+"/restricted/tmp.txt")

		if err == nil {
			t.Errorf("Writing to none existing file")
		}
	})
}

func TestDirectory(t *testing.T) {
	t.Parallel()
	t.Run("make private directory should add private dir", func(t *testing.T) {
		localConfig := disk.LocalConfig{PermModeFilePrivate: 0700, PermModeDirectoryPrivate: 0700}
		storage := disk.NewLocal(localConfig)
		result := storage.MakeDirectory("local_private", fs.PRIVATE)
		defer storage.DeleteDirectory("local_private")

		if result != nil {
			t.Errorf("could not created dir")
		}

		f, err := os.Open("local_private")
		defer f.Close()

		if err != nil {
			t.Errorf("failed to get dir")
		}

		stat, err := f.Stat()

		if err != nil {
			t.Errorf("failed to get dir")
		}

		if stat.Mode().Perm() != 0700 {
			t.Errorf("wrong file mode set")
		}
	})

	t.Run("get files should return no files if files are not accessible", func(t *testing.T) {
		localConfig := disk.LocalConfig{PermModeFilePrivate: 0400, PermModeDirectoryPublic: 0777, PermModeDirectoryPrivate: 0400}
		storage := disk.NewLocal(localConfig)
		result := storage.MakeDirectory("local_list", fs.PUBLIC)
		defer storage.DeleteDirectory("local_list")

		if result != nil {
			t.Errorf("could not created dir")
		}

		storage.Put("local_list/test/tmp.txt", []byte("test"), fs.PRIVATE)

		listResult := storage.Files("local_list/test")

		if len(listResult) > 0 {
			t.Errorf("could read all of it")
		}
	})
}
