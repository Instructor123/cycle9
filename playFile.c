#include <stdlib.h>
#include <stdio.h>


void normalFunction(void){
    unsigned char array[10] = {'\xc0', '\x90'};
    // for(int i =0; i < 10; ++i){
    //     array[i] = '\0';
    // }

//    array[0] = '\xc0';

    for(int i =0; i < 10; ++i){
        printf("%02x", array[i]);
    }

}

void __attribute__((naked)) assemblyFunction(void){
    asm volatile(
        "push rbp;"
        "mov rbp, rsp;"
        "pop rbp;"
        "ret;"
    );
}

int main(void){

    assemblyFunction();
    normalFunction();

    return 0;
}