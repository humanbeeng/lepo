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
	// RepoID string
	// TODO: Add Source[] once made available
}

func main() {
	messages := make([]openai.ChatCompletionMessage, 0)
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
			Conversation: messages,
			ChatCommand:  Ask,
		}
		fmt.Printf("Sending request %+v\n", chatReq)

		b, err := json.Marshal(chatReq)
		if err != nil {
			panic(err)
		}

		resp, err := http.Post("http://localhost:3000/v1/chat", "application/json", bytes.NewBuffer(b))
		if err != nil {
			panic(err)
		}

		var chatResp ChatResponse

		body, _ := io.ReadAll(resp.Body)

		err = json.Unmarshal([]byte(body), &chatResp)
		if err != nil {
			fmt.Println("Error unmarshaling data from request.")
		}

		fmt.Println(chatResp.Response)

		messages = chatResp.Conversation
	}
}
