# Template Enabler
This CLI tool can be used to generate files and folders using a reference, similar to how a template is used to generate fils, but much easier. 

It expects as input, the reference folder path and a text file containing the words to be replaced and the corresponding new word in "," separated pair in each line. The tool replaces the occurance of the word in the contents of each file and also the file and folder names themselves and yes, you can have as many nested folders as you want. 

It also has a "code_mode" option which can be used to replace the lower case ,upper case and camel case occurances of all the entries in the replace pairs file without having to repeat them!

## Usage

### Option 1:

Download and execute the binary from the [Releases](https://github.com/Nandu96/File-String-Replacer/releases) corresponding to the system being used: 

```go
// example 1 - macOs

./fsr "path/to/reference/folder" "path/to/replacement_pairs_file" code_mode

// example 2 - Windows

.\fsr.exe "path\to\reference\folder" "path\to\replacement_pairs_file" code_mode
```

### Option 2:

If you have go installed, you can directly run the main.go file by passing the arguments:

```go
go run main.go "path/to/reference/folder" "path/to/replacement_pairs_file" code_mode
```

A sample replacement_pairs_file contains the word to replace and the new word separated by `,` character with `no spaces after it`.

code_mode has 2 valid options:

`false` - String replacement is done for the replacement pairs in a case-sensitive format (exact match found for the word will be replaced)

`true`  - String replacement is done for the replacement pairs in a case-insensitive format, replacing the upper case, lower case and camel case occurances of the word.


```
NOTE:
1. Make sure the folder contents have proper read accesses. 
2. In addition to file contents, file names and folder names are also replaced in case a match is found.
```

## Example

Sample replcaement_pairs_file:

```go
//key_word.txt
MyExistingWord,MyNewWord
MyExistingAppName,MyNewAppName
I want this whole sentence replaced!,This is the sentence I want to replace it with!
```

### code_mode Usage explained:

When code_mode is used as true, in addition to normally replacing the pairs, the tool replaces the below occurances of the example as well:
```go
MYEXISTINGWORD -> MYNEWWORD
myexistingword -> mynewword
myExistingWord -> myNewWord //useful for changing variable names as well.
```

## Limitations

1. snake_case occurance replacement is not handled for the time being but if the word is part of a snake_case variable, its occurance is replaced with the lower case of the pair.
Example: my_existing_word is not replaced with my_new_word, but myexistingword_is_awesome is replaced wirth mynewword_is_awesome
2. When using code_mode = true, camel case occurance is ignored and only lower case occurance happens in case the "word to be replaced" starts does not contain any UpperCase alphabets after the starting letter. 
Example: If the pair entered is: `Myexistingword,MyNewWord` then `myExistingWord` is not replaced as it is not identifyable from the input.
3. While the code works for any number of nested folders and files, the length of the arguments you are passing may be limited by the system being used. This shouldn't be a problem for most practical use cases.
4. Replacing entire sentences will not work if they have `,` in them.
