package jiraui

import (
	"strings"
)

type CommandBarFragment struct {
	commandBar  *CommandBar
	commandMode bool
}

func (p *CommandBarFragment) ExecuteCommand() {
	command := string(p.commandBar.text)
	if command == "" {
		return
	}
	commandMode := string([]rune(command)[0])
	switch commandMode {
	case "/":
		log.Debugf("Search down: %q", command)
		if obj, ok := currentPage.(Searcher); ok {
			obj.SetSearch(command)
			obj.Search()
		}
	case "?":
		log.Debugf("Search up: %q", command)
		if obj, ok := currentPage.(Searcher); ok {
			obj.SetSearch(command)
			obj.Search()
		}
	case ":":
		log.Debugf("Command: %q", command)
		handleCommand(command)
	}
}

func handleCommand(command string) {
	if len(command) < 2 {
		// must be :something
		return
	}
	fields := strings.Fields(string(command[1:]))
	action := fields[0]
	var args []string
	if len(fields) > 1 {
		args = fields[1:]
	}
	log.Debugf("handleCommand: action %q, args %s", action, args)
	switch {
	case action == "q" || action == "quit":
		handleQuit()
	case action == "label" || action == "labels":
		handleLabelCommand(args)
	case action == "help":
		handleHelp()
	case action == "watch":
		handleWatchCommand(args)
	case action == "vote":
		handleVoteCommand(true)
	case action == "unvote":
		handleVoteCommand(false)
	case action == "assign":
		handleAssignCommand(args[0])
	case action == "unassign":
		handleAssignCommand("-1")
	case action == "take":
		opts := getJiraOpts()
		handleAssignCommand(opts["user"].(string))
	case action == "comment":
		if len(command) > 10 {
			handleCommentCommand(string(command[9:]))
		}
	case action == "query":
		n := len(":query ")
		if len(command) > n {
			handleQueryCommand(string(command[(n - 1):]))
		}
	case action == "view":
		if len(args) > 0 {
			handleViewCommand(args[0])
		}
	}
}

func handleLabelCommand(args []string) {
	log.Debugf("handleLabelCommand: args %s", args)
	if obj, ok := currentPage.(TicketCommander); ok {
		ticketId := obj.ActiveTicketId()
		if ticketId == "" || args == nil {
			return
		}
		action := "add"
		var labels []string
		switch args[0] {
		case "add":
			action = "add"
			if len(args) > 1 {
				labels = args[1:]
			}
		case "remove":
			action = "remove"
			if len(args) > 1 {
				labels = args[1:]
			}
		default:
			labels = args
		}
		runJiraCmdLabels(ticketId, action, labels)
		obj.Refresh()
	}
}

func handleCommentCommand(comment string) {
	log.Debugf("handleCommentCommand: comment %s", comment)
	if obj, ok := currentPage.(TicketCommander); ok {
		ticketId := obj.ActiveTicketId()
		if ticketId == "" || comment == "" {
			return
		}
		log.Debugf("handleCommentCommand: ticket: %s, comment %s", ticketId, comment)
		runJiraCmdCommentNoEditor(ticketId, comment)
		obj.Refresh()
	}
}

func handleAssignCommand(user string) {
	log.Debugf("handleAssignCommand: user %s", user)
	if obj, ok := currentPage.(TicketCommander); ok {
		ticketId := obj.ActiveTicketId()
		if ticketId == "" || user == "" {
			return
		}
		log.Debugf("handleAssignCommand: ticket: %s, user %s", ticketId, user)
		runJiraCmdAssign(ticketId, user)
		obj.Refresh()
	}
}

func handleViewCommand(ticket string) {
	log.Debugf("handleViewCommand: ticket %s", ticket)
	if ticket == "" {
		return
	}
	q := new(TicketShowPage)
	q.TicketId = ticket
	previousPages = append(previousPages, currentPage)
	currentPage = q
	changePage()
}

func handleQueryCommand(query string) {
	log.Debugf("handleQueryCommand: query %q", query)
	if query == "" {
		return
	}
	q := new(TicketListPage)
	q.ActiveQuery.Name = "adhoc query"
	q.ActiveQuery.JQL = query
	previousPages = append(previousPages, currentPage)
	currentPage = q
	changePage()
}

func handleVoteCommand(up bool) {
	log.Debugf("handleVoteCommand: up %q", up)
	if obj, ok := currentPage.(TicketCommander); ok {
		ticketId := obj.ActiveTicketId()
		if ticketId == "" {
			return
		}
		runJiraCmdVote(ticketId, up)
		obj.Refresh()
	}
}

func handleWatchCommand(args []string) {
	log.Debugf("handleWatchCommand: args %s", args)
	if obj, ok := currentPage.(TicketCommander); ok {
		ticketId := obj.ActiveTicketId()
		if ticketId == "" {
			return
		}
		log.Debugf("handleWatchCommand: ticket: %s, args %s", ticketId, args)
		if len(args) == 0 {
			runJiraCmdWatch(ticketId, "", false) // watch issue
		} else if args[0] == "add" {
			if len(args) > 1 {
				runJiraCmdWatch(ticketId, args[1], false) // add any user as watcher
			} else {
				runJiraCmdWatch(ticketId, "", false) // add self as watcher
			}
		} else if args[0] == "remove" {
			if len(args) > 1 {
				runJiraCmdWatch(ticketId, args[1], true) // remove any user as watcher
			} else {
				runJiraCmdWatch(ticketId, "", true) // remove self as watcher
			}
		} else {
			return
		}
		obj.Refresh()
	}
}

func (p *CommandBarFragment) SetCommandMode(mode bool) {
	p.commandMode = mode
}

func (p *CommandBarFragment) CommandMode() bool {
	return p.commandMode
}

func (p *CommandBarFragment) CommandBar() *CommandBar {
	return p.commandBar
}
