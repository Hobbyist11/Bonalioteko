package xattr

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/pkg/xattr"
)

var (
	homedir, _ = os.UserHomeDir()
	ebookdir   = filepath.Join(homedir, "Downloads/Ebooks")
)

const (
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

func GetXattrmap() map[string]string {
	filelist := find(ebookdir, ".epub")
	tags := make(map[string]string)

	for _, actualname := range filelist {
		value, err := xattr.Get(actualname, prefix)
		if err != nil {
			errors.New("got error")
		}

		if (string(value)) == "" {
			tags[actualname] = "untagged"
			continue
		}
		tags[actualname] = string(value)

	}
	return tags
}

func GetTagsFromPath(filePath string) ([]string, error) {
	tagsBytes, err := xattr.Get(filePath, prefix)
	if err != nil {
		// If the attribute doesn't exist, treat it as no tags, not as an error.
		if strings.Contains(err.Error(), "no such attribute") { // Check for attribute not found error
			// Changed to "" from untagged
			return []string{"untagged"}, nil
		}
		return nil, err
	}

	tagsString := string(tagsBytes)

	if tagsString == "" {
		return []string{"untagged"}, nil
	}
	tags := strings.Split(tagsString, ",")
	return tags, nil
}

func GetXattrMapFilePathToTag() map[string][]string {
	filelist := find(ebookdir, ".epub")

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

func GetXattrMapTagToFilePath() map[string][]string {
	filelist := find(ebookdir, ".epub")
	tagToFiles := make(map[string][]string)
	for _, fileNames := range filelist {
		tags, _ := GetTagsFromPath(fileNames)
		addTagAndFile(fileNames, tags, tagToFiles)

	}
	return tagToFiles
}

func addTagAndFile(filePath string, tags []string, mymap map[string][]string) {
	for _, tag := range tags {
		mymap[tag] = append(mymap[tag], filePath)
	}

	if len(tags) == 0 {
		mymap["untagged"] = append(mymap["untagged"], filePath)
	}
}

func GetUniqueTags(tagFiles map[string][]string) []string {
	uniqueTags := []string{}
	seenTags := make(map[string]bool)

	for tag := range tagFiles {
		if !seenTags[tag] {
			uniqueTags = append(uniqueTags, tag)
			seenTags[tag] = true
		}
	}

	return uniqueTags
}



func MultipleTagsFilter(selectedTags []string) []string {
	init := GetXattrMapTagToFilePath()

	result := init[selectedTags[0]]

	for i := 1; i < len(selectedTags); i++ {
		filesForNextTag := init[selectedTags[i]]
		result = GetIntersection(result, filesForNextTag)
	}

	return result
}

func GetIntersection(setA []string, setB []string) []string {
	var intersection []string
	if len(setA) > len(setB) {
		setB, setA = setA, setB
	}
	hashsetA := make(map[string]bool, len(setA))
	for _, filename := range setA {
		hashsetA[filename] = true
	}

	for _, item := range setB {
		if _, exists := hashsetA[item]; exists {
			intersection = append(intersection, item)
		}
	}

	return intersection
}

func Addtag(file string, tagname []byte) {
	xattr.Set(file, prefix, tagname)
}
