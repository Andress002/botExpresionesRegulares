package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const botName = "ToBeBot"

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Chatbot: Hello! My name is %s. What's your name?\n", botName)
	name := readLine(reader)
	if name == "" {
		name = "friend"
	}

	fmt.Printf("Chatbot: Nice to meet you, %s! Please type an AFFIRMATIVE sentence in English using the verb TO BE in the PAST tense (was/were).\n", name)
	fmt.Println("Chatbot: Example: 'andres was sick last week.' — type 'exit' to stop")

	sentence := ""
	askForSentence := true

	for {
		if askForSentence {
			fmt.Print("User: ")
			sentence = readLine(reader)
		}
		trimmed := strings.TrimSpace(sentence)

		if trimmed == "" {
			fmt.Println("Chatbot: Please type a sentence, or 'exit' to stop.")
			askForSentence = true
			continue
		}
		if strings.EqualFold(trimmed, "exit") || strings.EqualFold(trimmed, "salir") ||
			strings.EqualFold(trimmed, "no") || strings.EqualFold(trimmed, "n") {
			break
		}

		result := Validate(trimmed)
		fmt.Println("Chatbot: " + result.Message)

		fmt.Println("Chatbot: Do you want to try another sentence? (yes/no)")
		fmt.Print("User: ")
		answer := strings.TrimSpace(readLine(reader))
		answerLower := strings.ToLower(answer)

		switch {
		case answerLower == "no" || answerLower == "n" || strings.EqualFold(answer, "salir"):
			fmt.Println("Chatbot: Thanks for practicing! See you next time.")
			return
		case answerLower == "yes" || answerLower == "y" || answerLower == "si" || answerLower == "sí":
			fmt.Println("Chatbot: Great! Type your next sentence.")
			askForSentence = true
		default:
			sentence = answer
			askForSentence = false
		}
	}

	fmt.Println("Chatbot: Thanks for practicing! See you next time.")
}

func readLine(reader *bufio.Reader) string {
	line, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}
	return strings.TrimSpace(line)
}
