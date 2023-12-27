package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

// todo 完成对数组的支持
type temp struct {
	A uint `json:"123" binding:"required,max=30"` // 说明一下
	B int  `binding:"omitempty"`
	C float64
	D string
	E string
}

func main() {
	fmt.Println("strcut 和 json 互相转换, 直接按回车即可, 将从剪贴板中读取内容并转换")
	fmt.Println("by linchpin1029@qq.com")

	l := 0
	for {
		_, _ = fmt.Scanln(&l)

		// 从剪贴板获取数据
		all, err2 := clipboard.ReadAll()
		if err2 != nil {
			fmt.Println(err2.Error())
			continue
		}
		fmt.Println("--- 读取 ---")
		fmt.Println(all)

		paras, err := parseContent(all)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("--- 输出 ---")
		fmt.Println(paras)
		// format, count := convertToPostmanFormat(paras)

		err = clipboard.WriteAll(paras)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		// fmt.Println(format + strconv.Itoa(count))
	}
}

func parseContent(content string) (string, error) {
	res := "{\n"
	// fmt.Println(strings.Split(content, "\r\n"))

	for _, o := range strings.Split(content, "\r\n") {

		o = strings.Trim(o, "\r\n\t ")
		// fmt.Println("->", o)
		if strings.HasPrefix(o, "type") || strings.HasPrefix(o, "}") {
			continue
		}
		// fmt.Println("->>", o)
		// 处理每一行

		remarkNew := "//"
		key, typ, def, remark, err := parseLine(o)
		if err != nil {
			return "", err
		}

		// 键
		if def != "" {
			t, err := parseTag(def)
			if err != nil {
				return "", err
			}
			if t != "" {
				key = t
			}
		}
		res += "\"" + key + "\": "

		// 类型
		switch typ {
		case "int8":
			res += "1"
			remarkNew += " int8"
		case "int16":
			res += "1"
			remarkNew += " int16"
		case "int32":
			res += "1"
			remarkNew += " int32"
		case "int64", "int":
			res += "1"
			remarkNew += " int64"
		case "uint8":
			res += "1"
			remarkNew += " uint8"
		case "uint16":
			res += "1"
			remarkNew += " uint16"
		case "uint32":
			res += "1"
			remarkNew += " uint32"
		case "uint64", "uint":
			res += "1"
			remarkNew += " uint64"
		case "float32":
			res += "1.1"
			remarkNew += " float32"
		case "float64":
			res += "2.64"
			remarkNew += " float64"
		case "bool":
			res += "false"
			remarkNew += " bool"
		case "string":
			res += "\"test\""
			remarkNew += " string"
		}
		res += ", " + remarkNew + ". "
		if len(remark) > 2 {
			res += remark[2:]
		}

		res += "\n"
	}

	res += "}"

	return res, nil
}

// 从一行结构体文本中取出四个部分
func parseLine(line string) (key, typ, def, remark string, err error) {
	// 找到注释
	if i := strings.Index(line, "//"); i >= 0 {
		remark = strings.TrimLeft(line[i:], " \r\n\t")
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
			return
		}
	}

	err = errors.New("invalid line: " + line)

	return
}

// 解析 tag 并返回 json 的名字
// 没有找到 json 字样时, 返回空
func parseTag(tag string) (string, error) {
	i := strings.Index(tag, "json:\"")
	if i == -1 {
		return "", nil
	}

	// 从 i 处开始取第二个结束标识
	found := false
	j := i + 6
	for l := len(tag); j < l; j++ {
		if tag[j] == '"' {
			found = true
			break
		}
	}
	if !found || i+6 == j {
		return "", errors.New("invalid tag: " + tag)
	}
	tag = tag[i+6 : j]
	param := strings.Split(tag, ",")
	return param[0], nil
}
