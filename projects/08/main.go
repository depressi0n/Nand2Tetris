package main

import (
	"fmt"
	"os"
)


type VMTranslator struct{
	files []string
	cursor int
	parser *Parser
	symbolTable *SymbolTable
}

func NewVMTranslator(filenames []string) *VMTranslator{
	return &VMTranslator{
		files: filenames,
		cursor: 0,
	}
}
func(vm *VMTranslator)Run(){
	for ;vm.cursor<len(vm.files);vm.cursor++{
		// Create a New Parser
		parser:=NewParser(vm.files[vm.cursor])
		parser.Run()
	}
}
const UsageHelp=`
Usage:
	vmTranslator [soure.vm]
`
func main(){
	if len(os.Args)<2{
		fmt.Printf("Invalid parameter: length less than 2")
		fmt.Println(UsageHelp)
		os.Exit(1)
	}
	vm:=NewVMTranslator(os.Args[1:])
	vm.Run()
}