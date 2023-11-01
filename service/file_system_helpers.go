package service

import (
	"errors"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func getFoldersIn(rootPath string) ([]string, error) {
	var folders []string

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folderPath := filepath.Join(rootPath, entry.Name())
			folders = append(folders, folderPath)
		}
	}

	return folders, nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false // File does not exist
		}
		log.Printf("ERROR: Failed to check file existence for '%s': %v", filePath, err)
		return false
	}
	return true // File exists
}

func folderExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Folder does not exist
		return false, nil
	}
	// Some other error occurred
	if err != nil {
		return false, err
	}

	// Path exists but is not a folder
	if !info.IsDir() {
		return false, nil
	}

	return true, nil
}

func isFolderEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read item in the folder
	_, err = f.Readdirnames(1)

	// Folder is not empty
	if err == nil {
		return false, nil

	}
	// Folder is empty
	if err == io.EOF {
		return true, nil
	}

	// Some other error occurred
	return false, err
}

func removeFolder(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Printf("ERROR: Failed to remove folder '%s': %v", path, err)
		return
	}
	log.Printf("INFO: Removed folder %s", path)
}
