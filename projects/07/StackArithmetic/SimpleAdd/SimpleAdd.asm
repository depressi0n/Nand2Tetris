@7
D=A

@SP
A=M
M=D
@SP
M=M+1
@8
D=A

@SP
A=M
M=D
@SP
M=M+1

@SP
M=M-1
A=M
D=M // 第一个操作数
@SP
M=M-1
A=M
D=D+M
@SP
A=M
M=D
@SP
M=M+1

(END)
	@END
	0;JMP