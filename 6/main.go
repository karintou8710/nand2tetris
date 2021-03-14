package main		// mainパッケージであることを宣言

import (
	"fmt"
	"m/module"
	"strconv"
	"strings"
	"io/ioutil"
	"os"
	"path"
)

func ten2byte(n int) string {
	res := ""
	cnt := 0
	for n>0 {
		res = strconv.Itoa(n%2) + res
		n/=2
		cnt++
	}
	
	for i:=0;i<15-cnt;i++ {
		res = "0" + res
	}
	return res
}



func main() {

	if len(os.Args) < 2 {
		fmt.Println("ERROR: 引数を指定してください。")
		os.Exit(1)
	}

	dir, file := path.Split(os.Args[1])


	var parser *module.Parser = module.NewParser(os.Args[1])
	var codeModule *module.Code

	var output_arr []string
	symbolMap := make(map[string]int)

	// シンボルマップ初期化
	symbolMap["SP"] = 0
	symbolMap["LCL"] = 1
	symbolMap["ARG"] = 2
	symbolMap["THIS"] = 3
	symbolMap["THAT"] = 4
	for i:=0;i<16;i++ {
		str_i := strconv.Itoa(i)
		symbolMap["R"+str_i] = i
	}

	symbolMap["SCREEN"] = 16384
	symbolMap["KBD"] = 24576

	for {
		if (!parser.HasMoreCommands()){
			break
		}

		parser.Advance()

		commandType := parser.CommandType()
		if commandType == "L_COMMAND" {
			v, _,  counter := parser.Symbol()
			symbolMap[v] = counter
		}
	}

	parser.Reset()
	memoryCnt := 16

	for {
		if (!parser.HasMoreCommands()){
			break
		}

		parser.Advance()

		commandType := parser.CommandType()
		var machineLn string
		if commandType == "A_COMMAND" {
			machineLn = "0"
			value, isNum, _ := parser.Symbol()
			num := 0
			
			if isNum {
				num, _ = strconv.Atoi(value)
			}else{
				if v, ok := symbolMap[value]; ok {
					num = v
				}else{
					num = memoryCnt
					symbolMap[value] = memoryCnt
					memoryCnt++
				}
			}
			machineLn += ten2byte(num)
		
		}else if commandType == "C_COMMAND" {
			machineLn = "111"
			machineLn += codeModule.Comp(parser.Comp())
			machineLn += codeModule.Dest(parser.Dest())
			machineLn += codeModule.Jump(parser.Jump())
		}else if commandType == "L_COMMAND" {
			continue
		}

		output_arr = append(output_arr, machineLn)
	}

	output_str := strings.Join(output_arr, "\r\n")
	output := []byte(output_str)
	name := strings.Split(file, ".")[0]
	err := ioutil.WriteFile(path.Join(dir, name)+".hack", output, 0666)
    if err != nil {
        fmt.Println(os.Stderr, err)
        os.Exit(1)
    }

	fmt.Println("OK!")
}