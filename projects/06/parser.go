// Parser负责对输入文件进行语法分析
package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type Parser struct{
	symbolTable *SymbolTable
	fileName string // file path
	currentCommand string
	instructions []*Command
	symbolCnt int
	variableCnt int
	pc int
}
func NewParser(fileName string,symbolTable *SymbolTable)*Parser{
	return &Parser{
		symbolTable:symbolTable,
		fileName: fileName,
		currentCommand:"",
		symbolCnt: 0,
		variableCnt: 0,
	}
}
// hasMoreCommands show if there are more command
func(p *Parser)hasMoreCommands()bool{
	return false
}

// read read the next command
func(p *Parser)read(){

}

func(p *Parser)BuildSymbolTable(){
	// open the assemble file
	f,err:=os.OpenFile(p.fileName,os.O_RDONLY,0600)
	if err!=nil{
		log.Fatalf("os.OpenFile has error with %s:%v",p.fileName,err)
	}
	defer f.Close()
	input:=bufio.NewReader(f)
	p.instructions=make([]*Command,0,100)
	for {
		// Read a line
		tmp, _, err := input.ReadLine()
		if err==io.EOF{
			break
		}
		line:=string(tmp)
		// remove white space
		line= strings.TrimSpace(line)
		// remove empty lines and the comment lines
		if len(line)==0 || strings.HasPrefix(line,"//")  {
			continue
		}
		// remove comment after instruction
		commentIndex:=strings.LastIndex(line,"//")
		if commentIndex != -1{
			line=line[:commentIndex]
		}
		// remove white space
		line= strings.TrimSpace(line)
		// Pre-process: Now do nothing but there will handle something such as Macro
		command:=NewCommand(line)
		commandType:=command.Type()
		if commandType!=L_COMMAND{
			p.pc++
		}
		if command.Type() == L_COMMAND {
			p.symbolTable.AddSymbolEntry(command.symbol,p.pc)
			p.symbolCnt++
		}
		p.instructions=append(p.instructions,command)
	}
}

// Run parse the asm file, and translate it to machine code
func(p *Parser)Run()[]string{
	p.BuildSymbolTable()
	// Something such as optimization can be do here, but now nothing

	// resolve symbol in the instruction
	// and transfer to machine code
	outputHack:=make([]string,0,len(p.instructions))
	for i:=0;i<len(p.instructions);i++ {
		instruction:=p.instructions[i]
		switch instruction.Type(){
		case A_COMMAND:
			if instruction.AInstructionType() {
				name:=instruction.Symbol()
				t:=p.symbolTable.GetAddress(name)
				if t==-1{
					// variable
					p.symbolTable.AddVariableEntry(name,p.variableCnt)
					t=p.variableCnt+VARIABLE_START_ADDRESS
					p.variableCnt++
				}
				outputHack=append(outputHack, "0"+Decimal2Binary(t))
			}else{
				outputHack=append(outputHack, "0"+Decimal2Binary(instruction.AValue()))
			}
			
		case C_COMMAND:
			dBin,ok:=dest2Binary[instruction.Dest()]
			if !ok{
				log.Fatalf("Invalid instruction '%s' in dest=comp;jump mode: invalid dest field\n",instruction.Origin())
			}
			cBin,ok:=comp2Binary[instruction.Comp()]
			if !ok{
				log.Fatalf("Invalid instruction '%s' in dest=comp;jump mode: invalid comp field\n",instruction.Origin())
			}
			jBin,ok:=jump2Binary[instruction.Jump()]
			if !ok{
				log.Fatalf("Invalid instruction '%s' in dest=comp;jump mode: invalid jump field\n",instruction.Origin())
			}
			outputHack=append(outputHack, "111"+cBin+dBin+jBin)
		case L_COMMAND:
			// do-nothing
		default:
			log.Fatal("Invalid Command Type")
	}
	
		// Check the validity of machine code

		// Append the got machine code to result
	}

	// Do something basic check for got Hack program
	if len(p.instructions) == 0 {
		log.Printf("No instructions found in file named %s\n",p.fileName)
	}else if p.instructions[len(p.instructions)-1].Type()!=C_COMMAND{
		log.Printf("Last instruction should be a JUMP instruction")
	}
	
	return outputHack
}