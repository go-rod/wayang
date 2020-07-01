package wayang

import (
	"context"
	"log"

	"github.com/go-rod/rod"
)

type Action map[string]interface{}

type Program struct {
	Selectors map[string]string `json:"selectors"`
	Actions   map[string]Action `json:"actions"`
	Steps     []Action          `json:"steps"`
}

type Runner struct {
	B         *rod.Browser
	P         *rod.Page
	ENV       map[string]interface{}
	Context   context.Context
	Canceller context.CancelFunc
	Logger    *log.Logger
	program   Program
}

type RuntimeError struct {
	parent *Runner
	source string
	stack  []byte
	action Action
	err    interface{}
}
