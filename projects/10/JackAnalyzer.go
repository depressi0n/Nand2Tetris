package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type JackAnalyzer struct{
	Tokenizer *JackTokenizer
	CompilationEngine *CompilationEngine
	InputFileNames []string
	OutputFileNames []string
}

func NewJackAnalyzer(s string)*JackAnalyzer{
	info,err:=os.Stat(s)
	if err!=nil{
		log.Fatalf("Error when call os.Stat:%v",err)
	}
	inputFileNames:=make([]string,0,100)
	if info.IsDir() {
		files,err:=ioutil.ReadDir(s)
		if err!=nil{
			log.Fatalf("Error when call ioutil.ReadDir:%v",err)
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(),".jack"){
				inputFileNames=append(inputFileNames, info.Name()+"/" +f.Name())
			}
		}
	}else{
		inputFileNames=append(inputFileNames, s)
	}

	outputFileNames:=make([]string,0,len(inputFileNames))
	for i := 0; i < len(inputFileNames); i++ {
		outputFileNames=append(outputFileNames,strings.TrimSuffix(inputFileNames[i],".jack")+"_cmp.xml")
	}
	return &JackAnalyzer{
		// Tokenizer: NewJackTokenizer(s),
		InputFileNames: inputFileNames,
		OutputFileNames: outputFileNames,
	}
}

func(a *JackAnalyzer)Analyze(){
	// tokenize the files
	a.Tokenizer=NewJackTokenizer()
	a.CompilationEngine=NewCompilationEngine(a.Tokenizer)
	for i := 0; i < len(a.InputFileNames); i++ {
		a.Tokenizer.FileName=a.InputFileNames[i]
		a.CompilationEngine.FileName=a.OutputFileNames[i]
		a.CompilationEngine.Compile()
	}
}