package main

import (
	"encoding/xml"
	"os"
)

type Test struct{
	XMLName     xml.Name `xml:"servers"`
    Version     string   `xml:"version,attr"`
    Svs         []server `xml:"server"`
    Description string   `xml:",innerxml"`
}
type server struct {
    XMLName    xml.Name `xml:"server"`
    ServerName string   `xml:"serverName"`
    ServerIP   string   `xml:"serverIP"`
	InServers []inServer `xml:"inServer"`
}
type inServer struct {
    XMLName    xml.Name `xml:"inServer"`
    Name string   `xml:"inServerName"`
    IP   string   `xml:"inServerIP"`
}
func main(){
	analyzer:=NewJackAnalyzer(os.Args[1])
	analyzer.Analyze()
}