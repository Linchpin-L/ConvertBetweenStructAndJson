package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

func main() {
	fmt.Println("strcut 和 json 互相转换, 直接按回车即可, 将从剪贴板中读取内容并转换")
	fmt.Println("v 0.2.3")
	fmt.Println("by linchpin1029@qq.com")
	var err error

	for {
		l := ""
		_, _ = fmt.Scanln(&l)

		// 从剪贴板获取数据
		all, err2 := clipboard.ReadAll()
		if err2 != nil {
			fmt.Println(err2.Error())
			continue
		}
		fmt.Println("--- 读取 ---")

		var paras string
		// 如果是 json 文件, 那么转换成 struct
		if json.Valid([]byte(all)) {
			fmt.Println("--- json -> struct ---")
			paras, err = jsonToStruct(all)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else if l == "dot" {
			fmt.Println("--- struct -> dot ---")
			paras, err = structToDot(all)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			fmt.Println("--- struct -> json ---")
			paras, err = structToJson(all)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		fmt.Println(paras)

		err = clipboard.WriteAll(paras)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
	}
}

func jsonToStruct(sss string) (string, error) {
	var res interface{}
	json.Unmarshal([]byte(sss), &res)

	plain := "type Apple "
	s, err := parseJsonSub(res)
	if err != nil {
		return "", err
	}
	plain += s

	return plain, nil
}

// 此函数只考虑 值部分
func parseJsonSub(o interface{}) (string, error) {
	line := ""

	switch t := o.(type) {
	case bool:
		line += "bool"
	case float64:
		// todo 0.0 的识别会成为 int
		if o == math.Trunc(o.(float64)) {
			line += "int"
		} else {
			line += "float64"
		}
	case string:
		line += "string"
	case []interface{}:
		if len(t) == 0 {
			line += "[]interface{}"
		} else {
			tt, err := parseJsonSub(t[0])
			if err != nil {
				return "", err
			}
			line += "[]" + tt
		}
	case map[string]interface{}:
		// 这里相当于一个顶级入口

		line += "struct {\n"

		for k, v := range t {
			ttt, err := parseJsonSub(v)
			if err != nil {
				return "", err
			}
			line += fmt.Sprintf("%s %s `json:\"%s\"`", caseToUpper(k), ttt, k)

			// 将原值追加为注释，只追加字符串类型和数值类型
			switch ttt {
			case "string":
				line += " // " + v.(string)
			case "int":
				line += " // " + strconv.Itoa(int(v.(float64)))
			case "float64":
				line += " // " + strconv.FormatFloat(v.(float64), 'f', 2, 64)
			}
			line += "\n"
		}

		line += "}"
	case nil:
		line += "interface{}"
	default:
		return "", fmt.Errorf("不支持的类型: %s, 值为: %s", t, o)
	}

	return line, nil
}

// 输入完整的结构体字符串
//
// 返回格式化好的 json 字符串
func structToJson(content string) (string, error) {
	res := "{\n"
	// fmt.Println(strings.Split(content, "\r\n"))

	// 按行处理
	lines := strings.Split(content, "\n")
	for i, l := 0, len(lines); i < l; i++ {
		line := lines[i]
		// fmt.Printf("处理行: '%s'\n", line)
		line = strings.Trim(line, "\r\n\t ")
		if line == "" {
			continue
		}
		// 如果以其为开头, 说明是中间字段有结构体, 那么跳过这一行的检查
		if strings.HasPrefix(line, "struct") {
			continue
		}
		// 如果包含了 type 开头的, 那么说明是文件头, 删除他们
		if strings.HasPrefix(line, "type") {
			continue
		}
		if strings.HasPrefix(line, "}") {
			continue
		}
		// fmt.Println("->>", o)
		// 处理每一行

		typeWithRemarkStart := "//" // 注释的双斜线加上值类型
		var remarkReq tagBinding    // 备注: 对参数的要求
		key, typ, def, remark, err := parseLine(line)
		// fmt.Printf("分析行:%s\n", line)
		// fmt.Printf("键:%s, 类别:%s, 标识:%s, 备注:%s\n", key, typ, def, remark)
		if err != nil {
			return "", err
		}
		if key == "" { // 如果键为空, 那么忽略这一行。比如行被注释掉了
			continue
		}

		// 确定字段名
		var t tagJson
		if def != "" {
			t, remarkReq, err = findKeyAndRemark(def)
			if err != nil {
				return "", err
			}
			if t.Key != "" {
				key = t.Key
			}
			// 如果 json 指示不输出, 那么忽略改行
			if t.Key == "-" {
				continue
			}
		}
		res += "\"" + key + "\": "

		// 类型
		isArray := false
		if typ[0] == '*' {
			// 如果是指针类型, 那么忽略指针即可
			typ = typ[1:]
		}
		if typ[0] == '[' {
			// 如果是切片, 那么返回值前后加上方括号.
			if typ[1] != ']' {
				return "", errors.New("unknown type")
			}
			isArray = true
			typ = typ[2:]
			res += "["
		}
		switch typ {
		case "int8":
			res += "1"
			typeWithRemarkStart += " int8"
		case "int16":
			res += "1"
			typeWithRemarkStart += " int16"
		case "int32":
			res += "1"
			typeWithRemarkStart += " int32"
		case "int64", "int":
			res += "1"
			typeWithRemarkStart += " int64"
		case "uint8":
			res += "1"
			typeWithRemarkStart += " uint8"
		case "uint16":
			res += "1"
			typeWithRemarkStart += " uint16"
		case "uint32":
			res += "1"
			typeWithRemarkStart += " uint32"
		case "uint64", "uint":
			res += "1"
			typeWithRemarkStart += " uint64"
		case "float32":
			res += "1.1"
			typeWithRemarkStart += " float32"
		case "float64":
			res += "2.64"
			typeWithRemarkStart += " float64"
		case "bool":
			res += "false"
			typeWithRemarkStart += " bool"
		case "string":
			res += "\"test\""
			typeWithRemarkStart += " string"
		case "struct":
			// 先不考虑复杂情况.
			// 先判断这一行是不是有闭合 {}
			if strings.Contains(line, "}") {
				res += "{}"
				typeWithRemarkStart += " struct"
			} else {
				// 向下寻找到 }
				child := "struct {\n" // 还保持行的模式
				for j := i + 1; j < l; j++ {
					currentLine := strings.TrimSpace(lines[j])
					if idx := strings.Index(currentLine, "}"); idx > -1 {
						child += currentLine[:idx]
						child += "\n"
						child += "}"
						// fmt.Println("$1", child, "$2")
						temp, err := structToJson(child)
						if err != nil {
							fmt.Println("PC error:", err)
						}
						// fmt.Println("#1", temp, "#2")
						res += temp
						i = j // 接下来的几行都不用看了
						break
					} else {
						child += currentLine
						child += "\n"
						// fmt.Println("add1", currentLine)
						// fmt.Println("add2", child)
					}
				}
				typeWithRemarkStart += " struct"
			}
		default:
			fmt.Printf("unknown type: \"%s\"\n", typ)
		}
		if isArray {
			res += "]"
		}
		res += ", " + typeWithRemarkStart + "."
		if remarkReq != nil {
			res += " " + remarkReq.value() + "."
		}
		if remark != "" {
			res += " " + remark
		}

		res += "\n"
	}

	res += "}"

	return res, nil
}

// 从一行结构体文本中取出四个部分
//
//	key: 字段名
//	typ: 字段类型
//	def: 标注
//	remark: 备注. 不含 "//". 去重两端空格.
func parseLine(line string) (key, typ, def, remark string, err error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	// 找到注释
	if i := strings.Index(line, "//"); i >= 0 {
		remark = strings.Trim(line[i+2:], " \r\n\t")
		line = line[:i]
	}

	// 找到标注
	if i, j := strings.Index(line, "`"), strings.LastIndex(line, "`"); j > i && i > -1 {
		def = line[i : j+1]
		line = line[:i]
	}

	// 找到第一个非空字符的空字符
	s := -1
	for i := range line {
		if line[i] != ' ' {
			if s == -1 {
				s = i
			}
			continue
		}
		// 是空格
		if s != -1 {
			key = line[s:i]
			typ = strings.Trim(line[i:], " \r\n\t")
			if strings.HasPrefix(typ, "struct") { // typ: "struct {}"
				typ = "struct"
			}
			return
		}
	}

	return
}

type tagJson struct {
	Key       string
	OmitEmpty bool // 是否含有 omitempty
}

type tagBinding []tag

func (b tagBinding) isRequired() bool {
	if b == nil {
		return false
	}

	for _, o := range b {
		if o.Key == "required" {
			return true
		}
	}
	return false
}

// 获取解析前的 binding 规则字符串
func (b tagBinding) value() string {
	if b == nil {
		return ""
	}

	var res string
	for _, o := range b {
		if o.Value == "" {
			res += o.Key + ","
		} else {
			res += o.Key + "=" + o.Value + ","
		}
	}
	return res[:len(res)-1]
}

// 解析 tag 并返回 json 和 binding 的内容
// 没有找到 json 字样时, 返回空
func findKeyAndRemark(tag1 string) (key tagJson, bindings tagBinding, err error) {
	tags, err := parseTag(tag1)
	if err != nil {
		return tagJson{}, nil, err
	}

	for _, o := range tags {
		switch o.Key {
		case "json":
			// 判断值是否含有逗号, 如果有, 则前边为 key, 后边为 omitempty
			if idx := strings.Index(o.Value, ","); idx > -1 {
				key.Key = o.Value[:idx]
				if o.Value[idx+1:] == "omitempty" {
					key.OmitEmpty = true
				}
			} else {
				key.Key = o.Value
			}
		case "binding":
			for _, o := range strings.Split(o.Value, ",") {
				if o == "" {
					continue
				}
				temp := strings.Split(o, "=")
				if len(temp) == 1 {
					bindings = append(bindings, tag{o, ""})
				} else {
					bindings = append(bindings, tag{temp[0], temp[1]})
				}
			}
		}
	}

	return key, bindings, nil
}

type tag struct {
	Key, Value string
}

// 解析 tag 并返回 json 的名字
// json:"123" binding:"required,max=30,oneof=1 2 3"
func parseTag(str string) ([]tag, error) {
	// 先找到第一个不为空格的字符
	// 然后找到第二个引号作为结束, 以此标注他的一部分
	quotaFound := false
	tags := make([]tag, 0)
	l := len(str)
	if l == 0 {
		return tags, nil
	}

	// 如果第一个字符是反引号, 那么忽略他. 最后一个反引号无需处理
	s := -1
	i := 0
	if str[0] == '`' {
		i = 1
	}
	for ; i < l; i++ {
		o := str[i]
		if o != ' ' {
			if s == -1 {
				s = i
				continue
			}
			if o == '"' {
				// 第一个引号
				if !quotaFound {
					quotaFound = true
					continue
				}
				// 第二个引号
				temp := strings.Split(str[s:i+1], ":")
				tags = append(tags, tag{temp[0], temp[1][1 : len(temp[1])-1]})

				s, quotaFound = -1, false
			}
		}
	}

	return tags, nil
}

// 下划线模式转驼峰式
func caseToUpper(camel string) string {
	l := len(camel)
	if l == 0 {
		return camel
	}

	res := make([]byte, l)

	n := true
	ri := 0
	for i := 0; i < l; i++ {
		if camel[i] == '_' {
			n = true
			continue
		}

		if n && camel[i] >= 'a' && camel[i] <= 'z' {
			res[ri] = camel[i] - 32
		} else {
			res[ri] = camel[i]
		}
		n = false
		ri++
	}

	return string(res[0:ri])
}

type dot struct {
	open     bool   // 启用
	name     string // 参数名
	value    string // 参数值
	ty       string // 类型, 如 string
	required bool   // 必须
	remark   string // 说明
}

// 将结构体解析为 dot 模式
// 即: 启用,参数名,参数值,类型,必需,固定参数值,说明
func structToDot(content string) (string, error) {
	dts := make([]*dot, 0)

	// 按行处理
	// 此模式应用于 get 请求, 所以结构体简单, 不涉及嵌套. 且暂不处理数组.
	lines := strings.Split(content, "\n")
	for i, l := 0, len(lines); i < l; i++ {
		line := lines[i]
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 如果包含了 type 开头的, 那么说明是文件头, 删除他们
		if strings.HasPrefix(line, "type") {
			continue
		}
		// 如果以其为开头, 说明是中间字段有结构体, 那么跳过这一行的检查
		if strings.HasPrefix(line, "struct") {
			continue
		}
		if strings.HasPrefix(line, "}") {
			continue
		}

		dt := new(dot)
		var binding tagBinding // 备注: 对参数的要求
		var def, typ string
		var err error
		dt.name, typ, def, dt.remark, err = parseLine(line)
		// fmt.Printf("分析行:%s\n", line)
		// fmt.Printf("键:%s, 类别:%s, 标识:%s, 备注:%s\n", key, typ, def, remark)
		if err != nil {
			return "", err
		}
		if dt.name == "" {
			continue
		}

		// 查询是否必填
		if def != "" {
			_, binding, err = findKeyAndRemark(def)
			if err != nil {
				return "", err
			}
			// 如果 binding 中包含了 required, 说明参数是必需的
			if binding.isRequired() {
				dt.required = true
				dt.open = true
			}
		}

		// 类型
		if typ[0] == '*' {
			// 如果是指针类型, 那么忽略指针即可
			typ = typ[1:]
		}
		remarkNew := "" // 完整的 remark 字符串
		switch typ {
		case "int8":
			dt.value = "1"
			dt.ty = "number"
			remarkNew += "int8"
		case "int16":
			dt.value = "1"
			dt.ty = "number"
			remarkNew += "int16"
		case "int32":
			dt.value = "1"
			dt.ty = "number"
			remarkNew += "int32"
		case "int64", "int":
			dt.value = "1"
			dt.ty = "number"
			remarkNew += "int64"
		case "uint8":
			dt.value = "1"
			dt.ty = "number"
			remarkNew += "uint8"
		case "uint16":
			dt.value = "1"
			dt.ty = "number"
			remarkNew += "uint16"
		case "uint32":
			dt.value = "1"
			dt.ty = "number"
			remarkNew += "uint32"
		case "uint64", "uint":
			dt.value = "1"
			dt.ty = "number"
			remarkNew += "uint64"
		case "float32":
			dt.value = "1.1"
			dt.ty = "number"
			remarkNew += "float32"
		case "float64":
			dt.value = "2.64"
			dt.ty = "number"
			remarkNew += "float64"
		case "bool":
			dt.value = "0"
			dt.ty = "boolean"
		case "string":
			dt.value = "test"
			dt.ty = "string"
		case "struct":
			return "", errors.New("该模式不支持嵌套结构体")
		default:
			return "", fmt.Errorf("unknown type: \"%s\"", typ)
		}

		// 拼装备注. 把备注需要的几个字段拼接起来
		// 在前后加双引号, 避免逗号识别错误, 因为是 csv 格式
		temp := []string{}
		if remarkNew != "" {
			temp = append(temp, remarkNew)
		}
		if t := binding.value(); t != "" {
			temp = append(temp, t)
		}
		if dt.remark != "" {
			temp = append(temp, dt.remark)
		}
		dt.remark = "\"" + strings.Join(temp, ", ") + "\""
		dts = append(dts, dt)
	}

	res := ""
	for _, o := range dts {
		// res += fmt.Sprintf("启用,参数名,参数值,类型,必需,固定参数值,说明\n", o.name, o.ty, o.value)
		res += fmt.Sprintf("%t,%s,%s,%s,%t,,%s\n", o.open, o.name, o.value, o.ty, o.required, o.remark)
	}

	return res, nil
}
