package main

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func init() {

}

func TestGetBetweenStr(t *testing.T) {
	var err error
	var str string
	var start, end int
	var between string

	//成功情况
	str = "a{123}b"
	start, end, between, err = GetBetweenStr(str, "{", "}")
	if start != 1 || end != 6 || between != "123" || err != nil {
		t.Error("测试失败")
	}

	//成功情况
	str = "a{{123}}b"
	start, end, between, err = GetBetweenStr(str, "{{", "}}")
	if start != 1 || end != 8 || between != "123" || err != nil {
		t.Error("测试失败")
	}

	//错误情况
	str = "a123}a"
	_, _, between, _ = GetBetweenStr(str, "", "")
	if between != "" {
		t.Error("测试失败")
	}
	str = "a123}a"
	_, _, _, err = GetBetweenStr(str, "{", "}")
	if err == nil {
		t.Error("测试失败")
	}
	str = "a{123a"
	_, _, _, err = GetBetweenStr(str, "{", "}")
	if err == nil {
		t.Error("测试失败")
	}
}

func TestParseNumRuleString(t *testing.T) {
	var err error
	var v LengthRange
	v, err = parseNumRuleString("1")
	if err != nil || v.min != 1 || v.max != 1 {
		t.Error("测试失败")
	}
	v, err = parseNumRuleString("1,2")
	if err != nil || v.min != 1 || v.max != 2 {
		t.Error("测试失败")
	}

	v, err = parseNumRuleString(",2")
	if err != nil || v.min != 1 || v.max != 2 {
		t.Error("测试失败")
	}

	//错误情况
	v, err = parseNumRuleString("1,")
	if err == nil {
		t.Error("测试失败")
	}
	v, err = parseNumRuleString("1,2,3")
	if err == nil {
		t.Error("测试失败")
	}
	v, err = parseNumRuleString("")
	if err == nil {
		t.Error("测试失败")
	}
	v, err = parseNumRuleString("a")
	if err == nil {
		t.Error("测试失败")
	}

	//小数

	//负数

	//从长到短
}

func Test手机号案例(t *testing.T) {
	rules := parseRuleString("130/d{8}")
	if len(rules) != 0 {
		t.Error("测试失败")
	}
}

func Test解析空字符(t *testing.T) {
	rules := parseRuleString("")
	if len(rules) != 0 {
		t.Error("测试失败")
	}
}

// /d
func Test解析数字和普通规则(t *testing.T) {
	var rules []ruleUnit

	rules = parseRuleString("/d")

	if len(rules) != 1 ||
		rules[0].typ != TypeRuleUnitDict ||
		strings.Join(rules[0].chars, "") != strings.Join(NUMBER, "") {
		t.Error("测试失败，数据错误")
	}

	rules = parseRuleString("//d")

	log.Println(rules)
	if len(rules) != 1 ||
		rules[0].typ != TypeRuleUnitNormal ||
		len(rules[0].chars) != 1 ||
		rules[0].chars[0] != "/d" {
		t.Error("测试失败，数据错误")
	}

	rules = parseRuleString("/d{2}")

	log.Println(rules)
	if len(rules) != 1 ||
		rules[0].typ != TypeRuleUnitDict ||
		strings.Join(rules[0].chars, "") != strings.Join(NUMBER, "") ||
		rules[0].lengthRange.min != 2 ||
		rules[0].lengthRange.max != 2 {
		t.Error("测试失败，数据错误")
	}

	rules = parseRuleString("/d{1,2}")

	log.Println(rules)
	if len(rules) != 1 ||
		rules[0].typ != TypeRuleUnitDict ||
		strings.Join(rules[0].chars, "") != strings.Join(NUMBER, "") ||
		rules[0].lengthRange.min != 1 ||
		rules[0].lengthRange.max != 2 {
		t.Error("测试失败，数据错误")
	}

	rules = parseRuleString("/d{,2}")
	log.Println(rules)
	if len(rules) != 1 ||
		rules[0].typ != TypeRuleUnitDict ||
		strings.Join(rules[0].chars, "") != strings.Join(NUMBER, "") ||
		rules[0].lengthRange.min != 1 ||
		rules[0].lengthRange.max != 2 {
		t.Error("测试失败，数据错误")
	}
}

func Test解析数字和普通规则2(t *testing.T) {
	var rules []ruleUnit
	rules = parseRuleString("a/d{1}")
	log.Println(rules)
	if len(rules) != 2 ||
		rules[0].typ != TypeRuleUnitNormal ||
		len(rules[0].chars) != 1 ||
		rules[0].chars[0] != "a" ||
		rules[1].typ != TypeRuleUnitDict ||
		strings.Join(rules[1].chars, "") != strings.Join(NUMBER, "") ||
		rules[1].lengthRange.min != 1 ||
		rules[1].lengthRange.max != 1 {
		t.Error("测试失败，数据错误")
	}
}

func Test解析数字和普通规则3(t *testing.T) {
	var rules []ruleUnit

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
	rules = parseRuleString("a/d{1")
	log.Println(rules)

	t.Error("测试失败，错误")
}

func Test时间字符串(t *testing.T) {
	rules := parseRuleString("[00-23]:[00-59]:[00-59]")

	t.Log(rules)
	if len(rules) != 2 {
		t.Error("测试失败")
	}
}
func Test解析自定义字符组合(t *testing.T) {
	CalcString("[one|two|three]{2}", 4)

	t.Error("测试失败")
}
func Test解析自定义字符序列1(t *testing.T) {
	rules := parseRuleString("[a-z]")

	t.Log(len(rules))
	if len(rules) != 2 {
		t.Error("测试失败")
	}
}
func Test解析自定义字符序列2(t *testing.T) {
	rules := parseRuleString("[A-Z]")

	t.Log(len(rules))
	if len(rules) != 2 {
		t.Error("测试失败")
	}
}

func Test解析自定义字符序列3(t *testing.T) {
	rules := parseRuleString("[0-9]")

	t.Log(rules)
	if len(rules) != 2 {
		t.Error("测试失败")
	}
}

func Test解析自定义字符序列4(t *testing.T) {
	rules := parseRuleString("[130-139]{0,1}")

	t.Log(rules)
	if len(rules) != 2 {
		t.Error("测试失败")
	}
}
func TestParseNumRuleString2(t *testing.T) {
	var chars []string
	chars = parseCharRuleString("甲乙丙")
	if len(chars) != 3 {
		log.Println(chars)
		t.Error("")
	}
	chars = parseCharRuleString("甲乙|丙")
	log.Println(chars)
	if len(chars) != 2 {
		log.Println(chars)
		t.Error("")
	}
}
func TestParseNumRuleString1(t *testing.T) {
	var chars []string

	chars = parseCharRuleString("0-1")
	log.Println(chars)
	chars = parseCharRuleString("0-9")
	log.Println(chars)
	chars = parseCharRuleString("130-139")
	log.Println(chars)
	chars = parseCharRuleString("a-z")
	log.Println(chars)
	chars = parseCharRuleString("01-12") //首位补0
	log.Println(chars)
	chars = parseCharRuleString("12-01") //首位补0
	log.Println(chars)
	//chars = ParseCharRuleString("1300000-1390000") //超大序列，太大
	//log.Println(chars)
	// chars = ParseCharRuleString("0-9a-e") //错误语法
	// log.Println(chars)
	// chars = ParseCharRuleString("一-九")
	// log.Println(chars)
	// chars = ParseCharRuleString("!-*")
	// log.Println(chars)

	chars = parseCharRuleString("	")
	log.Println(chars)
	t.Error("..")

}

func TestParseCharRuleString(t *testing.T) {
	//a-z
	//0-9
	//A-Z
	//a-z0-9A-Z-_!@#$%&*()=

	//[0-9|a-e]  不那么重要的特性
	//[0123456789abcde]

	//[130-139]/d{8} 不那么重要的特性
	//[13000000000-14000000000] 不支持这样的写法，子串数量太大了

	//[甲乙丙]
	//[one|two|three]

	//[0-24][0-60][0-60]不那么重要的特性

	//词组和单字的界定,有`|`以这个分割为准，而不是字符

}

//性能测试
func Benchmark性能测试(b *testing.B) {
	for index := 0; index < b.N; index++ {
		CalcString("[0-9|a-z|A-Z]{4}", index)
	}
}

func TestCalcString(t *testing.T) {

	for index := 0; index < 128; index++ {
		CalcString("code[100-999]", index)
		// log.Println("----------")
	}
	// for index := 0; index < 65; index++ {
	// 	CalcString("	", index)
	// 	//log.Println("----------")
	// }
	t.Error("测试失败")

}
func TestCalcCount(t *testing.T) {
	//count := CalcCount("[00-23]:[00-59]:[00-59]")
	var count int

	count = CalcCount("/d")
	t.Log(count)
	if count != 10 {
		t.Error("测试失败")
	}

	count = CalcCount("/d{,1}")
	t.Log(count)
	if count != 11 {
		t.Error("测试失败")
	}

	count = CalcCount("/d{1,2}")
	t.Log(count)
	if count != 90 {
		t.Error("测试失败")
	}

	count = CalcCount("/d{0,2}")
	t.Log(count)
	if count != 90 {
		t.Error("测试失败")
	}

	count = CalcCount("([01]{2,3})")
	t.Log(count)
	if count != 90 {
		t.Error("测试失败")
	}
}

//TODO 多线程
//字符串
