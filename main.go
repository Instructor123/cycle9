package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
)

func checkErr(e error){
	if e != nil{
		panic(e)
	}
}

func findRet(myBytes []byte) int{

	returnOpcodes := make([]byte,4)
	returnOpcodes[0] = byte('\xC3')
	returnOpcodes[1] = byte('\xCB')
	returnOpcodes[2] = byte('\xC2')
	returnOpcodes[3] = byte('\xCA')
	stop := false

	i := len(myBytes)-1
	j := 0
	for !stop{
		if myBytes[i] == returnOpcodes[j]{
			stop = true
		} else {
			i -= 1
			j += 1
		}
		if i <= 0 || j > 3{
			stop = true
		}
	}

	return i
}

func printMenu(myBytes []byte, endLocation int){

	keepGoing := true
	userInput := 0
	for keepGoing{
		fmt.Println("Please choose an option below: ")
		fmt.Println("\t1.)Print raw bytes")
		fmt.Println("\t2.)Print bytes with leading '\\x'.")
		fmt.Println("\t3.)Print Python2 code.")
		fmt.Println("\t4.)Print Python3 code.")
		fmt.Println("\t5.)Print C code.")
		fmt.Printf("Please enter your choice: ")
		fmt.Scan(&userInput)

		if userInput >= 1 && userInput <= 5{
			keepGoing = false
		}
	}

	switch userInput {
	case 1:
		printRaw(myBytes, endLocation)
		break
	case 2:
		printFormat(myBytes, endLocation)
		break
	case 3:
		printPython(myBytes, endLocation, 2)
		break
	case 4:
		printPython(myBytes, endLocation, 3)
		break
	case 5:
		printC(myBytes, endLocation)
	default:
		fmt.Println("ERROR: Shouldn't get here.")
		break
	}
}

func printRaw(myBytes []byte, endLocation int){
	for i := range myBytes[:endLocation+1]{
		fmt.Printf("%02x", myBytes[i] )
	}
	println("\n")
}

func printFormat(myBytes []byte,   endLocation int){
	for i := range myBytes[:endLocation+1]{
		fmt.Printf("\\x%02x", myBytes[i] )
	}
	println("\n")
}

func printPython(myBytes []byte,  endLocation int, version int){
	
}

func printC(myBytes []byte, endLocation int){

}

func retrieveInfo(sym []elf.Symbol, function *string)(uint64,uint64){
	functionStart := uint64(0)
	functionSize := uint64(0)
	for x := range sym{
		if sym[x].Name == *function{
			functionStart = sym[x].Value
			functionSize = sym[x].Size
			break
		}
	}
	return functionStart, functionSize
}

func retrieveBytes(file *elf.File, start, size uint64)[]byte{
	textSection := file.Section(".text")
	funcOffsetStart := start - textSection.Offset
	funcOffsetEnd := funcOffsetStart + size
	output, err := textSection.Data()
	checkErr(err)
	funcBytes := output[funcOffsetStart:funcOffsetEnd]

	return funcBytes
}

func main(){
	//validate file exists
	flagFilename := flag.String("f", "", "file name: Provide the full path to the file containing your shellcode.")
	flagFunctionName := flag.String("t", "", "function name: Provide the function name that contains your assembly.")
	flag.Parse()

	retValue := 1

	if "" == *flagFilename || "" == *flagFunctionName {
		flag.PrintDefaults()
		retValue = -1
	}

	if 1 == retValue{
		fileStat, err := os.Stat(*flagFilename)
		checkErr(err)

		if nil != fileStat{
			shellFile, err := elf.Open(*flagFilename)
			checkErr(err)

			mySymb, err := shellFile.Symbols()
			checkErr(err)

			functionStart, functionSize := retrieveInfo(mySymb, flagFunctionName)

			funcBytes := retrieveBytes(shellFile, functionStart, functionSize)

			retLocation := findRet(funcBytes)

			if retLocation > 0{
				printMenu(funcBytes, retLocation)
			}
		}
	}
}
