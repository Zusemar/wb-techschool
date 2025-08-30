package main

import "testing"

func TestNamePrinter_embedding(t *testing.T) {
	human := Human{Age: 55, Name: "human_name"}
	action := Action{Name: "action_name", Human: human}

	if human.NamePrinter() == action.NamePrinter() {
		t.Errorf("NamePrinter error")
	}

	human1 := Human{Age: 55, Name: "human_name"}
	action1 := Action{Human: human1}

	if human1.NamePrinter() == action1.NamePrinter() {
		t.Errorf("NamePrinter error")
	}

	if human1.NamePrinter() != action1.Human.NamePrinter() {
		t.Errorf("NamePrinter error")
	}
}

func TestPoop_embedding(t *testing.T) {
	human := Human{Age: 55, Name: "human_name"}
	action := Action{Name: "action_name", Human: human}

	valHuman, _ := human.Poop()
	valAction, _ := action.Poop()
	if valHuman != valAction {
		t.Errorf("Poop error")
	}
}
