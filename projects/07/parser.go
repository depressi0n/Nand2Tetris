package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Parser struct{
	VMFileName string
	SymbolTable *SymbolTable
	VMCommands []*Command
	Current int
	// some variables that changed in the parse process such as some counter for different segment 

}

func NewParser(filename string)*Parser{
	res:=&Parser{
		VMFileName: filename,
	}
	res.SymbolTable=NewSymbolTable()
	return res
}

func(p *Parser)Prepare(){
	p.VMCommands=make([]*Command, 0,100)
	// Open the .vm file
	f, err := os.OpenFile(p.VMFileName, os.O_RDONLY, 0600)
	if err!=nil{
		log.Fatalf("Error in os.OpenFile with %s:%v",p.VMFileName,err)
	}
	defer f.Close()
	// Remove the comment lines and whitespace
	r := bufio.NewReader(f)
	for{
		line, _, err := r.ReadLine()
		if err!=nil{
			if err==io.EOF{
				break
			}
			log.Fatalf("Error in r.ReadLine:%v",err)
		}
		s:=string(line)
		// Remove the suffix and the prefix whitespace
		s = strings.TrimSpace(s)
		// Check the validity
		if len(line)==0 || strings.HasPrefix(s,"//") {
			continue
		}
		p.VMCommands=append(p.VMCommands, NewCommand(s))
	}
}
// xxx.vm文件中静态变量j被翻译为xxx.j，由Hack汇编编译器分配RAM空间
func(p *Parser)Translate(c *Command)string{
	var res string
	switch c.Type() {
	case C_ARITHMETIC:
		switch c.OperatorType(){
		case OP_ADD:
				res=`@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=D+M
@SP
A=M
M=D
@SP
M=M+1
`
		case OP_SUB:
				res=`@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@SP
A=M
M=D
@SP
M=M+1
`
		case OP_NEG:
				res=`@SP
M=M-1
A=M
D=-M
@SP
A=M
M=D
@SP
M=M+1
`
		case OP_EQ:
				res=`@SP
M=M-1
A=M
D=M // 第二个操作数
@SP
M=M-1
A=M
D=M-D // 第一个操作数减去第二个操作数
`+fmt.Sprintf("@TRUE%v\n",p.Current)+`D;JEQ
@SP
A=M
M=0
`+fmt.Sprintf("@SKIP%v\n",p.Current)+`0;JMP
`+fmt.Sprintf("(TRUE%v)\n",p.Current)+`  @SP
  A=M
  M=-1
  `+fmt.Sprintf("@SKIP%v\n",p.Current)+`  0;JMP
`+fmt.Sprintf("(SKIP%v)\n",p.Current)+`  @SP
  M=M+1
`
		case OP_GT:
			// 涉及到标签不重复？使用一个递增量
			// 另一种处理方式使用重复的标签值，因为所有的基本设定是一致
			res=`@SP
M=M-1
A=M
D=M // 第二个操作数
@SP
M=M-1
A=M
D=M-D // 第一个操作数减去第二个操作数
`+fmt.Sprintf("@TRUE%v\n",p.Current)+`D;JGT
@SP
A=M
M=0
`+fmt.Sprintf("@SKIP%v\n",p.Current)+`0;JMP
`+fmt.Sprintf("(TRUE%v)\n",p.Current)+`  @SP
  A=M
  M=-1
  `+fmt.Sprintf("@SKIP%v\n",p.Current)+`  0;JMP
`+fmt.Sprintf("(SKIP%v)\n",p.Current)+`  @SP
  M=M+1
`
		case OP_LT:
			res=`@SP
M=M-1
A=M
D=M // 第二个操作数
@SP
M=M-1
A=M
D=M-D // 第一个操作数减去第二个操作数
`+fmt.Sprintf("@TRUE%v\n",p.Current)+`D;JLT
@SP
A=M
M=0
`+fmt.Sprintf("@SKIP%v\n",p.Current)+`0;JMP
`+fmt.Sprintf("(TRUE%v)\n",p.Current)+`  @SP
  A=M
  M=-1
  `+fmt.Sprintf("@SKIP%v\n",p.Current)+`  0;JMP
`+fmt.Sprintf("(SKIP%v)\n",p.Current)+`  @SP
  M=M+1
`
		case OP_AND:
				res=`@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=D&M
@SP
A=M
M=D
@SP
M=M+1
`
		case OP_OR:
			res=`@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=D|M
@SP
A=M
M=D
@SP
M=M+1
`
		case OP_NOT:
				res=`
@SP
M=M-1
A=M
D=!M
@SP
A=M
M=D
@SP
M=M+1
`
		default:
			log.Fatalf("Invalid Operator Type")

		}
	case C_PUSH:
		// push segment index
		// 暂时只支持push constant x，将x放入栈中
		switch c.SegmentType(){
		case ARGUMENT:
		case LOCAL:
			// 
		case STATIC:
		case CONSTANT:
			// 将常数压入栈中
				res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@SP
A=M
M=D
@SP
M=M+1
`
		case THIS:
		case THAT:
		case POINTER:
		case TEMP:
		default:
			log.Fatalf("Invalid Segment Type")
		}

	
		
		
	case C_POP:

	case C_LABEL:

	case C_GOTO:

	case C_IF:
		panic("implement me!")
	case C_FUNCTION:
		panic("implement me!")
	case C_RETURN:
		panic("implement me!")
	case C_CALL:
		panic("implement me!")
	default:
		log.Fatalf("Invalid Command Type:%v",c.Type)
	}
	return res
}

func(p *Parser)Run(){
	p.Prepare()
	outputFileName:=strings.TrimSuffix(p.VMFileName, ".vm")+".asm"
	f, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err!=nil{
		log.Fatalf("Error in os.OpenFile with %s:%v",outputFileName,err)
	}
	defer f.Close()
	for i:=0;i<len(p.VMCommands);i++{
		// handle each command
		f.WriteString(p.Translate(p.VMCommands[i]))
		p.Current++
	}
	f.WriteString(`
(END)
	@END
	0;JMP`)
}