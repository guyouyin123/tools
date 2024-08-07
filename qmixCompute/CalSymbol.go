package qmixCompute

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
)

var magic_num int32 = 0xF0000 //Unicode Supplementary Private Use Area-A
type Calculate struct {
	temp            []rune //Algo String
	temp_Cp         []rune //Algo String Copy
	tempStack       []rune
	top             int
	magic_num_shift int32
	Max_Str         int
	MapValue        map[rune]float64
}

var OparateMap = make(map[rune]int)

func IniOperateDic() {
	OparateMap['('] = 0
	OparateMap[')'] = 4
	OparateMap['^'] = 3
	OparateMap['*'] = 2
	OparateMap['/'] = 2
	OparateMap['+'] = 1
	OparateMap['-'] = 1
}
func (p *Calculate) IniValue() {
	p.MapValue = make(map[rune]float64)
	for i := range p.tempStack {
		p.temp[i] = '\u0000'
		p.tempStack[i] = '\u0000'
	}
	p.magic_num_shift = 0
}
func (p *Calculate) IniSize() {
	p.temp = make([]rune, p.Max_Str)
	p.temp_Cp = make([]rune, p.Max_Str)
	p.tempStack = make([]rune, p.Max_Str)
}
func NewStruct(max_str_size ...interface{}) Calculate {
	new_one := Calculate{}
	if len(max_str_size) > 0 {
		new_one.Max_Str = max_str_size[0].(int)
	} else {
		new_one.Max_Str = 20
	}
	new_one.IniSize()
	new_one.IniValue()
	return new_one
}

func IsOperate(index rune) bool {
	if OparateMap[index] == 0 {
		return false
	}
	return true
}
func (p *Calculate) push_post(index rune) {
	p.tempStack[p.top] = index
	p.top++
}
func (p *Calculate) pop_post(pos_in int) int {
	if p.top > 0 {
		p.temp[pos_in] = p.tempStack[p.top-1]
		p.top--
		return pos_in + 1
	}
	return pos_in
}
func (p *Calculate) pop_val() rune {
	return p.tempStack[p.top-1]
}
func (p *Calculate) pop_able() bool {
	if p.top > 0 {
		return true
	}
	return false
}
func (p *Calculate) in2post_rune(a []rune) {
	p.top = 0
	i, j := 0, 0

	for ; i < len(a); i++ {

		switch a[i] {
		case '(':
			p.push_post(rune(a[i]))
		case '+', '-', '*', '/', '^':
			for p.pop_able() {
				if OparateMap[rune(p.pop_val())] >= OparateMap[rune(a[i])] {
					j = p.pop_post(j)
				} else {
					break
				}
			}
			p.push_post(rune(a[i]))
		case ')':
			for p.pop_able() {
				if rune(p.pop_val()) == '(' {
					p.top = p.top - 1
					break
				} else {
					j = p.pop_post(j)
				}
			}
		default:
			p.temp[j] = rune(a[i])
			j++

		}

	}
	for p.pop_able() {
		j = p.pop_post(j)
	}
}
func (p *Calculate) Set(a rune, b float64) { //a:target word ,b:target word 's value
	p.MapValue[a] = b
}
func (p *Calculate) SetByMap(a map[rune]float64) { //a:target word ,b:target word 's value
	for k, v := range a {
		p.MapValue[k] = v
	}
}
func (p *Calculate) Compute() float64 { //
	for i := 0; i < p.Max_Str; i++ {
		p.temp_Cp[i] = p.temp[i]
	}
	var uni_valid_symbol rune = p.magic_num_shift + magic_num
	var uni_start_pos int32 = p.magic_num_shift + magic_num + 1
	var compute_num [2]int32 = [2]int32{uni_valid_symbol, uni_valid_symbol} /*Map What Real Value*/
	var compute_pos [2]int = [2]int{0, 0}                                   /*Place In temp_Cp*/
	for i := 0; p.temp_Cp[i] != '\u0000'; i++ {
		if IsOperate(p.temp_Cp[i]) {
			for j, k := i-1, 2; k != 0 && j > -1; j-- {
				if p.temp_Cp[j] != uni_valid_symbol && (!IsOperate(p.temp_Cp[j])) {
					compute_num[k-1] = p.temp_Cp[j]
					compute_pos[k-1] = j
					k--
				}
			}
			switch p.temp_Cp[i] {
			case '^':
				p.MapValue[rune(uni_start_pos)] = math.Pow(p.MapValue[rune(compute_num[0])], p.MapValue[rune(compute_num[1])])
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			case '+':
				p.MapValue[rune(uni_start_pos)] = p.MapValue[rune(compute_num[0])] + p.MapValue[rune(compute_num[1])]
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			case '-':
				p.MapValue[rune(uni_start_pos)] = p.MapValue[rune(compute_num[0])] - p.MapValue[rune(compute_num[1])]
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			case '*':
				p.MapValue[rune(uni_start_pos)] = p.MapValue[rune(compute_num[0])] * p.MapValue[rune(compute_num[1])]
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			case '/':
				p.MapValue[rune(uni_start_pos)] = p.MapValue[rune(compute_num[0])] / p.MapValue[rune(compute_num[1])]
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			}
		}
	}
	if p.magic_num_shift+magic_num == uni_start_pos {
		return p.MapValue[p.temp_Cp[0]]
	}
	return p.MapValue[rune(uni_start_pos-1)]

}
func (p *Calculate) ComputeAndShow() float64 { //
	for i := 0; i < p.Max_Str; i++ {
		p.temp_Cp[i] = p.temp[i]
	}
	var uni_valid_symbol rune = p.magic_num_shift + magic_num
	var uni_start_pos int32 = p.magic_num_shift + magic_num + 1
	var compute_num [2]int32 = [2]int32{uni_valid_symbol, uni_valid_symbol} /*Map What Real Value*/
	var compute_pos [2]int = [2]int{0, 0}                                   /*Place In temp_Cp*/
	for i := 0; p.temp_Cp[i] != '\u0000'; i++ {
		if IsOperate(p.temp_Cp[i]) {
			for j, k := i-1, 2; k != 0 && j > -1; j-- {
				if p.temp_Cp[j] != uni_valid_symbol && (!IsOperate(p.temp_Cp[j])) {
					compute_num[k-1] = p.temp_Cp[j]
					compute_pos[k-1] = j
					k--
				}
			}
			switch p.temp_Cp[i] {
			case '^':
				fmt.Println("Let U+", fmt.Sprintf("%x", uni_start_pos), "=", p.MapValue[rune(compute_num[0])], "^", p.MapValue[rune(compute_num[1])])
				p.MapValue[rune(uni_start_pos)] = math.Pow(p.MapValue[rune(compute_num[0])], p.MapValue[rune(compute_num[1])])
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			case '+':
				fmt.Println("Let U+", fmt.Sprintf("%x", uni_start_pos), "=", p.MapValue[rune(compute_num[0])], "+", p.MapValue[rune(compute_num[1])])
				p.MapValue[rune(uni_start_pos)] = p.MapValue[rune(compute_num[0])] + p.MapValue[rune(compute_num[1])]
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			case '-':
				fmt.Println("Let U+", fmt.Sprintf("%x", uni_start_pos), "=", p.MapValue[rune(compute_num[0])], "-", p.MapValue[rune(compute_num[1])])
				p.MapValue[rune(uni_start_pos)] = p.MapValue[rune(compute_num[0])] - p.MapValue[rune(compute_num[1])]
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			case '*':
				fmt.Println("Let U+", fmt.Sprintf("%x", uni_start_pos), "=", p.MapValue[rune(compute_num[0])], "*", p.MapValue[rune(compute_num[1])])
				p.MapValue[rune(uni_start_pos)] = p.MapValue[rune(compute_num[0])] * p.MapValue[rune(compute_num[1])]
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			case '/':
				fmt.Println("Let U+", fmt.Sprintf("%x", uni_start_pos), "=", p.MapValue[rune(compute_num[0])], "/", p.MapValue[rune(compute_num[1])])
				p.MapValue[rune(uni_start_pos)] = p.MapValue[rune(compute_num[0])] / p.MapValue[rune(compute_num[1])]
				p.temp_Cp[compute_pos[1]] = rune(uni_start_pos)
				p.temp_Cp[compute_pos[0]] = uni_valid_symbol
				uni_start_pos++
			}
		}
	}
	if p.magic_num_shift+magic_num == uni_start_pos {
		return p.MapValue[p.temp_Cp[0]]
	}
	return p.MapValue[rune(uni_start_pos-1)]

}
func RegexNormal(a string) string {
	var re = regexp.MustCompile(`\(-`)
	s := re.ReplaceAllString(a, `(0-`)
	if s[0] == '-' {
		s = string('0') + s
	}
	return s
}

type ByLen []string

func (a ByLen) Len() int           { return len(a) }
func (a ByLen) Less(i, j int) bool { return len(a[i]) > len(a[j]) }
func (a ByLen) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func RegexFindNumberReplace(a string, b map[rune]float64, c int32) (string, int32) { //a:Ori , b:Contain , c:Start
	var re = regexp.MustCompile(`[\d]*[.]{0,1}[\d]+`)
	s := re.FindAllString(a, -1)
	d := a
	sort.Sort(ByLen(s))
	for _, j := range s {
		b[c], _ = strconv.ParseFloat(j, 64)
		var re2 = regexp.MustCompile(j)
		d = re2.ReplaceAllString(d, string(rune(c)))
		c++
	}
	return d, c
}
func (p *Calculate) GiveRule(a string) {
	b := RegexNormal(a)
	c, d := RegexFindNumberReplace(b, p.MapValue, magic_num)
	p.in2post_rune([]rune(c))
	p.magic_num_shift = d - magic_num
}

func (p *Calculate) DebugMap() {
	fmt.Println(">>Each Step Temp Keeping<<")
	for i, j := range p.MapValue {
		fmt.Println(string(i), "(U+", fmt.Sprintf("%x", i), "):", j)
	}
	for i, j := range p.MapValue {
		fmt.Println(string(i), "(U+", fmt.Sprintf("%x", i), "):", j)
	}
	fmt.Println("<<Each Step Temp Keeping>>")
}
func (p *Calculate) DebugPostfix() {
	fmt.Println(">>Postfix<<")
	for _, i := range p.temp {
		if val, exist := p.MapValue[i]; exist && !IsOperate(i) {
			fmt.Print(val, " ")
		} else {
			fmt.Print(string(i), " ")
		}
	}
	fmt.Print("\n")
	fmt.Println("<<Postfix>>")
}
func init() {
	IniOperateDic()
}
