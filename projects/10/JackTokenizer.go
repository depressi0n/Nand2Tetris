package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)



type JackTokenizer struct{
	FileName string
	Scanner *bufio.Reader
	CurrentLine string
	CurrentToken *Token
	TokenStream []*Token
}


func NewJackTokenizer()*JackTokenizer{
	return &JackTokenizer{
		// FileName:s,
		// Scanner: bufio.NewReader(scanner),
		TokenStream: nil,
		CurrentLine: "",
	}
}
// readLine process the original character stream :
// - Remove the comment lines and whitespace lines
// - Remove the comment after statements
func (t *JackTokenizer)readLine()bool{
	line, _, err := t.Scanner.ReadLine()
	if err!=nil{
			if  err==io.EOF {
				t.CurrentLine=""
				return false
			}
			log.Fatalf("Error in reading the file named %s:%v",t.FileName,err)
		}
	t.CurrentLine=string(line)
	t.CurrentLine=strings.TrimSpace(t.CurrentLine)
	// for whitespace line, comment line begin with "*", "/*" ,"*/","//","/**"
	// Now just support the single comment line
	for t.CurrentLine == "" || t.CurrentLine[0:1]=="*" || (len(t.CurrentLine)>=2 && (t.CurrentLine[0:2] == "//" || t.CurrentLine[0:2] == "*/")) || ( (len(t.CurrentLine)>=3 && t.CurrentLine[0:3]=="/**")){
		// read the next line
		line, _, err := t.Scanner.ReadLine()
		if err!=nil{
			if  err==io.EOF {
				t.CurrentLine=""
				return false
			}
			log.Fatalf("Error in reading the file named %s:%v",t.FileName,err)
		}
		t.CurrentLine=string(line)
		t.CurrentLine=strings.TrimSpace(t.CurrentLine)
	}
	t.CurrentLine=strings.TrimSpace(t.CurrentLine)
	// Remove the comment after statements
	commentStart:=strings.Index(t.CurrentLine,"//")
	if commentStart!=-1 {
		t.CurrentLine=t.CurrentLine[:commentStart]
	}
	return true
}

func(t *JackTokenizer)Tokenize(){
	scanner,err:=os.OpenFile(t.FileName,os.O_RDONLY,0600)
	if err!=nil {
		log.Fatalf("Error in opening the file named %s with os.OpenFile:%v",t.FileName,err)
	}
	defer scanner.Close()
	t.Scanner=bufio.NewReader(scanner)
	for t.readLine() {
		// handle the t.CurrentLine
		s:=&strings.Builder{}
		for i:=0;i<len(t.CurrentLine);i++{
			// add whitespace to split the symbol with other token
			if _,ok:=Rune2Symbol[rune(t.CurrentLine[i])];ok{
				s.WriteRune(' ')
				s.WriteRune(rune(t.CurrentLine[i]))
				s.WriteRune(' ')

			}else if t.CurrentLine[i] == '"'{
				// add '\' to difference the '"' for tokenize
				s.WriteByte('\\')
				s.WriteByte('"')
				s.WriteByte(' ')
			}else{
				s.WriteByte(t.CurrentLine[i])
			}
		}
		
		// split s by whitespace
		tokens:=strings.Split(s.String()," ")
		// handle all token
		for i:=0;i<len(tokens);i++{
			token:=tokens[i]
			if token == ""{
				continue
			}
			// TODO: keep the constant origin
			if strings.HasPrefix(token,`\"`) {
				i++
				token+=" "
				for i<len(tokens) && strings.Index(tokens[i],`\"`)==-1{
					token+=tokens[i]+" "
					i++
				}
				token+=tokens[i]+" "
			}
			t.TokenStream=append(t.TokenStream,NewToken(token))
		}
	}
}
func(t *JackTokenizer)HasMoreToken() bool {
	if len(t.TokenStream)>0 {
		return true
	}
	return false
}

// Advance 从输入中读取下一个Token，使其成为当前Token
func(t *JackTokenizer)Advance(){
	if t.HasMoreToken(){
		t.CurrentToken=t.TokenStream[0]
		t.TokenStream=t.TokenStream[1:]
	}
}

func(t *JackTokenizer)TokenType() TokenType {
	return t.CurrentToken.Type
}

func(t *JackTokenizer)Keyword()KeywordType{
	if t.CurrentToken.Type != T_KEYWORD {
		log.Fatalf("Current Token [%v] is not a keyword",t.CurrentToken)
	}
	return t.CurrentToken.Keyword
}

func(t *JackTokenizer)Symbol() rune {
	if t.CurrentToken.Type != T_SYMBOL {
		log.Fatalf("Current Token [%v] is not a keyword",t.CurrentToken)
	}
	return t.CurrentToken.Symbol
}

func(t *JackTokenizer)Identifier()string{
	if t.CurrentToken.Type != T_IDENTIFIER {
		log.Fatalf("Current Token [%v] is not a keyword",t.CurrentToken)
	}
	return t.CurrentToken.Identifier
}

func(t *JackTokenizer)IntVal()int{
	if t.CurrentToken.Type != T_INT_CONST {
		log.Fatalf("Current Token [%v] is not a keyword",t.CurrentToken)
	}
	return t.CurrentToken.IntVal
}

func(t *JackTokenizer)StringVal()string{
	if t.CurrentToken.Type != T_STRING_CONST {
		log.Fatalf("Current Token [%v] is not a keyword",t.CurrentToken)
	}
	return t.CurrentToken.StringVal
}


type Token struct{
	Type TokenType
	Keyword KeywordType
	Symbol rune
	Identifier string
	IntVal int
	StringVal string
}
func (t *Token)String()string{
	switch t.Type {
	case T_SYMBOL:
		return string(t.Symbol)
	case T_KEYWORD:
		return KeywordType2Str[t.Keyword]
	case T_IDENTIFIER:
		return t.Identifier
	case T_INT_CONST:
		return strconv.Itoa(t.IntVal)
	case T_STRING_CONST:
		return t.StringVal
	default:
		log.Fatalf("error for Token.String()")
		return ""
	}
	
}
func NewToken(s string)*Token{
	res:=&Token{}
	// consider the single character as symbol
	if len(s) == 1 {
		r, _ := utf8.DecodeRuneInString(s)
		if _,ok:= Rune2Symbol[r];ok{
			res.Type=T_SYMBOL
			res.Symbol=r
			return res
		}
	}
	// check if this token is an keyword
	if v,ok:= Str2KeywordType[s];ok{
		res.Type=T_KEYWORD
		res.Keyword=v
		return res
	}
	// view as constant or identifier
	
	if strings.Index(s,"\\\"")!=-1 && strings.LastIndex(s,"\\\"")!=strings.Index(s,"\\\""){
		res.Type=T_STRING_CONST
		res.StringVal=s[strings.Index(s,"\\\"")+2:strings.LastIndex(s,"\\\"")]
		return res
	} else if b, _ := regexp.MatchString("^[0-9]+",s);b{
		res.Type=T_INT_CONST
		v,_:=strconv.Atoi(s)
		res.IntVal=v
		return res
	}else{
		res.Type=T_IDENTIFIER
		res.Identifier=s
		return res
	}
}
var Str2KeywordType = map[string]KeywordType{
	"class":K_CLASS,
	"method":K_METHOD,
	"int":K_INT,
	"function":K_FUNCTION,
	"boolean":K_BOOLEAN,
	"constructor":K_CONSTRUCTOR,
	"char":K_CHAR,
	"void":K_VOID,
	"var":K_VAR,
	"static":K_STATIC,
	"field":K_FIELD,
	"let":K_LET,
	"do":K_DO,
	"if":K_IF,
	"else":K_ELSE,
	"while":K_WHILE,
	"return":K_RETURN,
	"true":K_TRUE,
	"false":K_FALSE,
	"null":K_NULL,
	"this":K_THIS,
}
var KeywordType2Str = map[KeywordType]string{
	K_CLASS:"class",
	K_METHOD:"method",
	K_INT:"int",
	K_FUNCTION:"function",
	K_BOOLEAN:"boolean",
	K_CONSTRUCTOR:"constructor",
	K_CHAR:"char",
	K_VOID:"void",
	K_VAR:"var",
	K_STATIC:"static",
	K_FIELD:"field",
	K_LET:"let",
	K_DO:"do",
	K_IF:"if",
	K_ELSE:"else",
	K_WHILE:"while",
	K_RETURN:"return",
	K_TRUE:"true",
	K_FALSE:"false",
	K_NULL:"null",
	K_THIS:"this",
}

// "":T_IDENTIFIER, 一系列字母、数字、下划线组成，不能以数字开头
// "":T_INT_CONST, 0-32767
// "":T_STRING_CONST, 使用双引号包含的一系列的Unicode字符
var Rune2Symbol = map[rune]struct{} {
	rune('{'):{},
	rune('}'):{},
	rune('('):{},
	rune(')'):{},
	rune('['):{},
	rune(']'):{},
	rune('.'):{},
	rune(','):{},
	rune(';'):{},
	rune('+'):{},
	rune('-'):{},
	rune('*'):{},
	rune('/'):{},
	rune('&'):{},
	rune('|'):{},
	rune('<'):{},
	rune('>'):{},
	rune('='):{},
	rune('~'):{},
}

type TokenType int

const (
	T_KEYWORD TokenType= iota
	T_SYMBOL
	T_IDENTIFIER
	T_INT_CONST
	T_STRING_CONST
)

type KeywordType int

const (
	K_CLASS KeywordType=iota
	K_METHOD
	K_INT
	K_FUNCTION
	K_BOOLEAN
	K_CONSTRUCTOR
	K_CHAR
	K_VOID
	K_VAR
	K_STATIC
	K_FIELD
	K_LET
	K_DO
	K_IF
	K_ELSE
	K_WHILE
	K_RETURN
	K_TRUE
	K_FALSE
	K_NULL
	K_THIS
)
