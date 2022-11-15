package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type Label struct {
	Name  string
	Value string
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

func main() {
	flag.Parse()
	fmt.Println("Hi from the test")
	fmt.Printf("Your project: %+v\n", project)
	fmt.Printf("Your task: %+v\n", task)
	if task == "" {
		panic("task parameter cannot be empty")
	}
	fmt.Printf("Your labels: %+v\n", labels)
	action, err := getAction()
	if err != nil {
		fmt.Printf("Action error: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Your action: %+v\n", action)

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
