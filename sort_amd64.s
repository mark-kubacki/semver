// Copyright 2014 The Semver Package Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !purego
// +build !go1.17

#include "go_asm.h"
#include "textflag.h"

TEXT ·twoFieldKey(SB),NOSPLIT,$0-32
    MOVQ    v+0(FP), SI
    MOVQ    fieldAdjustment+8(FP), M1    // Two packed L.
    MOVBQZX keyIndex+16(FP), AX

    // Contains the two relevant fields. They need to be swapped, though.
    MOVQ    (SI)(AX*4), M0
    PADDL   M1, M0

    // This function is comprised of two interleaved calculations (to hide within latencies)
    // which use register as follows.
    // M5, M6, M7: Select value of fields <=11, or NN to add to the result below.
    // M0 to M4: Number of digits (bytes) used. Will eventually be 0 for any <=11.
    //           Calculated without LZCNT or POPCNT.

    MOVQ    $0x0000000b0000000b, DX // {11, 11}
   MOVQ     $0x0000000100000010, CX // shift the rightmost byte by <<4.
    MOVQ    DX, M1
    // PSHUFW  $0xe4, M0, M4 // Just a fancy MOVQ M0, M4 // The assembler throws "invalid instruction".
    BYTE $0x0f; BYTE $0x70; BYTE $0xe0; BYTE $0xe4
   MOVQ     DX, M7
    PCMPGTL M1, M4          // Holds the "greater-than-11"-mask.
    PXOR    M1, M1
   BYTE $0x0f; BYTE $0x70; BYTE $0xe8; BYTE $0xe4 // PSHUFW // MOVQ M0, M5
    PCMPEQB M0, M1          // This saturates whole bytes, and leaves gaps inbetween.
   BYTE $0x0f; BYTE $0x70; BYTE $0xf0; BYTE $0xe4 // PSHUFW // MOVQ M0, M6
    // M0, the input, is no longer needed.
    PCMPEQB M3, M3
    PXOR    M3, M1          // ~M1
    MOVQ    $0x0101010110101010, BX // {0x01…, 0x10…}
    MOVQ    BX, M2
    PAND    M2, M1          // 0xff → 0x01 and so forth
   MOVQ     CX, M0
   PCMPGTL  M7, M5
    PAND    M1, M2          // MOVQ M1, M2; M1 and M2 are the same.
   PAND     M5, M7
   PANDN    M6, M5
    PSRLL   $8, M1          // XXX(mark): can be achieved with one less shift.
   POR      M5, M7
   PMULLW   M7, M0
    POR     M1, M2
    PSRLL   $8, M1
    POR     M1, M2
    PSRLL   $8, M1
    POR     M1, M2          // The final mask. (Due to the chosen bits PSADBW will flip nibbles.)

    // Combine the two results.
    PAND    M4, M2
   PADDW    M0, M2
    PXOR    M3, M3
    // This packs everything into an octet, swapping nibbles. Hence the *16 or <<4.
    BYTE $0x0f; BYTE $0xf6; BYTE $0xd3 // PSADBW  M3, M2

    // Converting this to a uint8 here would be inefficient.
    MOVQ    M2, ret+24(FP)
    EMMS
    RET
