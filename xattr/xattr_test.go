package xattr_test

import (
	"testing"

	"Bonalioteko/xattr"

	"github.com/google/go-cmp/cmp"
)

func TestGetXattr(t *testing.T) {
	// I want to see these tags that are found on a folder
	want := []string{"philosophy", "religion"}

	got := xattr.GetXattr()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}


//	// A map has an ID which should be faster for retrieval
//
// Tags should be unique, there should just be tag 1 and tag 2
// Selected tag must show the files associated with these tags

func TestGetxattrMap(t *testing.T) {
	want := map[string]string{
		"/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub": "religion",
		"/home/dd/Downloads/Ebooks/laozi_tao-te-ching_james-legge_advanced.epub":                 "philosophy",
		"/home/dd/Downloads/Ebooks/fyodor-dostoevsky_demons_constance-garnett_advanced.epub":     "philosophy",
	}
	got := xattr.GetXattrmap()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
