@i // i = 1
M=1
@sum // sum = 0
M=0
(LOOP)
@i // if i-100 = 0 goto END
D=M
@100
D=D-A
@END
D;JGT
@i
D=M
@sum
M=D+M
@i
M=M+1
@LOOP
0;JMP
(END)
@END
0;JMP
