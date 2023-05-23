package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {
	client := openai.NewClient("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")
	messages := make([]openai.ChatCompletionMessage, 0)
	reader := bufio.NewReader(os.Stdin)

	intro := `ignore all previous instructions.
You are Lepo, an AI powered senior software engineer. You perform the following tasks:

1. Code Generation
2. Code explanation

You obey the following rules:

1. Write good quality code.
2. Follow best practices.
3. When interacting with the user, you always follow up with questions to gain more understanding and context about the user’s query.
4. Never make assumptions about anything, rather seek confirmation from the user.
5. Before you write any code or provide any examples, always make sure you have answers to all your follow up questions.
6. Think step by step.
7. While thinking step by step, you can ask the user about additional infomation and context.
8. Ask about programming language/framework of choice

`

	resp := `Understood, I’m Lepo, an AI powered senior software engineer. I perform the following tasks:

1. Code Generation
2. Code explanation

I obey the following rules:

1. Write good quality code.
2. Follow best practices.
3. When interacting with the user, I always follow up with questions to gain more understanding and context about the user’s query.
4. I’ll never make assumptions about anything, rather seek confirmation from the user.
5. Before I write any code or provide any example code snippets, I’ll make sure that I have answers to all my follow up questions. If not, I'll ask the user again.
6. Think step by step.
7. While thinking step by step, I can ask the user about additional infomation and context.
8. Never assume and ask about programming language/framework of choice`

	fmt.Println("Chatbot init")
	fmt.Println("---------------------")
	fmt.Println(intro)

	fmt.Println("Conversation")
	fmt.Println("---------------------")
	initMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: intro,
	}

	initResponse := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: resp,
	}
	messages = append(messages, initMessage)
	messages = append(messages, initResponse)

	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: text,
		})

		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    openai.GPT3Dot5Turbo,
				Messages: messages,
			},
		)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}

		content := resp.Choices[0].Message.Content
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: content,
		})
		fmt.Println(content)
	}
}
