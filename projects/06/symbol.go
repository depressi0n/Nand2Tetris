// 对编译过程中的符号建表，并赋予实际地址
package main

type SymbolTable struct{
	table map[string]int
	baseSymbolAddress int
	baseVariableAddress int
}
func NewSymbolTable()*SymbolTable{
	res:=&SymbolTable{
		table:make(map[string]int),
		baseVariableAddress: VARIABLE_START_ADDRESS,
		baseSymbolAddress: SYMBOL_START_ADDRESS,
	}
	// pre-defined symbol
	res.table["SP"]=0
	res.table["LCL"]=1
	res.table["ARG"]=2
	res.table["THIS"]=3
	res.table["THAT"]=4
	
	res.table["R0"]=0
	res.table["R1"]=1
	res.table["R2"]=2
	res.table["R3"]=3
	res.table["R4"]=4
	res.table["R5"]=5
	res.table["R6"]=6
	res.table["R7"]=7
	res.table["R8"]=8
	res.table["R9"]=9
	res.table["R10"]=10
	res.table["R11"]=11
	res.table["R12"]=12
	res.table["R13"]=13
	res.table["R14"]=14
	res.table["R15"]=15
	
	res.table["SCREEN"]=16384
	res.table["KBD"]=24576
	return res
}
// addEntry add a entry (symbol,address) to symbol table
func(st *SymbolTable)AddVariableEntry(symbol string,address int){
	st.table[symbol]=address+st.baseVariableAddress
}
func(st *SymbolTable)AddSymbolEntry(symbol string,address int){
	st.table[symbol]=address+st.baseSymbolAddress
}

func(st *SymbolTable)GetAddress(symbol string)int{
	address,ok:=st.table[symbol]
	if !ok{
		return -1
	}
	return address
}