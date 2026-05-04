package xattr

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"Bonalioteko/config"

	"github.com/pkg/xattr"
)

const (
	prefix = "user.xdg.tags"
)

func GetEbookDir() (string, error) {
	cfg, err := config.ParseConfig()
	if err != nil {
		return "", err
	}
	return cfg.Settings.EbookDir, nil
}

func InitEbookdir() (string, error) {
	ebookdir, err := GetEbookDir()
	if err != nil {
		return "", err
	}
	return ebookdir, nil
}

// XattrClient defines the contract for extended attribute operations.
type XattrClient interface {
	Get(path, name string) ([]byte, error)
	Set(path, name string, data []byte) error
	Remove(path, name string) error
}

// RealXattr implements the interface using the actual OS calls.
type RealXattr struct{}

func (r RealXattr) Get(path, name string) ([]byte, error)    { return xattr.Get(path, name) }
func (r RealXattr) Set(path, name string, data []byte) error { return xattr.Set(path, name, data) }
func (r RealXattr) Remove(path, name string) error           { return xattr.Remove(path, name) }

type TagManager struct {
	Client XattrClient
	Prefix string
}

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

func GetXattrmap(directory string) map[string]string {
	filelist := find(directory, ".epub")
	tags := make(map[string]string)

	for _, actualname := range filelist {
		value, err := xattr.Get(actualname, prefix)
		if err != nil {
			log.Printf("error:%v", err)
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

func GetXattrMapFilePathToTag(directory string) map[string][]string {
	filelist := find(directory, ".epub")

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

func GetXattrMapTagToFilePath(directory string) map[string][]string {
	filelist := find(directory, ".epub")
	tagToFiles := make(map[string][]string)
	for _, fileNames := range filelist {
		tags, _ := GetTagsFromPath(fileNames)
		AddTagAndFile(fileNames, tags, tagToFiles)

	}
	return tagToFiles
}

func AddTagAndFile(filePath string, tags []string, mymap map[string][]string) {
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

func MultipleTagsFilter(selectedTags []string, tags map[string][]string) []string {
	if len(selectedTags) == 0 {
		return nil
	}

	result := tags[selectedTags[0]]

	for i := 1; i < len(selectedTags); i++ {
		filesForNextTag := tags[selectedTags[i]]
		result = GetIntersection(result, filesForNextTag)
	}

	return result
}

func GetIntersection(setA []string, setB []string) []string {
	var intersection []string
	if len(setA) > len(setB) {
		setB, setA = setA, setB
	}
	hashsetA := CreateHashSet(setA)

	for _, item := range setB {
		if _, exists := hashsetA[item]; exists {
			intersection = append(intersection, item)
		}
	}

	return intersection
}

func GetUnion(setA []string, setB []string) []string {
	var result []string
	if len(setA) > len(setB) {
		setB, setA = setA, setB
	}

	hashsetA := CreateHashSet(setA)
	for key := range hashsetA {
		if key == "" || key == " " || key == "untagged" {
			continue
		}

		result = append(result, key)
	}
	hashsetB := CreateHashSet(setB)
	for key := range hashsetB {
		if _, exists := hashsetA[key]; exists {
			continue
		}
		if key == "" || key == " " || key == "untagged" {
			continue
		}
		result = append(result, key)
	}

	return result
}

func CreateHashSet(set []string) map[string]bool {
	hashsetA := make(map[string]bool, len(set))
	for _, filename := range set {
		hashsetA[filename] = true
	}
	return hashsetA
}

// Add tag adds an xattr tag to a selected file
func Addtag(filepath string, newTags []byte) error {
	existingTags, err := xattr.Get(filepath, prefix)
	// Empty tags case
	if err != nil {
		return xattr.Set(filepath, prefix, newTags)
	}

	currentString := string(existingTags)
	if currentString == "" || currentString == "untagged" {
		return xattr.Set(filepath, prefix, newTags)
	} else {

		merged := GetUnion(strings.Split(string(existingTags), ","), strings.Split(string(newTags), ","))
		return xattr.Set(filepath, prefix, []byte(strings.Join(merged, ",")))
	}
}

// Remove tag removes  xattr tags on the file
func RemoveTag(filepath string, tagToRemove string) error {
	tagbyte, err := xattr.Get(filepath, prefix)
	if err != nil {
		return err
	}

	var newTags []string
	for tag := range strings.SplitSeq(string(tagbyte), ",") {
		if tag != "" && tag != tagToRemove {
			newTags = append(newTags, tag)
		}
	}

	if len(newTags) == 0 {
		return xattr.Remove(filepath, prefix)
	}

	return xattr.Set(filepath, prefix, []byte(strings.Join(newTags, ",")))
}
