package main

import "net/url"

type Webhook struct {
	Events []Event `json:"events"`
}
type Event struct {
	Action  string  `json:"action"` //push, delete
	Target  Target  `json:"target"`
	Source  Source  `json:"source"`
	Request Request `json:"request"`
}

type Target struct {
	Repository string  `json:"repository"`
	URL        url.URL `json:"url"`
	Tag        string  `json:"tag"` //only on push
}

type Source struct {
	Address string `json:"addr"`
}

type Request struct {
	Host string `json:"host"`
}
