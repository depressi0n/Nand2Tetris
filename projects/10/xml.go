package main

import (
	"encoding/xml"
	"strconv"
)



type XMLSymbol struct {
	Name string `xml:",innerxml"`
}

func (k *XMLSymbol) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "symbol"
	e.EncodeToken(start)

	e.EncodeToken(xml.CharData(" "+k.Name+" "))

	start.Name.Local = "symbol"
	e.EncodeToken(xml.EndElement{start.Name})

	return nil
}

type XMLKeyword struct {
	Name string `xml:",innerxml"`
}

func (k *XMLKeyword) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "keyword"
	e.EncodeToken(start)
	e.EncodeToken(xml.CharData(" "+k.Name+" "))
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

type XMLIdentifier struct {
	Name string `xml:",innerxml"`
}

func (id *XMLIdentifier) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "identifier"
	e.EncodeToken(start)
	e.EncodeToken(xml.CharData(" "+id.Name+" "))
	//e.EncodeElement(Keyword2XML[K_CLASS], start)
	//e.EncodeToken(c.Name)
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

type XMLStringConstant struct {
	Name string `xml:",innerxml"`
}

func (k *XMLStringConstant) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start = xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "stringConstant",
		},
		Attr: nil,
	}
	e.EncodeToken(start)
	e.EncodeToken(xml.CharData(" "+k.Name+" "))
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

type XMLIntegerConstant struct {
	Value int `xml:",innerxml"`
}

func (k *XMLIntegerConstant) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start = xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "integerConstant",
		},
		Attr: nil,
	}
	e.EncodeToken(start)
	s := strconv.Itoa(k.Value)
	e.EncodeToken(xml.CharData(" "+s+" "))
	e.EncodeToken(xml.EndElement{start.Name})
	return nil
}

var Keyword2XML = map[KeywordType]*XMLKeyword{
	K_CLASS:       {Name: "class"},
	K_METHOD:      {Name: "method"},
	K_INT:         {Name: "int"},
	K_FUNCTION:    {Name: "function"},
	K_BOOLEAN:     {Name: "boolean"},
	K_CONSTRUCTOR: {Name: "constructor"},
	K_CHAR:        {Name: "char"},
	K_VOID:        {Name: "void"},
	K_VAR:         {Name: "var"},
	K_STATIC:      {Name: "static"},
	K_FIELD:       {Name: "field"},
	K_LET:         {Name: "let"},
	K_DO:          {Name: "do"},
	K_IF:          {Name: "if"},
	K_ELSE:        {Name: "else"},
	K_WHILE:       {Name: "while"},
	K_RETURN:      {Name: "return"},
	K_TRUE:        {Name: "true"},
	K_FALSE:       {Name: "false"},
	K_NULL:        {Name: "null"},
	K_THIS:        {Name: "this"},
}

var Symbol2XML = map[rune]*XMLSymbol{
	rune('{'): {Name: "{"},
	rune('}'): {Name: "}"},
	rune('('): {Name: "("},
	rune(')'): {Name: ")"},
	rune('['): {Name: "["},
	rune(']'): {Name: "]"},
	rune('.'): {Name: "."},
	rune(','): {Name: ","},
	rune(';'): {Name: ";"},
	rune('+'): {Name: "+"},
	rune('-'): {Name: "-"},
	rune('*'): {Name: "*"},
	rune('/'): {Name: "/"},
	rune('&'): {Name: "&"},
	rune('|'): {Name: "|"},
	rune('<'): {Name: "<"},
	rune('>'): {Name: ">"},
	rune('='): {Name: "="},
	rune('~'): {Name: "~"},
}



// type XMLClass struct {
// 	XMLName xml.Name `xml:"class"`
// 	Keyword string `xml:"keyword"`
// 	Name XMLIdentifier  `xml:"identifier"`
// 	LeftBraces XMLSymbol  //`xml:">left,omitempty"`
// 	Variables []XMLClassVarDec `xml:",omitempty"`
// 	Subroutines []XMLSubroutineDec  `xml:",omitempty"`
// 	RightBraces XMLSymbol  //`xml:">right,omitempty"`
// }
// func(c *XMLClass)MarshalXML(e *xml.Encoder,start xml.StartElement)error{

// }



// func NewXMLClass(name string,variables []XMLClassVarDec,subroutines []XMLSubroutineDec)*XMLClass{
// 	return &XMLClass{
// 		Keyword: KeywordType2Str[K_CLASS],
// 		Name:*NewXMLIdentifier(name),
// 		LeftBraces:*Symbol2XML[rune('{')],
// 		Variables:variables,
// 		Subroutines:subroutines,
// 		RightBraces:Symbol2XML[rune('}')],
// 	}
// }

// type XMLClassVarDec struct{
// 	XMLName xml.Name `xml:"classVarDec"`
// 	Static *XMLKeyword 
// 	Type *XMLKeyword
// 	Variables []*XMLIdentifier 
// 	Simicolon *XMLSymbol
// }

// func NewXMLClassVarDec(staticFlag bool,fieldFlag bool,varType KeywordType,varNames []string)*XMLClassVarDec{
// 	res:=&XMLClassVarDec{
// 		Type:Keyword2XML[varType],
// 		Simicolon:Symbol2XML[rune(';')],
// 	}
// 	if staticFlag && fieldFlag {
// 		log.Fatalf("Error for use 'static' and 'feild' meanwhile")
// 	}
// 	if staticFlag {
// 		res.Static=Keyword2XML[K_STATIC]
// 	}
// 	if fieldFlag {
// 		res.Static=Keyword2XML[K_FIELD]
// 	}
// 	res.Variables=make([]*XMLIdentifier, 0,len(varNames))
// 	for i := 0; i < len(varNames); i++ {
// 		res.Variables=append(res.Variables, NewXMLIdentifier(varNames[i]))
// 	}
// 	return res
// }

// type XMLSubroutineDec struct{
// 	XMLName xml.Name `xml:"subroutineDec"`

// }

// func NewXMLSubroutineDec()*XMLSubroutineDec{
// 	return nil
// }

// type XMLSymbol  struct{
// 	XMLName xml.Name `xml:"symbol"`
// 	Name string `xml:",innerxml"`
// }

// var Symbol2XML =map[rune]*XMLSymbol{
// 	rune('{'):{Name:"{"},
// 	rune('}'):{Name:"}"},
// 	rune('('):{Name:"("},
// 	rune(')'):{Name:")"},
// 	rune('['):{Name:"["},
// 	rune(']'):{Name:"]"},
// 	rune('.'):{Name:"."},
// 	rune(','):{Name:","},
// 	rune(';'):{Name:";"},
// 	rune('+'):{Name:"+"},
// 	rune('-'):{Name:"-"},
// 	rune('*'):{Name:"*"},
// 	rune('/'):{Name:"/"},
// 	rune('&'):{Name:"&"},
// 	rune('|'):{Name:"|"},
// 	rune('<'):{Name:"<"},
// 	rune('>'):{Name:">"},
// 	rune('='):{Name:"="},
// 	rune('~'):{Name:"~"},
// }

// type XMLKeyword struct{
// 	XMLName xml.Name `xml:"keyword"`
// 	Name string `xml:",innerxml"`
// }

// var Keyword2XML = map[KeywordType]*XMLKeyword{
// 	K_CLASS:{Name:"class"},
// 	K_METHOD:{Name:"method"},
// 	K_INT:{Name:"int"},
// 	K_FUNCTION:{Name:"function"},
// 	K_BOOLEAN:{Name:"boolean"},
// 	K_CONSTRUCTOR:{Name:"constructor"},
// 	K_CHAR:{Name:"char"},
// 	K_VOID:{Name:"void"},
// 	K_VAR:{Name:"var"},
// 	K_STATIC:{Name:"static"},
// 	K_FIELD:{Name:"field"},
// 	K_LET:{Name:"let"},
// 	K_DO:{Name:"do"},
// 	K_IF:{Name:"if"},
// 	K_ELSE:{Name:"else"},
// 	K_WHILE:{Name:"while"},
// 	K_RETURN:{Name:"return"},
// 	K_TRUE:{Name:"true"},
// 	K_FALSE:{Name:"false"},
// 	K_NULL:{Name:"null"},
// 	K_THIS:{Name:"this"},
// }

// type XMLIdentifier struct{
// 	XMLName xml.Name `xml:"identifier"`
// 	Name string `xml:",innerxml"`
// }
// func NewXMLIdentifier(name string) *XMLIdentifier{
// 	return &XMLIdentifier{
// 		Name:name,
// 	}
// }