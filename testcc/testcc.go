package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type TestChainCode struct {
}

func (t *TestChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("GPTestChainCode: Init")
	return shim.Success(nil)
}

func (t *TestChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("GPTestChainCode: Invoke: %s", function)

	if function == "set" {
		return t.set(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}

	return shim.Error(fmt.Sprintf("Invalid invoke function name: %s", function))
}

func (t *TestChainCode) set(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key string
	var val int
	var err error

	if len(args) != 2 {
		return shim.Error("Wrong arguments")
	}
	fmt.Println("SET: %s, %s", args[0], args[1])

	// Get Key
	key = args[0]

	// Get Value
	val, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Failed to get value: " + err.Error())
	}

	// Save Key-Value
	err = stub.PutState(key, []byte(strconv.Itoa(val)))
	if err != nil {
		return shim.Error("Failed to put value to ledger: " + err.Error())
	}

	return shim.Success(nil)
}

func (t *TestChainCode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (t *TestChainCode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key string
	var err error

	if len(args) != 1 {
		return shim.Error("Wrong arguments")
	}
	fmt.Println("QUERY: %s", args[0])

	key = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(key)
	if err != nil {
		return shim.Error("Failed to get state from ledger")
	}

	resp := "" + string(Avalbytes)
	fmt.Printf("QUERY Response:%s\n", resp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(TestChainCode))
	if err != nil {
		fmt.Printf("Error starting test chaincode: %s", err)
	}
}
