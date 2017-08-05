package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
    	"strings"
	//"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var loginPrefix = "LOGIN"

type SimpleChaincode struct {

}

type user struct {
	LoginID 		string 	`json:"user_id"`
	UserType		string 	`json:"user_type"`
	CompanyName 	string	`json:"company_name"`
	CompanyLocation	string	`json:"company_location"`
	BankAccountNum		int	`json:"bank_account_num"`
	BankBalance		float64	`json:"bank_balance"`
}

type userLogin struct {
	LoginName		string 	`json:"login_name"`
	Password 		string	`json:"password"`
}

type tradeRequest struct {
	TradeRequestID int
	ShipperID string
	ProducerID string
	EnergyKWH float64
	GasPrice float64
	EntryLocation string
	TradeRequestStartDate string
	TradeRequestEndDate string
	TradeRequestStatus string
	TradeRequestInvoiceID int
	TradeRequestIncidentID int
}

type UserIDList []string

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    
	//create Maps for Each Type of User
	var producerInfoMap UserIDList
    t.addUser(stub, producerInfoMap, "producer", "Producer", "Producer Company 1", "Producer Company Location", "producer", 3456, 10000.0)    	
    
	var shipperInfoMap UserIDList
    t.addUser(stub, shipperInfoMap, "shipper", "Shipper", "Shipper Company 1", "Shipper Company Location", "shipper", 1234, 10000.0)	
    
	var buyerInfoMap UserIDList
    t.addUser(stub, buyerInfoMap, "buyer", "Buyer", "Buyer Company 1", "Buyer Company Location", "buyer", 4567, 10000.0)	
    
	var transporterInfoMap UserIDList
    t.addUser(stub, transporterInfoMap, "transporter", "Transporter", "Transporter Company 1", "Transporter Company Location", "transporter", 6789, 10000.0)

	return nil, nil

}

func (t *SimpleChaincode) addUser (stub shim.ChaincodeStubInterface, userIDArr UserIDList, userName string, 
				       userType string, compName string, compLoc string, password string, 
				       bankAccountNum int, bankBalance float64 ) bool {
	
	var newUser user
	var newUserLogin userLogin

	newUser = user{LoginID: userName, UserType: userType, CompanyName: compName, 
	CompanyLocation: compLoc, BankAccountNum: bankAccountNum, BankBalance: bankBalance}
	userObjBytes, err := json.Marshal(&newUser)
	if err != nil {
		fmt.Println(err)
		return false
	}

	err1 := stub.PutState(userName, userObjBytes)
	if err1 != nil {
		fmt.Println(err1)
	}

	newUserLogin =	userLogin{LoginName: userName, Password: password} 
	userObjLoginBytes, err := json.Marshal(&newUserLogin)
	err2 := stub.PutState(loginPrefix + userName, userObjLoginBytes)
	if err2 != nil {
		fmt.Println(err2)
	}
        
    //Add the user IDs into array of user types
    var arrKey = strings.ToLower(userType) + "InfoMap"
    userIDArr = append(userIDArr, userName)
    infoMapBytes, _ := json.Marshal(userIDArr)
    _ = stub.PutState(arrKey, infoMapBytes)      
    
    	return true
}

func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userName, userType, compName, compLoc, password, arrKey string
	var bankAccountNum int
	var bankBalance float64

	var userArr UserIDList
	fmt.Println("Running function Register")

	if len (args) != 7 {
		return nil, errors.New("Incorrect number of argumets. Expecting 7")
	}
	
	userName = args[0]
	userType = args[1]
	compName = args[2]
	compLoc = args[3]
	bankAccountNum, _ = strconv.Atoi(args[4])
	bankBalance, _ = strconv.ParseFloat(args[5], 64)
	password = args[6]

	arrKey = strings.ToLower(userType) + "InfoMap"
    	userArrObj, _ := stub.GetState(arrKey)
    	_ = json.Unmarshal(userArrObj, &userArr)
    
    t.addUser (stub, userArr, userName, userType, compName, compLoc, password, bankAccountNum, bankBalance )
    	
	return nil, nil

}

func (t *SimpleChaincode) getUserInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userNameGuess, returnMessage string
	var userSample user
	fmt.Println("Getting User Credentials")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	userNameGuess = args[0]
	
	verifyBytes, err3 := t.verifyUser(stub, args)
	if err3 != nil {
		return nil, err3
	}
	if testEqualSlice(verifyBytes, []byte("Valid")) {
		userInfo, err := stub.GetState(userNameGuess)
		if err != nil {
			return nil, errors.New("User was not properly registered")
		}
		err1 := json.Unmarshal(userInfo, &userSample)
		if err1 != nil {
			return nil, err1
		}
		
		returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" :" + string(userInfo) + "}"
		return []byte(returnMessage), nil
	} else {
        returnMessage = "{\"statusCode\" : \"FAIL\", \"body\" : \"ERROR: Invalid user !\"}"
		return []byte(returnMessage), nil
	}
	return nil, nil

}

func (t *SimpleChaincode) verifyUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userNameGuess, returnMessage, passwordGuess string
	var loginObj userLogin

	fmt.Println("Verifying User")
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	userNameGuess = args[0]
	passwordGuess = args[1]

	userLoginInfo, err := stub.GetState(loginPrefix + userNameGuess)
	if userLoginInfo == nil {
		returnMessage = "Invalid Username"
		return []byte(returnMessage), nil
	}

	err1 := json.Unmarshal(userLoginInfo, &loginObj)
	if err1 != nil {
		return nil, err
	}

	if passwordGuess == loginObj.Password {
		returnMessage = "Valid"
		return []byte(returnMessage), nil
	} else {
		returnMessage = "Invalid Password"
		return []byte(returnMessage), nil
	}
	return nil, nil
}

func (t *SimpleChaincode) getProducerList(stub shim.ChaincodeStubInterface) ([]byte, error) {
	//var userSample user
	var lenMap int

	var mapProducerInfo UserIDList
	var returnMessage string
	fmt.Println("Getting Producer List")
    mapProducerInfoBytes, _ := stub.GetState("producerInfoMap")
	_ = json.Unmarshal(mapProducerInfoBytes, &mapProducerInfo)
    fmt.Println("Printing the map")
    fmt.Println(&mapProducerInfo)
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["
	lenMap = len(mapProducerInfo)
	for _, k := range mapProducerInfo {
		fmt.Println(k)
		userStructInfo, _ := stub.GetState(k)
        fmt.Println(string(userStructInfo))
		returnMessage = returnMessage + string(userStructInfo) 
		lenMap = lenMap - 1 
		if (lenMap != 0) {
			returnMessage = returnMessage + ","
		} 
	} 
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil

}

func (t *SimpleChaincode) updateUserInfo(stub shim.ChaincodeStubInterface, argsUpdated[] string, argsVerify[] string) ([]byte, error) {
	var userName, userType, compName, compLoc, password string
	var bankAccountNum int
	var bankBalance float64

	var userObj user
	var userLoginObj userLogin


	userName = argsUpdated[0]
	userType = argsUpdated[1]
	compName = argsUpdated[2]
	compLoc = argsUpdated[3]
	bankAccountNum, _ = strconv.Atoi(argsUpdated[4])
	bankBalance, _ = strconv.ParseFloat(argsUpdated[5], 64)
	password = argsUpdated[6]

	verifyBytes, _ := t.verifyUser(stub, argsVerify)

	if testEqualSlice(verifyBytes, []byte("Valid")) {
		userObj = user{LoginID: userName, UserType: userType, 
		CompanyName: compName, CompanyLocation: compLoc, BankAccountNum: bankAccountNum, 
		BankBalance: bankBalance}
		userObjBytes, err := json.Marshal(&userObj)
		if err != nil {
			fmt.Println("Failed to save user credentials. UserObj")
		}
		err3 := stub.PutState(userName, userObjBytes)
		if err3 != nil {
			return nil, errors.New("Failed to save User credentials")
		} 

		userLoginObj = userLogin{LoginName: userName, Password: password}
		userLoginBytes, err1 := json.Marshal(&userLoginObj)
		if err1 != nil {
			fmt.Println("Failed to save user credentials. UserObj")
		}

		err2 := stub.PutState(loginPrefix + userName, userLoginBytes)
		if err2 != nil {
			fmt.Println("Failed to save user credentials. UserLoginObj")
		}
		return nil, nil
	} else {
		return []byte("Not authorized to change user details"), nil
	}
	return nil, nil
}

func (t *SimpleChaincode) createTradeRequest(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var shipperID, tradeRequestIDString, producerID, entryLocation, tradeRequestStartDate, tradeRequestEndDate, tradeRequestStatus string
	var tradeRequestID, tradeRequestInvoiceID, tradeRequestIncidentID int
	var energyKWH, gasPrice float64
	var tradeRequestObj tradeRequest

	var tradeRequestShipperMap map[string][]byte
	var tradeRequestProducerMap map[string][]byte

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. 8 expected")
	}

	tradeRequestIDString = args[0]
	tradeRequestID, _ = strconv.Atoi(args[0])
	shipperID = args[1]
	producerID = args[2]
	energyKWH, _ = strconv.ParseFloat(args[3], 64)
	gasPrice, _ = strconv.ParseFloat(args[4], 64)
	entryLocation = args[5]
	tradeRequestStartDate = args[6]
	tradeRequestEndDate = args[7]
	tradeRequestStatus = "New"	
	tradeRequestInvoiceID = 0
	tradeRequestIncidentID = 0

	tradeRequestObj = tradeRequest{TradeRequestID: tradeRequestID, ShipperID: shipperID, ProducerID: producerID,
	EnergyKWH: energyKWH, GasPrice: gasPrice, EntryLocation: entryLocation, TradeRequestStartDate: tradeRequestStartDate,
	TradeRequestEndDate: tradeRequestEndDate, TradeRequestStatus: tradeRequestStatus, TradeRequestInvoiceID: tradeRequestInvoiceID,
	TradeRequestIncidentID: tradeRequestIncidentID}

	//Putting on RocksDB database.
	tradeRequestObjBytes, err1 := json.Marshal(tradeRequestObj)
	if err1 != nil {
		return nil, err1
	}
	_ = stub.PutState(tradeRequestIDString, tradeRequestObjBytes)

	//Putting in Maps for Shipper
	tradeRequestShipperMapObjBytes, err2 := stub.GetState(shipperID + "TradeRequestShipperMap")
	if err2 != nil {
		return nil, err2
	}
	if tradeRequestShipperMapObjBytes == nil {
		tradeRequestShipperMap = make(map[string][]byte)
		tradeRequestShipperMap[tradeRequestIDString] = tradeRequestObjBytes
		tradeRequestShipperMapObjBytes, _ = json.Marshal(&tradeRequestShipperMap)
		_ = stub.PutState(shipperID + "TradeRequestShipperMap", tradeRequestShipperMapObjBytes)
	} else {
		_ = json.Unmarshal(tradeRequestShipperMapObjBytes, &tradeRequestShipperMap)
		tradeRequestShipperMap[tradeRequestIDString] = tradeRequestObjBytes
		tradeRequestShipperMapObjBytes, _ = json.Marshal(&tradeRequestShipperMap)
		_ = stub.PutState(shipperID + "TradeRequestShipperMap", tradeRequestShipperMapObjBytes)
	}

	//Putting in Maps for Prodcuer
	tradeRequestProducerMapObjBytes, err3 := stub.GetState(producerID + "TradeRequestProducerMap")
	if err3 != nil {
		return nil, err3
	}
	if tradeRequestProducerMapObjBytes == nil {
		tradeRequestProducerMap = make(map[string][]byte)
		tradeRequestProducerMap[tradeRequestIDString] = tradeRequestObjBytes
		tradeRequestProducerMapObjBytes, _ = json.Marshal(&tradeRequestProducerMap)
		_ = stub.PutState(producerID + "TradeRequestProducerMap", tradeRequestProducerMapObjBytes)
	} else {
		_ = json.Unmarshal(tradeRequestProducerMapObjBytes, &tradeRequestProducerMap)
		tradeRequestProducerMap[tradeRequestIDString] = tradeRequestObjBytes
		tradeRequestProducerMapObjBytes, _ = json.Marshal(&tradeRequestProducerMap)
		_ = stub.PutState(producerID + "TradeRequestProducerMap", tradeRequestProducerMapObjBytes)
	}

	return nil, nil
}

func (t *SimpleChaincode) updateTradeRequestStatus(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var tradeRequestIDString string
	var tradeRequestObj tradeRequest
	
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. 2 expected (TradeRequestID and TradeRequestStatus)")
	}
	
	tradeRequestIDString = args[0]
	tradeRequestObjBytes, _ := stub.GetState(tradeRequestIDString)
	err1 := json.Unmarshal(tradeRequestObjBytes, &tradeRequestObj)
	if err1 != nil {
		return nil, err1
	}
	
	//Update the status
	tradeRequestObj.TradeRequestStatus = args[1];
	
	//Save the updated trade request
	tradeRequestBytes, err2 := json.Marshal(&tradeRequestObj)
	if err2 != nil {
		return nil, err2
	}
	err3 := stub.PutState(tradeRequestIDString, tradeRequestBytes)
	if err3 != nil {
		return nil, err3
	}
	
	return nil, nil;
}

/*func (t *SimpleChaincode) getTradeRequest(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var tradeRequestID string
	
	tradeRequestID = args[0]
	tradeRequestObjBytes, _ := stub.GetState(tradeRequestID)
	return []byte(string(tradeRequestObjBytes)), nil
}*/

func (t *SimpleChaincode) getShipperTradeRequestList(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var shipperID, returnMessage string
	var lenMap int
	mapShipperRequestInfo := make(map[string][]byte)
	fmt.Println("Getting Trade Requests for one shipper")

	shipperID = args[0]
	mapShipperRequestInfoBytes, _ := stub.GetState(shipperID + "TradeRequestShipperMap")
	_ = json.Unmarshal(mapShipperRequestInfoBytes, &mapShipperRequestInfo)
	lenMap = len(mapShipperRequestInfo)
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["

	for k, _ := range mapShipperRequestInfo {
		tradeRequestInfo, _ := stub.GetState(k)
		returnMessage = returnMessage + string(tradeRequestInfo)
		lenMap = lenMap - 1
		if (lenMap!= 0) {
			returnMessage = returnMessage + ","
		}
	}
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil
}

func (t *SimpleChaincode) getProducerTradeRequestList(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var producerID, returnMessage string
	var lenMap int
	mapProducerRequestInfo := make(map[string][]byte)
	fmt.Println("Getting Trade Requests for one Producer")

	producerID = args[0]
	mapProducerRequestInfoBytes, _ := stub.GetState(producerID + "TradeRequestProducerMap")
	_ = json.Unmarshal(mapProducerRequestInfoBytes, &mapProducerRequestInfo)
	lenMap = len(mapProducerRequestInfo)
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["

	for k, _ := range mapProducerRequestInfo {
		tradeRequestInfo, _ := stub.GetState(k)
		returnMessage = returnMessage + string(tradeRequestInfo)
		lenMap = lenMap - 1
		if (lenMap!= 0) {
			returnMessage = returnMessage + ","
		}
	}
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil
}

func testEqualSlice (a []byte, b []byte) bool {

	if a == nil && b == nil { 
        return true; 
    } else if a == nil || b == nil { 
        return false; 
    } 
	
	if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Running Invoke function")

	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "register" {
		return t.register(stub, args)
	} else if function == "createTradeRequest" {
		return t.createTradeRequest(stub, args)
	}
 
	fmt.Println("Invoke did not find function:" + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("Querying function: " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "verifyUser" {
		return t.verifyUser(stub, args)
	} else if function == "getUserInfo" {
		return t.getUserInfo(stub, args)
	} else if function == "getProducerList" {
		return t.getProducerList(stub)
	} else if function == "getProducerTradeRequestList" {
		return t.getProducerTradeRequestList(stub, args)
	} else if function == "getShipperTradeRequestList" {
		return t.getShipperTradeRequestList(stub, args)
	}
	fmt.Println("Query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("Running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0]                            //rename for fun
	value = args[1]
	err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    key = args[0]
    valAsbytes, err := stub.GetState(key)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }
    return valAsbytes, nil
}
