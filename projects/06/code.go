// 提供所有汇编命令所对应的二进制代码
package main

import (
	"strings"
)

type CommandType int

const (
	// @xxx中xxx是符号或者十进制数字时
	A_COMMAND CommandType=iota
	// dest=comp:jump
	C_COMMAND
	// 伪指令 当(xxx)中xxx是符号时 
	L_COMMAND
) 
type Command struct{
	commandType CommandType
	origin string
	aInstructionType bool // limited usage to A-instruction, true for symbol, false for numeric
	aValue int
	symbol string
	dest string
	comp string
	jump string
}
func(c *Command)Type()CommandType{
	return c.commandType
}
func(c *Command)Origin()string{
	return c.origin
}

func(c *Command)AInstructionType()bool{
	return c.aInstructionType
}
func(c *Command)AValue()int{
	return c.aValue
}
func(c *Command)Symbol()string{
	return c.symbol
}
func(c *Command)Dest()string{
	return c.dest
}
func(c *Command)Comp()string{
	return c.comp
}
func(c *Command)Jump()string{
	return c.jump
}

func NewCommand(input string)*Command{

	switch input[0] {
	case '@':
		// if the string after '@' is a number, it means using the memory
		value,ok:=parseNumberConstant(input[1:])
		if ok{
			return &Command{
				commandType: A_COMMAND,
				origin: input,
				aInstructionType: false,
				aValue:value,
			}
		}
		return &Command{
			commandType: A_COMMAND,
			origin: input,
			aInstructionType: true,
			symbol: input[1:],
		}
		
	case '(':
		// got the symbol
		return &Command{
			commandType: L_COMMAND,
			origin: input,
			symbol: input[1:len(input)-1],
		}
	default:
		// parse the C-instruction to format dest=comp:jump
		var tail string
		dest,comp,jump:="null","null","null"
		if strings.Contains(input,"="){
			tmp:=strings.SplitAfter(input,"=")
			dest,tail=tmp[0],tmp[1]
			dest=dest[:len(dest)-1]
		}else{
			tail=input
		}
		if strings.Contains(tail,";"){
			tmp:=strings.SplitAfter(tail,";")
			comp,jump=tmp[0],tmp[1]
			comp=comp[:len(comp)-1]
		}else{
			comp=tail
		}
		return &Command{
			commandType: C_COMMAND,
			origin: input,
			dest:dest,
			comp:comp,
			jump:jump,
		}
	}
}
// parseNumberConstant try to parse the token to numeric constant. Now support decimal, hex, and binary format.
func parseNumberConstant(token string)(int,bool){
	// eliminate human-friendly formate such XXX_XXX
	if strings.Contains(token,"_"){
		strings.ReplaceAll(token,"_","")
	}

	// the sign of number default is postive
	sign:=true
	if token[0]=='+' || token[0]=='-'{
		if token[0]=='-' {
			sign=false
		}
		token=token[1:]
	}
	res:=0
	// hex or binary format
	if len(token)>=2 && token[0:2]=="0x"{
		// hex-format
		token=token[2:]
		for i:=2;i<len(token);i++{
			v,ok:=HexMap[token[i]]
			if !ok{
				return -1,false
			}
			res=res*16+v
		}
	}else if len(token)>=2 && token[0:2]=="0b"{
		// binary-format
		token=token[0:2]
		for i:=2;i<len(token);i++{
			if token[i]=='0' || token[i]=='1'{
				res=res*2+int(token[i]-'0')
			}else{
				return -1,false
			}
		}
	}else{
		// decimal-format
		for i:=0;i<len(token);i++{
			if 0<=token[i]-'0' && token[i]-'0'<=9 {
				res=res*10+int(token[i]-'0')
			}else{
				return -1,false
			}
			
		}
	}
	if !sign{
		res=-res
	}
	return res,true
} 

// 助记符符号到二进制码的表格
var dest2Binary = map[string]string{
	"null":"000",
	"M":"001",
	"D":"010",
	"MD":"011",
	"A":"100",
	"AM":"101",
	"AD":"110",
	"AMD":"111",
}
var comp2Binary = map[string]string{
		"0"     : "0101010",
        "1"     : "0111111",
        "-1"    : "0111010",
        "D"     : "0001100",
        "A"     : "0110000",
        "!D"    : "0001101",
        "!A"    : "0110001",
        "-D"    : "0001111",
        "-A"    : "0110011",
        "D+1"   : "0011111",
        "A+1"   : "0110111",
        "D-1"   : "0001110",
        "A-1"   : "0110010",
        "D+A"   : "0000010",
        "D-A"   : "0010011",
        "A-D"   : "0000111",
        "D&A"   : "0000000",
        "D|A"   : "0010101",
		
		"M"		: "1110000",
		"!M"	: "1110001",
		"-M"	: "1110011",
		"M+1"	: "1110111",
		"M-1"	: "1110010",
		"D+M"	: "1000010",
		"D-M"	: "1010011",
		"M-D"	: "1000111",
		"D&M"	: "1000000",
		"D|M"	: "1010101",
}
var jump2Binary = map[string]string{
	"null":"000",
	"JGT":"001",
	"JEQ":"010",
	"JGE":"011",
	"JLT":"100",
	"JNE":"101",
	"JLE":"110",
	"JMP":"111",
}