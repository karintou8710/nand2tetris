package module

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"regexp"
)

type Parser struct {
	code []string
	nowLine string
	nowIndex int
	counter int
}

func NewParser(filePath string) (*Parser) {
	parser := new(Parser)
	parser.nowIndex = 0
	parser.counter = 0

	f, err := os.Open(filePath)
	if err!=nil {
		fmt.Println("can't open file !")
		os.Exit(1)
	}
	defer f.Close()

	code, err := ioutil.ReadAll(f)
	if err!=nil {
		fmt.Println("can't read file !")
		os.Exit(1)
	}
	codeSlice := strings.Split(string(code),"\r\n")
	for _, line := range codeSlice {
		line = strings.Replace(line, " ", "", -1)
		// コメント
		commentReg := regexp.MustCompile(`//.*`)
		line = commentReg.ReplaceAllString(line, "")
		// 空行
		if (len(line)<=1){
			continue
		}

		parser.code = append(parser.code, line)
	}

	return parser
}

func (p *Parser) HasMoreCommands() bool {
	if p.nowIndex+1 <= len(p.code) {
		return true
	}else {
		return false
	}
}

func (p *Parser) Advance() {
	p.nowLine = p.code[p.nowIndex]
	p.nowIndex++
}

// 例外処理は対応しない
func (p *Parser) CommandType() string{
	if p.nowLine[0]=='@' {
		p.counter++
		return "A_COMMAND"
	}else if p.nowLine[0]=='(' {
		return "L_COMMAND"
	}else{
		p.counter++
		return "C_COMMAND"
	}
}

func (p *Parser) Symbol() (string, bool, int) {
	rep := regexp.MustCompile(`[@\(\)]`)
	rep2 := regexp.MustCompile(`^[0-9]+$`)
	res := rep.ReplaceAllString(p.nowLine, "")
	return res, rep2.MatchString(res), p.counter
}

func (p *Parser) Dest()string {
	indexEqual := strings.Index(p.nowLine, "=")
	if indexEqual!=-1 {
		return p.nowLine[0:indexEqual]
	}else{
		return ""
	}
}

func (p *Parser) Comp() string {
	indexEqual := strings.Index(p.nowLine, "=")
	if indexEqual==-1 {
		indexSemi := strings.Index(p.nowLine, ";")
		if (indexSemi!=-1){
			return p.nowLine[:indexSemi]
		}else{
			return "error"
		}
	}else{
		return p.nowLine[indexEqual+1:]
	}
}

func (p *Parser) Jump() string {
	indexSemi := strings.Index(p.nowLine, ";")
	if (indexSemi!=-1){
		return p.nowLine[indexSemi+1:]
	}else{
		return ""
	}
}

func (p *Parser) Reset() {
	p.nowIndex = 0
	p.counter = 0
	p.nowLine = ""
}