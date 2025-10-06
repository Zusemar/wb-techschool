// Применимость - когда необходимо чтобы клиент мог поьлзоваться методами интерфейса сервиса,
// если это напрямую невозможно по разным причинам
//
// Плюсы - теперь клиент !совместим! с сервисом, с которым был не совместим
//
// Минусы - новые абстракции, классы
//
// Реальные примеры использования - когда нужно уметь работать с внешней либой,
// но мой код этого не может условно по причине разных форматов.
// Хороший пример - один сервис выдает данные в xml а сервис с анализом данных
// работает только с json. Вот тут и нужно написать адаптер переводящий данные из xml->json

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
