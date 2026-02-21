#include<cstdint>
//provides fixed size integer types, such as uint8_t, uint16_t, etc.

class Chip8 {
    public:
        Chip8();
    private:
        uint8_t memory[4096]; //4K memory
        uint8_t disp[64 * 32]; //display buffer
        uint16_t pc; //program counter
        uint16_t I; //index register
        uint16_t stack[16]; //stack for 16 bit addresses
        uint8_t delay; //delay timer
        uint8_t sound; //sound timer
        uint8_t V[16]; //16 8-bit general purpose registers

};