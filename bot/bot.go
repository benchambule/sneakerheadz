package bot

import (
	"slices"
	"strings"
)

type Request struct {
	Msisdn     string
	Prompt     string
	Parameters map[string]string
}

type Menu struct {
	Name string
	Body string
}

func (m Menu) String() string {
	result := strings.ReplaceAll(m.Body, "\n{{errors}}\n", "")
	return result
}

type Session struct {
	Location    string
	Parameters  map[string]string
	CurrentMenu *Menu
}

var sessions map[string]*Session = make(map[string]*Session)

func HasSession(msisdn string) bool {
	_, ok := sessions[msisdn]
	return ok
}

func ProcessRequest(r *Request) Menu {
	menu := getMenu("000_init", r.Parameters)

	return menu
}

func ProcessRequest1(r *Request) Menu {
	session, ok := sessions[r.Msisdn]

	if session == nil || !ok {
		menu := getMenu("000_init", r.Parameters)
		// TODO: Uncomment this when sessions are implemented

		// sessions[r.Msisdn] = &Session{
		// 	Location:    menu.Name,
		// 	Parameters:  r.Parameters,
		// 	CurrentMenu: &menu,
		// }

		return menu
	}

	menu := getMenu(session.Location, r.Parameters)
	switch session.Location {
	case "000_init":
		menu = processInit(r)

	case "002_recommend":
		menu = getMenu("020_recommended", r.Parameters)

	case "004_evaluate":
		menu = getMenu("040_evaluated", r.Parameters)

	default:
		menu = getMenu("000_init", r.Parameters)
	}

	if slices.Contains([]string{"001_talktoadmin", "003_info", "020_recommended", "040_evaluated"}, menu.Name) {
		delete(sessions, r.Msisdn)

		return menu
	}

	sessions[r.Msisdn].CurrentMenu = &menu
	sessions[r.Msisdn].Location = menu.Name
	return menu
}
