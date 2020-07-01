package wayang

import (
	"fmt"
	"log"
	"os"

	"github.com/go-rod/rod"
	"github.com/ysmood/kit"
)

func (parent *Runner) RunProgram(program Program) (interface{}, *RuntimeError) {
	parent.program = program

	var res interface{}
	for i, action := range parent.program.Steps {
		source := fmt.Sprintf("root[%d]", i)
		res = parent.runAction(action, source)
		if err, ok := res.(RuntimeError); ok {
			return nil, &err
		}
	}

	return res, nil
}

func RunProgram(program Program) (interface{}, *RuntimeError) {
	browser := rod.New().Connect()
	page := browser.Page("")
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	parent := &Runner{
		B:      browser,
		P:      page,
		Logger: logger,
	}

	return parent.RunProgram(program)
}

func RunActions(actions []Action) (interface{}, *RuntimeError) {
	return RunProgram(Program{
		Steps: actions,
	})
}

func (parent *Runner) RunActions(actions []Action) (interface{}, *RuntimeError) {
	return parent.RunProgram(Program{
		Steps: actions,
	})
}

func RunAction(action Action) (interface{}, *RuntimeError) {
	return RunProgram(Program{
		Steps: []Action{action},
	})
}

func (parent *Runner) RunAction(action Action) (interface{}, *RuntimeError) {
	return parent.RunProgram(Program{
		Steps: []Action{action},
	})
}

func (re *RuntimeError) Action() Action {
	return re.action
}

func (re *RuntimeError) Source() string {
	return re.source
}

func (re *RuntimeError) ErrorRaw() interface{} {
	return re.err
}

func (re *RuntimeError) Error() string {
	return fmt.Sprintln(re.err)
}

func (re *RuntimeError) Dump() string {
	return kit.Sdump(re)
}

func (re *RuntimeError) Log() {
	msg := kit.Sdump(re.err)
	re.parent.Logger.Printf(`level="error" msg="%s"`, msg)
}

func (re *RuntimeError) Print() {
	fmt.Println(kit.Sdump(re.err))
}

func (re *RuntimeError) Stack() string {
	return string(re.stack)
}

func (re *RuntimeError) LogStack() {
	re.parent.Logger.Printf(`level="error" msg="%s"`, string(re.stack))
}

func (re *RuntimeError) PrintStack() {
	fmt.Println(string(re.stack))
}
