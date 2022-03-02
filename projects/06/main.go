package main

import (
	"fmt"
	"os"
	"strings"
)

const USAGEHELP =`
Usage:
	asm [input path] [output path]
`

// 无符号汇编编译器的实现过程：
// 1. 打开输出文件
// 2. 处理输入文件的每一行（汇编指令）：
// 2.1 对于C-指令，程序翻译后的指令域的二进制连接到一个单一的16位字上，写到输出文件中
// 2.2 对于A-指令，语法分析器返回的十进制常数翻译成对应的二进制表示，然后得到16位字，写到输出文件中

// 有符号汇编编译器的实现过程：允许在符号被定义之前使用符号标签，两次遍历输入文件，第一次读取时构建符号表并分配内存地址，第二次读取遇到的所有标签的内存地址
// 1. 逐行处理整个程序，构建符号表，并利用数字来记录ROM地址——当前命令最终被加载到这个地址中，从0开始，遇到A-指令或C指令则增1。每次遇到伪指令或注释时不发生变化，但遇到伪指令时，符号表上增加一个新条目
// 2. 每一行进行语法分析，遇到A-指令时如果内容是符号中，则在符号表中查找相应符号，成功则使用对应的地址代替符号，失败则表示变量，将这个变量和RAM地址加入到符号表中，分配的RAM地址是连续数字，从地址16开始
func main(){
	if len(os.Args) < 2{
		fmt.Println("Invalid parameters: length less than 3")
		fmt.Println(USAGEHELP)
		os.Exit(1)
	}
	inputASMPath:=os.Args[1]
	outputHackPath:=strings.TrimSuffix(inputASMPath,".asm")
	asm:=NewAssemble(inputASMPath,outputHackPath)
	asm.Run()
}