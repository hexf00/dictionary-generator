package main

import (
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
)

//LengthRange 重复规则
type LengthRange struct {
	min int
	max int
}

type ruleUnit struct {
	typ         string
	lengthRange LengthRange
	chars       []string
}

// TypeRuleUnitNormal rule Unit Type
var TypeRuleUnitNormal = "normal"

// TypeRuleUnitDict rule Unit Type
var TypeRuleUnitDict = "dict"

// NUMBER 内置序列
var NUMBER []string
var lowercaseLetters []string
var capitalLetter []string

func init() {
	for i := 'a'; i <= 'z'; i++ {
		lowercaseLetters = append(lowercaseLetters, string(i))
	}
	for i := 'A'; i <= 'Z'; i++ {
		capitalLetter = append(capitalLetter, string(i))
	}
	for i := '0'; i <= '9'; i++ {
		NUMBER = append(NUMBER, string(i))
	}
}

func main() {

}

//GetBetweenStr 通过左右两边字符串来获取中间字符串，并返回下标信息
func GetBetweenStr(str, left, right string) (start int, end int, between string, err error) {
	start = strings.Index(str, left)
	if start == -1 {
		err = errors.New("err: need char " + left)
		return
	}

	str = string([]byte(str)[start+len(left):])

	end = strings.Index(str, right)
	if end == -1 {
		err = errors.New("err: need char " + right)
		return
	}

	str = string([]byte(str)[:end])
	return start, start + +len(left) + end + len(right), str, nil
}

//parseNumRuleString 解析字符串里面的1,2这种字符串，如{1,2}
func parseNumRuleString(ruleString string) (lg LengthRange, err error) {

	nums := strings.Split(ruleString, ",")

	if len(nums) == 1 {
		var num int
		num, err = strconv.Atoi(nums[0])
		if err != nil {
			return
		}
		lg = LengthRange{num, num}
	} else if len(nums) == 2 {
		var min, max int
		min, err = strconv.Atoi(nums[0])
		if err != nil && nums[0] != "" {
			return
		}
		max, err = strconv.Atoi(nums[1])
		if err != nil {
			return
		}
		lg = LengthRange{min, max}

	} else {
		err = errors.New("num params number error")
		return
	}
	// if lg.min == 0 {
	// 	lg.min = 1
	// }
	return lg, err
}

//parseCharRuleString 将字典规则语句解析为具体的内容
func parseCharRuleString(ruleString string) (chars []string) {

	// 按 | 分割
	isWord := false

	subRules := strings.Split(ruleString, "|")
	if len(subRules) > 1 {
		isWord = true
	}

	for index := 0; index < len(subRules); index++ {
		rule := subRules[index]
		if strings.Index(rule, "-") != -1 {
			//序列
			if subs := strings.Split(rule, "-"); len(subs) == 2 {
				var min, max int
				var err error
				isNum := true
				min, err = strconv.Atoi(subs[0])
				if err != nil {
					isNum = false
					//utf8.RuneCountInString(subs[0])
					// log.Println(subs[0], subs[0][0], ((subs[0][0] >= 'A' && subs[0][0] <= 'Z') || (subs[0][0] >= 'a' && subs[0][0] <= 'z')))
					if len(subs[0]) == 1 &&
						((subs[0][0] >= 'A' && subs[0][0] <= 'Z') || (subs[0][0] >= 'a' && subs[0][0] <= 'z')) {
						min = int(subs[0][0])
					} else {
						log.Panicln("序列不支持数字和字母以外的符号")
					}

				}

				max, err = strconv.Atoi(subs[1])
				if err != nil {
					isNum = false
					if len(subs[1]) == 1 &&
						((subs[1][0] >= 'A' && subs[1][0] <= 'Z') || (subs[1][0] >= 'a' && subs[1][0] <= 'z')) {
						max = int(subs[1][0])
					} else {
						log.Panicln("序列不支持数字和字母以外的符号")
					}
				}

				//log.Println(isNum, min, max)
				if isNum {
					//数字序列
					if math.Abs(float64(min-max)) > 10000 {
						log.Panicln("序列数量超过限制，请尝试更换表达", math.Abs(float64(min-max)))
					}
					zeroLength := 0

					//小的数第一位是0则补0
					if len(subs[0]) > 1 && subs[0][0] == '0' {
						zeroLength = len(subs[0])
					}
					if len(subs[1]) > 1 && subs[1][0] == '0' {
						if len(subs[1]) >= len(subs[0]) {
							zeroLength = len(subs[1])
						}
					}

					if max > min {
						for index := min; index <= max; index++ {
							if zeroLength > 0 {
								tmp := "000000" + strconv.Itoa(index)
								chars = append(chars, tmp[len(tmp)-zeroLength:])
							} else {
								chars = append(chars, strconv.Itoa(index))
							}
						}
					} else {
						for index := min; index >= max; index-- {
							if zeroLength > 0 {
								tmp := "000000" + strconv.Itoa(index)
								chars = append(chars, tmp[len(tmp)-zeroLength:])
							} else {
								chars = append(chars, strconv.Itoa(index))
							}
						}
					}
				} else {
					//字母序列
					if max > min {
						for index := min; index <= max; index++ {
							chars = append(chars, string(index))
						}
					} else {
						for index := min; index >= max; index-- {
							chars = append(chars, string(index))
						}
					}
				}
			} else {
				log.Panic("序列语法错误", rule)
			}
		} else {
			if isWord {
				//单词切割
				chars = append(chars, subRules[index])
			} else {
				//单字切割
				chars = append(chars, strings.Split(subRules[index], "")...)

			}
		}
	}

	return
}

func parseRuleString(ruleString string) []ruleUnit {
	var result []ruleUnit
	lastString := ""
	for {
		if len(ruleString) == 0 {
			if len(lastString) > 0 {
				result = append(result, ruleUnit{
					typ:   TypeRuleUnitNormal,
					chars: []string{lastString},
				})
				lastString = ""
			}
			break
		}

		if ruleString[:1] == "/" {
			if len(lastString) > 0 {
				result = append(result, ruleUnit{
					typ:   TypeRuleUnitNormal,
					chars: []string{lastString},
				})
				lastString = ""
			}

			if string(ruleString[:2]) == "/d" {
				//ParseNumRuleString(ruleString)

				if len(ruleString) > 2 && ruleString[2:3] == "{" {

					_, end, betwwen, err := GetBetweenStr(ruleString, "{", "}")
					if err != nil {
						log.Panicln(err)
					}
					lg, err := parseNumRuleString(betwwen)
					if err != nil {
						log.Panicln(err)
					}
					ruleString = ruleString[end:]
					result = append(result, ruleUnit{
						typ:         TypeRuleUnitDict,
						lengthRange: lg,
						chars:       NUMBER,
					})
				} else {
					ruleString = ruleString[2:]
					result = append(result, ruleUnit{
						typ:         TypeRuleUnitDict,
						lengthRange: LengthRange{1, 1},
						chars:       NUMBER,
					})
				}
			} else {
				// // /{ /[
				lastString = lastString + ruleString[1:2]
				ruleString = ruleString[2:]
			}
		} else if ruleString[0] == '[' {
			if len(lastString) > 0 {
				result = append(result, ruleUnit{
					typ:   TypeRuleUnitNormal,
					chars: []string{lastString},
				})
				lastString = ""
			}

			_, end, betwwen, err := GetBetweenStr(ruleString, "[", "]")
			if err != nil {
				log.Panicln(err)
			}
			chars := parseCharRuleString(betwwen)
			ruleString = ruleString[end:]

			if len(ruleString) > 1 && ruleString[0:1] == "{" {

				_, end, betwwen, err := GetBetweenStr(ruleString, "{", "}")
				if err != nil {
					log.Panicln(err)
				}
				lg, err := parseNumRuleString(betwwen)
				if err != nil {
					log.Panicln(err)
				}
				ruleString = ruleString[end:]
				result = append(result, ruleUnit{
					typ:         TypeRuleUnitDict,
					lengthRange: lg,
					chars:       chars,
				})
			} else {
				result = append(result, ruleUnit{
					typ:         TypeRuleUnitDict,
					lengthRange: LengthRange{1, 1},
					chars:       chars,
				})
			}
			//log.Println(chars)
			//parse num
			//break
		} else {
			//普通文字
			lastString = lastString + ruleString[:1]
			ruleString = ruleString[1:]
		}

	}

	return result
}

func pow(x, n int) int {
	ret := 1 // 结果初始为0次方的值，整数0次方为1。如果是矩阵，则为单元矩阵。
	for n != 0 {
		if n%2 != 0 {
			ret = ret * x
		}
		n /= 2
		x = x * x
	}
	return ret
}

//CalcCount 计算字典结果条目大小
func CalcCount(ruleString string) (count int) {
	rules := parseRuleString(ruleString)

	count = 0
	for index := 0; index < len(rules); index++ {
		rule := rules[index]

		if rule.typ == TypeRuleUnitDict {
			currCount := 0

			lg := rule.lengthRange

			for index := lg.min; index <= lg.max; index++ {
				currCount += pow(len(rule.chars), index)
			}

			if count == 0 {
				count = currCount
			} else {
				count = count * currCount
			}
		}
	}
	return
}

var rulesCache map[string][]ruleUnit

//CalcString 计算字典结果条目大小
func CalcString(ruleString string, i int) {
	var rules []ruleUnit

	//此处必须缓存，不然太浪费
	if _, ok := rulesCache[ruleString]; ok {
		rules = rulesCache[ruleString]
	} else {
		rules = parseRuleString(ruleString)
		rulesCache = make(map[string][]ruleUnit)
		rulesCache[ruleString] = rules
	}

	//序号
	//i := 11

	//对应的字符串
	text := ""

	//从右到左遍历规则单元
	for index := len(rules) - 1; index >= 0; index-- {
		rule := rules[index]

		//计算当前单元的具体值
		if rule.typ == TypeRuleUnitDict {

			//当前单元最大组合数量
			currCount := 0
			var currRanage [][]int //组合阶梯
			lg := rule.lengthRange
			//因为长度不同，所以我们要把这些都累加在一起，顺序是由短到长
			for index := lg.min; index <= lg.max; index++ {
				currCount += pow(len(rule.chars), index)
				currRanage = append(currRanage, []int{index, pow(len(rule.chars), index)})
			}
			// log.Println("当前单元最大组合数量", currCount)
			// log.Println("当前组合阶梯", currRanage)

			//通过计算余数可得当前单元对应的序列游标
			currIndex := i % currCount
			// log.Println("当前规则单元的游标", currIndex)

			for index := lg.min; index <= lg.max; index++ {
				阶梯长度 := pow(len(rule.chars), index)
				if currIndex >= 阶梯长度 { //这里要用等于判断
					//减掉当前阶梯，进入下一级的阶梯

					currIndex = currIndex - pow(len(rule.chars), index)

					//log.Println("减掉当前阶梯，进入下一级的阶梯", currIndex)
					continue
				} else {
					//无法进入下一级阶梯，计算当前组合

					//log.Println("计算字符长度为", index, "索引为", currIndex, "阶梯长度", 阶梯长度)

					总共长度位数 := index //长度 ，此处还需要补齐长度如000

					// log.Println(index, currIndex)
					tmptext := ""
					for 当前第几位 := 0; 当前第几位 < 总共长度位数; 当前第几位++ {

						//log.Println(currIndex % len(rule.chars))
						tmptext = rule.chars[currIndex%len(rule.chars)] + tmptext

						currIndex = currIndex / len(rule.chars)

					}
					// log.Println(len(tmptext), long)
					if len(tmptext) < 总共长度位数 {
						for i := 0; i < (总共长度位数 - len(tmptext)); i++ {
							tmptext = rule.chars[0] + tmptext
						}
					}

					text = tmptext + text
					// log.Println(i, "总共长度位数", 总共长度位数, "当前表达", tmptext)
					//得出结果后中断停止计算
					break
				}
			}

			i = i / currCount //下一个规则单元循环时候使用

		} else {
			//普通文本直接追加
			text = rules[index].chars[0] + text
		}
	}

	//log.Println(text)

}
