package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

type (
	ChatCommand string
	Role        string
)

const (
	User Role = "user"
	Lepo Role = "lepo"
)

const (
	Ask    ChatCommand = "/ask"
	Search ChatCommand = "/search"
)

type ChatRequest struct {
	// TODO: Add other metadata once made available
	Query        string                         `json:"query"`
	Conversation []openai.ChatCompletionMessage `json:"conversation"`
	ChatCommand  ChatCommand                    `json:"chat_command"`
	RepoID       string                         `json:"repo_id"`
}

type ChatResponse struct {
	Response     string                         `json:"response"`
	Conversation []openai.ChatCompletionMessage `json:"conversation"`
	RepoID       string                         `json:"repo_id"`
	Sources      []string                       `json:"sources"`
}

func main() {
	conversation := make([]openai.ChatCompletionMessage, 0)
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Chatbot init")
	fmt.Println("---------------------")

	fmt.Println("Conversation")
	fmt.Println("---------------------")

	for {

		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		chatReq := ChatRequest{
			RepoID:       "1",
			Query:        text,
			Conversation: conversation,
		}

		b, err := json.Marshal(chatReq)
		if err != nil {
			panic(err)
		}

		resp, err := http.Post(
			"http://localhost:3000/v1/chat",
			"application/json",
			bytes.NewBuffer(b),
		)
		if err != nil {
			panic(err)
		}

		var chatResp ChatResponse

		body, _ := io.ReadAll(resp.Body)

		err = json.Unmarshal([]byte(body), &chatResp)
		if err != nil {
			fmt.Println("Error unmarshaling data from request.")
			continue
		}

		yellow := color.New(color.FgYellow)
		fmt.Println("-----------")
		yellow.Println(chatResp.Sources)
		fmt.Println("-----------")
		yellow.Println(chatResp.Response)
		conversation = chatResp.Conversation
	}
}
