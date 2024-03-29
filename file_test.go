package storage

import (
	"github.com/evolidev/storage/disk"
	"github.com/evolidev/storage/fs"
	"testing"
)

func TestFilePut(t *testing.T) {
	t.Parallel()
	storage := disk.NewMemory(disk.MemoryConfig{})
	file := fs.NewFile(storage, "put", "test.txt")

	file.Put([]byte("test"), fs.PUBLIC)

	f, err := storage.Get("put/test.txt")

	if err != nil {
		t.Errorf("failed to get file")
	}

	if string(f) != "test" {
		t.Errorf("content mismatch")
	}
}

func TestWrite(t *testing.T) {
	t.Parallel()
	storage := disk.NewMemory(disk.MemoryConfig{})
	file := fs.NewFile(storage, "write", "test.txt")

	file.Write([]byte("test"))

	f, err := storage.Get("write/test.txt")

	if err != nil {
		t.Errorf("failed to get file")
	}

	if string(f) != "test" {
		t.Errorf("content mismatch")
	}
}

func TestFileDelete(t *testing.T) {
	storage := disk.NewMemory(disk.MemoryConfig{})
	file := fs.NewFile(storage, "delete", "test.txt")

	file.Put([]byte("test"), fs.PUBLIC)
	file.Delete()

	if storage.Exists("delete/test.txt") {
		t.Errorf("failed to delete file")
	}
}
