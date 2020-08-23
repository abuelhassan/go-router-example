package trie

import (
	"reflect"
	"testing"
)

func TestPathTrie(t *testing.T) {
	tests := []struct {
		name    string
		paths   []string
		values  []interface{}
		queries []string
		want    []interface{}
	}{
		{
			name:    "path trie",
			paths:   []string{"/a", "/a/b", "/a/b/c/d", "/a/b"},
			values:  []interface{}{1, 2, 3, 4},
			queries: []string{"/a", "/a/b", "/a/b/c", "/a/b/c/e"},
			want:    []interface{}{1, 4, nil, nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trie := NewPathTrie()
			for i, p := range tt.paths {
				trie.Put(p, tt.values[i])
			}
			got := make([]interface{}, 0)
			for _, q := range tt.queries {
				got = append(got, trie.Get(q))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_nextPathSegment(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantSegs []string
		wantNxts []int
	}{
		{
			name:     "good path",
			path:     "/a/b/c",
			wantSegs: []string{"/a", "/b", "/c"},
			wantNxts: []int{2, 4, -1},
		},
		{
			name:     "path ends with a slash",
			path:     "/a/b/c/",
			wantSegs: []string{"/a", "/b", "/c", "/"},
			wantNxts: []int{2, 4, 6, -1},
		},
		{
			name:     "empty path",
			path:     "",
			wantSegs: []string{""},
			wantNxts: []int{-1},
		},
		{
			name:     "consecutive slashes //",
			path:     "//",
			wantSegs: []string{"/", "/"},
			wantNxts: []int{1, -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSegs := make([]string, 0)
			gotNxts := make([]int, 0)
			st := 0
			for st != -1 {
				seg, nxt := nextPathSegment(tt.path, st)
				gotSegs = append(gotSegs, seg)
				gotNxts = append(gotNxts, nxt)
				st = nxt
			}
			if !reflect.DeepEqual(gotSegs, tt.wantSegs) {
				t.Errorf("nextPathSegment() gotSegs = %v, want %v", gotSegs, tt.wantSegs)
			}
			if !reflect.DeepEqual(gotNxts, tt.wantNxts) {
				t.Errorf("nextPathSegment() gotNxts = %v, want %v", gotNxts, tt.wantNxts)
			}
		})
	}
}
