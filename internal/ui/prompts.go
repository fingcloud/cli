package ui

import (
	"github.com/AlecAivazis/survey/v2"
)

func PromptInput(msg string, v interface{}) error {
	p := &survey.Input{
		Message: msg,
	}

	return survey.AskOne(p, v)
}

func PromptEmail(v interface{}) error {
	return PromptInput("Enter your email", v)
}

func PromptPassword(v interface{}) error {
	p := &survey.Password{
		Message: "Enter your password:",
	}

	return survey.AskOne(p, v)
}
