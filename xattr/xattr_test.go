package xattr_test

import (
	"testing"

	"Bonalioteko/xattr"

	"github.com/google/go-cmp/cmp"
)

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
		{
			tags: []string{"religion"},
			want: []string{
				"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
				"/var/home/dd/Downloads/Ebooks/r-h-tawney_religion-and-the-rise-of-capitalism_advanced.epub",
			},
		},
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
	want := []string{
		"/var/home/dd/Downloads/Ebooks/g-k-chesterton_the-everlasting-man_advanced.epub",
	}
	got := xattr.MultipleTagsFilter([]string{"unread", "religion"})

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
