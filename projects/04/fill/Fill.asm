// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.
// 初始值为白屏
@FLAG
M=0

@SCREEN
D=A
@arr // 使用一个变量记录屏幕基地址
M=D 

// 256行，每行32个寄存器则n=256*32=8192个寄存器
@8192
D=A
@n
M=D // n=8192

@i
M=0

// 通过KBD寄存器检查键盘是否有输入（无限循环）
(LOOP)
    @KBD
    D=M
    @BLACK
    D;JNE // 如果不等于0，则跳转到@BLACK
    @WHITE
    0;JMP // 如果等于0，则跳转到@WHITE

    @LOOP
    0;JMP

(BLACK)
    @FLAG
    D=M-1
    @LOOP
    D;JEQ // 如果当前状态为黑，则无需设置，跳回LOOP继续监控
    @FLAG
    M=1 // 设置当前状态为白
    @i
    M=0 // i = 0
    (SETBLACKLOOP)
        @i
        D=M
        @n
        D=D-M // 判断当前循环次数
        @LOOP
        D;JEQ // 完成

        @arr
        D=M // 取出基地址
        @i
        A=D+M
        M=-1 // 设置为黑色

        @i
        M=M+1 // 计数器递增

        @SETBLACKLOOP
        0;JMP // 跳转

(WHITE)
@FLAG
    D=M
    @LOOP
    D;JEQ // 如果当前状态为白，则无需设置，跳回LOOP继续监控
    @FLAG
    M=0 // 否则设置为白
    @i
    M=0 
    (SETWHITELOOP)
        @i
        D=M
        @n
        D=D-M
        @LOOP
        D;JEQ // 完成

        @arr
        D=M
        @i
        A=D+M
        M=0 // 设置为白色

        @i
        M=M+1 // 计数器递增

        @SETWHITELOOP
        0;JMP