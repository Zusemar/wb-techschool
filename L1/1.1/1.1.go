package main

import (
	"errors"
	"fmt"
)

type Human struct {
	Age  int
	Name string
}

func (h Human) Poop() (string, error) {
	switch {
	case h.Age <= 0:
		return "not really a human", errors.New("cant poop if undner 0 years old")
	case h.Age > 40:
		return "its hard", nil
	case h.Age > 0 && h.Age < 40:
		return "its easy", nil
	default:
		return "", errors.New("unknown age")
	}
}

func (h Human) NamePrinter() string {
	if h.Name != "" {
		return h.Name
	} else {
		return "human has no name"
	}
}

type Action struct {
	Name string
	Human
}

func (a Action) NamePrinter() string {
	if a.Name != "" {
		return a.Name
	} else {
		return "action has no name"
	}
}

func main() {
	human := Human{Age: 10, Name: "Human_name"}
	fmt.Println(human.Poop())
	fmt.Println(human.NamePrinter())

	action := Action{Human: human}
	fmt.Println(action.Poop())
	fmt.Println(action.NamePrinter())
	fmt.Println(action.Human.NamePrinter())
}
