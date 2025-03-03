package xattr_test

import (
	"testing"

	"Bonalioteko/xattr"

	"github.com/google/go-cmp/cmp"
)

/*
func TestGetXattr(t *testing.T) {
	// I want to see these tags that are found on a folder
	want := []string{
		"untagged",
		"philosophy",
		"untagged",
		"untagged",
		"untagged",
		"untagged",
		"untagged",
		"untagged",
		"untagged",
		"philosophy",
		"untagged",
		"untagged",
		"religion",
		"untagged",
		"untagged",
	}

	got := xattr.GetXattr()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
*/

// A map has an ID which should be faster for retrieval
// Selected tag must show the files associated with these tags
// The IDs are unique, however the values aren't but we only need to render them once in the UI

func TestGetxattrMap(t *testing.T) {
	want := map[string]string{
		"/var/home/dd/Downloads/Ebooks/bertrand-russell_roads-to-freedom_advanced.epub":                               "untagged",
		"/var/home/dd/Downloads/Ebooks/fyodor-dostoevsky_demons_constance-garnett_advanced.epub":                      "philosophy",
		"/var/home/dd/Downloads/Ebooks/g-k-chesterton_heretics_advanced.epub":                                         "untagged",
		"/var/home/dd/Downloads/Ebooks/g-k-chesterton_orthodoxy_advanced.epub":                                        "untagged",
		"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub":                              "unread,religion,philosophy",
		"/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub":                              "untagged",
		"/var/home/dd/Downloads/Ebooks/john-dewey_human-nature-and-conduct_advanced.epub":                             "untagged",
		"/var/home/dd/Downloads/Ebooks/john-locke_two-treatises-of-government_advanced.epub":                          "untagged",
		"/var/home/dd/Downloads/Ebooks/karl-marx_friedrich-engels_the-communist-manifesto_samuel-moore_advanced.epub": "untagged",
		"/var/home/dd/Downloads/Ebooks/laozi_tao-te-ching_james-legge_advanced.epub":                                  "philosophy",
		"/var/home/dd/Downloads/Ebooks/leo-tolstoy_the-kingdom-of-god-is-within-you_leo-wiener_advanced.epub":         "untagged",
		"/var/home/dd/Downloads/Ebooks/liam-oflaherty_the-informer_advanced.epub":                                     "untagged",
		"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub":                  "religion",
		"/var/home/dd/Downloads/Ebooks/rene-descartes_philosophical-works_john-veitch_advanced.epub":                  "untagged",
		"/var/home/dd/Downloads/Ebooks/william-james_pragmatism_advanced.epub":                                        "untagged",
	}
	got := xattr.GetXattrmap()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetXattrPathtoTags(t *testing.T) {
	want := map[string][]string{
		"/var/home/dd/Downloads/Ebooks/bertrand-russell_roads-to-freedom_advanced.epub":                               {"untagged"},
		"/var/home/dd/Downloads/Ebooks/fyodor-dostoevsky_demons_constance-garnett_advanced.epub":                      {"philosophy"},
		"/var/home/dd/Downloads/Ebooks/g-k-chesterton_heretics_advanced.epub":                                         {"untagged"},
		"/var/home/dd/Downloads/Ebooks/g-k-chesterton_orthodoxy_advanced.epub":                                        {"untagged"},
		"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub":                              {"unread", "religion", "philosophy"},
		"/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub":                              {"untagged"},
		"/var/home/dd/Downloads/Ebooks/john-dewey_human-nature-and-conduct_advanced.epub":                             {"untagged"},
		"/var/home/dd/Downloads/Ebooks/john-locke_two-treatises-of-government_advanced.epub":                          {"untagged"},
		"/var/home/dd/Downloads/Ebooks/karl-marx_friedrich-engels_the-communist-manifesto_samuel-moore_advanced.epub": {"untagged"},
		"/var/home/dd/Downloads/Ebooks/laozi_tao-te-ching_james-legge_advanced.epub":                                  {"philosophy"},
		"/var/home/dd/Downloads/Ebooks/leo-tolstoy_the-kingdom-of-god-is-within-you_leo-wiener_advanced.epub":         {"untagged"},
		"/var/home/dd/Downloads/Ebooks/liam-oflaherty_the-informer_advanced.epub":                                     {"untagged"},
		"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub":                  {"religion"},
		"/var/home/dd/Downloads/Ebooks/rene-descartes_philosophical-works_john-veitch_advanced.epub":                  {"untagged"},
		"/var/home/dd/Downloads/Ebooks/william-james_pragmatism_advanced.epub":                                        {"untagged"},
	}
	got := xattr.GetXattrMapFilePathToTag()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetXattrTagtoFilepath(t *testing.T) {
	want := map[string][]string{
		"philosophy": {
			"/var/home/dd/Downloads/Ebooks/fyodor-dostoevsky_demons_constance-garnett_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/laozi_tao-te-ching_james-legge_advanced.epub",
		},
		"religion": {
			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub",
		},
		"untagged": {
			"/var/home/dd/Downloads/Ebooks/bertrand-russell_roads-to-freedom_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_heretics_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_orthodoxy_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/john-dewey_human-nature-and-conduct_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/john-locke_two-treatises-of-government_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/karl-marx_friedrich-engels_the-communist-manifesto_samuel-moore_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/leo-tolstoy_the-kingdom-of-god-is-within-you_leo-wiener_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/liam-oflaherty_the-informer_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/rene-descartes_philosophical-works_john-veitch_advanced.epub",
			"/var/home/dd/Downloads/Ebooks/william-james_pragmatism_advanced.epub",
		},

		// This has multiple tags
		"unread": {
			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
		},
	}

	got := xattr.GetXattrMapTagToFilePath()

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestMultipleTagsFilter(t *testing.T) {

	type testCase struct {
		tags []string
		want []string
	}

	testCases := []testCase{
		// Test single tag
		{
			tags: []string{"religion"},
			want: []string{
				"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
				"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub",
			},
		},
		// Test two out of three tags
		{
			tags: []string{"unread", "religion"},
			want: []string{"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub"},
		},
	}

	for _, tc := range testCases {
		got := xattr.MultipleTagsFilter(tc.tags)
		if !cmp.Equal(tc.want, got) {
			t.Errorf("MultipleTagsFilter(%v): want %v, got %v", tc.tags, tc.want, got)
		}
	}
	// Want is the file names with these tags
	// Make test case with single tag
	want := []string{
		"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
	}
	got := xattr.MultipleTagsFilter([]string{"unread", "religion"})

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

// func TestGetfiles(t *testing.T) {
// 	// Getting files with the selected tag/ value in map reverse of getting value from key
// 	// I want the slice of file names
// 	want := []string{"Demons", "Tao Te Ching"}
// 	got := xattr.Getfiles("philosophy")
//
// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}
// }

// func TestAddtags(t *testing.T) {
// 	// I want to be able to add a tag to a certain file
// 	want := map[string]string{
// 		"/var/home/dd/Downloads/Ebooks/leo-tolstoy_the-kingdom-of-god-is-within-you_leo-wiener_advanced.epub":         "religion",
// 	}
// 	got := xattr.Addtag("/var/home/dd/Downloads/Ebooks/leo-tolstoy_the-kingdom-of-god-is-within-you_leo-wiener_advanced.epub","religion")
// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}
// }
