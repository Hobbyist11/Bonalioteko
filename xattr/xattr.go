package xattr

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pirmd/epub"
	"github.com/pkg/errors"
	"github.com/pkg/xattr"
)

var (
	homedir, _ = os.UserHomeDir()
	ebookdir   = filepath.Join(homedir, "Downloads/Ebooks")
)

const (
	// ebookdir = "$HOME/Downloads/Ebooks/"
	// ebookdir = "/var/home/dd/Downloads/Ebooks/"
	// The extended attribute we want
	prefix = "user.xdg.tags"
)

func find(root, ext string) []string {
	var filename []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			filename = append(filename, s)
		}
		return nil
	})
	return filename
}

// Get xattr map
func GetXattrmap() map[string]string {
	filelist := find(ebookdir, ".epub")
	tags := make(map[string]string)

	for _, actualname := range filelist {
		value, err := xattr.Get(actualname, prefix)
		if err != nil {
			errors.New("got error")
		}

		if (string(value)) == "" {
			// append as "untagged"
			tags[actualname] = "untagged"
			continue
		}
		tags[actualname] = string(value)

	}
	return tags
}

// Gets slice of Tags from a filePath
func GetTagsFromPath(filePath string) ([]string, error) {
	tagsBytes, err := xattr.Get(filePath, prefix)
	if err != nil {
		// If the attribute doesn't exist, treat it as no tags, not as an error.
		if strings.Contains(err.Error(), "no such attribute") { // Check for attribute not found error
			// Changed to "" from untagged
			return []string{"untagged"}, nil
		}
		// FIX:  Having this as "untagged" breaks the untagged selector
		// But  "" makes the choice "untagged disappear"
		return nil, err
	}

	// Assuming tags are comma-separated strings
	tagsString := string(tagsBytes)

	if tagsString == "" { // Handle empty tags attribute
		return []string{"untagged"}, nil
	}
	tags := strings.Split(tagsString, ",")
	return tags, nil
}

// GetXattrMapFilePathToTag Gets a map of key filepath to tag values
func GetXattrMapFilePathToTag() map[string][]string {
	filelist := find(ebookdir, ".epub")

	// File path to tag
	fileToTag := make(map[string][]string)

	for _, fileNames := range filelist {
		tags, _ := GetTagsFromPath(fileNames)
		if tags == nil {
			tags = append(tags, "untagged")
		}
		addFileAndTag(fileNames, tags, fileToTag)
	}
	return fileToTag
}

func addFileAndTag(filePath string, tags []string, mymap map[string][]string) {
	mymap[filePath] = tags
}

// Gets the key tags and  filepath values
func GetXattrMapTagToFilePath() map[string][]string {
	filelist := find(ebookdir, ".epub")
	// Tags to file path
	tagToFiles := make(map[string][]string)
	for _, fileNames := range filelist {
		tags, _ := GetTagsFromPath(fileNames)
		addTagAndFile(fileNames, tags, tagToFiles)

	}
	return tagToFiles
}

// Adds the found tag to the file
func addTagAndFile(filePath string, tags []string, mymap map[string][]string) {
	for _, tag := range tags {
		// tags are the keys, filepath is the []value
		mymap[tag] = append(mymap[tag], filePath)
	}

	if len(tags) == 0 {
		mymap["untagged"] = append(mymap["untagged"], filePath)
	}
}

// GetUniqueTags Gets Unique tags from key tags and filepath map
func GetUniqueTags(tagFiles map[string][]string) []string {
	uniqueTags := []string{}
	seenTags := make(map[string]bool) // Use a map to track seen tags

	for tag := range tagFiles {
		if !seenTags[tag] {
			uniqueTags = append(uniqueTags, tag)
			seenTags[tag] = true
		}
	}

	return uniqueTags
}

// GetTitlesSlice Gets the files associated with a tag
func GetTitlesSlice(tag string) []string {
	filelist := find(ebookdir, ".epub")
	// store files here
	var files []string
	// Loop over th epub files
	for _, actualname := range filelist {
		value, err := xattr.Get(actualname, prefix)
		if err != nil {
			errors.New("got error")
		}
		// This is what was causing the untagged files not to show up
		// if (string(value)) == "" {
		// 	continue
		// }
		if (string(value)) == tag {
			actualname, err := epub.GetMetadataFromFile(actualname)
			if err != nil {
				errors.New("got an error")
			}
			files = append(files, actualname.Title...)
		}
	}
	return files
}

func MultipleTagsFilter(selectedTags []string) []string {
	init := GetXattrMapTagToFilePath()

	result := init[selectedTags[0]]

	// Find the intersection between the remaining tags
	for i := 1; i < len(selectedTags); i++ {
		// nextTag := GetTitlesSlice(selectedTags[i])
		// nextTag :=selectedTags[i] 
		filesForNextTag := init[selectedTags[i]]
		// If intesect, turn result = intersection
		result = GetIntersection(result,filesForNextTag)
	}

	return result
}

// GetIntersection finds the intersection between two slices
func GetIntersection(setA []string, setB []string) []string {
	var intersection []string
	// Compare the sets find which one is smaller
	if len(setA) > len(setB) {
		setB, setA = setA, setB
	}
	// Create hash map of the smaller set
	hashsetA := make(map[string]bool, len(setA))
	for _, filename := range setA {
		hashsetA[filename] = true
	}

	// Loop over the second set and find the intersection
	// for _, elements := range setB {
	// 	if(hashsetA[elements] == true){
	// 		intersection = append(intersection,elements)
	// 	}
	// }
	for _, item := range setB{
		if _, exists := hashsetA[item]; exists{
			intersection = append(intersection, item)
		}
	}

	return intersection
}

// Addtag adds a tag to a file
func Addtag(file string, tagname []byte) {
	xattr.Set(file, prefix, tagname)
}

// Unused, Map is now used
//Gets slice of xattr tags
// func GetXattr() []string {
// 	// We get the epub files
// 	filelist := find(ebookdir, ".epub")
// 	// We will store the tags that we find here
// 	var tags []string
// 	// We loop over the epub files, actualname being the name of the files
// 	for _, actualname := range filelist {
// 		// We pass the actual name of the file to xattr.Get to get it's tags
// 		value, err := xattr.Get(actualname, prefix)
// 		if err != nil {
// 			errors.New("got error")
// 		}
//
//
// 		// Changed "" from "untagged to """
// 		if (string(value)) == "" {
// 			tags = append(tags, "untagged")
// 			continue
// 		}
// 		// We can append the actual name here to a filewith tags list
// 		tags = append(tags, string(value))
// 	}
// 	return tags
// }

// // Unused
// func ListEpub(directory string) []string {
// 	sr2, err := epub.GetMetadataFromFile(directory)
// 	if err != nil {
// 		errors.Cause(err)
// 	}
// 	return sr2.Title
// }
//
// // Defined in main
// func ListEpubs(directory string) []string {
// 	var sr []string
// 	for _, sr2 := range find(directory, ".epub") {
// 		sr2, err := epub.GetMetadataFromFile(sr2)
// 		if err != nil {
// 			errors.Cause(err)
// 		}
// 		sr = append(sr, sr2.Title...)
// 	}
// 	return sr
// }

// Gets the file/s associated with the selectedTag
// Almost duplicate of GetTitlesSlice
// func GetTagsMaps(selectedTag string, tagFiles map[string][]string) string {
// 	if files, ok := tagFiles[selectedTag]; ok {
// 		for _, file := range files {
// 			return file
// 		}
//
