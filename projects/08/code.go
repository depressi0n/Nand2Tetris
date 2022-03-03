package main

import (
	"log"
	"strings"
)

type SymbolTable struct{
	Table map[string]int
}

func NewSymbolTable()*SymbolTable{
	res:=&SymbolTable{}
	res.Table=map[string]int{
		"R0":0,
		"R1":1,
		"R2":2,
		"R3":3,
		"R4":4,
		"R5":5, // 用于保存temp段的内容
		"R6":6, // 用于保存temp段的内容
		"R7":7, // 用于保存temp段的内容
		"R8":8, // 用于保存temp段的内容
		"R9":9, // 用于保存temp段的内容
		"R10":10, // 用于保存temp段的内容
		"R11":11, // 用于保存temp段的内容
		"R12":12, // 用于保存temp段的内容
		"R13":13, // 可用作通用寄存器
		"R14":14, // 可用作通用寄存器
		"R15":15, // 可用作通用寄存器
		"SP":0, // 指向栈顶的下一个位置
		"LCL":1, // 指向当前VM函数的local段基址
		"ARG":2, // 指向当前VM函数的argument段基址
		"THIS":3, // 指向当前this段（堆中）的基址
		"THAT":4, //指向当前that段（堆中）的基地
	}
	return res
}

type CommandType int
const (
	C_ARITHMETIC CommandType=iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)
type Operator int
const (
	OP_ADD Operator=iota
	OP_SUB
	OP_NEG
	OP_EQ
	OP_GT
	OP_LT
	OP_AND
	OP_OR
	OP_NOT
)

var Str2CommandType= map[string]CommandType{
	"add":C_ARITHMETIC,
	"sub":C_ARITHMETIC,
	"neg":C_ARITHMETIC,
	"eq":C_ARITHMETIC,
	"gt":C_ARITHMETIC,
	"lt":C_ARITHMETIC,
	"and":C_ARITHMETIC,
	"or":C_ARITHMETIC,
	"not":C_ARITHMETIC,

	"push":C_PUSH,
	"pop":C_POP,
	"label":C_LABEL,
	"goto":C_GOTO,
	"if-goto":C_IF,
	"function":C_FUNCTION,
	"return":C_RETURN,
	"call":C_CALL,
}
var CommandType2Str= map[CommandType]string{
	C_ARITHMETIC:"arithmetic",
	C_PUSH:"push",
	C_POP:"pop",
	C_LABEL:"label",
	C_GOTO:"goto",
	C_IF:"if-goto",
	C_FUNCTION:"function",
	C_RETURN:"return",
	C_CALL:"call",
}

var Str2Op = map[string]Operator {
	"add":OP_ADD,
	"sub":OP_SUB,
	"neg":OP_NEG,
	"eq":OP_EQ,
	"gt":OP_GT,
	"lt":OP_LT,
	"and":OP_AND,
	"or":OP_OR,
	"not":OP_NOT,
}
var Op2Str= map[Operator]string{
	OP_ADD:"ADD",
	OP_SUB:"SUB",
	OP_NEG:"NEG",
	OP_EQ:"EQ",
	OP_GT:"GT",
	OP_LT:"LT",
	OP_AND:"AND",
	OP_OR:"OR",
	OP_NOT:"NOT",
}

type SegmentType int
const(
	ARGUMENT SegmentType=iota
	LOCAL
	STATIC
	CONSTANT
	THIS
	THAT
	POINTER
	TEMP
)
var Str2SegmentType = map[string]SegmentType{
	"argument":ARGUMENT,
	"local":LOCAL,
	"static":STATIC,
	"constant":CONSTANT,
	"this":THIS,
	"that":THAT,
	"pointer":POINTER,
	"temp":TEMP,
}
var SegmentType2Str = map[SegmentType]string{
	ARGUMENT:"argument",
	LOCAL:"local",
	STATIC:"static",
	CONSTANT:"constant",
	THIS:"this",
	THAT:"that",
	POINTER:"pointer",
	TEMP:"temp",
}

type Command struct{
	commandType CommandType
	operatorType Operator
	segmentType SegmentType
	args []string
}
func NewCommand(s string)*Command{
	var commandType CommandType
	for prefix,v  := range Str2CommandType {
		if strings.HasPrefix(s, prefix){
			commandType=v
			if commandType != C_ARITHMETIC {
				s=strings.TrimPrefix(s,prefix)
			}
			break
		}
	}
	var operator Operator
	if commandType == C_ARITHMETIC{
		for prefix,v := range Str2Op {
			if strings.HasPrefix(s,prefix){
				operator=v
				s=strings.TrimPrefix(s,prefix)
			}
		}
	}
	s=strings.TrimSpace(s)
	args:=strings.Split(s," ")
	for i:=0;i<len(args);i++{
		args[i]=strings.TrimSpace(args[i])
	}
	return &Command{
		commandType: commandType,
		operatorType: operator,
		args:args,
	}

}

func(c *Command)Type()CommandType{
	return c.commandType
}
func(c *Command)OperatorType()Operator{
	return c.operatorType
}
func(c *Command)SegmentType()SegmentType{
	if c.Type()!=C_PUSH && c.Type()!=C_POP{
		log.Fatalf("Command.SegmentType shoud not be call by type:%s",CommandType2Str[c.Type()])
	}

	return Str2SegmentType[c.Arg1()]
}
func(c *Command)Arg1()string{
	switch c.commandType {
	case C_ARITHMETIC:
		return Op2Str[c.operatorType]
	case C_RETURN:
		log.Fatalf("C_RETURN should not call the Arg1 method")
		return c.args[0]
	default:
		return c.args[0]
	}
	
}
func(c *Command)Arg2()string{
	switch c.commandType {
	case C_PUSH:
		return c.args[1]
	case C_POP:
		return c.args[1]
	case C_FUNCTION|C_CALL:
		return c.args[1]
	default:
		log.Fatalf("Should not call the Arg2 method")
		return c.args[1]
	}
}

var seg2Addr =  map[string]int{
	"argument":2, // 基础地址存放在RAM[2]中
	"local":1, // 基址存放在RAM[1]中
	"static":16, // 从RAM单元地址16开始分配
	// "constant": null,  // 真正虚拟的段
	"this":3, // RAM[3]
	"that":4, // RAM[4]
	// "pointer":, // 暂时未知
	// "temp":, // RAM[5-12]
}