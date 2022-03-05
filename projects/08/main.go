package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)


type VMTranslator struct{
	files []string
	// parser *Parser
	symbolTable *SymbolTable
	outputFileName string
}

func NewVMTranslator(s string) *VMTranslator{
	files:=make([]string,0,100)
	if !strings.HasSuffix(s,".vm"){
		// Support directory
		all, err := ioutil.ReadDir(s)
		if err!=nil{
			log.Fatal("Error when read the directory:%v",err)
		}
		// filter the .vm file
		for _, name := range all {
			if !name.IsDir() && strings.HasSuffix(name.Name(),".vm") {
				files = append(files, s+"/"+name.Name())
			}
		}
	}
	fileNameWithSuffix:=path.Base(s)
   	fileType:=path.Ext(fileNameWithSuffix)
   	outputFileName:=strings.TrimSuffix(fileNameWithSuffix, fileType)+".asm"
	return &VMTranslator{
		files: files,
		outputFileName: s+"/"+outputFileName,
	}
}
func(vm *VMTranslator)Run(){
	// Create a New Parser
	parser:=NewParser(vm.files,vm.outputFileName)
	parser.Run()
}
const UsageHelp=`
Usage:
	vmTranslator [soure]
`
func main(){
	if len(os.Args)<2{
		fmt.Printf("Invalid parameter: length less than 2")
		fmt.Println(UsageHelp)
		os.Exit(1)
	}
	vm:=NewVMTranslator(os.Args[1])
	vm.Run()
}