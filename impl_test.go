package wayang_test

import (
	"time"

	"github.com/ysmood/kit"

	"github.com/go-rod/wayang"
)

func (s *S) TestDo() {
	res, _ := s.singleAction(action(
		"action", "do",
		"statements", []wayang.Action{
			{
				"action": "navigate",
				"link":   "http://example.com/",
			},
			{
				"action":  "has",
				"element": "//div",
			},
		},
	))
	s.Equal(true, res)
}

func (s *S) TestIf() {
	s.page.Navigate("http://example.com")
	res, _ := s.singleAction(action(
		"action", "if",
		"condition", action(
			"action", "has",
			"element", "//a",
		),
		"statement", action(
			"action", "text",
			"element", "//a",
		),
		"otherwise", action(
			"action", "error",
			"message", "expected website to have anchor element",
		),
	))
	s.Equal("More information...", res)
}

func (s *S) TestAttribute() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "attribute",
		"name", "cols",
		"element", "//textarea",
	))
	s.Equal("30", res)

	res, _ = s.singleAction(action(
		"action", "attribute",
		"name", "cols2",
		"element", "//textarea",
	))
	s.Nil(res)
}

func (s *S) TestHTML() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "html",
		"element", "//input[@type='submit']",
	))
	s.Equal(`<input type="submit" value="submit">`, res)
}

func (s *S) TestText() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "text",
		"element", "//input[@type='submit']",
	))
	s.Equal("submit", res)
}

func (s *S) TestHas() {
	s.page.Navigate(srcFile("fixtures/click.html"))

	res, _ := s.singleAction(action(
		"action", "do",
		"statements", []wayang.Action{
			{
				"action":  "click",
				"element": "//button",
			},
			{
				"action":  "has",
				"element": "//button[@a='ok']",
			},
		},
	))
	s.Equal(true, res)
}

func (s *S) TestNot() {
	s.page.Navigate(srcFile("fixtures/click.html"))

	res, _ := s.singleAction(action(
		"action", "not",
		"statement", action(
			"action", "has",
			"element", "//buton[@a='ok']",
		),
	))
	s.Equal(true, res)
}

func (s *S) TestTextEqual() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "textEqual",
		"expected", "submit",
		"statement", action(
			"action", "text",
			"element", "//input[@type='submit']",
		),
	))
	s.Equal(true, res)
}

func (s *S) TestTextNotEqual() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "textNotEqual",
		"expected", "submit",
		"statement", action(
			"action", "text",
			"element", "//input[@type='submit']",
		),
	))
	s.Equal(false, res)
}

func (s *S) TestWaitInvisible() {
	s.page.Navigate(srcFile("fixtures/click.html"))

	res, _ := s.singleAction(action(
		"action", "visible",
		"element", "//h4",
	))

	s.Equal(true, res)

	go func() {
		kit.Sleep(0.03)
		s.page.ElementX("//h4").Eval("() => this.style.visibility = 'hidden'")
	}()

	s.page = s.page.Timeout(time.Second * 3)
	res, _ = s.singleAction(action(
		"action", "waitInvisible",
		"element", "//h4",
	))
	s.Equal(nil, res)
}

func (s *S) TestBlur() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "do",
		"statements", []wayang.Action{
			{
				"action":  "focus",
				"element": "//input[@id='blur']",
			},
			{
				"action":  "blur",
				"element": "//input[@id='blur']",
			},
		},
	))
	el := s.page.ElementX("//input[@id='blur']")
	s.Nil(res)
	s.Equal("ok", *el.Attribute("a"))
}

func (s *S) TestClear() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "do",
		"statements", []wayang.Action{
			{
				"action":  "input",
				"element": "//input[@id='blur']",
				"text":    "test",
			},
			{
				"action":  "clear",
				"element": "//input[@id='blur']",
			},
		},
	))
	el := s.page.ElementX("//input[@id='blur']")
	s.Nil(res)
	s.Empty(el.Text())
}

func (s *S) TestClick() {
	s.page.Navigate(srcFile("fixtures/click.html"))
	el := s.page.ElementX("//button")

	s.Nil(el.Attribute("a"))
	res, _ := s.singleAction(action(
		"action", "click",
		"element", "//button",
	))
	s.Nil(res)
	s.Equal("ok", *el.Attribute("a"))
}

func (s *S) TestError() {
	_, err := s.singleAction(action(
		"action", "error",
		"message", "text",
	))
	s.NotNil(err)
	s.Contains(s.Output.String(), "level=error msg=text")
	s.Equal([]interface{}{"text"}, err.ErrorRaw())
}

func (s *S) TestEval() {
	res, _ := s.singleAction(action(
		"action", "eval",
		"expression", `() => window.location.href = "http://example.com/"`,
	))
	s.page.WaitLoad()
	s.Nil(res)
	s.Equal("http://example.com/", s.pageURL())
}

func (s *S) TestFocus() {
	p := s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "focus",
		"element", "//input",
	))

	p.Keyboard.Press('A')
	p.Keyboard.InsertText(" Test")
	el := p.ElementX("//input")
	s.Nil(res)
	s.Equal("A Test", el.Text())
}

func (s *S) TestInput() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "input",
		"element", "//input",
		"text", "A Test",
	))
	el := s.page.ElementX("//input")
	s.Nil(res)
	s.Equal("A Test", el.Text())
}

func (s *S) TestLog() {
	res, _ := s.singleAction(action(
		"action", "log",
		"message", "text",
	))
	s.Nil(res)
	s.Contains(s.Output.String(), "level=info msg=text")
}

func (s *S) TestNavigate() {
	url := "http://example.com/"
	res, _ := s.singleAction(action(
		"action", "navigate",
		"link", url,
	))
	s.Nil(res)
	s.Equal(url, s.pageURL())
}

func (s *S) TestScrollIntoView() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "scrollIntoView",
		"element", "//form",
	))
	s.Nil(res)
}

func (s *S) TextSelectAll() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.singleAction(action(
		"action", "do",
		"statements", []wayang.Action{
			{
				"action":  "input",
				"element": "//input",
				"text":    "test",
			},
			{
				"action":  "selectAll",
				"element": "//input",
			},
			{
				"action":  "clear",
				"element": "//input",
			},
		},
	))
	s.Nil(res)
	el := s.page.ElementX("//input")
	s.Empty(el.Text())
}

func (s *S) TestCustomSelector() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.execute(`{
		"selectors": {
			"input": "//input"
		},
		"steps": [
			{
				"action": "input",
				"element": "$input",
				"text": "A Test"
			}
		]
	}`)
	s.Nil(res)
	el := s.page.ElementX("//input")
	s.Equal("A Test", el.Text())
}

func (s *S) TestCustomAction() {
	s.page.Navigate(srcFile("fixtures/input.html"))

	res, _ := s.execute(`{
	"selectors": {
		"input": "//input[@id='blur']"
	},
	"actions": {
		"inputBlur": {
			"action": "do",
			"statements": [
				{
					"action": "input",
					"element": "$input",
					"text": "A Test"
				},
				{
					"action": "blur",
					"element": "$input"
				}
			]
		}
	},
	"steps": [
		{
			"action": "$inputBlur"
		}
	]
}`)
	kit.Dump(res)
	s.Nil(res)

	el := s.page.ElementX("//input[@id='blur']")
	s.Equal("A Test", el.Text())
	s.Equal("ok", *el.Attribute("a"))
}
