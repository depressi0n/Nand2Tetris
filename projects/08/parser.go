package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Parser struct{
	VMFileNames []string
	Cursor int // file cursor
	SymbolTable *SymbolTable
	VMCommands []*Command
	Current int // command cursor
	OutputFileName string
	// some variables that changed in the parse process such as some counter for different segment
	CallCnt int // record the call count for unique the return name when translating the call instruction
	RetCnt int

}

func NewParser(filenames []string,outputFileName string)*Parser{
	res:=&Parser{
		VMFileNames: filenames,
		OutputFileName: outputFileName,
		CallCnt: 0,
		RetCnt: 0,
	}
	res.SymbolTable=NewSymbolTable()
	return res
}
func(p *Parser)Prepare(){
	p.VMCommands=make([]*Command, 0,100)
	// Open all .vm file
	for ;p.Cursor<len(p.VMFileNames);p.Cursor++ {
		f, err := os.OpenFile(p.VMFileNames[p.Cursor], os.O_RDONLY, 0600)
		if err!=nil{
			log.Fatalf("Error in os.OpenFile with %s:%v\n",p.VMFileNames[p.Cursor],err)
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
			cmd:=NewCommand(s)
			if cmd.Type() == C_PUSH ||  cmd.Type() == C_POP {
				// check if the segment is static, if so, add the filename as prefix
				// Due to the static variable will not be used by different file but just used in same file
				if cmd.SegmentType() == SEG_STATIC {
					cmd.AddPrefixForArg2(p.VMFileNames[p.Cursor])
				}
			}
			p.VMCommands=append(p.VMCommands,cmd )
		}
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
		case SEG_ARGUMENT:
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
		case SEG_LOCAL:
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
		case SEG_STATIC:
			// 使用p.OutputFileName_[vairiableName]表示静态变量，借助了汇编程序为程序中变量提供RAM单元的特性
			staticVariableName:=strings.TrimPrefix(strings.ReplaceAll(p.OutputFileName,"/","."),".asm")+"."+c.Arg2()
			res=fmt.Sprintf("@%s\n",staticVariableName)+`D=M
@SP
A=M
M=D
@SP
M=M+1 // 压入栈中
`
		case SEG_CONSTANT:
			// 将常数压入栈中
				res=fmt.Sprintf("@%v\nD=A\n",c.Arg2())+`@SP
A=M
M=D
@SP
M=M+1
`
		case SEG_THIS:
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
		case SEG_THAT:
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
		case SEG_POINTER:
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
		case SEG_TEMP:
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
		case SEG_ARGUMENT:
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
		case SEG_LOCAL:
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
		case SEG_STATIC:
			staticVariableName:=strings.TrimPrefix(strings.ReplaceAll(p.OutputFileName,"/","."),".asm")+"."+c.Arg2()
			res=`@SP
A=M-1
D=M
@SP
M=M-1 // 栈中数据出栈
`+fmt.Sprintf("@%s\n",staticVariableName)+
`M=D
`
		case SEG_CONSTANT:
			// 不应该支持pop constant n
			log.Fatalf("Invalid pop command with segment constant")
		case SEG_THIS:
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
		case SEG_THAT:
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
		case SEG_POINTER:
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
		case SEG_TEMP:
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
		res=fmt.Sprintf("(%s)\n",c.Arg1())
	case C_GOTO:
		res=fmt.Sprintf("@%s\n0;JMP\n",c.Arg1())
	case C_IF:
		res=`@SP
A=M-1
D=M
@SP
M=M-1
`+fmt.Sprintf("@%s\n",c.Arg1())+
`
D;JNE
`
	case C_FUNCTION:
		res=fmt.Sprintf("(%s)\n",c.Arg1())
		// 声明k个局部变量
		k,_:=strconv.Atoi(c.Arg2())
		for i:=0;i<k;i++{
			res+=`@SP
A=M
M=0
@SP
M=M+1
`
		}
	case C_RETURN:
		res=`@R1 // FRAME=LCL
D=M
@FRAME
M=D
@5
A=D-A
D=M
`+fmt.Sprintf("@RET%v\n",p.RetCnt)+
`M=D
@SP // 重置调用者的返回值，从栈里面弹出一个放到ARG的位置
M=M-1
A=M
D=M
@R2
A=M
M=D
@R2 // 恢复调用者SP
D=M
@R0
M=D+1
@FRAME // 恢复调用者THAT
D=M
@1
A=D-A
D=M
@R4
M=D
@FRAME // 恢复调用者THIS
D=M
@2
A=D-A
D=M
@R3
M=D
@FRAME // 恢复调用者ARG
D=M
@3
A=D-A
D=M
@R2
M=D
@FRAME // 恢复调用者LCL
D=M
@4
A=D-A
D=M
@R1
M=D
`+fmt.Sprintf("@RET%v\n",p.RetCnt)+
`A=M
0;JMP
`
		p.RetCnt++
	case C_CALL:
		res=fmt.Sprintf("@%s.return-address%v // 返回地址压栈\n",c.Arg1(),p.CallCnt)+
`D=A
@SP
A=M
M=D
@SP
M=M+1
@R1 // LCL压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R2 // ARG压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R3 // THIS压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R4 // THAT压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP // ARG=SP-n-5
D=M
@5
D=D-A
`+fmt.Sprintf("@%v\n",c.Arg2())+
`D=D-A
@R2
M=D
@SP // 重置LCL=SP
D=M
@R1
M=D
`+fmt.Sprintf("@%s // 跳转控制\n",c.Arg1())+
`0;JMP
`+fmt.Sprintf("(%s.return-address%v) // 返回地址\n",c.Arg1(),p.CallCnt)
	p.CallCnt++
	default:
		log.Fatalf("Invalid Command Type:%v",c.Type)
	}
	return res
}

func(p *Parser)Run(){
	f, err := os.OpenFile(p.OutputFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err!=nil{
		log.Fatalf("Error in os.OpenFile with %s:%v",p.OutputFileName,err)
	}
	defer f.Close()
	p.Prepare()
	// Write the bootstrap code
	f.WriteString(`@256  // SP=256
D=A
@R0
M=D
@Sys.return-address // 返回地址压栈
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
@R1 // LCL压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R2 // ARG压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R3 // THIS压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R4 // THAT压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP // ARG=SP-5
D=M
@5
D=D-A
@R2
M=D
@SP // 重置LCL=SP
D=M
@R1
M=D
@Sys.init // call Sys.init
0;JMP
(Sys.return-address)
`)
f.WriteString(`(END)
	@END
	0;JMP
`)
	for i:=0;i<len(p.VMCommands);i++{
		// handle each command
		f.WriteString(p.Translate(p.VMCommands[i]))
		p.Current++
	}
	
}