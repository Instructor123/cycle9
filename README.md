# Cycle9 - GoShell

## Description
This tool takes two arguments: a file name and a function name. The file (should) contain assembly instructions that you want to retrieve the opcodes for: the tool will work for any function, but is intended for opcode extraction. The second argument is the function name you want to extract from. You will be presented with 6 different menu options:
* Print the raw codes.
* Print the "\x" codes.
* Print Python2 code that prints the hex codes.
* Print Python3 code that prints the hex codes.
* Print C code thta prints the hex codes.
* Encode the bytes and save it to a file.

## Dependencies
This tool depends on sgn being installed in "home/user/go/bin". This can be changed if your sgn is in a different directory. Directions for installing sgn can be found here - https://github.com/EgeBalci/sgn

## Usage
The two flags are "-f" and "-t" for file and function appropriately. An example use would be:
* ./main -f a.out -t assemblyFunction