package types

import (
	"fmt"
	"strings"

	"github.com/gobigbang/validator/utils"
)

type (
	// Interface to be implemented for all the translators
	ITranslator interface {
		Translate(message string, args map[string]interface{}) string
	}

	// A simple translator struct
	Translator struct {
		Messages map[string]string
	}
)

// Creates a new instance of translator
func NewTranslator(messages map[string]interface{}) *Translator {
	t := &Translator{}
	return t.SetMessages(messages)
}

// Set the messages
func (t *Translator) SetMessages(messages map[string]interface{}) *Translator {
	t.Messages = utils.ParseMessages(messages)
	return t
}

// Translates a message
func (t *Translator) Translate(message string, params map[string]interface{}) string {
	m, ok := t.Messages[message]
	if !ok {
		m = message
	}
	args := make([]string, 0)
	if params != nil {
		for k, v := range params {
			args = append(args, "{"+fmt.Sprint(k)+"}")
			args = append(args, fmt.Sprintf("%v", v))
		}
	}
	m = strings.NewReplacer(args...).Replace(m)

	return m
}
