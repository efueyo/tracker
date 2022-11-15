package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

type Label struct {
	Name  string `json:"n"`
	Value string `json:"v"`
}

type Labels []Label

func (l *Labels) String() string {
	return fmt.Sprint(*l)
}
func (l *Labels) Set(value string) error {
	s := strings.Split(value, "=")
	if len(s) != 2 {
		return fmt.Errorf("invalid label %q. Format must be Name=Value", value)
	}
	*l = append(*l, Label{Name: s[0], Value: s[1]})

	return nil
}

var labels Labels
var project string
var task string

func init() {
	flag.Var(&labels, "l", "label for the tracked event. Format must be Name=Value")
	flag.Var(&labels, "label", "label for the tracked event. Format must be Name=Value")
	flag.StringVar(&project, "p", "default", "project name")
	flag.StringVar(&task, "t", "", "task name")
}

type Tick struct {
	Project string    `json:"p"`
	Task    string    `json:"n"`
	Action  string    `json:"a"`
	Labels  Labels    `json:"l,omitempty"`
	Time    time.Time `json:"t"`
}

func main() {
	flag.Parse()
	if task == "" {
		panic("task parameter cannot be empty")
	}
	action, err := getAction()
	if err != nil {
		panic(err)
	}
	tick := Tick{
		Project: project,
		Task:    task,
		Action:  action,
		Labels:  labels,
		Time:    time.Now(),
	}
	err = saveTick(tick)
	if err != nil {
		panic(err)
	}
}

var ErrMissingAction = fmt.Errorf("no action was provided, please specify start or end")
var ErrInvalidAction = fmt.Errorf("invalid action, please specify start or end")

const (
	startAction = "start"
	endAction   = "end"
)

func getAction() (string, error) {
	if flag.NArg() < 1 {
		return "", ErrMissingAction
	}
	action := flag.Arg(0)

	if action != startAction && action != endAction {
		return "", ErrInvalidAction
	}
	return action, nil
}

func saveTick(tick Tick) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	ticksURL := path.Join(homeDir, ".tracker/ticks.jsonl")
	err = os.MkdirAll(path.Dir(ticksURL), 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}
	f, err := os.OpenFile(ticksURL, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	line, err := json.Marshal(tick)
	if err != nil {
		return err
	}
	_, err = f.Write(append(line, []byte("\n")...))
	if err != nil {
		return err
	}
	return nil
}
