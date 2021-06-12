// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !purego
// +build !go1.17

#include "go_asm.h"
#include "textflag.h"

TEXT ·twoFieldKey(SB),NOSPLIT,$0-20
    MOVL    v+0(FP), SI
    MOVQ    fieldAdjustment+4(FP), X1    // Two packed L.
    MOVBLZX keyIndex+12(FP), AX

    // Contains the two relevant fields. They need to be swapped, though.
    MOVQ    (SI)(AX*4), X0
    PADDL   X1, X0

    // This function is comprised of two interleaved calculations (to hide within latencies)
    // which use register as follows.
    // X5, X6, X7: Select value of fields <=11, or NN to add to the result below.
    // X0 to X4: Number of digits (bytes) used. Will eventually be 0 for any <=11.
    //           Calculated without LZCNT or POPCNT.

    MOVQ    elevenPL(SB), X1
    MOVQ    X0, X4
   MOVQ     elevenPL(SB), X7
    PCMPGTL X1, X4          // Holds the "greater-than-11"-mask.
    PXOR    X1, X1
   MOVQ     X0, X5
    PCMPEQB X0, X1          // This saturates whole bytes, and leaves gaps inbetween.
   MOVQ     X0, X6
    // X0, the input, is no longer needed.
    PCMPEQB X3, X3
    PXOR    X3, X1          // ~X1
    MOVQ    movmaskAndSwapPB(SB), X2
    PAND    X2, X1          // 0xff → 0x01 and so forth
   MOVQ     shiftRightByFourPL(SB), X0
   PCMPGTL  X7, X5
    PAND    X1, X2          // MOVQ X1, X2; X1 and X2 are the same.
   PAND     X5, X7
   PANDN    X6, X5
    PSRLL   $8, X1          // XXX(mark): can be achieved with one less shift.
   POR      X5, X7
   PMULLW   X7, X0
    POR     X1, X2
    PSRLL   $8, X1
    POR     X1, X2
    PSRLL   $8, X1
    POR     X1, X2          // The final mask. (Due to the chosen bits PSADBW will flip nibbles.)

    // Combine the two results.
    PAND    X4, X2
   PADDW    X0, X2
    PXOR    X3, X3
    // This packs everything into an octet, swapping nibbles. Hence the *16 or <<4.
    PSADBW  X3, X2

    // Converting this to a uint8 here would be inefficient.
    MOVL    X2, ret+16(FP)
    RET

DATA elevenPL+0x00(SB)/8,           $0x0000000b0000000b // {11, 11}
GLOBL elevenPL(SB), (RODATA+NOPTR), $8

DATA shiftRightByFourPL+0x00(SB)/8, $0x0000000100000010 // shift the rightmost byte by <<4.
GLOBL shiftRightByFourPL(SB), (RODATA+NOPTR), $8

DATA movmaskAndSwapPB+0x00(SB)/8,   $0x0101010110101010 // {0x01…, 0x10…}
GLOBL movmaskAndSwapPB(SB), (RODATA+NOPTR), $8
