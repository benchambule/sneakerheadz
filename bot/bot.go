package bot

type Request struct {
	Msisdn     string
	Prompt     string
	Parameters map[string]string
}

var location map[string]string

type Menu struct {
	Name string
	Body string
}

func ProcessRequest(r *Request) (string, string) {
	location, err := location[r.Msisdn]

	if !err {
		menu := getMenu("init", r.Parameters)

		return menu.Body, menu.Name
	}

	switch location {
	case "init":
	}

	return "", ""
}
