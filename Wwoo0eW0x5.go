package main

import (
	"github.com/Sirupsen/logrus"
)

type Hook struct{}

func (h Hook) Fire(e *logrus.Entry) error {
	e.Data["test"] = "changed"
	return nil
}

func (h Hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.Error,
	}
}

func main() {
	l := logrus.New()
	l.Hooks.Add(Hook{})
	l.WithField("test", "test").Error("test")
}
