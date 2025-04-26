package main

import (
	"fmt"

	"github.com/benchambule/lottus"
)

func main() {
	app := lottus.New("main", lottus.InMemorySessionStorage{Sessions: map[string]lottus.Session{}})

	app.At("main",
		func(r *lottus.Request, res lottus.Message) lottus.Message {
			return lottus.Message{
				Text: "Hello, what's your name",
				Input: lottus.TextInput{
					NextMessage:   "final",
					ParameterName: "name",
				},
			}
		},
		lottus.DefaultProcessor(app),
	)

	app.At("final",
		func(r *lottus.Request, res lottus.Message) lottus.Message {
			return lottus.Message{
				Text: "Hello {{name}}",
				Input: lottus.TextInput{
					NextMessage: "main",
				},
			}
		},
		lottus.DefaultProcessor(app),
	)

	_, msg := app.ProcRequest(lottus.Request{Msisdn: "25884XXXXXXX", Prompt: "Hello World"})
	fmt.Println(msg)

	_, msg = app.ProcRequest(lottus.Request{Msisdn: "25884XXXXXXX", Prompt: "Hello World"})
	fmt.Println(msg)
}
