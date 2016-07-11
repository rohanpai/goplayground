package main

import (
	"os"
	"os/exec"

	"gopkg.in/qml.v1"
)

func main() {
	if err := qml.Run(run); err != nil {
		panic(err)
	}
}

type Foo struct{}

// DoClick runs the command "echo FOO!" and copies it's stdout to
// this programs' stdout.
func (Foo) DoClick() {
	cmd := exec.Command("echo", "FOO!")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// GoDoClick runs DoClick in a new goroutine
func (f Foo) GoDoClick() {
	go f.DoClick()
}

func run() error {
	engine := qml.NewEngine()

	ctx := engine.Context()
	ctx.SetVar("foo", Foo{})

	obj, err := engine.LoadString("foo.qml", QML)
	if err != nil {
		return err
	}

	window := obj.CreateWindow(nil)
	window.Show()
	window.Wait()

	return nil
}

const QML = `
import QtQuick 2.0
import QtQuick.Controls 1.2

Rectangle {
	width: 300
	height: 300

	Button {
		anchors.top: parent.top
		anchors.left: parent.left
		text: "DoClick"
		onClicked: foo.doClick()
	}

	Button {
		anchors.bottom: parent.bottom
		anchors.right: parent.right
		text: "GoDoClick"
		onClicked: foo.goDoClick()
	}
}
`
