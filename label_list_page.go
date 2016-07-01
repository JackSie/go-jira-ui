package jiraui

import (
	"fmt"

	ui "github.com/gizak/termui"
)

type LabelListPage struct {
	BaseListPage
	CommandBarFragment
	StatusBarFragment
	labelCounts map[string]int
	ActiveQuery Query
}

func (p *LabelListPage) Search() {
	s := p.ActiveSearch
	n := len(p.cachedResults)
	if s.command == "" {
		return
	}
	increment := 1
	if s.directionUp {
		increment = -1
	}
	// we use modulo here so we can loop through every line.
	// adding 'n' means we never have '-1 % n'.
	startLine := (p.uiList.Cursor + n + increment) % n
	for i := startLine; i != p.uiList.Cursor; i = (i + increment + n) % n {
		if s.re.MatchString(p.cachedResults[i]) {
			p.uiList.SetCursorLine(i)
			p.Update()
			break
		}
	}
}

func (p *LabelListPage) labelsAsSortedList() []string {
	return sortedKeys(p.labelCounts)
}

func (p *LabelListPage) labelsAsSortedListWithCounts() []string {
	data := p.labelsAsSortedList()
	ret := make([]string, len(data))
	for i, v := range data {
		ret[i] = fmt.Sprintf("%s (%d found)", v, p.labelCounts[v])
	}
	return ret
}

func (p *LabelListPage) SelectItem() {
	label := p.cachedResults[p.uiList.Cursor]
	// calling TicketList page is the last 'previousPage'. We leave it on the stack,
	// as we will likely want to GoBack to it.
	switch oldTicketListPage := previousPages[len(previousPages)-1].(type) {
	default:
		return
	case *TicketListPage:
		q := new(TicketListPage)
		if label == "NOT LABELLED" {
			q.ActiveQuery.Name = oldTicketListPage.ActiveQuery.Name + " (unlabelled)"
			q.ActiveQuery.JQL = oldTicketListPage.ActiveQuery.JQL + " AND labels IS EMPTY"
		} else {
			q.ActiveQuery.Name = oldTicketListPage.ActiveQuery.Name + "+" + label
			q.ActiveQuery.JQL = oldTicketListPage.ActiveQuery.JQL + " AND labels = " + label
		}
		// our label list is useful, prob want to come back to it to look at a different
		// label
		previousPages = append(previousPages, currentPage)
		currentPage = q
		changePage()
	}
}

func (p *LabelListPage) markActiveLine() {
	for i, v := range p.cachedResults {
		p.displayLines[i] = fmt.Sprintf("%-40s -- %d tickets", v, p.labelCounts[v])
	}
}

func (p *LabelListPage) GoBack() {
	if len(previousPages) == 0 {
		currentPage = new(QueryPage)
	} else {
		currentPage, previousPages = previousPages[len(previousPages)-1], previousPages[:len(previousPages)-1]
	}
	changePage()
}

func (p *LabelListPage) Update() {
	ls := p.uiList
	p.markActiveLine()
	ls.Items = p.displayLines
	ui.Render(ls)
	p.statusBar.Update()
	p.commandBar.Update()
}

func (p *LabelListPage) Create() {
	ui.Clear()
	ls := NewScrollableList()
	p.uiList = ls
	if p.statusBar == nil {
		p.statusBar = new(StatusBar)
	}
	if p.commandBar == nil {
		p.commandBar = commandBar
	}
	queryName := p.ActiveQuery.Name
	queryJQL := p.ActiveQuery.JQL
	p.labelCounts = countLabelsFromQuery(queryJQL)
	p.cachedResults = p.labelsAsSortedList()
	p.isPopulated = true
	p.displayLines = make([]string, len(p.cachedResults))
	ls.ItemFgColor = ui.ColorYellow
	ls.BorderLabel = fmt.Sprintf("Label view -- %s: %s", queryName, queryJQL)
	ls.Height = ui.TermHeight() - 2
	ls.Width = ui.TermWidth()
	ls.Y = 0
	p.statusBar.Create()
	p.commandBar.Create()
	p.Update()
}
