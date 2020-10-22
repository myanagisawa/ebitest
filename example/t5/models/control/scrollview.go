package control

import "github.com/myanagisawa/ebitest/example/t5/models/input"

type listRowView struct {
	UIControlImpl
}

func (r *listRowView) GetRowHeight() int {
	_, h := r.bg.Size()
	return h
}

type listView struct {
	UIControlImpl
	rows []listRowView
}

func (l *listView) GetListHeight() int {
	listHeight := 0
	for i := range l.rows {
		row := l.rows[i]
		listHeight += row.GetRowHeight()
	}
	return listHeight
}

// UIScrollViewImpl ...
type UIScrollViewImpl struct {
	UIControlImpl
	list   *listView
	stroke *input.Stroke
}
