package main

type Webhook struct {
	Events []Event `json:"events"`
}
type Event struct {
	Action string `json:"action"` //push, delete
	Target Target `json:"target"`
}

type Target struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"` //only on push
}
