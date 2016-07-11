package main

import (
	&#34;github.com/Sirupsen/logrus&#34;
)

type Hook struct{}

func (h Hook) Fire(e *logrus.Entry) error {
	e.Data[&#34;test&#34;] = &#34;changed&#34;
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
	l.WithField(&#34;test&#34;, &#34;test&#34;).Error(&#34;test&#34;)
}
