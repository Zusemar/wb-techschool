package main

import (
	"fmt"
	"strings"
)

type client struct {
	Data string
}

type clientInterface interface {
	poop(str string) string
}

type service struct {
	Data []rune
}

type serviceInterface interface {
	piss(str []rune) []rune
}

type ClientToServiceAdapter struct {
	service serviceInterface
}

// реализуем интерфейс клиента в адаптере
func (adapter *ClientToServiceAdapter) poop(str string) string {
	res := adapter.service.piss([]rune(str))
	return string(res)
}

func (client *client) poop(str string) string {
	var builder strings.Builder
	builder.WriteString(str)
	builder.WriteString(" жопа")
	return builder.String()
}

func (service *service) piss(str []rune) []rune {
	str = append(str, []rune(" попа")...)
	return str
}

func main() {
	input := "тест"

	c := &client{}

	var d clientInterface
	d = &ClientToServiceAdapter{service: &service{}}

	fmt.Println("Client:", c.poop(input))
	fmt.Println("Adapter:", d.poop(input))
}
