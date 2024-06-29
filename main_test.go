package main

import (
	"testing"
)

func TestParseJson(t *testing.T) {
	res, err := parseJson("{\"a\":1,\"b\":2,\"c\":3}")
	if err != nil {
		t.Errorf("parseJson error: %v", err)
		return
	}
	if res != "type Apple struct {\nA int `json:\"a\"`\nB int `json:\"b\"`\nC int `json:\"c\"`\n}" {
		t.Errorf("parseJson error: %v", res)
		return
	}
}

// type temp struct {
// 	A uint    `json:"123" binding:"required,max=30,oneof=1 2 3"` // 说明一下
// 	B int     `binding:"omitempty"`
// 	C float64 `json:"cc" binding:"required,max=30,oneof=1 2 3"`
// 	D *string
// 	E []string
// 	F struct {
// 		FA string
// 		FB uint
// 	}
// 	G uint `json:"-"`
// }

func TestParseContent(t *testing.T) {
	example := "type temp struct {\n"+
	"A uint    `json:\"123\" binding:\"required,max=30,oneof=1 2 3\"` // 说明一下\n"+
	"B int     `binding:\"omitempty\"`\n"+
	"C float64 `json:\"cc\" binding:\"required,max=30,oneof=1 2 3\"`\n"+
	"D *string\n"+
	"E []string\n"+
	"F struct {\n"+
		"FA string\n"+
		"FB uint\n"+
	"}\n"+
	"G uint `json:\"-\"`\n"+
"}"

	res, err := parseContent(example)
	if err != nil {
		t.Errorf("parseContent error: %v", err)
		return
	}
	if res != `{
"123": 1, // uint64. required,max=30,oneof=1 2 3.  说明一下
"B": 1, // int64. omitempty. 
"cc": 2.64, // float64. required,max=30,oneof=1 2 3. 
"D": "test", // string. 
"E": ["test"], // string. 
"F": {
"FA": "test", // string. 
"FB": 1, // uint64. 
}, // struct. 
}` {
		t.Errorf("parseContent error: %v", res)
		return
	}
}