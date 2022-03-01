// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

// Put your code here.
// R0*R1 等价于 R0 + R0 + ... + R0
// 初始化R2作为和，初始值为0
@R2
M=0
// 计数器 i = 0
@i
M=0

// 循环R1次
(LOOP)
// 条件判断
@i
D=M
@R1
D=D-M // i-R1=0 则JMP
D;JEQ

// 循环执行加法
@R0
D=M
@R2
M=D+M

// 计数器
@i
M=M+1

@LOOP
0;JMP // 无条件跳转

(END)
@END
0;JMP // 无条件跳转
