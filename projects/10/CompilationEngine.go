package main

import (
	"encoding/xml"
	"log"
	"os"
)
 
type CompilationEngine struct{
	FileName string
	outputWriter *os.File
	encoder *xml.Encoder
	*JackTokenizer
	start xml.StartElement
	indentLevel int
}

func NewCompilationEngine(tokenizer *JackTokenizer) *CompilationEngine{
	return &CompilationEngine{
		JackTokenizer:tokenizer,
	}
}
func(c *CompilationEngine)Compile(){
	// Start For testing 

	// c.tokenizer.Tokenize()
	// f, err := os.OpenFile(c.FileName+".txt", os.O_CREATE|os.O_RDWR, 0600)
	// if err!=nil{
	// 	log.Fatalf("Error in opening the file named %s with os.OpenFile",c.FileName)
	// }
	// defer f.Close()
	// for i := 0; i < len(c.tokenizer.TokenStream); i++ {
	// 	s:=c.tokenizer.TokenStream[i].String()+"\n"
	// 	f.WriteString(s)
	// }

	// End For testing

	f, err := os.OpenFile(c.FileName, os.O_CREATE|os.O_RDWR, 0600)
	if err!=nil{
		log.Fatalf("Error in opening the file named %s with os.OpenFile",c.FileName)
	}
	defer f.Close()
	c.outputWriter=f
	c.encoder=xml.NewEncoder(f)
	c.Tokenize()
	c.CompileClass()
	c.encoder.Flush()
}
// class className { classVarDec* subroutineDec* }
func(c *CompilationEngine)CompileClass(){
	c.Advance()
	if !(c.CurrentToken.Type == T_KEYWORD && c.CurrentToken.Keyword == K_CLASS) {
		log.Fatalf("Syntax error on class")
	}
	c.start.Name.Local = "class"
	c.encoder.EncodeToken(c.start)
	Keyword2XML[K_CLASS].MarshalXML(c.encoder,c.start)


	// read the className
	c.Advance()
	if c.CurrentToken.Type != T_IDENTIFIER {
		log.Fatalf("Syntax error on class")
	}
	className:=&XMLIdentifier{
		Name:c.CurrentToken.Identifier,
	}
	className.MarshalXML(c.encoder,c.start)

	// read the {
	c.Advance()
	if !(c.CurrentToken.Type == T_SYMBOL && c.CurrentToken.Symbol  == '{' ){
		log.Fatalf("Syntax error on class")
	}
	Symbol2XML['{'].MarshalXML(c.encoder,c.start)
	
	
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error:")
	}
	// class body
	for !(c.HasMoreToken()&&(c.TokenStream[0].Type== T_SYMBOL && c.TokenStream[0].Symbol=='}')){
		// Distinguish between subroutineDec and calssVarDec
		nextToken:=c.TokenStream[0]
		if nextToken.Type != T_KEYWORD {
			log.Fatalf("Syntax error on class:%s",className.Name)
		}
		switch nextToken.Keyword{
		case K_STATIC:
			fallthrough
		case K_FIELD:
			c.CompileClassVarDec()
		case K_CONSTRUCTOR:
			fallthrough
		case K_FUNCTION:
			fallthrough
		case K_METHOD:
			c.CompileSubroutineDec()
		default:
			log.Fatalf("Syntax error on class:%s",className.Name)
		}
	}

	// read the }
	c.Advance()
	if !(c.CurrentToken.Type == T_SYMBOL && c.CurrentToken.Symbol  == '}' ){
		log.Fatalf("Syntax error on class")
	}
	Symbol2XML['}'].MarshalXML(c.encoder,c.start)

	c.start.Name.Local = "class"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}
// static|fiele type varName [,varName]*;
func(c *CompilationEngine)CompileClassVarDec(){
	c.start.Name.Local = "classVarDec"
	c.encoder.EncodeToken(c.start)

	// read the 'static' or 'field'
	c.Advance()
	Keyword2XML[c.CurrentToken.Keyword].MarshalXML(c.encoder,c.start)
	
	// read the return type
	c.Advance()
	switch c.CurrentToken.Type {
	case T_KEYWORD:
		if c.CurrentToken.Keyword == K_VOID ||c.CurrentToken.Keyword == K_INT ||c.CurrentToken.Keyword == K_CHAR ||c.CurrentToken.Keyword == K_BOOLEAN {
			Keyword2XML[c.CurrentToken.Keyword].MarshalXML(c.encoder,c.start)
		}else{
			log.Fatalf("Syntax error on CompileClassVarDec")
		}
	case T_IDENTIFIER:
		retType:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
		retType.MarshalXML(c.encoder,c.start)
	default:
		log.Fatalf("Syntax error on CompileClassVarDec")
	}
	
	// varName
	c.Advance()
	if c.CurrentToken.Type != T_IDENTIFIER{
		log.Fatalf("Syntax error on CompileClassVarDec")
	}
	varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
	varName.MarshalXML(c.encoder,c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileClassVarDec")
	}
	if c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol == ',' {
		for !(c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol == ';'){
			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on CompileClassVarDec")
			}
			c.Advance()
			Symbol2XML[','].MarshalXML(c.encoder,c.start)
			if c.TokenStream[0].Type != T_IDENTIFIER {
				log.Fatalf("Syntax error on CompileClassVarDec")
			}
			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on CompileClassVarDec")
			}
			c.Advance()
			varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
			varName.MarshalXML(c.encoder,c.start)
			
			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on CompileClassVarDec")
			}
		}
	}
	
	c.Advance()
	Symbol2XML[';'].MarshalXML(c.encoder,c.start)

	c.start.Name.Local = "classVarDec"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}
// constructor|function|method void|type subroutineName ( parameterList ) subroutineBody
func(c *CompilationEngine)CompileSubroutineDec(){
	c.start.Name.Local = "subroutineDec"
	c.encoder.EncodeToken(c.start)
	
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	c.Advance()
	// constructor|function|method
	Keyword2XML[c.CurrentToken.Keyword].MarshalXML(c.encoder,c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	c.Advance()
	switch c.CurrentToken.Type{
	case T_IDENTIFIER:
		retTypeName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
		retTypeName.MarshalXML(c.encoder,c.start)
	case T_KEYWORD:
		if c.CurrentToken.Keyword == K_VOID ||c.CurrentToken.Keyword == K_INT ||c.CurrentToken.Keyword == K_CHAR ||c.CurrentToken.Keyword == K_BOOLEAN {
			Keyword2XML[c.CurrentToken.Keyword].MarshalXML(c.encoder,c.start)
		}else{
			log.Fatalf("Syntax error on CompileSubroutineDec")
		}
	default:
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	// name
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	c.Advance()
	if c.CurrentToken.Type != T_IDENTIFIER{
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	suroutineName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
	suroutineName.MarshalXML(c.encoder,c.start)

	// (
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	c.Advance()
	if !(c.CurrentToken.Type == T_SYMBOL && c.CurrentToken.Symbol == '('  ){
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	Symbol2XML['('].MarshalXML(c.encoder,c.start)

	// parameterList | )
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	switch c.TokenStream[0].Type {
	case T_SYMBOL :
		c.start.Name.Local = "parameterList"
		c.encoder.EncodeToken(c.start)
		c.start.Name.Local = "parameterList"
		c.encoder.EncodeToken(xml.EndElement{c.start.Name})
	case T_IDENTIFIER:
		fallthrough
	case T_KEYWORD:
		c.CompileParameterList()
	default:
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	c.Advance()
	if c.CurrentToken.Type == T_SYMBOL && c.CurrentToken.Symbol==')'{
		Symbol2XML[')'].MarshalXML(c.encoder,c.start)
	}else{
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	if !(c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol == '{'){
		log.Fatalf("Syntax error on CompileSubroutineDec")
	}
	c.CompileSubroutineBody()

	c.start.Name.Local = "subroutineDec"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

func (c *CompilationEngine)CompileSubroutineBody() {
	c.start.Name.Local = "subroutineBody"
	c.encoder.EncodeToken(c.start)

	// {
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileSubroutineBody")
	}
	c.Advance()
	if !(c.CurrentToken.Type == T_SYMBOL && c.CurrentToken.Symbol == '{'){
		log.Fatalf("Syntax error on CompileSubroutineBody")
	}
	Symbol2XML['{'].MarshalXML(c.encoder,c.start)

	if c.TokenStream[0].Type == T_KEYWORD {
		for c.TokenStream[0].Type == T_KEYWORD && c.TokenStream[0].Keyword == K_VAR{
			// classVarDec*
			c.compileVarDec()
		}
		switch c.TokenStream[0].Keyword {
			case K_LET:
				fallthrough
			case K_IF:
				fallthrough
			case K_WHILE:
				fallthrough
			case K_DO:
				fallthrough
			case K_RETURN:
				// subroutineDec*
				c.compileStatements()
			default:
				log.Fatalf("Syntax error on CompileSubroutineBody")
		}
	}

	// }
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileSubroutineBody")
	}
	if c.TokenStream[0].Symbol != '}'{
		log.Fatalf("Syntax error on CompileSubroutineBody")
	}else{
		c.Advance()
		Symbol2XML['}'].MarshalXML(c.encoder,c.start)
	}
	c.start.Name.Local = "subroutineBody"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
	return
}
// [type varName [,type varName]*] ?
func(c *CompilationEngine)CompileParameterList(){
	c.start.Name.Local = "parameterList"
	c.encoder.EncodeToken(c.start)
	// type
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileParameterList")
	}
	c.Advance()
	switch c.CurrentToken.Type {
	case T_KEYWORD:
		if c.CurrentToken.Keyword == K_VOID ||c.CurrentToken.Keyword == K_INT ||c.CurrentToken.Keyword == K_CHAR ||c.CurrentToken.Keyword == K_BOOLEAN {
			Keyword2XML[c.CurrentToken.Keyword].MarshalXML(c.encoder,c.start)
		}else{
			log.Fatalf("Syntax error on CompileParameterList")
		}
	case T_IDENTIFIER:
		typeName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
		typeName.MarshalXML(c.encoder,c.start)
	default:
		log.Fatalf("Syntax error on CompileParameterList")
	}

	// varName
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileParameterList")
	}
	c.Advance()
	if !(c.CurrentToken.Type == T_IDENTIFIER){
		log.Fatalf("Syntax error on CompileParameterList")
	}
	varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
	varName.MarshalXML(c.encoder,c.start)

	// [, type varName]*
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileParameterList")
	}
	for c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol == ',' {
		c.Advance()
		Symbol2XML[','].MarshalXML(c.encoder,c.start)
		
		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on CompileParameterList")
		}
		c.Advance()
		switch c.CurrentToken.Type {
		case T_KEYWORD:
			if c.CurrentToken.Keyword == K_VOID ||c.CurrentToken.Keyword == K_INT ||c.CurrentToken.Keyword == K_CHAR ||c.CurrentToken.Keyword == K_BOOLEAN {
				Keyword2XML[c.CurrentToken.Keyword].MarshalXML(c.encoder,c.start)
			}else{
				log.Fatalf("Syntax error on CompileParameterList")
			}
		case T_IDENTIFIER:
			typeName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
			typeName.MarshalXML(c.encoder,c.start)
		default:
			log.Fatalf("Syntax error on CompileParameterList")
		}

		// varName
		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on CompileParameterList")
		}
		c.Advance()
		if !(c.CurrentToken.Type == T_IDENTIFIER){
			log.Fatalf("Syntax error on CompileParameterList")
		}
		varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
		varName.MarshalXML(c.encoder,c.start)
	}

	c.start.Name.Local = "parameterList"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}
// var type varName (,varName)* ;
func(c *CompilationEngine)compileVarDec(){
	c.start.Name.Local = "varDec"
	c.encoder.EncodeToken(c.start)
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileVarDec")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_KEYWORD && c.CurrentToken.Keyword==K_VAR){
		log.Fatalf("Syntax error on compileVarDec")
	}
	Keyword2XML[K_VAR].MarshalXML(c.encoder,c.start)
	
	// type
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileVarDec")
	}
	c.Advance()
	switch c.CurrentToken.Type {
	case T_KEYWORD:
		if c.CurrentToken.Keyword == K_VOID ||c.CurrentToken.Keyword == K_INT ||c.CurrentToken.Keyword == K_CHAR ||c.CurrentToken.Keyword == K_BOOLEAN {
			Keyword2XML[c.CurrentToken.Keyword].MarshalXML(c.encoder,c.start)
		}else{
			log.Fatalf("Syntax error on compileVarDec")
		}
	case T_IDENTIFIER:
		typeName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
		typeName.MarshalXML(c.encoder,c.start)
	default:
		log.Fatalf("Syntax error on compileVarDec")
	}

	// varName
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileVarDec")
	}
	c.Advance()
	if !(c.CurrentToken.Type == T_IDENTIFIER){
		log.Fatalf("Syntax error on compileVarDec")
	}
	varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
	varName.MarshalXML(c.encoder,c.start)

	// (, varName)*
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileVarDec")
	}
	for c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol == ',' {
		c.Advance()
		Symbol2XML[','].MarshalXML(c.encoder,c.start)

		// varName
		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on compileVarDec")
		}
		c.Advance()
		if !(c.CurrentToken.Type == T_IDENTIFIER){
			log.Fatalf("Syntax error on compileVarDec")
		}
		varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
		varName.MarshalXML(c.encoder,c.start)
	}

	// ; 
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileVarDec")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL&& c.CurrentToken.Symbol==';'){
		log.Fatalf("Syntax error on compileVarDec")
	}
	Symbol2XML[';'].MarshalXML(c.encoder,c.start)
	c.start.Name.Local = "varDec"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}
// statement*
func(c *CompilationEngine)compileStatements(){
	c.start.Name.Local = "statements"
	c.encoder.EncodeToken(c.start)
	
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileVarDec")
	}
	for !(c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol == '}'){
		switch c.TokenStream[0].Keyword {
		case K_LET:
			fallthrough
		case K_DO:
			fallthrough
		case K_IF:
			fallthrough
		case K_WHILE:
			fallthrough
		case K_RETURN:
			c.compileStatement()
		default:
			log.Fatalf("Syntax error on compileStatements")
		}
	}
	
	c.start.Name.Local = "statements"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

// letStatement|ifStatement|whileStatement|doStatement|returnStatement
func(c *CompilationEngine)compileStatement(){
	// c.start.Name.Local = "statement"
	// c.encoder.EncodeToken(c.start)
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileStatement ")
	}
	switch c.TokenStream[0].Keyword {
		case K_LET:
			c.compileLetStatement()
		case K_DO:
			c.compileDoStatement()
		case K_IF:
			c.compileIfStatement()
		case K_WHILE:
			c.compileWhileStatement()
		case K_RETURN:
			c.compileReturnStatement()
		default:
			log.Fatalf("Syntax error on compileStatement")
	}
	// c.start.Name.Local = "statement"
	// c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

// let varName [[ expression ]]? = expression ;
func(c *CompilationEngine)compileLetStatement(){
	c.start.Name.Local = "letStatement"
	c.encoder.EncodeToken(c.start)
	
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileLetStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_KEYWORD && c.CurrentToken.Keyword==K_LET){
		log.Fatalf("Syntax error on compileLetStatement")
	}
	Keyword2XML[K_LET].MarshalXML(c.encoder,c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileLetStatement")
	}
	c.Advance()
	if c.CurrentToken.Type!=T_IDENTIFIER{
		log.Fatalf("Syntax error on compileLetStatement")
	}
	varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
	varName.MarshalXML(c.encoder,c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileLetStatement")
	}
	if c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol=='['{
		c.Advance()
		Symbol2XML['['].MarshalXML(c.encoder,c.start)

		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on compileLetStatement")
		}
		c.CompileExpression()

		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on compileLetStatement")
		}
		c.Advance()
		if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol==']'){
			log.Fatalf("Syntax error on compileLetStatement")
		}
		Symbol2XML[']'].MarshalXML(c.encoder,c.start)
	}

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileLetStatement")
	}
	if  !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol=='=') {
		log.Fatalf("Syntax error on compileLetStatement")
	}
	c.Advance()
	Symbol2XML['='].MarshalXML(c.encoder,c.start)

	c.CompileExpression()


	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileLetStatement")
	}
	c.Advance()
	if  !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol==';') {
		log.Fatalf("Syntax error on compileLetStatement")
	}
	Symbol2XML[';'].MarshalXML(c.encoder,c.start)

	c.start.Name.Local = "letStatement"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

// do subroutineCall ;
func(c *CompilationEngine)compileDoStatement(){
	c.start.Name.Local = "doStatement"
	c.encoder.EncodeToken(c.start)
	// do
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_KEYWORD && c.CurrentToken.Keyword==K_DO){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	Keyword2XML[K_DO].MarshalXML(c.encoder,c.start)

	// id
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	c.Advance()
	if c.CurrentToken.Type!=T_IDENTIFIER{
		log.Fatalf("Syntax error on compileDoStatement")
	}
	name:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
	name.MarshalXML(c.encoder,c.start)
	// .id
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	if !(c.TokenStream[0].Symbol=='(' || c.TokenStream[0].Symbol=='.'){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	if c.TokenStream[0].Symbol=='.' {
		c.Advance()
		Symbol2XML['.'].MarshalXML(c.encoder,c.start)

		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on compileDoStatement")
		}
		c.Advance()
		if !(c.CurrentToken.Type == T_IDENTIFIER){
			log.Fatalf("Syntax error on compileDoStatement")
		}
		varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
		varName.MarshalXML(c.encoder,c.start)

		if !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol=='('){
			log.Fatalf("Syntax error on compileDoStatement")
		}
	}
	// ( expressionList )
	c.Advance()
	Symbol2XML['('].MarshalXML(c.encoder,c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	c.CompileExpressionList()

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	if !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol==')'){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	c.Advance()
	Symbol2XML[')'].MarshalXML(c.encoder,c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	if !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol==';'){
		log.Fatalf("Syntax error on compileDoStatement")
	}
	c.Advance()
	Symbol2XML[';'].MarshalXML(c.encoder,c.start)

	c.start.Name.Local = "doStatement"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

// while ( expression ) { statements }
func(c *CompilationEngine)compileWhileStatement(){
	c.start.Name.Local = "whileStatement"
	c.encoder.EncodeToken(c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compilewhileStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_KEYWORD && c.CurrentToken.Keyword==K_WHILE){
		log.Fatalf("Syntax error on compilewhileStatement")
	}
	Keyword2XML[K_WHILE].MarshalXML(c.encoder,c.start)

	// ( expression )
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileWhileStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol=='('){
		log.Fatalf("Syntax error on compileWhileStatement")
	}
	Symbol2XML['('].MarshalXML(c.encoder,c.start)

	c.CompileExpression()

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileWhileStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol==')'){
		log.Fatalf("Syntax error on compileWhileStatement")
	}
	Symbol2XML[')'].MarshalXML(c.encoder,c.start)

	// { statements }
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileWhileStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol=='{'){
		log.Fatalf("Syntax error on compileWhileStatement")
	}
	Symbol2XML['{'].MarshalXML(c.encoder,c.start)

	c.compileStatements()

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileWhileStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol=='}'){
		log.Fatalf("Syntax error on compileWhileStatement")
	}
	Symbol2XML['}'].MarshalXML(c.encoder,c.start)


	c.start.Name.Local = "whileStatement"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

// if ( expression ) { statements } [else { statements }]?
func(c *CompilationEngine)compileIfStatement(){
	c.start.Name.Local = "ifStatement"
	c.encoder.EncodeToken(c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_KEYWORD && c.CurrentToken.Keyword==K_IF){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	Keyword2XML[K_IF].MarshalXML(c.encoder,c.start)

	// ( expression )
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol=='('){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	Symbol2XML['('].MarshalXML(c.encoder,c.start)

	c.CompileExpression()

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol==')'){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	Symbol2XML[')'].MarshalXML(c.encoder,c.start)

	// { statements }
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol=='{'){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	Symbol2XML['{'].MarshalXML(c.encoder,c.start)

	c.compileStatements()

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol=='}'){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	Symbol2XML['}'].MarshalXML(c.encoder,c.start)


	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileIfStatement")
	}
	if c.TokenStream[0].Type==T_KEYWORD && c.TokenStream[0].Keyword==K_ELSE {
		c.Advance()
		Keyword2XML[K_ELSE].MarshalXML(c.encoder,c.start)
		// { statements }
		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on compileIfStatement")
		}
		c.Advance()
		if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol=='{'){
			log.Fatalf("Syntax error on compileIfStatement")
		}
		Symbol2XML['{'].MarshalXML(c.encoder,c.start)

		c.compileStatements()

		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on compileIfStatement")
		}
		c.Advance()
		if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol=='}'){
			log.Fatalf("Syntax error on compileIfStatement")
		}
		Symbol2XML['}'].MarshalXML(c.encoder,c.start)
	}

	c.start.Name.Local = "ifStatement"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

// return expression? ;
func(c *CompilationEngine)compileReturnStatement(){
	c.start.Name.Local = "returnStatement"
	c.encoder.EncodeToken(c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileReturnStatment")
	}
	c.Advance()
	if !(c.CurrentToken.Type==T_KEYWORD && c.CurrentToken.Keyword==K_RETURN){
		log.Fatalf("Syntax error on compileReturnStatment")
	}
	Keyword2XML[K_RETURN].MarshalXML(c.encoder,c.start)

	// expression?
	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileReturnStatment")
	}

	if (c.TokenStream[0].Type==T_INT_CONST) || 
	(c.TokenStream[0].Type==T_STRING_CONST) ||
	(c.TokenStream[0].Type==T_IDENTIFIER)||
	(c.TokenStream[0].Type==T_KEYWORD && (c.TokenStream[0].Keyword == K_THIS ||c.TokenStream[0].Keyword == K_TRUE|| c.TokenStream[0].Keyword == K_FALSE || c.TokenStream[0].Keyword == K_NULL )) ||
	(c.TokenStream[0].Type==T_SYMBOL && (c.TokenStream[0].Symbol == '(' || c.TokenStream[0].Symbol == '~' ||c.TokenStream[0].Symbol == '-')) {
		c.CompileExpression()
	}

	c.Advance()
	if !(c.CurrentToken.Type==T_SYMBOL && c.CurrentToken.Symbol==';'){
		log.Fatalf("Syntax error on compileReturnStatment")
	}
	Symbol2XML[';'].MarshalXML(c.encoder,c.start)

	c.start.Name.Local = "returnStatement"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

// term [op term]*
func(c *CompilationEngine)CompileExpression(){
	c.start.Name.Local = "expression"
	c.encoder.EncodeToken(c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileExpression")
	}
	if (c.TokenStream[0].Type==T_INT_CONST) || 
	(c.TokenStream[0].Type==T_STRING_CONST) ||
	(c.TokenStream[0].Type==T_IDENTIFIER)||
	(c.TokenStream[0].Type==T_KEYWORD && (c.TokenStream[0].Keyword == K_THIS ||c.TokenStream[0].Keyword == K_TRUE|| c.TokenStream[0].Keyword == K_FALSE || c.TokenStream[0].Keyword == K_NULL )) ||
	(c.TokenStream[0].Type==T_SYMBOL && (c.TokenStream[0].Symbol == '(' || c.TokenStream[0].Symbol == '~' ||c.TokenStream[0].Symbol == '-')) {
		c.CompileTerm()
	}else{
		log.Fatalf("Syntax error on CompileExpression")
	}

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileExpression")
	}
	if c.TokenStream[0].Type==T_SYMBOL && (c.TokenStream[0].Symbol == ']' || c.TokenStream[0].Symbol == ')'){
		c.start.Name.Local = "expression"
		c.encoder.EncodeToken(xml.EndElement{c.start.Name})
		return
	}
	for c.TokenStream[0].Type==T_SYMBOL && 
	(c.TokenStream[0].Symbol == '+' || c.TokenStream[0].Symbol == '-' ||
	c.TokenStream[0].Symbol == '*' || c.TokenStream[0].Symbol == '/' ||
	c.TokenStream[0].Symbol == '&' || c.TokenStream[0].Symbol == '|' ||
	c.TokenStream[0].Symbol == '<' || c.TokenStream[0].Symbol == '>' ||
	c.TokenStream[0].Symbol == '='){
		c.Advance()
		Symbol2XML[c.CurrentToken.Symbol].MarshalXML(c.encoder,c.start)

		if (c.TokenStream[0].Type==T_INT_CONST) || 
		(c.TokenStream[0].Type==T_STRING_CONST) ||
		(c.TokenStream[0].Type==T_IDENTIFIER)||
		(c.TokenStream[0].Type==T_SYMBOL && (c.TokenStream[0].Symbol == '(' || c.TokenStream[0].Symbol == '~' ||c.TokenStream[0].Symbol == '-')) {
			c.CompileTerm()
		}else{
			log.Fatalf("Syntax error on CompileExpression")
		}

	}
	c.start.Name.Local = "expression"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
	return
	
}
// integerConstant|stringConstant|keywordConstant|varName|varName [ expression ] |subroutineCall| ( expression ) | unaryOp term
func(c *CompilationEngine)CompileTerm(){
	c.start.Name.Local = "term"
	c.encoder.EncodeToken(c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on compileTermStatement")
	}
	c.Advance()
	switch c.CurrentToken.Type{
	case T_INT_CONST:
		value:=&XMLIntegerConstant{Value: c.CurrentToken.IntVal}
		value.MarshalXML(c.encoder,c.start)

	case T_STRING_CONST:
		strValue:=&XMLStringConstant{Name: c.CurrentToken.StringVal}
		strValue.MarshalXML(c.encoder,c.start)
		
	case T_IDENTIFIER:
		name:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
		name.MarshalXML(c.encoder,c.start)

		if !c.HasMoreToken(){
			log.Fatalf("Syntax error on compileTermStatement")
		}
		if c.TokenStream[0].Type !=  T_SYMBOL {
			log.Fatalf("Syntax error on compileTermStatement")
		}else if c.TokenStream[0].Symbol=='['{
			c.Advance()
			Symbol2XML['['].MarshalXML(c.encoder,c.start)
			
			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.CompileExpression()

			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			if !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol==']'){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.Advance()
			Symbol2XML[']'].MarshalXML(c.encoder,c.start)
		}else if c.TokenStream[0].Symbol=='('{
			c.Advance()
			
			Symbol2XML['('].MarshalXML(c.encoder,c.start)

			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.CompileExpressionList()

			if !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol==')'){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.Advance()
			Symbol2XML[')'].MarshalXML(c.encoder,c.start)
		}else if c.TokenStream[0].Symbol=='.'{
			// . id 
			c.Advance()
			Symbol2XML['.'].MarshalXML(c.encoder,c.start)

			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.Advance()
			if !(c.CurrentToken.Type == T_IDENTIFIER){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			varName:=&XMLIdentifier{Name: c.CurrentToken.Identifier}
			varName.MarshalXML(c.encoder,c.start)
			
			// (
			if !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol=='('){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.Advance()
			
			Symbol2XML['('].MarshalXML(c.encoder,c.start)

			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.CompileExpressionList()

			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			
			if !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol==')'){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.Advance()
			Symbol2XML[')'].MarshalXML(c.encoder,c.start)
		}

	case T_KEYWORD:
		switch c.CurrentToken.Keyword {
		case K_TRUE:
			fallthrough
		case K_FALSE:
			fallthrough
		case K_NULL:
			fallthrough
		case K_THIS:
			Keyword2XML[c.CurrentToken.Keyword].MarshalXML(c.encoder,c.start)
		default:
			log.Fatalf("Syntax error on compileTermStatement")
			
		}
	case T_SYMBOL:
		switch c.CurrentToken.Symbol {
		case '(':
			Symbol2XML['('].MarshalXML(c.encoder,c.start)
			
			c.CompileExpression()

			if !c.HasMoreToken(){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			
			if !(c.TokenStream[0].Type==T_SYMBOL && c.TokenStream[0].Symbol==')'){
				log.Fatalf("Syntax error on compileTermStatement")
			}
			c.Advance()
			Symbol2XML[')'].MarshalXML(c.encoder,c.start)
		case '-':
			fallthrough
		case '~':
			Symbol2XML[c.CurrentToken.Symbol].MarshalXML(c.encoder,c.start)
			c.CompileTerm()
		default:
			log.Fatalf("Syntax error on compileTermStatement")
		}
	default:
		log.Fatalf("Syntax error on compileTermStatement")
	}
	c.start.Name.Local = "term"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
}

// [expression [, expression]*]?
func(c *CompilationEngine)CompileExpressionList(){
	c.start.Name.Local = "expressionList"
	c.encoder.EncodeToken(c.start)

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileExpressionList")
	}
	if c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol==')'{
		c.start.Name.Local = "expressionList"
		c.encoder.EncodeToken(xml.EndElement{c.start.Name})
		return
	}

	c.CompileExpression()

	if !c.HasMoreToken(){
		log.Fatalf("Syntax error on CompileExpressionList")
	}
	if !(c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol==')'){
		for !(c.TokenStream[0].Type == T_SYMBOL && c.TokenStream[0].Symbol==')'){
			c.Advance()
			if !(c.CurrentToken.Type == T_SYMBOL && c.CurrentToken.Symbol==','){
				log.Fatalf("Syntax error on CompileExpressionList")
			}
			Symbol2XML[','].MarshalXML(c.encoder,c.start)

			c.CompileExpression()
		}
	}

	c.start.Name.Local = "expressionList"
	c.encoder.EncodeToken(xml.EndElement{c.start.Name})
	return
}