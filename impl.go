package wayang

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

type runtimeAction struct {
	runner *Runner
	act    Action
	source string
}

type actionFunc func(ra runtimeAction, act Action) interface{}

var actions map[string]actionFunc

func init() {
	actions = map[string]actionFunc{
		"do":             doAction,
		"forEach":        forEachAction,
		"if":             ifAction,
		"store":          storeAction,
		"attribute":      attributeAction,
		"html":           htmlAction,
		"text":           textAction,
		"has":            hasAction,
		"not":            notAction,
		"textContains":   textContainsAction,
		"textEqual":      textEqualAction,
		"textNotEqual":   textNotEqualAction,
		"visible":        visibleAction,
		"blur":           blurAction,
		"clear":          clearAction,
		"click":          clickAction,
		"error":          errorAction,
		"eval":           evalAction,
		"focus":          focusAction,
		"input":          inputAction,
		"log":            logAction,
		"logStore":       logStoreAction,
		"navigate":       navigateAction,
		"press":          pressAction,
		"scrollIntoView": scrollIntoViewAction,
		"selectAll":      selectAllAction,
		"sleep":          sleepAction,
		"waitIdle":       waitIdleAction,
		"waitInvisible":  waitInvisibleAction,
		"waitLoad":       waitLoadAction,
		"waitStable":     waitStableAction,
		"waitVisible":    waitVisibleAction,
	}
}

func (parent *Runner) runAction(act Action, source string) interface{} {
	source = source + "." + act["action"].(string)

	ra := runtimeAction{
		runner: parent,
		act:    act,
		source: source,
	}

	if len(source) > 1000 {
		return ra.err("action chain longer than 1000 chars, expected to be inside recursive loop")
	}

	action, ok := act["action"].(string)
	if !ok {
		return ra.err("could not convert action to type string")
	}

	actFunc, ok := actions[action]
	if ok {
		return actFunc(ra, act)
	}
	if !strings.HasPrefix(action, "$") {
		return ra.err("could not find a defined action with the requested name (" + source + ")")
	}

	cAction := strings.TrimPrefix(action, "$")
	cActionFunc, ok := parent.program.Actions[cAction]
	if !ok {
		return ra.err("could not find a custom action with the requested name (" + source + ")")
	}
	return parent.runAction(cActionFunc, source)
}

func doAction(ra runtimeAction, act Action) interface{} {
	run := ra.runner
	stmts, ok := act["statements"].([]interface{})
	var res interface{}
	if !ok {
		typed, ok := act["statements"].([]Action)
		if !ok {
			return ra.err("expected statements to be able to be parsed as []Action")
		}

		for _, action := range typed {
			res = run.runAction(action, ra.source)
			if _, ok = res.(RuntimeError); ok {
				return res
			}
		}
	} else {
		for _, rawAction := range stmts {
			stmt := Action(rawAction.(map[string]interface{}))
			res = run.runAction(stmt, ra.source)
			if _, ok = res.(RuntimeError); ok {
				return res
			}
		}
	}

	return res
}

func forEachAction(ra runtimeAction, act Action) interface{} {
	//TODO not implemented properly
	run := ra.runner

	action, ok := act["execute"].(Action)
	if !ok {
		return ra.err("expected execute to be able to be parsed as Action")
	}
	elements, ok := act["elements"].(string)
	if !ok {
		return ra.err("could not convert elements to type string")
	}

	sel, ok := run.sel(elements)
	if !ok {
		return ra.err("could not find a custom selector defined with the specified value")
	}

	for _, element := range run.P.ElementsX(sel) {
		action["element"] = element // needs to be the xpath
		if res, ok := run.runAction(action, ra.source).(RuntimeError); ok {
			return res
		}
	}
	return nil
}

func ifAction(ra runtimeAction, act Action) interface{} {
	run := ra.runner

	condition := run.makeAction(act["condition"])
	if condition == nil {
		return ra.err("a condition is required to be present")
	}

	res := run.runAction(*condition, ra.source)
	toBool, ok := res.(bool)
	if !ok {
		if _, ok := res.(RuntimeError); ok {
			return res
		}
		return ra.err("expected condition to return a bool type, got", res)
	}

	if toBool {
		stmt := run.makeAction(act["statement"])
		if stmt == nil {
			return ra.err("could not transform execute to action")
		}
		return run.runAction(*stmt, ra.source)
	}
	action := run.makeAction(act["otherwise"])
	if action == nil {
		return nil
	}
	return run.runAction(*action, ra.source)
}

func storeAction(ra runtimeAction, act Action) interface{} {
	run := ra.runner

	items, ok := act["items"].(map[string]interface{})
	if !ok {
		items = make(map[string]interface{})
	}

	source := ra.source
	for s, field := range items {
		switch item := field.(type) {
		case map[string]interface{}:
			source = fmt.Sprintf("%s.field[%s]", source, s)
			res := run.runAction(item, source)
			if err, ok := res.(RuntimeError); ok {
				return err
			}
			run.ENV[s] = res

		case Action:
			source = fmt.Sprintf("%s.field[%s]", source, s)
			res := run.runAction(item, source)
			if err, ok := res.(RuntimeError); ok {
				return err
			}
			run.ENV[s] = res

		case string:
			if !strings.HasPrefix(item, "$") {
				run.ENV[s] = item
				break
			}

			value := strings.TrimPrefix(item, "$")
			if selector, ok := run.program.Selectors[value]; ok {
				run.ENV[s] = selector
				break
			}

			if action, ok := run.program.Actions[value]; ok {
				res := run.runAction(action, source)
				if err, ok := res.(RuntimeError); ok {
					return err
				}
				run.ENV[s] = res
				break
			}
			run.ENV[s] = item

		default:
			run.ENV[s] = item
		}
	}

	return nil
}

func attributeAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	attrName, ok := act["name"].(string)
	if !ok {
		return ra.err("could not find name to retrieve attribute")
	}

	attr := element.Attribute(attrName)
	if attr == nil {
		return nil
	}
	return *attr
}

func htmlAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	return element.HTML()
}

func textAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	return element.Text()
}

func hasAction(ra runtimeAction, act Action) interface{} {
	run := ra.runner

	attrSel, ok := act["element"].(string)
	if !ok {
		return ra.err("could not find an element key of type string")
	}

	sel, ok := run.sel(attrSel)
	if !ok {
		return ra.err("could not find a custom selector defined with the specified value")
	}
	return run.P.HasX(sel)
}

func notAction(ra runtimeAction, act Action) interface{} {
	run := ra.runner

	stmtAttr := run.makeAction(act["statement"])
	if stmtAttr == nil {
		if asBool, ok := act["statement"].(bool); ok {
			return !asBool
		}

		return ra.err("a statement is required to be present")
	}

	res := run.runAction(*stmtAttr, ra.source)
	toBool, ok := res.(bool)
	if !ok {
		if _, ok := res.(RuntimeError); ok {
			return res
		}
		return ra.err("expected statement to return a bool type, got", res)
	}

	return !toBool
}

func textContainsAction(ra runtimeAction, act Action) interface{} {
	run := ra.runner

	text, ok := act["text"].(string)
	if !ok {
		return ra.err("a 'text' key (type string) is required to present")
	}

	stmt := run.makeAction(act["statement"])
	if stmt == nil {
		return ra.err("a statement key (type action) is required to be present")
	}

	stmtRes := run.runAction(*stmt, ra.source)
	switch stmtRes.(type) {
	case RuntimeError:
		return stmtRes
	case string:
		break
	default:
		return ra.err("unexpected result from action: ", stmtRes)
	}
	stmtStr := stmtRes.(string)

	ignoreCase, ok := act["ignoreCase"].(bool)
	if !ok {
		ignoreCase = false
	}
	if ignoreCase {
		stmtStr = strings.ToLower(stmtStr)
		text = strings.ToLower(text)
	}

	return strings.Contains(stmtStr, text)
}

func textEqualAction(ra runtimeAction, act Action) interface{} {
	run := ra.runner

	expected, ok := act["expected"].(string)
	if !ok {
		return ra.err("an 'expected' key (type string) is required to present")
	}

	stmt := run.makeAction(act["statement"])
	if stmt == nil {
		return ra.err("an actual key (type action) is required to be present")
	}

	actualRes := run.runAction(*stmt, ra.source)
	switch actualRes.(type) {
	case RuntimeError:
		return actualRes
	case string:
		break
	default:
		return ra.err("unexpected result from action: ", actualRes)
	}
	actual := actualRes.(string)

	ignoreCase, ok := act["ignoreCase"].(bool)
	if !ok {
		ignoreCase = false
	}
	if ignoreCase {
		expected = strings.ToLower(expected)
		actual = strings.ToLower(actual)
	}

	return expected == actual
}

func textNotEqualAction(ra runtimeAction, act Action) interface{} {
	res := textEqualAction(ra, act)
	if boolRes, ok := res.(bool); ok {
		return !boolRes
	}

	return res
}

func visibleAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	return element.Visible()
}

func blurAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.Blur()
	return nil
}

func clearAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.SelectAllText().Input("")
	return nil
}

func clickAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.Click()
	return nil
}

func errorAction(ra runtimeAction, act Action) interface{} {
	message, ok := act["message"].(string)
	if !ok {
		return ra.err("a 'message' key (type string) is required to present")
	}

	ra.runner.Logger.Printf(`level=error msg=%s`, message)
	return ra.err(message)
}

func evalAction(ra runtimeAction, act Action) interface{} {
	expression, ok := act["expression"].(string)
	if !ok {
		return ra.err("an 'expression' key (type string) is required to be present")
	}

	if _, ok := act["element"]; !ok {
		return ra.runner.P.Eval(expression).Raw
	}

	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	return element.Eval(expression).Raw
}

func focusAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.Focus()
	return nil
}

func inputAction(ra runtimeAction, act Action) interface{} {
	text, ok := act["text"].(string)
	if !ok {
		return ra.err("a 'text' key (type string) is required to be present")
	}

	element, err := ra.createElem(act)
	if err != nil {
		ra.runner.P.Keyboard.InsertText(text)
		return nil
	}
	element.Input(text)
	return nil
}

func logAction(ra runtimeAction, act Action) interface{} {
	message, ok := act["message"].(string)
	if !ok {
		return ra.err("a 'message' key (type string) is required to present")
	}

	ra.runner.Logger.Printf(`level=info msg=%s`, message)
	return nil
}

func logStoreAction(ra runtimeAction, act Action) interface{} {
	keyStmt, ok := act["key"].(string)
	if !ok {
		return ra.err("a 'key' key (type string) is required to present")
	}
	if !strings.HasPrefix(keyStmt, "$") {
		return ra.err("querying the store requires a $ prefix")
	}

	key := strings.TrimPrefix(keyStmt, "$")
	exists, ok := ra.runner.ENV[key]
	if !ok {
		return ra.err("the specified key is not in the program store")
	}
	return exists
}

func navigateAction(ra runtimeAction, act Action) interface{} {
	link, ok := act["link"].(string)
	if !ok {
		return ra.err("a 'link' key (type string) is required to present")
	}
	ra.runner.P.Navigate(link)
	return nil
}

func pressAction(ra runtimeAction, act Action) interface{} {
	keyAttr, ok := act["key"].(string)
	var press rune
	if !ok {
		return ra.err("a 'key' key (type string) is required to be present")
	}

	if len(keyAttr) == 1 {
		press = []rune(keyAttr)[0]
	} else {
		for rne, key := range input.Keys {
			if strings.ToLower(key.Code) == strings.ToLower(keyAttr) {
				press = rne
				break
			}
		}
	}

	element, err := ra.createElem(act)
	if err != nil {
		ra.runner.P.Keyboard.Press(press)
		return nil
	}
	element.Press(press)
	return nil
}

func scrollIntoViewAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.ScrollIntoView()
	return nil
}

func selectAllAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.SelectAllText()
	return nil
}

func sleepAction(ra runtimeAction, act Action) interface{} {
	duration, ok := act["duration"].(float64)
	if !ok {
		return ra.err("a 'duration' key (type float) is required to present")
	}
	if duration < 0 {
		return ra.err("duration value must be greater than or equal to 0")
	}

	ctx := ra.runner.P.GetContext()
	asDuration := time.Duration(duration)
	t := time.NewTimer(asDuration * time.Second)

	select {
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	case <-t.C:
	}
	return nil
}

func waitIdleAction(ra runtimeAction, _ Action) interface{} {
	run := ra.runner
	done := make(chan bool)
	ctx := run.P.GetContext()

	go func() {
		run.P.WaitRequestIdle()()
		done <- true
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func waitInvisibleAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.WaitInvisible()
	return nil
}

func waitLoadAction(ra runtimeAction, _ Action) interface{} {
	ra.runner.P.WaitLoad()
	return nil
}

func waitStableAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.WaitStable()
	return nil
}

func waitVisibleAction(ra runtimeAction, act Action) interface{} {
	element, err := ra.createElem(act)
	if err != nil {
		return *err
	}
	element.WaitVisible()
	return nil
}

func (ra runtimeAction) err(errs ...interface{}) RuntimeError {
	return RuntimeError{
		parent: ra.runner,
		source: ra.source,
		stack:  debug.Stack(),
		action: ra.act,
		err:    errs,
	}
}

func (ra runtimeAction) createElem(act Action) (*rod.Element, *RuntimeError) {
	run := ra.runner

	attrElement := act["element"]
	var element *rod.Element
	switch typed := attrElement.(type) {
	case string:
		strElement := typed
		sel, ok := run.sel(strElement)
		if !ok {
			err := ra.err("could not find a custom selector defined with the specified value")
			return nil, &err
		}

		element = run.P.ElementX(sel)
	case *rod.Element:
		element = typed
	default:
		err := ra.err("could not find element key to retrieve")
		return nil, &err
	}
	return element, nil
}

func (parent *Runner) makeAction(act interface{}) *Action {
	switch act.(type) {
	case map[string]interface{}:
		action := Action(act.(map[string]interface{}))
		return &action
	case Action:
		action := act.(Action)
		return &action
	default:
		return nil
	}
}

func (parent *Runner) sel(element string) (string, bool) {
	if strings.HasPrefix(element, "$") {
		res, ok := parent.program.Selectors[strings.TrimPrefix(element, "$")]
		return res, ok
	}
	return element, true
}
