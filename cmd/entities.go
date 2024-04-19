package cmd

type History struct {
    UserMessage string `bson:"userMessage"`
    GPTResponse string `bson:"gptResponse"`
}