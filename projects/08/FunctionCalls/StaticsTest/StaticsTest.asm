@256  // SP=256
D=A
@R0
M=D
@Sys.return-address // 返回地址压栈
A=M
D=M
@SP
A=M
M=D
@SP
M=M+1
@R1 // LCL压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R2 // ARG压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R3 // THIS压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R4 // THAT压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP // ARG=SP-5
D=M
@5
D=D-A
@R2
M=D
@SP // 重置LCL=SP
D=M
@R1
M=D
@Sys.init // call Sys.init
0;JMP
(Sys.return-address)
(END)
	@END
	0;JMP
(Class1.set)
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
A=M-1
D=M
@SP
M=M-1 // 栈中数据出栈
@FunctionCalls.StaticsTest.StaticsTest.asm.FunctionCalls.StaticsTest.Class1.vm0
M=D
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
A=M-1
D=M
@SP
M=M-1 // 栈中数据出栈
@FunctionCalls.StaticsTest.StaticsTest.asm.FunctionCalls.StaticsTest.Class1.vm1
M=D
@0
D=A
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
A=D-A
D=M
@RET0
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
@RET0
A=M
0;JMP
(Class1.get)
@FunctionCalls.StaticsTest.StaticsTest.asm.FunctionCalls.StaticsTest.Class1.vm0
D=M
@SP
A=M
M=D
@SP
M=M+1 // 压入栈中
@FunctionCalls.StaticsTest.StaticsTest.asm.FunctionCalls.StaticsTest.Class1.vm1
D=M
@SP
A=M
M=D
@SP
M=M+1 // 压入栈中
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
A=D-A
D=M
@RET1
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
@RET1
A=M
0;JMP
(Class2.set)
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
A=M-1
D=M
@SP
M=M-1 // 栈中数据出栈
@FunctionCalls.StaticsTest.StaticsTest.asm.FunctionCalls.StaticsTest.Class2.vm0
M=D
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
A=M-1
D=M
@SP
M=M-1 // 栈中数据出栈
@FunctionCalls.StaticsTest.StaticsTest.asm.FunctionCalls.StaticsTest.Class2.vm1
M=D
@0
D=A
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
A=D-A
D=M
@RET2
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
@RET2
A=M
0;JMP
(Class2.get)
@FunctionCalls.StaticsTest.StaticsTest.asm.FunctionCalls.StaticsTest.Class2.vm0
D=M
@SP
A=M
M=D
@SP
M=M+1 // 压入栈中
@FunctionCalls.StaticsTest.StaticsTest.asm.FunctionCalls.StaticsTest.Class2.vm1
D=M
@SP
A=M
M=D
@SP
M=M+1 // 压入栈中
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
A=D-A
D=M
@RET3
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
@RET3
A=M
0;JMP
(Sys.init)
@6
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
@Class1.set.return-address0 // 返回地址压栈
D=A
@SP
A=M
M=D
@SP
M=M+1
@R1 // LCL压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R2 // ARG压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R3 // THIS压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R4 // THAT压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP // ARG=SP-n-5
D=M
@5
D=D-A
@2
D=D-A
@R2
M=D
@SP // 重置LCL=SP
D=M
@R1
M=D
@Class1.set // 跳转控制
0;JMP
(Class1.set.return-address0) // 返回地址
@0
D=A
@5
D=D+A
@SP
A=M
M=D // 使用栈顶存放位置信息
@SP
A=M-1
D=M
@SP
A=M
A=M
M=D
@SP
M=M-1
@23
D=A
@SP
A=M
M=D
@SP
M=M+1
@15
D=A
@SP
A=M
M=D
@SP
M=M+1
@Class2.set.return-address1 // 返回地址压栈
D=A
@SP
A=M
M=D
@SP
M=M+1
@R1 // LCL压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R2 // ARG压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R3 // THIS压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R4 // THAT压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP // ARG=SP-n-5
D=M
@5
D=D-A
@2
D=D-A
@R2
M=D
@SP // 重置LCL=SP
D=M
@R1
M=D
@Class2.set // 跳转控制
0;JMP
(Class2.set.return-address1) // 返回地址
@0
D=A
@5
D=D+A
@SP
A=M
M=D // 使用栈顶存放位置信息
@SP
A=M-1
D=M
@SP
A=M
A=M
M=D
@SP
M=M-1
@Class1.get.return-address2 // 返回地址压栈
D=A
@SP
A=M
M=D
@SP
M=M+1
@R1 // LCL压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R2 // ARG压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R3 // THIS压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R4 // THAT压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP // ARG=SP-n-5
D=M
@5
D=D-A
@0
D=D-A
@R2
M=D
@SP // 重置LCL=SP
D=M
@R1
M=D
@Class1.get // 跳转控制
0;JMP
(Class1.get.return-address2) // 返回地址
@Class2.get.return-address3 // 返回地址压栈
D=A
@SP
A=M
M=D
@SP
M=M+1
@R1 // LCL压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R2 // ARG压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R3 // THIS压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@R4 // THAT压栈
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP // ARG=SP-n-5
D=M
@5
D=D-A
@0
D=D-A
@R2
M=D
@SP // 重置LCL=SP
D=M
@R1
M=D
@Class2.get // 跳转控制
0;JMP
(Class2.get.return-address3) // 返回地址
(WHILE)
@WHILE
0;JMP
