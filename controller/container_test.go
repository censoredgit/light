package controller

import (
	"fmt"
	"testing"
)

type printer interface {
	Print(msg string)
}

type printService struct {
	lastMsg string
}

func (t *printService) ReturnLastMsg() string {
	return t.lastMsg
}

func (t *printService) Print(msg string) {
	t.lastMsg = msg
	fmt.Println(msg)
}

func TestContainerInjectSuccess(t *testing.T) {
	container := Container{}
	err := container.Singleton(&printService{})
	if err != nil {
		t.Error(err)
	}

	container.inject(func(service *printService, printer printer) *Action {
		msg := "Hello world!!!"
		printer.Print(msg)
		if service.ReturnLastMsg() != msg {
			t.Error("ReturnLastMsg fail")
		}

		return &Action{}
	})
}

func TestContainerInjectNonStructFail(t *testing.T) {
	cont := Container{}
	a := 4
	err := cont.Singleton(&a)
	if err == nil {
		t.Error(err)
	}
}
