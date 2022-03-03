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
				res=`@SP
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
		switch c.SegmentType(){
		case ARGUMENT:
			// 基址直接映射到RAM[2]即ARG寄存器上
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@R2
A=D+M
D=M // 取出argument[index]的数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
`
		case LOCAL:
			// 基址直接映射到RAM[1]即LCL寄存器上
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@R1
A=D+M
D=M // 取出local[index]的数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
`
		case STATIC:
			panic("implement me!")
		case CONSTANT:
			// 将常数压入栈中
				res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@SP
A=M
M=D
@SP
M=M+1
`
		case THIS:
			// RAM[3]
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@R3
A=D+M
D=M // 取出local[index]的数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
`
		case THAT:
			// RAM[4]
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@R4
A=D+M
D=M // 取出local[index]的数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
`
		case POINTER:
			arg2:=c.Arg2()
			if arg2== "0"{
				res=`@R3
D=M // 取出this当前指向的地址
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
`
			}else if arg2 == "1"{
				res=`@R4
D=M // 取出this当前指向的地址
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
`
			}else{
				log.Fatalf("Invalid parameter with pointer segment:%s",arg2)
			}
		case TEMP:
			// RAM[5-12]
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@5
A=D+A
D=M // 取出数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
`
		default:
			log.Fatalf("Invalid Segment Type")
		}
	case C_POP:
		switch c.SegmentType(){
		case ARGUMENT:
			// 基址直接映射到RAM[2]即ARG寄存器上
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@R2
D=D+M
@SP
A=M
M=D // 使用栈顶存放位置信息
@SP
A=M-1
D=M
@SP
A=M
A=M
M=D
@SP
M=M-1
`
		case LOCAL:
			// 基址直接映射到RAM[1]即LCL寄存器上
res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@R1
D=D+M
@SP
A=M
M=D // 使用栈顶存放位置信息
@SP
A=M-1
D=M
@SP
A=M
A=M
M=D
@SP
M=M-1
`
		case STATIC:
			panic("implement me!")
		case CONSTANT:
			// 不应该支持pop constant n
			log.Fatalf("Invalid pop command with segment constant")
		case THIS:
			// RAM[3]
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@R3
D=D+M
@SP
A=M
M=D // 使用栈顶存放位置信息
@SP
A=M-1
D=M
@SP
A=M
A=M
M=D
@SP
M=M-1
`
		case THAT:
			// RAM[4]
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@R4
D=D+M
@SP
A=M
M=D // 使用栈顶存放位置信息
@SP
A=M-1
D=M
@SP
A=M
A=M
M=D
@SP
M=M-1
`
		case POINTER:
			arg2:=c.Arg2()
			if arg2== "0"{
				res=`@SP
A=M-1
D=M // 取出栈顶数据
@SP
M=M-1
@R3
M=D
`
			}else if arg2 == "1"{
				res=`@SP
A=M-1
D=M // 取出栈顶数据
@SP
M=M-1
@R4
M=D
`
			}else{
				log.Fatalf("Invalid parameter with pointer segment:%s",arg2)
			}
		case TEMP:
			// RAM[5-12]
			res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@5
D=D+A
@SP
A=M
M=D // 使用栈顶存放位置信息
@SP
A=M-1
D=M
@SP
A=M
A=M
M=D
@SP
M=M-1
`
		default:
			log.Fatalf("Invalid Segment Type")
		}

	case C_LABEL:
		panic("implement me!")
	case C_GOTO:
		panic("implement me!")
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
	f.WriteString(`(END)
	@END
	0;JMP`)
}