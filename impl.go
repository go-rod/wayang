package wayang

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

func (parent *Runner) runAction(act Action, source string) interface{} {
	source = source + "." + act["action"].(string)
	err := func(errs ...interface{}) RuntimeError {
		return RuntimeError{
			parent: parent,
			source: source,
			stack:  debug.Stack(),
			action: act,
			err:    errs,
		}
	}

	if len(source) > 1000 {
		return err("action chain longer than 1000 chars, expected to be inside recursive loop")
	}
	switch act["action"] {
	default:
		actAttr, ok := act["action"].(string)
		if !ok {
			return err("a 'action' key (type string) is required to be present")
		}
		if !strings.HasPrefix(actAttr, "$") {
			return err("could not find a defined action with the requested name")
		}
		value := strings.TrimPrefix(actAttr, "$")
		action, ok := parent.program.Actions[value]
		if !ok {
			return err("could not find a custom action with the requested name")
		}
		return parent.runAction(action, source)

	case "":
		return err("empty action type")

	case "do":
		stmts, ok := act["statements"].([]interface{})
		var res interface{}
		if !ok {
			typed, ok := act["statements"].([]Action)
			if !ok {
				return err("expected statements to be able to be parsed as []Action")
			}

			for _, action := range typed {
				res = parent.runAction(action, source)
				if _, ok = res.(RuntimeError); ok {
					return res
				}
			}
		} else {
			for _, rawAction := range stmts {
				stmt := Action(rawAction.(map[string]interface{}))
				res = parent.runAction(stmt, source)
				if _, ok = res.(RuntimeError); ok {
					return res
				}
			}
		}

		return res

	case "forEach":
		//TODO not implented properly

		action, ok := act["execute"].(Action)
		if !ok {
			return err("expected execute to be able to be parsed as Action")
		}
		elements, ok := act["elements"].(string)
		if !ok {
			return err("could not convert elements to type string")
		}

		for _, element := range parent.P.ElementsX(parent.sel(elements)) {
			action["element"] = element // needs to be the xpath
			if res, ok := parent.runAction(action, source).(RuntimeError); ok {
				return res
			}
		}

	case "if":
		condition := parent.makeAction(act["condition"])
		if condition == nil {
			return err("a condition is required to be present")
		}

		res := parent.runAction(*condition, source)
		toBool, ok := res.(bool)
		if !ok {
			if _, ok := res.(RuntimeError); ok {
				return res
			}
			return err("expected condition to return a bool type, got", res)
		}

		if toBool {
			stmt := parent.makeAction(act["statement"])
			if stmt == nil {
				return err("could not transform execute to action")
			}
			return parent.runAction(*stmt, source)
		}
		action := parent.makeAction(act["otherwise"])
		if action == nil {
			return nil
		}
		return parent.runAction(*action, source)

	case "attribute":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}

		attrName, ok := act["name"].(string)
		if !ok {
			return err("could not find name to retrieve attribute")
		}

		attr := element.Attribute(attrName)
		if attr == nil {
			return nil
		}
		return *attr

	case "html":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		return element.HTML()

	case "text":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		return element.Text()

	case "has":

		attrSel, ok := act["element"].(string)
		if !ok {
			return err("could not find an element key of type string")
		}
		sel := parent.sel(attrSel)
		return parent.P.HasX(sel)

	case "not":
		stmtAttr := parent.makeAction(act["statement"])
		if stmtAttr == nil {
			return err("a statement is required to be present")
		}
		res := parent.runAction(*stmtAttr, source)
		toBool, ok := res.(bool)
		if !ok {
			if _, ok := res.(RuntimeError); ok {
				return res
			}
			return err("expected statement to return a bool type, got", res)
		}
		return !toBool

	case "textEqual":
		expected, ok := act["expected"].(string)
		if !ok {
			return err("an 'expected' key (type string) is required to present")
		}

		stmt := parent.makeAction(act["statement"])
		if stmt == nil {
			return err("an actual key (type action) is required to be present")
		}
		actualRes := parent.runAction(*stmt, source)
		switch actualRes.(type) {
		case RuntimeError:
			return actualRes
		case string:
			break
		default:
			return err("unexpected result from action: ", actualRes)
		}

		return expected == actualRes

	case "textNotEqual":
		expected, ok := act["expected"].(string)
		if !ok {
			return err("an 'expected' key (type string) is required to present")
		}

		stmt := parent.makeAction(act["statement"])
		if stmt == nil {
			return err("an actual key (type action) is required to be present")
		}
		actualRes := parent.runAction(*stmt, source)
		switch actualRes.(type) {
		case RuntimeError:
			return actualRes
		case string:
			break
		default:
			return err("unexpected result from action: ", actualRes)
		}

		return expected != actualRes

	case "visible":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}

		return element.Visible()

	case "blur":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.Blur()
		return nil

	case "clear":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.SelectAllText().Input("")
		return nil

	case "click":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.Click()
		return nil

	case "error":
		message, ok := act["message"].(string)
		if !ok {
			return err("a 'message' key (type string) is required to present")
		}

		parent.Logger.Printf(`level=error msg=%s`, message)
		return err(message)

	case "eval":
		expression, ok := act["expression"].(string)
		if !ok {
			return err("an 'expression' key (type string) is required to be present")
		}

		if _, ok := act["element"]; !ok {
			parent.P.Eval(expression)
			return nil
		}
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.Eval(expression)
		return nil

	case "focus":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.Focus()
		return nil

	case "input":
		text, ok := act["text"].(string)
		if !ok {
			return err("a 'text' key (type string) is required to be present")
		}

		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.Input(text)
		return nil

	case "log":
		message, ok := act["message"].(string)
		if !ok {
			return err("a 'message' key (type string) is required to present")
		}

		parent.Logger.Printf(`level=info msg=%s`, message)
		return nil

	case "logStore":
		keyStmt, ok := act["key"].(string)
		if !ok {
			return err("a 'key' key (type string) is required to present")
		}
		if !strings.HasPrefix(keyStmt, "$") {
			return err("querying the store re")
		}
		key := strings.TrimPrefix(keyStmt, "$")
		exists, ok := parent.ENV[key]
		if !ok {
			return err("the specified key is not in the program memory")
		}
		return exists

	case "navigate":
		link, ok := act["link"].(string)
		if !ok {
			return err("a 'link' key (type string) is required to present")
		}
		parent.P.Navigate(link)
		return nil

	case "press":
		//TODO implement keyboard input stuff

	case "scrollIntoView":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.ScrollIntoView()
		return nil

	case "selectAll":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.SelectAllText()
		return nil

	case "sleep":
		duration, ok := act["duration"].(float64)
		if !ok {
			return err("a 'duration' key (type float) is required to present")
		}
		if duration < 0 {
			return err("duration value must be greater than or equal to 0")
		}

		ctx := parent.P.GetContext()
		second := float64(time.Second)
		t := time.NewTimer(time.Duration(duration * second))

		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
		return nil

	case "store":
		items, ok := act["items"].(map[string]interface{})
		if !ok {
			items = make(map[string]interface{})
		}

		for s, field := range items {
			switch field.(type) {
			case map[string]interface{}:
				source = fmt.Sprintf("%s.field[%s]", source, s)
				res := parent.runAction(field.(map[string]interface{}), source)
				if err, ok := res.(RuntimeError); ok {
					return err
				}
				parent.ENV[s] = res
			case Action:
				source = fmt.Sprintf("%s.field[%s]", source, s)
				res := parent.runAction(field.(Action), source)
				if err, ok := res.(RuntimeError); ok {
					return err
				}
				parent.ENV[s] = res
			case string:
				if !strings.HasPrefix(field.(string), "$") {
					parent.ENV[s] = field
					break
				}

				value := strings.TrimPrefix(field.(string), "$")
				if selector, ok := parent.program.Selectors[value]; ok {
					parent.ENV[s] = selector
					break
				}

				if action, ok := parent.program.Actions[value]; ok {
					res := parent.runAction(action, source)
					if err, ok := res.(RuntimeError); ok {
						return err
					}
					parent.ENV[s] = res
					break
				}
				parent.ENV[s] = field
			default:
				parent.ENV[s] = field
			}
		}
		return nil

	case "waitIdle":
		parent.P.WaitRequestIdle()

	case "waitInvisible":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.WaitInvisible()
		return nil

	case "waitLoad":
		parent.P.WaitLoad()

	case "waitStable":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.WaitStable()
		return nil

	case "waitVisible":
		element, e := parent.createElem(act, err)
		if e != nil {
			return *e
		}
		element.WaitVisible()
	}
	return nil
}

func (parent *Runner) createElem(act Action, err func(...interface{}) RuntimeError) (*rod.Element, *RuntimeError) {
	page := parent.P

	attrElement := act["element"]
	var element *rod.Element
	switch attrElement.(type) {
	case string:
		strElement := attrElement.(string)
		element = page.ElementX(parent.sel(strElement))
	case *rod.Element:
		element = attrElement.(*rod.Element)
	default:
		e := err("could not find element key to retrieve")
		return nil, &e
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

func (parent *Runner) sel(element string) string {
	if strings.HasPrefix(element, "$") {
		return parent.program.Selectors[strings.TrimPrefix(element, "$")]
	}
	return element
}
