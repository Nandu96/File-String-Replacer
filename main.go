package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) != 4 {
		fmt.Println("Arguments mismatch!\nExpected Usage: program \"path_to_reference_folder\" \"path_to_replacement_pairs_file\" code_mode")
		os.Exit(1)
	}

	srcPath := os.Args[1]
	destPath := srcPath + "(generated)"
	replacementPair := os.Args[2]
	inputCodeMode := os.Args[3]

	performStringReplacement(srcPath, replacementPair, destPath, inputCodeMode)
}

func performStringReplacement(srcPath, replacementPair, destPath, inputCodeMode string) {

	isCodeMode, err := strconv.ParseBool(inputCodeMode)
	if err != nil {
		fmt.Println("Error: Third argument should be 'true' or 'false' to indicate code_mode")
		os.Exit(1)
	}

	fmt.Println("Path entered for reference: ", srcPath)
	fmt.Println("File entered for replacement pairs: ", replacementPair)
	fmt.Println("Code mode opted: ", isCodeMode)

	keyValuePairs := make(map[string]string)
	keyValueFile, err := os.Open(replacementPair)
	if err != nil {
		fmt.Println("Error opening file for replacement pairs:", err)
		return
	}
	defer keyValueFile.Close()

	scanner := bufio.NewScanner(keyValueFile)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ",")

		if len(words) == 2 {
			keyValuePairs[words[0]] = words[1]
		} else {
			fmt.Println("Skipping invalid line in replacement pairs file which does not follow convention: word_to_replace, new_word. \nLine contains entry:", line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if isCodeMode {
		fmt.Println("Code mode is activated. Enriching the template map!")
		keyValuePairs = enrichKeyValueMap(keyValuePairs)
	}

	err = duplicateFolderStructure(srcPath, destPath, keyValuePairs)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Folder structure duplicated successfully!")
}

func duplicateFolderStructure(srcPath, destPath string, keyValuePairs map[string]string) error {
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", srcPath)
	}

	err = os.MkdirAll(destPath, srcInfo.Mode())
	if err != nil {
		return err
	}

	contents, err := os.ReadDir(srcPath)
	if err != nil {
		return err
	}

	for _, content := range contents {
		srcFile := filepath.Join(srcPath, content.Name())
		destFile := filepath.Join(destPath, getReplacementFileNameOrDefault(keyValuePairs, content.Name()))

		if content.IsDir() {
			err = duplicateFolderStructure(srcFile, destFile, keyValuePairs)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcFile, destFile, keyValuePairs)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func replaceWordInFile(filePath string, keyValuePairs map[string]string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	newData := string(data)
	for key, value := range keyValuePairs {
		newData = strings.ReplaceAll(newData, key, value)
	}

	err = os.WriteFile(filePath, []byte(newData), 0644)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(srcFile, destFile string, keyValuePairs map[string]string) error {
	srcFileStat, err := os.Stat(srcFile)
	if err != nil {
		return err
	}

	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dest, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer dest.Close()

	err = dest.Chmod(srcFileStat.Mode())
	if err != nil {
		return err
	}

	_, err = dest.ReadFrom(src)
	if err != nil {
		return err
	}

	err = replaceWordInFile(destFile, keyValuePairs) // make file
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

func enrichKeyValueMap(myMap map[string]string) map[string]string {
	newMap := make(map[string]string)
	for key, value := range myMap {
		newMap[key] = value
		uppercaseKey := strings.ToUpper(key)
		uppercaseValue := strings.ToUpper(value)
		newMap[uppercaseKey] = uppercaseValue
		lowercaseKey := strings.ToLower(key)
		lowercaseValue := strings.ToLower(value)
		newMap[lowercaseKey] = lowercaseValue
		addCamelCasePairIfNeeded(newMap, key, value)
	}
	return newMap
}

func addCamelCasePairIfNeeded(newMap map[string]string, key string, value string) {
	if strings.HasPrefix(key, strings.ToUpper(key[:1])) && strings.HasPrefix(value, strings.ToUpper(value[:1])) {
		firstCharLowerKey := strings.ToLower(key[:1]) + key[1:]
		if _, exists := newMap[firstCharLowerKey]; !exists {
			// Only add if the key doesn't already exist
			firstCharLowerValue := strings.ToLower(value[:1]) + value[1:]
			newMap[firstCharLowerKey] = firstCharLowerValue
		}
	}
}

func getReplacementFileNameOrDefault(keyValuePairs map[string]string, key string) string {
	for mapKey, mapValue := range keyValuePairs {
		if strings.Contains(key, mapKey) {
			return strings.ReplaceAll(key, mapKey, mapValue)
		}
	}
	return key
}
