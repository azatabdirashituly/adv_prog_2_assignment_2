package cmd

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	openai "github.com/sashabaranov/go-openai"
)

var store = sessions.NewCookieStore([]byte("something-secret"))

func init() {
	// Error saving session: securecookie:
	// error - caused by: securecookie: error - caused by: gob: type not registered for interface: []openai.ChatCompletionMessage
	gob.Register([]openai.ChatCompletionMessage{})
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge: -1,
	}
}

func sendingMessageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "chat-session")
	if err != nil {
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		return
	}

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

	if session.IsNew {
		session.ID = uuid.New().String()
		session.Values["history"] = []openai.ChatCompletionMessage{}
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var history []openai.ChatCompletionMessage
	if h, ok := session.Values["history"].([]openai.ChatCompletionMessage); ok {
		history = h
	}

	userMsg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	}
	history = append(history, userMsg)

	client := openai.NewClient("api-key")
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
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	history = append(history, resp.Choices[0].Message)

	session.Values["history"] = history
	if err = session.Save(r, w); err != nil {
		http.Error(w, "Error saving session: "+err.Error(), http.StatusInternalServerError)
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
		Response:    resp.Choices[0].Message.Content,
		History:     string(historyJson),
	}

	renderTemplate(w, "home.html", data)
}
