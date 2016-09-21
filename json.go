package main

import (
	"encoding/json"
	"log"
	"strings"
)

/*
{
    "color": "green",
    "message": "It's going to be sunny tomorrow! (yey)",
    "notify": false,
    "message_format": "text"
}
*/

func jsonFormatReply(color string, text string) string {
	mapD := make(map[string]interface{})
	mapD["color"] = color
	mapD["notify"] = false
	mapD["message_format"] = "text"
	mapD["message"] = text
	data, _ := json.Marshal(mapD)
	return string(data)
}

type hipchatMessage struct {
	node    string
	name    string
	command string
	args    []string
}

func parseInputJson(data []byte) hipchatMessage {
	var out hipchatMessage
	var f interface{}
	err := json.Unmarshal(data, &f)

	if err != nil {
		log.Println("Malformed json input")
		return out
	}

	item := f.(map[string]interface{})["item"]
	m1 := item.(map[string]interface{})["message"]
	from := m1.(map[string]interface{})["from"]
	name := from.(map[string]interface{})["mention_name"]
	m2 := m1.(map[string]interface{})["message"]

	args := strings.Split(m2.(string), " ")

	switch {
	case len(args) == 0: //should never happen
	case len(args) == 1:
		out.node = args[0]
		out.name = name.(string)
		out.command = ""
		out.args = make([]string, 0)
	case len(args) == 2:
		out.node = args[0]
		out.name = name.(string)
		out.command = args[1]
		out.args = make([]string, 0)
	case len(args) > 2:
		out.node = args[0]
		out.name = name.(string)
		out.command = args[1]
		out.args = args[2:]
	}
	return out
}
