package whatsapp

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/benchambule/sneakerheadz/bot"
	_ "github.com/mattn/go-sqlite3"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"google.golang.org/protobuf/proto"
)

func extractProductMT(input string) (string, float64, bool) {
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", " ")    // Replace non-breaking spaces with regular spaces
	input = strings.ReplaceAll(input, " - ", " ")  // Remove hyphens surrounded by spaces
	input = strings.ReplaceAll(input, "MT", " MT") // Ensure "MT" is separated by a space for easier parsing

	if !strings.HasSuffix(input, "MT") {
		return "", 0, false
	}

	parts := strings.Fields(input)
	if len(parts) < 2 {
		return "", 0, false
	}

	if !strings.HasSuffix(parts[len(parts)-1], "MT") {
		return "", 0, false
	}

	re := regexp.MustCompile(`^\d+$`)
	price := strings.Replace(parts[len(parts)-2], ",", "", -1)

	productPrice, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return "", 0, false
	}

	if !re.MatchString(price) {
		return "", 0, false
	}

	productName := strings.Join(parts[:len(parts)-2], " ")

	return productName, productPrice, true

}

func GetEventHandler(client *whatsmeow.Client) func(interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			if v.Info.IsFromMe {
				fmt.Println("Message from me. Will not process")
				return
			}

			fmt.Println("----------------------------NEW MESSAGE---------------------------")
			fmt.Println("Chat: ", v.Info.Chat)
			fmt.Println("User: ", v.Info.Chat.User)
			fmt.Println("RawAgent: ", v.Info.Chat.RawAgent)
			fmt.Println("Device: ", v.Info.Chat.Device)
			fmt.Println("Integrator: ", v.Info.Chat.Integrator)
			fmt.Println("Server: ", v.Info.Chat.Server)
			fmt.Println("Id: ", v.Info.ID)
			fmt.Println("Sender: ", v.Info.Sender)
			fmt.Println("Name: ", v.Info.PushName)
			fmt.Println("Time: ", v.Info.Timestamp)
			fmt.Println("----------------------------BEGIN---------------------------------")

			if v.Message.Conversation != nil {
				var input = v.Message.GetConversation()
				fmt.Printf("Message: %s\n", v.Message.GetConversation())

				product_name, product_price, ok := extractProductMT(input)

				if !ok {
					fmt.Printf("[%s], is not a valid product name with price in MT format", input)
					return
				}

				client_name := v.Info.PushName
				fifty_percent := product_price / 2
				now := time.Now()

				request := &bot.Request{
					Msisdn: "258849902174",
					Prompt: "1",
					Parameters: map[string]string{
						"{{client_name}}":   client_name,
						"{{product_name}}":  product_name,
						"{{product_price}}": fmt.Sprintf("%.2f", product_price),
						"{{fifty_percent}}": fmt.Sprintf("%.2f", fifty_percent),
						"{{book_date}}":     now.Format("02/01/2006"),
						"{{delivery_date}}": now.Add(14 * 24 * time.Hour).Format("02/01/2006"),
					},
				}

				resp := bot.ProcessRequest(request)

				client.MarkRead([]types.MessageID{v.Info.ID}, v.Info.Timestamp.Add(time.Second), v.Info.Chat, v.Info.Sender)
				client.SendChatPresence(v.Info.Chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
				time.Sleep(time.Second * 5)

				client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{Conversation: proto.String(resp.Body)})

				client.SendChatPresence(v.Info.Chat, types.ChatPresencePaused, types.ChatPresenceMediaText)
			}

			var extend = v.Message.GetExtendedTextMessage()
			if extend != nil {
				if extend.ContextInfo.QuotedMessage == nil {
					fmt.Println("No quoted message found in the extended text message.")
					return
				}

				if extend.ContextInfo.QuotedMessage.Conversation == nil {
					fmt.Println("No conversation found in the quoted message.")
					return
				}

				input := extend.ContextInfo.QuotedMessage.GetConversation()
				product_name, product_price, ok := extractProductMT(input)

				if !ok {
					fmt.Printf("[%s], is not a valid product name with price in MT format", input)
					return
				}

				client_name := v.Info.PushName
				fifty_percent := product_price / 2
				now := time.Now()

				request := &bot.Request{
					Msisdn: "258849902174",
					Prompt: "1",
					Parameters: map[string]string{
						"{{client_name}}":   client_name,
						"{{product_name}}":  product_name,
						"{{product_price}}": fmt.Sprintf("%.2f", product_price),
						"{{fifty_percent}}": fmt.Sprintf("%.2f", fifty_percent),
						"{{book_date}}":     now.Format("02/01/2006"),
						"{{delivery_date}}": now.Add(14 * 24 * time.Hour).Format("02/01/2006"),
					},
				}

				resp := bot.ProcessRequest(request)

				client.MarkRead([]types.MessageID{v.Info.ID}, v.Info.Timestamp.Add(time.Second), v.Info.Chat, v.Info.Sender)
				client.SendChatPresence(v.Info.Chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
				time.Sleep(time.Second * 5)

				client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{Conversation: proto.String(resp.Body)})

				client.SendChatPresence(v.Info.Chat, types.ChatPresencePaused, types.ChatPresenceMediaText)
			}

			if v.Message.ReactionMessage != nil {
				fmt.Println("Reaction Message:", v.Message.GetReactionMessage())
				fmt.Println("Message Text:", v.Message.ReactionMessage.Text)
				fmt.Println("Message ID:", v.Message.ReactionMessage.GetKey().GetID())
				fmt.Println("Participant:", v.Message.ReactionMessage.GetKey().GetParticipant())
				fmt.Println("Message GetText:", v.Message.ReactionMessage.GetText())
				fmt.Println("Message GetGroupingKey:", v.Message.ReactionMessage.GetGroupingKey())
			}
		}
	}
}
