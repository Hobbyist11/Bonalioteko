package xattr_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"Bonalioteko/xattr"

	xattrpkg "github.com/pkg/xattr"

	"github.com/google/go-cmp/cmp"
)

//	func TestGetxattrMap(t *testing.T) {
//		want := map[string]string{
//			"/var/home/dd/Downloads/Ebooks/bertrand-russell_roads-to-freedom_advanced.epub":                               "untagged",
//			"/var/home/dd/Downloads/Ebooks/fyodor-dostoevsky_demons_constance-garnett_advanced.epub":                      "philosophy",
//			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_heretics_advanced.epub":                                         "untagged",
//			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_orthodoxy_advanced.epub":                                        "untagged",
//			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub":                              "unread,religion,philosophy",
//			"/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub":                              "untagged",
//			"/var/home/dd/Downloads/Ebooks/john-dewey_human-nature-and-conduct_advanced.epub":                             "untagged",
//			"/var/home/dd/Downloads/Ebooks/john-locke_two-treatises-of-government_advanced.epub":                          "untagged",
//			"/var/home/dd/Downloads/Ebooks/karl-marx_friedrich-engels_the-communist-manifesto_samuel-moore_advanced.epub": "untagged",
//			"/var/home/dd/Downloads/Ebooks/laozi_tao-te-ching_james-legge_advanced.epub":                                  "philosophy",
//			"/var/home/dd/Downloads/Ebooks/leo-tolstoy_the-kingdom-of-god-is-within-you_leo-wiener_advanced.epub":         "untagged",
//			"/var/home/dd/Downloads/Ebooks/liam-oflaherty_the-informer_advanced.epub":                                     "untagged",
//			"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub":                  "religion",
//			"/var/home/dd/Downloads/Ebooks/rene-descartes_philosophical-works_john-veitch_advanced.epub":                  "untagged",
//			"/var/home/dd/Downloads/Ebooks/william-james_pragmatism_advanced.epub":                                        "untagged",
//		}
//		got := xattr.GetXattrmap()
//		if !cmp.Equal(want, got) {
//			t.Error(cmp.Diff(want, got))
//		}
//	}
//
//	func TestGetXattrPathtoTags(t *testing.T) {
//		want := map[string][]string{
//			"/var/home/dd/Downloads/Ebooks/bertrand-russell_roads-to-freedom_advanced.epub":                               {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/fyodor-dostoevsky_demons_constance-garnett_advanced.epub":                      {"philosophy"},
//			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_heretics_advanced.epub":                                         {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_orthodoxy_advanced.epub":                                        {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub":                              {"unread", "religion", "philosophy"},
//			"/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub":                              {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/john-dewey_human-nature-and-conduct_advanced.epub":                             {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/john-locke_two-treatises-of-government_advanced.epub":                          {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/karl-marx_friedrich-engels_the-communist-manifesto_samuel-moore_advanced.epub": {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/laozi_tao-te-ching_james-legge_advanced.epub":                                  {"philosophy"},
//			"/var/home/dd/Downloads/Ebooks/leo-tolstoy_the-kingdom-of-god-is-within-you_leo-wiener_advanced.epub":         {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/liam-oflaherty_the-informer_advanced.epub":                                     {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub":                  {"religion"},
//			"/var/home/dd/Downloads/Ebooks/rene-descartes_philosophical-works_john-veitch_advanced.epub":                  {"untagged"},
//			"/var/home/dd/Downloads/Ebooks/william-james_pragmatism_advanced.epub":                                        {"untagged"},
//		}
//		got := xattr.GetXattrMapFilePathToTag()
//
//		if !cmp.Equal(want, got) {
//			t.Error(cmp.Diff(want, got))
//		}
//	}
//
//	func TestGetXattrTagtoFilepath(t *testing.T) {
//		want := map[string][]string{
//			"philosophy": {
//				"/var/home/dd/Downloads/Ebooks/fyodor-dostoevsky_demons_constance-garnett_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/laozi_tao-te-ching_james-legge_advanced.epub",
//			},
//			"religion": {
//				"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub",
//			},
//			"untagged": {
//				"/var/home/dd/Downloads/Ebooks/bertrand-russell_roads-to-freedom_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/g-k-chesterton_heretics_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/g-k-chesterton_orthodoxy_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/john-dewey_human-nature-and-conduct_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/john-locke_two-treatises-of-government_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/karl-marx_friedrich-engels_the-communist-manifesto_samuel-moore_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/leo-tolstoy_the-kingdom-of-god-is-within-you_leo-wiener_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/liam-oflaherty_the-informer_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/rene-descartes_philosophical-works_john-veitch_advanced.epub",
//				"/var/home/dd/Downloads/Ebooks/william-james_pragmatism_advanced.epub",
//			},
//
//			"unread": {
//				"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
//			},
//		}
//
//		got := xattr.GetXattrMapTagToFilePath()
//
//		if !cmp.Equal(want, got) {
//			t.Error(cmp.Diff(want, got))
//		}
//	}
//
//	func TestMultipleTagsFilter(t *testing.T) {
//		type testCase struct {
//			tags []string
//			want []string
//		}
//
//		testCases := []testCase{
//			{
//				tags: []string{"religion"},
//				want: []string{
//					"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
//					"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub",
//				},
//			},
//			{
//				tags: []string{"unread", "religion"},
//				want: []string{"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub"},
//			},
//		}
//
//		for _, tc := range testCases {
//			got := xattr.MultipleTagsFilter(tc.tags)
//			if !cmp.Equal(tc.want, got) {
//				t.Errorf("MultipleTagsFilter(%v): want %v, got %v", tc.tags, tc.want, got)
//			}
//		}
//		want := []string{
//			"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
//		}
//		got := xattr.MultipleTagsFilter([]string{"unread", "religion"})
//
//		if !cmp.Equal(want, got) {
//			t.Error(cmp.Diff(want, got))
//		}
//	}
//
//	func TestMultipleTagsFilter_EmptySelection(t *testing.T) {
//		// Test with nil input
//		gotNil := xattr.MultipleTagsFilter(nil)
//		if gotNil != nil {
//			t.Errorf("MultipleTagsFilter(nil): expected nil, got %v", gotNil)
//		}
//
//		// Test with empty slice input
//		gotEmpty := xattr.MultipleTagsFilter([]string{})
//		if gotEmpty != nil {
//			t.Errorf("MultipleTagsFilter([]string{}): expected nil, got %v", gotEmpty)
//		}
//	}
// func TestAddtag(t *testing.T) {
// 	want := []string{"education", "philosophy", "religion","unread" }
// 	xattr.Addtag("/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub", []byte("education,philosophy,religion,unread"))
// 	got, _ := xattr.GetTagsFromPath("/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub")
// 	slices.Sort(got)
// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}
// }

func TestAddTag_TempOS(t *testing.T) {
	tmpDir := t.TempDir()

	// Create untagged and "" test case
	testFile := filepath.Join(tmpDir, "test.epub")
	err := os.WriteFile(testFile, []byte("dummy content"), 0o644)
	if err != nil {
		t.Errorf("got error:%s", err)
	}
	err = xattr.Addtag(testFile, []byte("education,politics,religion"))
	if err != nil {
		t.Errorf("got error:%s", err)
	}

	want := []string{"education", "politics", "religion"}
	got, err := xattr.GetTagsFromPath(testFile)
	if err != nil {
		t.Errorf("got error:%s", err)
	}
	slices.Sort(got)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestRemoveTag_TempOS(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "remove.epub")
	os.WriteFile(testFile, []byte("dummy content"), 0o644)

	xattrpkg.Set(testFile, "user.xdg.tags", []byte("education,politics,religion"))

	xattr.RemoveTag(testFile, "education")

	want := []string{"politics", "religion"}

	got, _ := xattr.GetTagsFromPath(testFile)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func (f *FakeXattr) Set(path, name string, data []byte) error {
	if f.storage[path] == nil {
		f.storage[path] = make(map[string][]byte)
	}
	f.storage[path][name] = data
	return nil
}

//
// // // TODO: Remove tags
// // func TestRemovetags(t *testing.T) {
// // 	want := []string{"untagged", "philosophy", "unread", "religion", "education"}
// // 	xattr.RemoveTag("/var/home/dd/Downloads/Ebooks/john-dewey_democracy-and-education_advanced.epub")
// //
// // 	tagsMap := xattr.GetXattrMapTagToFilePath()
// // 	got := xattr.GetUniqueTags(tagsMap)
// //
// // 	if !cmp.Equal(want, got) {
// // 		t.Error(cmp.Diff(want, got))
// // 	}
// // }

type FakeXattr struct {
	storage map[string]map[string][]byte
}

func NewFakeXattr() *FakeXattr {
	return &FakeXattr{storage: make(map[string]map[string][]byte)}
}

func TestUnion(t *testing.T) {
	type testCase struct {
		a, b []string
		want []string
	}

	testcases := []testCase{
		{a: []string{"a", "b"}, b: []string{"a", "c"}, want: []string{"a", "b", "c"}},
		{a: []string{"unread", "b"}, b: []string{"unread", "philosophy"}, want: []string{"b", "philosophy", "unread"}},
		{a: []string{"1", "2", "3"}, b: []string{"4", "2", "3", "4"}, want: []string{"1", "2", "3", "4"}},
	}

	for _, tc := range testcases {
		got := xattr.GetUnion(tc.a, tc.b)
		slices.Sort(got)
		if !cmp.Equal(tc.want, got) {
			t.Error(cmp.Diff(tc.want, got))
		}

	}
}
