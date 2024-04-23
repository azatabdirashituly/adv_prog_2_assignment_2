package cmd

import (
	"Ex2_Week3/db"
	"context"
	"encoding/json"
	"net/http"
	
	"gopkg.in/mgo.v2/bson"

	openai "github.com/sashabaranov/go-openai"
)

func sendingMessageHandler(w http.ResponseWriter, r *http.Request) {

	userMessage := r.FormValue("userMessage")
	if !containsKeywords(userMessage) {
		data := struct {
			Error string
		}{
			Error: "Your request was declined because your question is not related to Football.",
		}
		renderTemplate(w, "error.html", data)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var history []openai.ChatCompletionMessage
	
	cursor, err := db.Client.Database("Chat").Collection("history").Find(context.Background(), bson.M{})
    if err != nil {
        http.Error(w, "Failed to fetch chat history: "+err.Error(), http.StatusInternalServerError)
        return
    }
	for cursor.Next(context.Background()) {
        var histEntry History
        if err = cursor.Decode(&histEntry); err != nil {
            http.Error(w, "Failed to decode chat history: "+err.Error(), http.StatusInternalServerError)
            return
        }
        history = append(history, openai.ChatCompletionMessage{
            Role:    openai.ChatMessageRoleUser,
            Content: histEntry.UserMessage,
        }, openai.ChatCompletionMessage{
            Role:    openai.ChatMessageRoleSystem,
            Content: histEntry.GPTResponse,
        })
    }

	history = append(history, openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleUser,
        Content: userMessage,
    })

	client := openai.NewClient("api")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     openai.GPT3Dot5Turbo,
			Messages:  history,
			MaxTokens: 512,
		},
	)
	if err != nil {
		http.Error(w, "Failed to get response from AI", http.StatusInternalServerError)
		return
	}

	aiResponse := resp.Choices[0].Message
	history = append(history, aiResponse)

	historyEntry := History{
		UserMessage: userMessage,
		GPTResponse: aiResponse.Content,
	}
	if _, err := db.Client.Database("Chat").Collection("history").InsertOne(context.Background(), historyEntry); err != nil {
		http.Error(w, "Failed to save chat history: "+err.Error(), http.StatusInternalServerError)
		return
	}

	historyJson, err := json.Marshal(history)
    if err != nil {
        http.Error(w, "Error marshaling history", http.StatusInternalServerError)
        return
    }

	data := struct {
		UserMessage string
		Response    string
		History     string
	}{
		UserMessage: userMessage,
		Response:    aiResponse.Content,
		History:     string(historyJson),
	}

	renderTemplate(w, "home.html", data)
}