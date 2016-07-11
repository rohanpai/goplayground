package main

import (
	&#34;os&#34;
	&#34;os/exec&#34;

	&#34;gopkg.in/qml.v1&#34;
)

func main() {
	if err := qml.Run(run); err != nil {
		panic(err)
	}
}

type Foo struct{}

// DoClick runs the command &#34;echo FOO!&#34; and copies it&#39;s stdout to
// this programs&#39; stdout.
func (Foo) DoClick() {
	cmd := exec.Command(&#34;echo&#34;, &#34;FOO!&#34;)
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
	ctx.SetVar(&#34;foo&#34;, Foo{})

	obj, err := engine.LoadString(&#34;foo.qml&#34;, QML)
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
		text: &#34;DoClick&#34;
		onClicked: foo.doClick()
	}

	Button {
		anchors.bottom: parent.bottom
		anchors.right: parent.right
		text: &#34;GoDoClick&#34;
		onClicked: foo.goDoClick()
	}
}
`
