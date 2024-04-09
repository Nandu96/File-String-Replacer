package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	reference_folder                               = "testdata/directory/my_folder"
	replacement_pairs_file                         = "testdata/directory/map_file.txt"
	actual_generated_folder_path                   = "testdata/directory/my_folder(generated)"
	expected_generated_folder_path_code_mode_false = "testdata/expected_directories/expected_code_mode_false"
	expected_generated_folder_path_code_mode_true  = "testdata/expected_directories/expected_code_mode_true"
)

func teardown() {
	if err := deleteRepository(actual_generated_folder_path); err != nil {
		fmt.Printf("Error deleting repository: %v\n", err)
	}
}

func TestValidInputFiles_CodeModeFalse_FolderStructureIsCreatedWithCaseSensitiveReplacement(t *testing.T) {

	performStringReplacement(reference_folder, replacement_pairs_file, actual_generated_folder_path, "false")

	assertFolderContentsMatch(actual_generated_folder_path, expected_generated_folder_path_code_mode_false, t)
	defer teardown()
}

func TestValidInputFiles_CodeModeTrue_FolderStructureIsCreatedWithAllCaseReplacement(t *testing.T) {

	performStringReplacement(reference_folder, replacement_pairs_file, actual_generated_folder_path, "true")

	assertFolderContentsMatch(actual_generated_folder_path, expected_generated_folder_path_code_mode_true, t)
	defer teardown()
}

func assertFolderContentsMatch(actual_generated_folder_path, expected_generated_folder_path string, t *testing.T) {
	files1, err := getAllFiles(actual_generated_folder_path)
	if err != nil {
		t.Fatalf("error getting files from first folder: %v", err)
	}

	files2, err := getAllFiles(expected_generated_folder_path)
	if err != nil {
		t.Fatalf("error getting files from second folder: %v", err)
	}

	if !reflect.DeepEqual(files1, files2) {
		t.Errorf("folder contents are not an exact match")
	}
}

func getAllFiles(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func deleteRepository(repoPath string) error {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return nil
	}
	if err := os.RemoveAll(repoPath); err != nil {
		return fmt.Errorf("error deleting repository: %v", err)
	}
	return nil
}
