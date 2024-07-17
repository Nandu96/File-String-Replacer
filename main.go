package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	srcPath := flag.String("s", "", "Path to reference folder - (mandatory).")
	replacementPair := flag.String("f", "", "Path to replacement pairs file (mandatory).")
	destPath := flag.String("d", "", "Path to place generated project, which by default is same as reference folder with (generated).")
	inputCodeMode := flag.Bool("code_mode", false, "Replace the UPPERCASE, lowercase and camelCase occurences of WordToReplace with that of NewWord.")
	help := flag.Bool("help", false, "Display help message.")
	verbose := flag.Bool("v", false, "Verbose mode to see detailed logs of processing of your files.")
	flag.Parse()

	if *help {
		displayHelp()
		return
	}

	if len(os.Args) < 2 || *srcPath == "" || *replacementPair == "" {
		fmt.Println("Arguments mismatch!")
		displayHelp()
		os.Exit(1)
	}

	if *destPath == "" {
		*destPath = *srcPath + "(generated)"
	}

	performStringReplacement(*srcPath, *replacementPair, *destPath, inputCodeMode, verbose)
}

func performStringReplacement(srcPath, replacementPair, destPath string, isCodeMode, verbose *bool) {

	if *verbose {
		fmt.Println("verbose mode enabled")
		fmt.Println("Path entered for reference: ", srcPath)
		fmt.Println("File entered for replacement pairs: ", replacementPair)
		fmt.Println("Code mode opted: ", *isCodeMode)
	}
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

	if *isCodeMode {
		fmt.Println("Code mode is activated. Enriching the template map!")
		keyValuePairs = enrichKeyValueMap(keyValuePairs)
	}

	err = duplicateFolderStructure(srcPath, destPath, keyValuePairs, verbose)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nFolder structure duplicated successfully!")
	fmt.Println("Find your generated project at: " + destPath)
}

func duplicateFolderStructure(srcPath, destPath string, keyValuePairs map[string]string, verbose *bool) error {
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
			err = duplicateFolderStructure(srcFile, destFile, keyValuePairs, verbose)
			if err != nil {
				return err
			}
		} else {
			if *verbose {
				fmt.Println("Generating file: " + destFile)
			}
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

func displayHelp() {
	fmt.Println("\nUsage:\n\texecutable -s \"path_to_reference_folder\" -f \"path_to_replacement_pairs_file\"")
	fmt.Println("Description:")
	fmt.Println("\texecutable: The name of the executable that is downloaded. For example, in Windows it should be used like fsr.exe.")
	fmt.Println("\tpath_to_reference_folder: The location of the folder to be replicated. Forward/backward slashes must be used appropriately.")
	fmt.Println("\tpath_to_replacement_pairs_file: The location of the replacement pairs file. File should contain the WordToReplace and NewWord in comma separated format without any spaces.")
	fmt.Println("Example:\n\t ./fsr -s \"Documents/My Project/Hello World\" -f \"Documents/My Project/keyfile.txt\"")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Print("\n\n")
}
