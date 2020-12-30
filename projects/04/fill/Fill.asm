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
(LOOP)
  @KBD
  D=M

  @BLACK
  D;JGT
  @WHITE
  0;JMP

(BLACK)
  @color
  M=-1
  @EXEC
  0;JMP
(WHITE)
  @color
  M=0
  @EXEC
  0;JMP

(EXEC)
  @i
  M=0

  (NEXT)
    // while(i < 8192)
    @i
    D=M
    // 512 / 16 * 256
    @8192
    D=D-A
    @END
    D;JGE

    @SCREEN
    D=A
    @i
    D=D+M
    @now
    M=D

    @color
    D=M
    @now
    A=M
    M=D

    // i++
    @i
    M=M+1
    @NEXT
    0;JMP

(END)
  @LOOP
  0;JMP
