package main
//some caveats to get the SGN package installed
/*
	https://github.com/keystone-engine/keystone/blob/master/docs/COMPILE-NIX.md
	https://github.com/keystone-engine/keystone/blob/master/docs/COMPILE-NIX.md
 */
import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func checkErr(e error){
	if e != nil{
		panic(e)
	}
}

const PYTHON3 = 3
const PYTHON2 = 2

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
		fmt.Println("\t1.)Print raw opcodes")
		fmt.Println("\t2.)Print opcodes with leading '\\x'.")
		fmt.Println("\t3.)Print Python2 code.")
		fmt.Println("\t4.)Print Python3 code.")
		fmt.Println("\t5.)Print C code.")
		fmt.Println("\t6.)Encode bytes.")
		fmt.Printf("Please enter your choice: ")
		fmt.Scan(&userInput)

		if userInput >= 1 && userInput <= 6{
			keepGoing = false
		}
	}

	switch userInput {
	case 1:
		printRaw(myBytes, endLocation)
		println("\n")
		break
	case 2:
		printFormat(myBytes, endLocation)
		println("\n")
		break
	case 3:
		printPython(myBytes, endLocation, 2)
		println("\n")
		break
	case 4:
		printPython(myBytes, endLocation, 3)
		println("\n")
		break
	case 5:
		printC(myBytes, endLocation)
		println("\n")
		break
	case 6:
		err := os.WriteFile("TEMP_plain.dat", myBytes, 0667)
		checkErr(err)
		cmd := exec.Command("/home/user/go/bin/sgn", "-a", "64", "-o", "test", "TEMP_plain.dat")
		err = cmd.Run()
		checkErr(err)
		cmd = exec.Command("rm", "TEMP_plain.dat")
		err = cmd.Run()
		checkErr(err)
	default:
		fmt.Println("ERROR: Shouldn't get here.")
		break
	}
}

func printRaw(myBytes []byte, endLocation int){
	for i := range myBytes[:endLocation+1]{
		fmt.Printf("%02x", myBytes[i] )
	}
}

func printFormat(myBytes []byte,   endLocation int){
	for i := range myBytes[:endLocation+1]{
		fmt.Printf("\\x%02x", myBytes[i] )
	}
}

func printPython(myBytes []byte,  endLocation int, version int){

	switch version{
	case PYTHON3:
		print("opCodes='")
		printFormat(myBytes, endLocation)
		print("'\n")
		print("print(opCodes, end=\"\")")
		break
	case PYTHON2:
		print("import sys\n\n")
		print("opCodes = '")
		printFormat(myBytes, endLocation)
		print("'\n")
		print("sys.stdout.write(opCodes)")
		break
	default:
		fmt.Println("Python version not supported")
		break
	}
}

func printC(myBytes []byte, endLocation int){
	print("unsigned char opCodes[", len(myBytes), "] = {")
	for i := range myBytes{
		fmt.Printf("'\\x%02x'", myBytes[i])
		if i != len(myBytes)-1{
			fmt.Printf(",")
		}
	}
	print("};\n\n")
	print("for(int i = 0; i < ",len(myBytes),"; ++i){\n")
	print("\tprintf(\"%02x\", opCodes[i]);\n")
	print("}")

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
