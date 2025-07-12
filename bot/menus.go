// Package bot
package bot

import (
	"embed"
	"strings"
)

var menus embed.FS

func getMenu(name string, params map[string]string) Menu {
	content, _ := menus.ReadFile("menus/" + name)

	menu := Menu{}

	menu.Body = string(content)
	menu.Name = name

	for key, value := range params {
		strings.ReplaceAll(menu.Body, key, value)
	}

	return menu
}

func ProcessInit(r *Request) Menu {
	if r == nil {
		panic("Request cannot be null")
	}

	switch prompt := r.Prompt; prompt {
	case "1":
		return getMenu("buy", r.Parameters)
	case "2":
		return getMenu("info", r.Parameters)
	case "3":
		return getMenu("recommend", r.Parameters)
	default:
		params := r.Parameters
		params["error"] = "Por favor seleccione uma das opcoes abaixo"
		return getMenu("init", params)
	}
}

func ProcessBuy(r *Request) Menu {
	return Menu{}
}
