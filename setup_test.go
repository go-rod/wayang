package wayang_test

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/stretchr/testify/suite"
	"github.com/ysmood/kit"

	"github.com/go-rod/wayang"
)

var slash = filepath.FromSlash

type S struct {
	suite.Suite

	browser *rod.Browser
	page    *rod.Page

	Logger *log.Logger
	Output *bytes.Buffer
}

func Test(t *testing.T) {
	s := new(S)

	url := launcher.New().Headless(true).Launch()
	s.browser = rod.New().ControlURL(url).Timeout(time.Second * 5).Connect()
	defer s.browser.Close()

	s.page = s.browser.Page("")
	s.page.Viewport(800, 600, 1, false)

	s.Output = bytes.NewBufferString("")
	s.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	s.Logger.SetOutput(s.Output)

	suite.Run(t, s)
}

func (s *S) execute(test string) (interface{}, *wayang.RuntimeError) {
	parent := &wayang.Runner{
		B:      s.browser,
		P:      s.page,
		Logger: s.Logger,
	}
	program := wayang.Program{}
	kit.E(json.Unmarshal([]byte(test), &program))

	return parent.RunProgram(program)
}

func (s *S) singleAction(action wayang.Action) (interface{}, *wayang.RuntimeError) {
	parent := &wayang.Runner{
		B:      s.browser,
		P:      s.page,
		Logger: s.Logger,
	}

	res, err := parent.RunAction(action)
	if err != nil {
		kit.Dump(err.ErrorRaw())
	}
	return res, err
}

func action(args ...interface{}) wayang.Action {
	action := wayang.Action{}
	for i := 0; i < len(args); i += 2 {
		action[args[i].(string)] = args[i+1]
	}

	return action
}

func (s *S) pageURL() string {
	info, err := proto.TargetGetTargetInfo{TargetID: s.page.TargetID}.Call(s.browser)
	kit.E(err)
	return info.TargetInfo.URL
}

// get abs file path from fixtures folder, return sample "file:///a/b/click.html"
func srcFile(path string) string {
	return "file://" + file(path)
}

// get abs file path from fixtures folder, return sample "/a/b/click.html"
func file(path string) string {
	f, err := filepath.Abs(slash(path))
	kit.E(err)
	return f
}
