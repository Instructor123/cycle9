#include <stdlib.h>
#include <stdio.h>


void normalFunction(void){

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

    return 0;
}