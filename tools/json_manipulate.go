package tools

import (
	//"fmt"
	"strconv"
	"strings"
)

type Token struct {
	start int
	end   int
}

type Node struct {
	token Token
	a     []*Node
}

func (node *Node) Insert(j ADJson, token Token, level int) {
	////fmt.Printf("before clear %v\n", j.jsonString[token.start:token.end])

	node.token = j.ADClear(token)
	////fmt.Printf("after clear  %v\n", j.jsonString[node.token.start:node.token.end])

	level++
	tokens := j.ADSplit(node.token, level)
	////fmt.Printf("len tokens: %d\n", len(tokens))
	for i, tok := range tokens {
		////fmt.Printf("level=%d i=%d tok: %v %s\n", level, i, tok, j.jsonString[tok.start:tok.end])
		////fmt.Printf("      %d   %d      %v %v\n", level, i, tok, j.arrbrackets[tok.start:tok.end])
		node.a = append(node.a, &Node{tok, nil})
		node.a[i].Insert(j, tok, level)
	}

}

type ADJson struct {
	jsonByte    []byte
	jsonString  string
	brackets    []byte
	arrbrackets []byte
}

type ADFind struct {
	start_idx   int
	end_idx     int
	start_value int
	end_value   int
}

func (j *ADJson) ADClear(token Token) Token {
	level := 0
	for _, x := range j.jsonString[token.start:token.end] {
		if string(x) == "[" {
			level++
		} else if string(x) == "]" {
			level--
		}
	}
	if level > 0 {
		l := level
		for k := token.start; k <= token.end; k++ {
			if string(j.jsonString[k]) == "[" {
				l--
				token.start = k + 1
			}
			if l == 0 {
				break
			}
		}
	}
	if level < 0 {
		l := level
		for k := token.end - 1; k >= token.start; k-- {
			if string(j.jsonString[k]) == "]" {
				l++
				token.end = k
			}
			if l == 0 {
				break
			}
		}
	}
	return token
}

func (j *ADJson) ADSplit(token Token, level int) []Token {
	if level > 5 {
		return nil
	}
	b_min := 100
	for k := token.start; k < token.end; k++ {
		if string(j.jsonString[k]) == "," {
			if int(j.arrbrackets[k]) < b_min {
				b_min = int(j.arrbrackets[k])
				////fmt.Printf("b_min=%d %v\n", b_min, string(j.jsonString[k]))
			}
		}
	}
	if b_min == 100 {
		return nil
	}
	tokens := make([]Token, 0)
	start := token.start
	end := token.start
	for k := token.start; k <= token.end; k++ {

		if (int(j.arrbrackets[k]) == b_min && string(j.jsonString[k]) == ",") || (k == token.end) {
			end = k
			if start < end {
				tokens = append(tokens, Token{start, end})
			}
			start = end + 1
		}
	}
	return tokens
}

func (j *ADJson) SetByte(jsonByte []byte) {
	j.jsonByte = jsonByte
	j.Byte2String()
	j.DeepCalculate()
	j.ArrCalculate()
}

func (j *ADJson) SetJson(jsonString string) {
	j.jsonString = jsonString
	j.String2Byte()
	j.DeepCalculate()
	j.ArrCalculate()
}

func (j *ADJson) DeepCalculate() {
	j.brackets = make([]byte, len(j.jsonString))
	level := 0
	for idx, x := range j.jsonString {
		if string(x) == "{" {
			level++
		} else if string(x) == "}" {
			level--
		}
		j.brackets[idx] = byte(level)
	}
}

func (j *ADJson) ArrCalculate() {
	j.arrbrackets = make([]byte, len(j.jsonString))
	level := 0
	for idx, x := range j.jsonString {
		if string(x) == "[" {
			level++
		} else if string(x) == "]" {
			level--
		}
		j.arrbrackets[idx] = byte(level)
	}
}

func (j *ADJson) Byte2String() {
	j.jsonString = string(j.jsonByte)
}

func (j *ADJson) String2Byte() {
	j.jsonByte = []byte(j.jsonString)
}

func (j *ADJson) ReadElement(path string, type_opt ...string) string {
	typ := "string"
	if len(type_opt) > 0 {
		typ = type_opt[0]
	}

	adFind := j.FindElement(path, typ)
	return j.jsonString[adFind.start_value:adFind.end_value]
}
func (j *ADJson) SetElement(path string, value string, type_opt ...string) {

	bull := false
	numer := false
	typ := "string"
	if len(type_opt) > 0 {
		switch type_opt[0] {
		case "string":
			bull = false
			numer = false
			typ = "string"
		case "number":
			bull = false
			numer = true
			typ = "number"
		case "bool":
			bull = true
			numer = false
			typ = "bool"
		case "json":
			bull = false
			numer = false
			typ = "json"
		}
	}

	adFind := j.FindElement(path, typ)
	//fmt.Println("adFind.start_value: " + //fmt.Sprint(adFind.start_value))
	//fmt.Println("set value: " + value)
	if !bull && !numer {
		j.jsonString = j.jsonString[:adFind.start_value] + "\"" + value + "\"" + j.jsonString[adFind.end_value:]
	} else {
		j.jsonString = j.jsonString[:adFind.start_value] + value + j.jsonString[adFind.end_value:]
	}
}

func (j *ADJson) FindElement(path string, type_opt ...string) ADFind {

	bull := false
	numer := false
	array := false
	json := false

	if len(type_opt) > 0 {
		switch type_opt[0] {
		case "string":
			bull = false
			numer = false
		case "number":
			bull = false
			numer = true
		case "bool":
			bull = true
			numer = false
		case "array":
			bull = false
			numer = false
			array = true
		case "json":
			bull = false
			numer = false
			json = true
		}
	}

	var start_idx int
	var end_idx int
	var start_value int
	var end_value int
	var adFind ADFind

	start_idx = 0
	end_idx = -2
	start_value = 0
	end_value = 0

	adFind.start_idx = 0
	adFind.end_idx = 0
	adFind.start_value = 0
	adFind.end_value = 0

	paths := strings.Split(path, ".")
	wherearrint := make([]int, 0)
	for k, element := range paths {
		where := ""
		wherearr := make([]string, 0)
		if strings.Index(element, "[") > 0 {
			where = element[strings.Index(element, "["):]
			where = strings.Replace(where, "[", "", -1)
			wherearr = strings.Split(where, "]")
			for _, w := range wherearr {
				////fmt.Printf("w= %s\n", w)
				wint, err := strconv.Atoi(w)
				if err == nil {
					wherearrint = append(wherearrint, wint)
				} else {
					//fmt.Printf("Error %v %v\n", err, w)
				}
			}
			element = element[:strings.Index(element, "[")]
		}
		for {
			//fmt.Println("Szukam w: " + j.jsonString[end_idx+2:])
			start_idx = strings.Index(j.jsonString[end_idx+2:], element+"\":")

			if start_idx != -1 {
				start_idx += end_idx + 2
			} else {
				return adFind
			}
			end_idx = start_idx + len(element)
			//fmt.Println("start_idx: " + //fmt.Sprint(start_idx))
			//fmt.Println("end_idx: " + //fmt.Sprint(end_idx))
			//fmt.Println(j.jsonString[start_idx:end_idx])
			bracket := strings.Index(j.jsonString[end_idx+2:end_idx+4], "}")
			doublebracket := strings.Index(j.jsonString[end_idx+2:end_idx+4], "{}")

			if k+1 == int(j.brackets[start_idx]) && doublebracket == -1 && bracket == -1 {
				//fmt.Println("NIe ma nawiasu i dobry poziom")
				break
			}
			if doublebracket != -1 {
				//fmt.Println("Jest nawias podwójny, szukaj dalej")
				continue
			}
			if j.jsonString[start_idx:end_idx] == element {
				break
			}
			if bracket != -1 && k == len(paths)-1 {
				//fmt.Println("Jest nawias ale koniec więc OK")
				break
			}
		}
	}
	//fmt.Println("Odczytuję wartość")
	if !bull && !numer && !array {

		//fmt.Println("Szukam wartości string w: " + j.jsonString[end_idx+2:])
		start_value = strings.Index(j.jsonString[end_idx+2:], "\"")
		if start_value != -1 {
			start_value += end_idx + 2
		} else {
			return adFind
		}

		//fmt.Println("Szukam końca wartości w: " + j.jsonString[start_value+1:])
		end_value = strings.Index(j.jsonString[start_value+1:], "\"")
		//fmt.Println(j.jsonString[start_value+1:])
		//fmt.Println("end_value: " + fmt.Sprint(end_value))
		if json {
			new_start_value := start_value
			for e := 1; e < 20; e++ {
				//fmt.Println(string(j.jsonString[new_start_value+1+end_value-1]))
				if string(j.jsonString[new_start_value+end_value]) == "\\" {

					new_start_value += end_value + 1
					end_value = strings.Index(j.jsonString[new_start_value+1:], "\"")
					//fmt.Println(j.jsonString[new_start_value+1:])
					//fmt.Println("inside end_value: " + fmt.Sprint(end_value))
				} else {
					end_value = new_start_value - start_value + 1
					break
				}
			}
		}
		if end_value != -1 {
			end_value += start_value + 2
		} else {
			return adFind
		}

		//fmt.Println("start_value: " + //fmt.Sprint(start_value))
		//fmt.Println("end_value: " + //fmt.Sprint(end_value))

		adFind.start_idx = start_idx
		adFind.end_idx = end_idx
		adFind.start_value = start_value
		adFind.end_value = end_value
	} else if numer {
		//fmt.Println("Szukam wartości number w: " + j.jsonString[end_idx+2:])
		bracket_value := strings.Index(j.jsonString[end_idx+2:], "}")
		comma_value := strings.Index(j.jsonString[end_idx+2:], ",")
		sep_value := -1
		if bracket_value != -1 && comma_value != -1 {
			if bracket_value < comma_value {
				sep_value = bracket_value
			} else {
				sep_value = comma_value
			}
		}
		if bracket_value != -1 && comma_value == -1 {
			sep_value = bracket_value
		}
		if bracket_value == -1 && comma_value != -1 {
			sep_value = comma_value
		}

		//fmt.Println("} or , " + //fmt.Sprint(sep_value))
		if sep_value != -1 && sep_value > 0 {
			start_value += end_idx + 2
			end_value = start_value + sep_value
		} else {
			return adFind
		}
		//fmt.Println("start_value: " + //fmt.Sprint(start_value))
		//fmt.Println("end_value: " + //fmt.Sprint(end_value))

		adFind.start_idx = start_idx
		adFind.end_idx = end_idx
		adFind.start_value = start_value
		adFind.end_value = end_value
	} else if array {
		//fmt.Println("Szukam wartości array w: " + j.jsonString[end_idx+2:])
		start_value = strings.Index(j.jsonString[end_idx+2:], "[")
		if start_value != -1 {
			start_value += end_idx + 2
		} else {
			return adFind
		}

		//fmt.Println("Szukam końca wartości array w: " + j.jsonString[start_value+1:])

		level := j.arrbrackets[start_value] - 1
		for k := start_value; k < len(j.jsonString); k++ {
			if j.arrbrackets[k] == level && string(j.jsonString[k]) == "]" {
				end_value = k + 1
				break
			}
		}

		node := Node{}
		token := Token{start_value, end_value}
		node.Insert(*j, token, 0)
		deep_node := node
		for _, d := range wherearrint {
			deep_node = *deep_node.a[d]
		}
		adFind.start_value = deep_node.token.start
		adFind.end_value = deep_node.token.end

		adFind.start_idx = start_idx
		adFind.end_idx = end_idx

	} else {
		//fmt.Println("Szukam wartości logicznej w: " + j.jsonString[end_idx+2:])
		start_value = strings.Index(j.jsonString[end_idx+2:], "true")
		if start_value != -1 && start_value < 2 {
			start_value += end_idx + 2
			end_value = start_value + 4
		} else {
			start_value = strings.Index(j.jsonString[end_idx+2:], "false")
			if start_value != -1 && start_value < 2 {
				start_value += end_idx + 2
				end_value = start_value + 5
			} else {
				return adFind
			}
		}
		//fmt.Println("start_value: " + //fmt.Sprint(start_value))
		//fmt.Println("end_value: " + //fmt.Sprint(end_value))

		adFind.start_idx = start_idx
		adFind.end_idx = end_idx
		adFind.start_value = start_value
		adFind.end_value = end_value
	}

	return adFind

}
