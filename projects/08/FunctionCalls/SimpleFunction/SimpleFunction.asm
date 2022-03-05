(SimpleFunction.test)
@SP
A=M
M=0
@SP
M=M+1
@SP
A=M
M=0
@SP
M=M+1
@0
D=A
@R1
A=D+M
D=M // 取出local[index]的数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
@1
D=A
@R1
A=D+M
D=M // 取出local[index]的数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=D+M
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=!M
@SP
A=M
M=D
@SP
M=M+1
@0
D=A
@R2
A=D+M
D=M // 取出argument[index]的数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=D+M
@SP
A=M
M=D
@SP
M=M+1
@1
D=A
@R2
A=D+M
D=M // 取出argument[index]的数据
@SP
A=M
M=D // 压入栈中
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
D=M-D
@SP
A=M
M=D
@SP
M=M+1
@R1 // FRAME=LCL
D=M
@FRAME
M=D
@5
D=D-A
@RET
M=D
@SP // 重置调用者的返回值，从栈里面弹出一个放到ARG的位置
M=M-1
A=M
D=M
@R2
A=M
M=D
@R2 // 恢复调用者SP
D=M
@R0
M=D+1
@FRAME // 恢复调用者THAT
D=M
@1
A=D-A
D=M
@R4
M=D
@FRAME // 恢复调用者THIS
D=M
@2
A=D-A
D=M
@R3
M=D
@FRAME // 恢复调用者ARG
D=M
@3
A=D-A
D=M
@R2
M=D
@FRAME // 恢复调用者LCL
D=M
@4
A=D-A
D=M
@R1
M=D
@RET
A=M
0;JMP
(END)
	@END
	0;JMP