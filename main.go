package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

// todo 完成对数组的支持
type temp struct {
	A uint    `json:"123" binding:"required,max=30,oneof=1 2 3"` // 说明一下
	B int     `binding:"omitempty"`
	C float64 `json:"cc" binding:"required,max=30,oneof=1 2 3"`
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

		err = clipboard.WriteAll(paras)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
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

		remarkNew := "//" // 完整的 remark 字符串
		remarkReq := ""   // 备注: 对参数的要求
		key, typ, def, remark, err := parseLine(o)
		if err != nil {
			return "", err
		}

		// 键
		var t string
		if def != "" {
			t, remarkReq, err = findKeyAndRemark(def)
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
		if remarkReq != "" {
			res += remarkReq + ". "
		}
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

// 解析 tag 并返回 json 和 binding 的内容
// 没有找到 json 字样时, 返回空
func findKeyAndRemark(tag string) (string, string, error) {
	tags, err := parseTag(tag)
	if err != nil {
		return "", "", err
	}

	var key, remark string
	for _, o := range tags {
		switch o.Key {
		case "json":
			key = o.Value
		case "binding":
			remark = o.Value
		}
	}

	return key, remark, nil
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
