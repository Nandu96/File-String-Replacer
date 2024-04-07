# Template Enabler
This CLI tool can be used to generate files and folders using a reference similar to how a template can be used, but much easier. 

It expects as input, the reference folder path and a text file containing the words to be replaced and the corresponding new word in "," separated pair in each line. The tool replaces the occurance of the word in the contents of each file and also the file and folder names themselves and yes, you can have as many nested folders as you want. 

It also has a "code_mode" option which can be used to replace the lower case ,upper case and camel case occurances of all the entries in the replace pairs file without having to repeat them!

## Usage

```go
./template_enabler "path_to_reference_folder" "path_to_replacement_pairs_file" code_mode
```

Where code_mode expects a boolean value if you want to enrich your replacement pairs file for different cases of same words.