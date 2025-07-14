package bot

import (
	"embed"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//go:embed menus/*
var menus embed.FS

func getMenu(name string, params map[string]string) Menu {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := fmt.Sprintf("%02d", r.Intn(20))
	filename := fmt.Sprintf("menus/%s/%s_%s.md", name, name, n)
	fmt.Println("Loading menu from:", filename)

	content, _ := menus.ReadFile(filename)

	menu := Menu{}

	menu.Body = string(content)
	menu.Body = strings.Split(menu.Body, "-------------")[0]
	menu.Name = name

	for key, value := range params {
		menu.Body = strings.ReplaceAll(menu.Body, key, value)
	}

	return menu
}

func processInit(r *Request) Menu {
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
