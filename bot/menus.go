package bot

import (
	"embed"
	"strings"
)

//go:embed menus/*.md
var menus embed.FS

func getMenu(name string, params map[string]string) Menu {
	content, _ := menus.ReadFile("menus/" + name + ".md")

	menu := Menu{}

	menu.Body = string(content)
	menu.Name = name

	for key, value := range params {
		menu.Body = strings.ReplaceAll(menu.Body, key, value)
	}

	return menu
}

func ProcessInit(r *Request) Menu {
	if r == nil {
		panic("Request cannot be null")
	}

	switch prompt := r.Prompt; prompt {
	case "1":
		return getMenu("001_talktoadmin", r.Parameters)
	case "2":
		return getMenu("002_recommend", r.Parameters)
	case "3":
		return getMenu("003_info", r.Parameters)
	case "4":
		return getMenu("004_evaluate", r.Parameters)
	default:
		params := r.Parameters
		params["error"] = "Por favor seleccione uma das opcoes abaixo"
		return getMenu("000_init", params)
	}
}
