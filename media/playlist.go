package media

import (
	"sort"
)

type Playlist []*Video

// By is the type of a "less" function that defines the ordering of its Playlist arguments.
type By func(p1, p2 *Video) bool

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(pl Playlist) {
	ps := &playlistSorter{
		pl: pl,
		by: by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(ps)
}

// playlistSorter joins a By function and a slice of Playlist to be sorted.
type playlistSorter struct {
	pl Playlist
	by func(p1, p2 *Video) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *playlistSorter) Len() int {
	return len(s.pl)
}

// Swap is part of sort.Interface.
func (s *playlistSorter) Swap(i, j int) {
	s.pl[i], s.pl[j] = s.pl[j], s.pl[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *playlistSorter) Less(i, j int) bool {
	return s.by(s.pl[i], s.pl[j])
}

func SortByTimestamp(v1, v2 *Video) bool {
	return v1.Timestamp.After(v2.Timestamp)
}

func SortByViews(v1, v2 *Video) bool {
	return v1.Views < v2.Views
}
