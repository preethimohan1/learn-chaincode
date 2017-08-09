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

var companyKey = "COMPANYIDLIST"
var userIDAffix = "USERLIST"
var tradeRequestKey = "TRADEREQUESTIDLIST"

type SimpleChaincode struct {

}

type company struct {
	CompanyID 		string 	`json:"company_id"`
	CompanyType		string 	`json:"company_type"`
	CompanyName 	string	`json:"company_name"`
	CompanyLocation	string	`json:"company_location"`
	BankBalance		float64	`json:"bank_balance"`
}

type user struct {
	UserID		string 	`json:"user_id"`
	Password 		string	`json:"user_password"`
    CompanyID 		string 	`json:"company_id"`
}

type userInfo struct {
	UserID		string 	`json:"user_id"`
    Company 	company `json:"company"`
}

type tradeRequest struct {
	TradeRequestID         int     `json:"tr_id"`
	ShipperID              string  `json:"tr_shipper_id"`
	ProducerID             string  `json:"tr_producer_id"`
	EnergyKWH              float64 `json:"tr_energy_kwh"`
	GasPrice               float64 `json:"tr_gas_price"`
	EntryLocation          string  `json:"tr_entry_location"`
	TradeRequestStartDate  string  `json:"tr_start_date"`
	TradeRequestEndDate    string  `json:"tr_end_date"`
	TradeRequestStatus     string  `json:"tr_status"`
	TradeRequestInvoiceID  int     `json:"tr_invoice_id"`
	TradeRequestIncidentID int     `json:"tr_incident_id"`
}

type CompanyIDList []string
type UserIDList []string
type TradeRequestIDList []string

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    
    var compIDArr CompanyIDList
    
    //Create default companies
    t.addCompany (stub, compIDArr, "PRODUCER1", "Producer", "Dong Energy", "Europe", 100000)
    t.addCompany (stub, compIDArr, "SHIPPER1", "Shipper", "RWE Supply and Trading", "Europe", 100000)
    t.addCompany (stub, compIDArr, "TRANSPORTER1", "Transporter", "Open Grid Europe", "Europe", 100000)
    t.addCompany (stub, compIDArr, "BUYER1", "Buyer", "EnBW", "Europe", 100000)
    
	//create Arrays for Each Type of User
	var producerInfoArr UserIDList
    t.addUser(stub, producerInfoArr, "producer", "producer", "PRODUCER1", "Producer")   
    
	var shipperInfoArr UserIDList
    t.addUser(stub, shipperInfoArr, "shipper", "shipper", "SHIPPER1", "Shipper")	
    
	var buyerInfoArr UserIDList
    t.addUser(stub, buyerInfoArr, "buyer", "buyer", "BUYER1", "Buyer")	
    
	var transporterInfoArr UserIDList
    t.addUser(stub, transporterInfoArr, "transporter", "transporter", "TRANSPORTER1", "Transporter")

	return nil, nil
}

func (t *SimpleChaincode) addCompany (stub shim.ChaincodeStubInterface, compIDArr CompanyIDList, compID string, 
				       compType string, compName string, compLoc string, bankBalance float64 ) bool {
    fmt.Println("Adding new company:"+ compName)
    
	var newCompany company
    
	newCompany = company{CompanyID: compID, CompanyType: compType, CompanyName: compName, 
	CompanyLocation: compLoc, BankBalance: bankBalance}
    
	compObjBytes, err := json.Marshal(&newCompany)
	if err != nil {
		fmt.Println(err)
		return false
	}

	err1 := stub.PutState(compID, compObjBytes)
	if err1 != nil {
		fmt.Println(err1)
        return false
	}    
    
    //Add companyID to a list
    compIDArr = append(compIDArr, compID)
    compIDArrBytes, _ := json.Marshal(compIDArr)
    _ = stub.PutState(companyKey, compIDArrBytes)     
    
    fmt.Println("Printing company ID List")
    fmt.Println(compIDArr)
    fmt.Println("Successfully added new company:"+ compName)
    	return true
}

func (t *SimpleChaincode) getCompanyList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var lenArr int
	var compIDArr CompanyIDList
	var companyType, returnMessage string
    var companyObj company
    
    if len(args) < 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1 (company type).")
	}
    
    companyType = args[0]    
	
	fmt.Println("Getting company list of type " + companyType)
    
    compIDArrBytes, _ := stub.GetState(companyKey)
	_ = json.Unmarshal(compIDArrBytes, &compIDArr)
    
   fmt.Println(compIDArr)
    
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["
	lenArr = len(compIDArr)
	for _, k := range compIDArr {
		compObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(compObjBytes, &companyObj)
        fmt.Println(companyObj)
        
        if(companyObj.CompanyType == companyType) {            
            returnMessage = returnMessage + string(compObjBytes) 
            lenArr = lenArr - 1 
            if (lenArr != 0) {
                returnMessage = returnMessage + ","
            } 
        }
	} 
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil

}

func (t *SimpleChaincode) addUser (stub shim.ChaincodeStubInterface, userIDArr UserIDList, userName string, 
				       password string, compID string, compType string ) bool {
    fmt.Println("Adding new user:"+ userName)
    
	var newUser user
    
    //Add user to login record
    newUser = user{UserID: userName, Password: password, CompanyID: compID} 
	userObjLoginBytes, _ := json.Marshal(&newUser)
	err2 := stub.PutState(userName, userObjLoginBytes)
	if err2 != nil {
		fmt.Println(err2)
	}
        
    //Add the user IDs into array of user types
    var arrKey = strings.ToLower(compType) + userIDAffix
    userIDArr = append(userIDArr, userName)
    userIDArrBytes, _ := json.Marshal(userIDArr)
    _ = stub.PutState(arrKey, userIDArrBytes)      
    
    fmt.Println("Successfully added new user:"+ userName)
    	return true
}

func (t *SimpleChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userName, password, companyJsonString, arrKey string
	var companyObj company
	var userArr UserIDList
	fmt.Println("Running function Register")

	if len (args) != 7 {
		return nil, errors.New("Incorrect number of argumets. Expecting 7")
	}
	
	userName = args[0]
	password = args[1]
	companyJsonString = args[2]
	err := json.Unmarshal([]byte(companyJsonString), &companyObj)
    if err != nil {
		return nil, err
	}
    fmt.Println(companyObj)
    
	arrKey = strings.ToLower(companyObj.CompanyType) + userIDAffix
    userArrObj, _ := stub.GetState(arrKey)
    _ = json.Unmarshal(userArrObj, &userArr)
    
    t.addUser (stub, userArr, userName, password, companyObj.CompanyID, companyObj.CompanyType )
    	
	return nil, nil

}

func (t *SimpleChaincode) getUserInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userName, returnMessage string
    var compStruct company
	var userInfoObj userInfo
    
	if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2 (userName and password).")
	}

    //Requires 2 arguments
	userName = args[0]
    //password = args[1]
    
    fmt.Println("Getting user info for user: "+userName)
    
	validUser, err1, compID := t.verifyUser(stub, args)
	if err1 != nil {
		return nil, err1
    }
    fmt.Println(compID)
    
	if validUser == true {
        fmt.Println("Valid user: "+userName)		
        
        //Get Company details
        compInfo, _ := stub.GetState(compID)	
        _ = json.Unmarshal(compInfo, &compStruct)
        
        userInfoObj.UserID = userName
        userInfoObj.Company = compStruct
        fmt.Println(userInfoObj)
        
        userInfoObjBytes, err2 := json.Marshal(userInfoObj)
        if err2 != nil {
		  return nil, err2
        }
        returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : " + string(userInfoObjBytes) + "} "
        
        
        /*returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : {\"user_id\" : \"" + userName + "\"," +
        "\"comp_id\" : \"" + compStruct.CompanyID + "\"," +
        "\"company_type\" : \"" + compStruct.CompanyType + "\"," +
        "\"company_name\" : \"" + compStruct.CompanyName + "\"," +
        "\"company_location\" : \"" + compStruct.CompanyLocation + "\"," +
        "\"bank_balance\" : \"" + compStruct.BankBalance + "\"} }"
           */     
        fmt.Println("User Info: "+ returnMessage)
		return []byte(returnMessage), nil
	} else {
        fmt.Println("Invalid user: "+userName)
        returnMessage = "{\"statusCode\" : \"FAIL\", \"body\" : \"ERROR: Invalid user !\"}"
		return []byte(returnMessage), nil
	}
	return nil, nil

}

func (t *SimpleChaincode) verifyUser(stub shim.ChaincodeStubInterface, args []string) (bool, error, string) {
	var userName, returnMessage, password string
	var loginObj user

	fmt.Println("Verifying User")
	if len(args) != 2 {
		return false, errors.New("Incorrect number of arguments. Expecting 2."), ""
	}

	userName = args[0]
	password = args[1]

	userInfo, err := stub.GetState(userName)
	if userInfo == nil {
        fmt.Println("Invalid Username")
		returnMessage = "Invalid Username"
        return false, errors.New(returnMessage), ""
	}

	err1 := json.Unmarshal(userInfo, &loginObj)
	if err1 != nil {
		return false, err, ""
	}
    fmt.Println(loginObj)
	if password == loginObj.Password {
		return true, nil, loginObj.CompanyID
	} else {        
        fmt.Println("Invalid Password")
		returnMessage = "Invalid Password"
		return false, errors.New(returnMessage), ""
	}
	return false, nil, ""
}


/*func (t *SimpleChaincode) updateUserInfo(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var userName, compType, compName, compLoc string
	var bankAccountNum int
	var bankBalance float64

	var userObj user

	userName = args[0]
	compType = args[1]
	compName = args[2]
	compLoc = args[3]
	bankAccountNum, _ = strconv.Atoi(args[4])
	bankBalance, _ = strconv.ParseFloat(args[5], 64)
		
    userObj = user{LoginID: userName, compType: compType, 
    CompanyName: compName, CompanyLocation: compLoc, BankAccountNum: bankAccountNum, BankBalance: bankBalance}
    userObjBytes, err := json.Marshal(&userObj)
    if err != nil {
        fmt.Println("Failed to marshal user info.")
        fmt.Println(err)
    }
    err3 := stub.PutState(userName, userObjBytes)
    if err3 != nil {
        return nil, errors.New("Failed to save User info")
        fmt.Println(err3)
    } 
	
	return nil, nil
}*/

func (t *SimpleChaincode) changePassword(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var userName, oldPassword, newPassword string	
	var userObj user

    userName = args[0]
    oldPassword = args[1]
    newPassword = args[2]    
    
    argsVerify := []string{userName, oldPassword}
	validUser, _ , compID := t.verifyUser(stub, argsVerify)
    
	if validUser == true {		
        userObj = user{UserID: userName, Password: newPassword, CompanyID: compID}
		userObjBytes, err1 := json.Marshal(&userObj)
		if err1 != nil {
            return []byte("Failed to marshal new password credentials."), err1
		}

		err2 := stub.PutState(userName, userObjBytes)
		if err2 != nil {
            return []byte("Failed to update password."), err2
		}
          
		return nil, nil
	} else {
		return []byte("ERROR! Not authorized to change password."), nil
	}
	return nil, nil
}

func (t *SimpleChaincode) createTradeRequest(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var shipperID, tradeRequestIDString, producerID, entryLocation, tradeRequestStartDate, tradeRequestEndDate, tradeRequestStatus string
	var tradeRequestID, tradeRequestInvoiceID, tradeRequestIncidentID int
	var energyKWH, gasPrice float64
	var tradeRequestObj tradeRequest

	var tradeRequestIDArr TradeRequestIDList

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

	//Putting in TR ID list    
	tradeRequestIDListObjBytes, err2 := stub.GetState(tradeRequestKey)
	if err2 != nil {
		return nil, err2
	}
	if tradeRequestIDListObjBytes != nil {
		_ = json.Unmarshal(tradeRequestIDListObjBytes, &tradeRequestIDArr)
	}
    tradeRequestIDArr = append(tradeRequestIDArr, tradeRequestIDString)	
    tradeRequestIDListObjBytes, _ = json.Marshal(&tradeRequestIDArr)
    _ = stub.PutState(tradeRequestKey, tradeRequestIDListObjBytes)

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
	tradeRequestObj.TradeRequestStatus = args[1]
	
	//Save the updated trade request
	tradeRequestBytes, err2 := json.Marshal(&tradeRequestObj)
	if err2 != nil {
		return nil, err2
	}
	err3 := stub.PutState(tradeRequestIDString, tradeRequestBytes)
	if err3 != nil {
		return nil, err3
	}
	
	return nil, nil
}

/*func (t *SimpleChaincode) getTradeRequest(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var tradeRequestID string
	
	tradeRequestID = args[0]
	tradeRequestObjBytes, _ := stub.GetState(tradeRequestID)
	return []byte(string(tradeRequestObjBytes)), nil
}*/

func (t *SimpleChaincode) getTradeRequestList(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var companyID, returnMessage string
	var lenMap int	
    var trList TradeRequestIDList
    var tradeRequestObj tradeRequest
    
    companyID = args[0]
    
    fmt.Println("Getting Trade Requests for company: "+ companyID)
	
	trLisObjBytes, _ := stub.GetState(tradeRequestKey)
	_ = json.Unmarshal(trLisObjBytes, &trList)
    
	lenMap = len(trList)
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["

	for _, k := range trList {
        
		tradeRequestObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(tradeRequestObjBytes, &tradeRequestObj)
        fmt.Println(tradeRequestObj)
        
        if(tradeRequestObj.ShipperID == companyID || tradeRequestObj.ProducerID == companyID) {
            returnMessage = returnMessage + string(tradeRequestObjBytes)
            
            lenMap = lenMap - 1
            if (lenMap!= 0) {
                returnMessage = returnMessage + ","
            }
        }        
	}
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil
}

func testEqualSlice (a []byte, b []byte) bool {

	if a == nil && b == nil { 
        return true
    } else if a == nil || b == nil { 
        return false
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
	} else if function == "changePassword" {
		return t.changePassword(stub, args)
	} /*else if function == "updateUserInfo" {
		return t.updateUserInfo(stub, args)
	}*/
 
	fmt.Println("Invoke did not find function:" + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    fmt.Println("Querying function: " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} else if function == "getUserInfo" {
		return t.getUserInfo(stub, args)
	} else if function == "getCompanyList" {
		return t.getCompanyList(stub, args)
	} else if function == "getTradeRequestList" {
		return t.getTradeRequestList(stub, args)
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
