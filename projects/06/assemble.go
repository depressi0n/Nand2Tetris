package main

import (
	"log"
	"os"
)

type Assemble struct{
	parser *Parser
	symbolTable *SymbolTable
	output string
}
func NewAssemble(inputASM string,output string)*Assemble{
	symbolTable:=NewSymbolTable()
	parser:=NewParser(inputASM,symbolTable)
	return &Assemble{
		parser:parser,
		symbolTable:symbolTable,
		output:output,
	}
}

func (asm *Assemble)Run(){
	// open the output file
	f,err:=os.OpenFile(asm.output+".hack",os.O_CREATE|os.O_RDWR,0600)
	if err!=nil{
		log.Fatalf("Error in os.OpenFile() with filename %s:%v",asm.output,err)
	}
	defer f.Close()
	machineCodes:=asm.parser.Run()
	for i:=0;i<len(machineCodes);i++{
		_, err = f.WriteString(machineCodes[i])
		if err!=nil{
			log.Fatalf("Error in file.Write() with filename %s:%v",asm.output,err)
		}
		_, err = f.Write([]byte{'\n'})
		if err!=nil{
			log.Fatalf("Error in file.Write() with filename %s:%v",asm.output,err)
		}

	}
	
}