package main

import (
	"fmt"

	"github.com/benchambule/sneakerheadz/bot"
)

func main() {
	menu, name := bot.ProcessRequest(&bot.Request{Msisdn: "258849902174", Prompt: ""})

	fmt.Println(menu, name)
}
