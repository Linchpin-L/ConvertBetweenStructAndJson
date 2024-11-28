package main

import (
	_ "embed"
	"testing"
)

func TestJsonToStruct(t *testing.T) {
	res, err := jsonToStruct(`{"a":1, "b": 1.1, "c":"text", "d":null}`)
	if err != nil {
		t.Errorf("parseJson error: %v", err)
		return
	}
	// 由于函数中使用了 map，这会导致此测试并不总是能通过
	if res != "type Apple struct {\n"+
		"A int `json:\"a\"` // 1\n"+
		"B float64 `json:\"b\"` // 1.10\n"+
		"C string `json:\"c\"` // text\n"+
		"D interface{} `json:\"d\"`\n"+
		"}" {
		t.Errorf("parseJson error: %v", res)
		return
	}
}

//go:embed test/struct.txt
var structExample string

func TestStructToJson(t *testing.T) {
	res, err := structToJson(structExample)
	if err != nil {
		t.Errorf("parseContent error: %v", err)
		return
	}
	if res != `{
"123": 1, // uint64. required,max=30,oneof=1 2 3. 注释
"B": 1, // int64. omitempty,required_if=Field1 foobar. 注释
"cc": 2.64, // float64. required,max=30,oneof=1 2 3.
"D": "test", // string.
"E": ["test"], // string.
"E1": ["test"], // string.
"F": {
"FA": "test", // string.
"FB": 1, // uint64.
}, // struct.
}` {
		t.Errorf("final result wrong: %v", res)
		return
	}
}

//go:embed test/structGet.txt
var structGetExample string

func TestStructToDot(t *testing.T) {
	// 因为 get 请求不支持嵌套, 所以这里无法与上边试用同一个 example
	res, err := structToDot(structGetExample)
	if err != nil {
		t.Errorf("function error: %v", err)
		return
	}
	if res != `true,ID,1,number,true,,"uint64, required"
false,U,1,number,false,,"uint64, required_if=Field1 foobar, [统计需要]如果用户登录, 需要将uid放置其中并传递"
false,Fav,0,boolean,false,,"是否获取收藏信息"
false,E,0,boolean,false,,""
` {
		t.Errorf("final result wrong: %v", res)
		return
	}
}
