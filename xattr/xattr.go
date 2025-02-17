package xattr

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pirmd/epub"
	"github.com/pkg/errors"
	"github.com/pkg/xattr"
)


var homedir, _ = os.UserHomeDir()
var ebookdir = filepath.Join(homedir,"Downloads/Ebooks")
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

func ListEpubs(directory string) []string {
	var sr []string
	for _, sr2 := range find(directory, ".epub") {
		sr2, err := epub.GetMetadataFromFile(sr2)
		if err != nil {
			errors.Cause(err)
		}
		sr = append(sr, sr2.Title...)
	}
	return sr
}

func GetXattr() []string {
	// We get the epub files
	filelist := find(ebookdir, ".epub")
	// We will store the tags that we find here
	var tags []string
	// We loop over the epub files, actualname being the name of the files
	for _, actualname := range filelist {
		// We pass the actual name of the file to xattr.Get to get it's tags
		value, err := xattr.Get(actualname, prefix)
		if err != nil {
			errors.New("got error")
		}
		if (string(value)) == "" {
			tags = append(tags, "untagged")
			continue
		}
		// We can append the actual name here to a filewith tags list
		tags = append(tags, string(value))
	}
	return tags
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

// Gets the list of files and their tags
func GetXattrMap2(){

}
// Gets the file/s associated with the selectedTag

// func GetTagsMaps(selectedTag string, tagFiles map[string][]string) {
//   fileTags := make(map[string][]string)
//   tagFiles := make(map[string][]string)
//   filelist := find(ebookdir, epub)
//   if files, ok := tagFiles[selectedTag]; ok {
//
//   }
//
// }

func Getfiles(tag string)[]string{
  filelist := find(ebookdir,".epub" )
  // store files here
  var files []string
  // Loop over th epub files
  for _, actualname := range filelist {
    value, err := xattr.Get(actualname, prefix)
    if err != nil{
      errors.New("got error")

    }
    if (string(value)) == ""{
      continue
    }
    if (string(value)) == tag {
      actualname , err := epub.GetMetadataFromFile(actualname)
      if err != nil {
        errors.New("got an error")
      }
      files = append(files, actualname.Title...)
    }
  }
  return files
}

// func Addtag(file string, tag string){
//   xattr.Set(file, prefix,tag)
// }


