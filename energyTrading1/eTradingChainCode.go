package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
    	"strings"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var companyKey = "COMPANYIDLIST"
var userIDAffix = "USERLIST"
var tradeRequestKey = "TRADEREQUESTIDLIST"
var transportRequestKey = "TRANSPORTREQUESTIDLIST"
var gasRequestKey = "GASREQUESTIDLIST"
var planKey = "PLANIDLIST"
var planIDPrefix = "PLAN_"

type SimpleChaincode struct {

}

type company struct {
	CompanyID 		string 	`json:"company_id"`
	CompanyType		string 	`json:"company_type"`
	CompanyName 	string	`json:"company_name"`
	CompanyLocation	string	`json:"company_location"`
	BankBalance		float64	`json:"bank_balance"`
    BalanceUpdatedDate		string	`json:"bank_balance_date"`
}

type user struct {
	UserID		    string 	`json:"user_id"`
	Password 		string	`json:"user_password"`
    CompanyID 		string 	`json:"company_id"`
}

type businessPlan struct {
	PlanID 		    string 	`json:"bp_plan_id"`
    PlanDate 		string 	`json:"bp_plan_date"`
	GasPrice		float64	`json:"bp_gas_price"`    
	EntryLocation 	string	`json:"bp_entry_location"`
	EntryCapacity	int	    `json:"bp_entry_capacity"`
	ExitLocation 	string	`json:"bp_exit_location"`
	ExitCapacity	int	    `json:"bp_exit_capacity"`
    CompanyID 		string 	`json:"bp_company_id"`
}

type userInfo struct {
	UserID		 string 	  `json:"user_id"`
    Company 	 company      `json:"company"`
    BusinessPlan businessPlan `json:"business_plan"`
}

type businessPlanInfo struct {    
    BusinessPlan businessPlan  `json:"business_plan"`
    Company 	 company       `json:"company"`
}

type contract struct {
	ContractID         int     `json:"contract_id"`
	InitiatorID        string  `json:"contract_initiator_id"`
	ReceiverID         string  `json:"contract_receiver_id"`
	EnergyMWH          float64 `json:"contract_energy_mwh"`
    EntryLocation      string  `json:"contract_entry_location"`
	ContractStartDate  string  `json:"contract_start_date"`
	ContractEndDate    string  `json:"contract_end_date"`
	ContractStatus     string  `json:"contract_status"`
	InvoiceID          int     `json:"contract_invoice_id"`
	IncidentID         int     `json:"contract_incident_id"`
}

type contractInfo struct {
    Contract     contract       `json:"contract"`
    InitiatorCompany company    `json:"initiator_company"`
    ReceiverCompany  company    `json:"receiver_company"`
}

type CompanyIDList []string
type UserIDList []string
type BusinessPlanIDList []string

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    var currentDate string    
    year, month, day := time.Now().Date()
    var monthInNumber int = int(month) //convert time.Month to integer
    currentDate = strconv.Itoa(day) + "/" + strconv.Itoa(monthInNumber) + "/" + strconv.Itoa(year)
    
    //Create default companies
    var compIDArr CompanyIDList
    t.addCompany (stub, compIDArr, "PRODUCER1", "Producer", "Dong Energy", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "PRODUCER1")
    t.addCompany (stub, compIDArr, "PRODUCER2", "Producer", "Gaz Promp", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "PRODUCER2")
    
    t.addCompany (stub, compIDArr, "SHIPPER1", "Shipper", "RWE Supply and Trading", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "SHIPPER1")
    t.addCompany (stub, compIDArr, "SHIPPER2", "Shipper", "UNIPER Energy Trading", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "SHIPPER2")
    
    t.addCompany (stub, compIDArr, "TRANSPORTER1", "Transporter", "Open Grid Europe", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "TRANSPORTER1")
    t.addCompany (stub, compIDArr, "TRANSPORTER2", "Transporter", "ONTRAS GMBH", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "TRANSPORTER2")
    t.addCompany (stub, compIDArr, "TRANSPORTER3", "Transporter", "Gasunie DTS", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "TRANSPORTER3")
    
    t.addCompany (stub, compIDArr, "BUYER1", "Buyer", "EnBW", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "BUYER1")
    t.addCompany (stub, compIDArr, "BUYER2", "Buyer", "Vattenfall", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "BUYER2")
    
	//create default users
	var userIDArr UserIDList
    t.addUser(stub, userIDArr, "producer1", "producer1", "PRODUCER1", "Producer")
    userIDArr = append(userIDArr, "producer1")
    t.addUser(stub, userIDArr, "producer2", "producer2", "PRODUCER2", "Producer")     
	userIDArr = append(userIDArr, "producer2")
    
    t.addUser(stub, userIDArr, "shipper1", "shipper1", "SHIPPER1", "Shipper")	
    userIDArr = append(userIDArr, "shipper1")
    t.addUser(stub, userIDArr, "shipper2", "shipper2", "SHIPPER2", "Shipper")	
    userIDArr = append(userIDArr, "shipper2")
	
    t.addUser(stub, userIDArr, "transporter1", "transporter1", "TRANSPORTER1", "Transporter")
    userIDArr = append(userIDArr, "transporter1")
    t.addUser(stub, userIDArr, "transporter2", "transporter2", "TRANSPORTER2", "Transporter")
    userIDArr = append(userIDArr, "transporter2")
    t.addUser(stub, userIDArr, "transporter3", "transporter3", "TRANSPORTER3", "Transporter")
    userIDArr = append(userIDArr, "transporter3")
	
    t.addUser(stub, userIDArr, "buyer1", "buyer1", "BUYER1", "Buyer")	
    userIDArr = append(userIDArr, "buyer1")
    t.addUser(stub, userIDArr, "buyer2", "buyer2", "BUYER2", "Buyer")	
    userIDArr = append(userIDArr, "buyer2")
    
    //Create business plan for producers
    var planID string
    var bpIDList BusinessPlanIDList
    planID = planIDPrefix + "PRODUCER1"
    t.createBusinessPlan(stub, bpIDList, planID, currentDate, 12.0, "Europe", 200, "Wardenburg", 200, "PRODUCER1")     
    bpIDList = append(bpIDList, planID)
    planID = planIDPrefix + "PRODUCER2"
    t.createBusinessPlan(stub, bpIDList, planID, currentDate, 10.0, "Europe", 300, "Ellund", 300, "PRODUCER2")
    bpIDList = append(bpIDList, planID)
    
    //Create business plan for trasporters
    planID = planIDPrefix + "TRANSPORTER1"
    t.createBusinessPlan(stub, bpIDList, planID, currentDate, 11.0, "Wardenburg", 200, "Bunder-Tief", 100, "TRANSPORTER1")  
    bpIDList = append(bpIDList, planID)
    
    planID = planIDPrefix + "TRANSPORTER2"
    t.createBusinessPlan(stub, bpIDList, planID, currentDate, 9.0, "Ellund", 300, "Steinbrink", 150, "TRANSPORTER2")
    bpIDList = append(bpIDList, planID)
    
    planID = planIDPrefix + "TRANSPORTER3"
    t.createBusinessPlan(stub, bpIDList, planID, currentDate, 8.0, "Ellund", 350, "Steinitz", 175, "TRANSPORTER3")
    bpIDList = append(bpIDList, planID)
    
	return nil, nil
}

func (t *SimpleChaincode) addCompany (stub shim.ChaincodeStubInterface, compIDArr CompanyIDList, compID string, 
				       compType string, compName string, compLoc string, bankBalance float64,  balanceDate string) bool {
    fmt.Println("Adding new company:"+ compName)
   
	var newCompany company
    
	newCompany = company{CompanyID: compID, CompanyType: compType, CompanyName: compName, 
                         CompanyLocation: compLoc, BankBalance: bankBalance, BalanceUpdatedDate: balanceDate}
    
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
        
        if(strings.ToLower(companyType) == "all" || strings.ToLower(companyObj.CompanyType) == strings.ToLower(companyType)) {        
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

	if len (args) < 3 {
        return nil, errors.New("Incorrect number of arguments. Expecting 3 (userName, password, company)")
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
    var busPlanStruct businessPlan
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
        fmt.Println(err1)
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
        
        //Get Business Plan info
        if compStruct.CompanyType == "Producer" || compStruct.CompanyType == "Transporter" {
            bpInfo, _ := stub.GetState(planIDPrefix + compID)	
            _ = json.Unmarshal(bpInfo, &busPlanStruct)
            userInfoObj.BusinessPlan = busPlanStruct
        }
        
        userInfoObjBytes, err2 := json.Marshal(userInfoObj)
        if err2 != nil {
		  return nil, err2
        }
        returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : " + string(userInfoObjBytes) + "} "
        
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

	userInfo, _ := stub.GetState(userName)
	if userInfo == nil {
		returnMessage = "Invalid Username"
        return false, errors.New(returnMessage), ""
	}

	err1 := json.Unmarshal(userInfo, &loginObj)
	if err1 != nil {
		return false, err1, ""
	}
    fmt.Println(loginObj)
	if password == loginObj.Password {
		return true, nil, loginObj.CompanyID
	} else {        
		returnMessage = "Invalid Password"
		return false, errors.New(returnMessage), ""
	}
	return false, nil, ""
}


func (t *SimpleChaincode) topupBankBalance(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var compID, topupDate string
	var topupAmount float64
	var companyObj company
    
    fmt.Println("Entered function topupBankBalance()")
    
    if len(args) < 3 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2 arguments (CompanyID, top-up amount, top-up date).")
	}

	compID = args[0]
	topupAmount, _ = strconv.ParseFloat(args[1], 64)
    topupDate = args[2]
    
    //Get the company object from DB
    compObjBytes, _ := stub.GetState(compID)
    _ = json.Unmarshal(compObjBytes, &companyObj)
    fmt.Println(companyObj)
    
    //Topup the amount   
    companyObj.BankBalance = companyObj.BankBalance + topupAmount   
    companyObj.BalanceUpdatedDate = topupDate
        
    companyObjBytes, err := json.Marshal(&companyObj)
    if err != nil {
        fmt.Println("Failed to marshal company info.")
        fmt.Println(err)
    }
    err3 := stub.PutState(compID, companyObjBytes)
    if err3 != nil {
        return nil, errors.New("Failed to save Company info")
        fmt.Println(err3)
    } 
	
	return nil, nil
}

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

func (t *SimpleChaincode) createBusinessPlan(stub shim.ChaincodeStubInterface, bpIDList BusinessPlanIDList, planID string, 
                                             planDate string, gasPrice float64, entryLocation string, entryCapacity int, exitLocation string, exitCapacity int, compID string) ([]byte, error) {
    fmt.Println("Creating new Business Plan: " + planID)
    
    var businessPlanObj businessPlan
    
    businessPlanObj = businessPlan{PlanID: planID, PlanDate: planDate, GasPrice: gasPrice, EntryLocation: entryLocation, EntryCapacity: entryCapacity, ExitLocation: exitLocation, ExitCapacity: exitCapacity, CompanyID: compID}
    
    businessPlanObjBytes, err1 := json.Marshal(businessPlanObj)
    if err1 != nil {
		return nil, err1
	}
    err2 := stub.PutState(planID, businessPlanObjBytes)
    if err2 != nil {
		return nil, err2
	}
    
    //Add the plan IDs into array of business plan ID
    if bpIDList != nil {
        bpIDList = append(bpIDList, planID)
        bpIDListBytes, _ := json.Marshal(bpIDList)
        _ = stub.PutState(planKey, bpIDListBytes) 
        
        fmt.Println(bpIDList)
    }
        
    return nil, nil
}

func (t *SimpleChaincode) getBusinessPlanList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var lenArr int
	var bpIDArr BusinessPlanIDList
	var returnMessage string
    var bpObj businessPlan
    var companyObj company
    var bpInfoObj businessPlanInfo
    
	fmt.Println("Getting all business plans.")
    
    bpIDArrBytes, _ := stub.GetState(planKey)
	_ = json.Unmarshal(bpIDArrBytes, &bpIDArr)
    
   fmt.Println(bpIDArr)
    
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["
	lenArr = len(bpIDArr)
	for _, k := range bpIDArr {
        //Fetch the Business plan
		bpObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(bpObjBytes, &bpObj)
        fmt.Println(bpObj)
        
        bpInfoObj.BusinessPlan = bpObj;
        
        //Fetch the company details
        companyObjBytes, _ := stub.GetState(bpObj.CompanyID)
        _ = json.Unmarshal(companyObjBytes, &companyObj)
        fmt.Println(companyObj)
        
        bpInfoObj.Company = companyObj;
        bpInfoObjBytes, _ := json.Marshal(bpInfoObj)
        
        returnMessage = returnMessage + string(bpInfoObjBytes) 
        lenArr = lenArr - 1 
        if (lenArr != 0) {
            returnMessage = returnMessage + ","
        } 
	} 
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil
}

func (t *SimpleChaincode) updateBusinessPlan(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
    fmt.Println("Entering function updateBusinessPlan()")
    
    var gasPrice float64
    var entryCapacity, exitCapacity int
    gasPrice, _ = strconv.ParseFloat(args[2], 64)
    entryCapacity, _ = strconv.Atoi(args[4])
    exitCapacity, _ = strconv.Atoi(args[6])
       
    _, err := t.createBusinessPlan(stub, nil, args[0], args[1], gasPrice, args[3], entryCapacity, args[5], exitCapacity, args[7])
    if err != nil {
		return nil, err
	}
	return nil, nil
}


func (t *SimpleChaincode) createContract(stub shim.ChaincodeStubInterface, idArrKey string, args[] string) ([]byte, error) {
    
	var initiatorID, contractIDString, receiverID, contractStartDate, contractEndDate, contractStatus string
	var contractID, contractInvoiceID, contractIncidentID int
	var energyMWH float64
	var contractObj contract
    var contractIDArr []string
    
	if len(args) < 6 {
		return nil, errors.New("Incorrect number of arguments. 6 expected")
	}
    
    fmt.Println("Creating new contract...")

	contractIDString = args[0]
	contractID, _ = strconv.Atoi(args[0])
	initiatorID = args[1]
	receiverID = args[2]
	energyMWH, _ = strconv.ParseFloat(args[3], 64)
	contractStartDate = args[4]
	contractEndDate = args[5]
	contractStatus = "New"	
	contractInvoiceID = 0
	contractIncidentID = 0

	contractObj = contract{ContractID: contractID, InitiatorID: initiatorID, ReceiverID: receiverID,
	EnergyMWH: energyMWH, ContractStartDate: contractStartDate, ContractEndDate: contractEndDate, 
	ContractStatus: contractStatus, InvoiceID: contractInvoiceID,
	IncidentID: contractIncidentID}

	//Putting on RocksDB database.
	contractObjBytes, err1 := json.Marshal(contractObj)
	if err1 != nil {
		return nil, err1
	}
	_ = stub.PutState(contractIDString, contractObjBytes)
    
    //Putting contract ID in Contract ID array  
	contractIDListObjBytes, err := stub.GetState(idArrKey)
	if err != nil {
		return nil, err
	}
	if contractIDListObjBytes != nil {
		_ = json.Unmarshal(contractIDListObjBytes, &contractIDArr)
	}    
    contractIDArr = append(contractIDArr, contractIDString)	
    contractIDListObjBytes, _ = json.Marshal(&contractIDArr)
    _ = stub.PutState(idArrKey, contractIDListObjBytes)
    
    fmt.Println(contractIDArr)
    
	return nil, nil
}

func (t *SimpleChaincode) createTradeRequest(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
    fmt.Println("Creating new trade request: " + args[0])
    return t.createContract(stub, tradeRequestKey, args)
}

func (t *SimpleChaincode) createTransportRequest(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
    fmt.Println("Creating new transport request: " + args[0])
    return t.createContract(stub, transportRequestKey, args)
}

func (t *SimpleChaincode) createGasRequest(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
    fmt.Println("Creating new gas request: " + args[0])
    return t.createContract(stub, gasRequestKey, args)
}

func (t *SimpleChaincode) updateContractStatus(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var contractIDString string
	var contractObj contract
	
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. 2 expected (ContractID and ContractStatus)")
	}
	
	contractIDString = args[0]
	contractObjBytes, _ := stub.GetState(contractIDString)
	err1 := json.Unmarshal(contractObjBytes, &contractObj)
	if err1 != nil {
		return nil, err1
	}
	
	//Update the status
	contractObj.ContractStatus = args[1]
	
	//Save the updated trade request
	contractBytes, err2 := json.Marshal(&contractObj)
	if err2 != nil {
		return nil, err2
	}
	err3 := stub.PutState(contractIDString, contractBytes)
	if err3 != nil {
		return nil, err3
	}
	
	return nil, nil
}

/*func (t *SimpleChaincode) getContract(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
	var contractID string
	
	contractID = args[0]
	contractObjBytes, _ := stub.GetState(contractID)
	return []byte(string(contractObjBytes)), nil
}*/

func (t *SimpleChaincode) getContractList(stub shim.ChaincodeStubInterface, idArrKey string, args[] string) ([]byte, error) {
	var companyID, returnMessage string
	var lenMap int	
    var contractIDList []string
    var contractObj contract
    var contractFullObj contractInfo
    var initiatorCompany, receiverCompany company
    
    companyID = args[0]
    
    fmt.Println("Getting Contracts for company: "+ companyID)
	
	contractListObjBytes, _ := stub.GetState(idArrKey)
	_ = json.Unmarshal(contractListObjBytes, &contractIDList)
    
	lenMap = len(contractIDList)
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["

	for _, k := range contractIDList {
        
		contractObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(contractObjBytes, &contractObj)
        fmt.Println(contractObj)
        
        if(contractObj.InitiatorID == companyID || contractObj.ReceiverID == companyID) {
            contractFullObj.Contract = contractObj
            
            //Add Initiator company object
            initiatorObjBytes, _ := stub.GetState(contractObj.InitiatorID)
            _ = json.Unmarshal(initiatorObjBytes, &initiatorCompany)
            contractFullObj.InitiatorCompany = initiatorCompany
            
            //Add Receiver company object
            receiverObjBytes, _ := stub.GetState(contractObj.ReceiverID)
            _ = json.Unmarshal(receiverObjBytes, &receiverCompany)
            contractFullObj.ReceiverCompany = receiverCompany
            
            contractFullObjBytes, err1 := json.Marshal(contractFullObj)
            if err1 != nil {
              return nil, err1
            }
            
            returnMessage = returnMessage + string(contractFullObjBytes)
            
            lenMap = lenMap - 1
            if (lenMap!= 0) {
                returnMessage = returnMessage + ","
            }
        }        
	}
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil
}

func (t *SimpleChaincode) getTradeRequestList(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
    return t.getContractList(stub, tradeRequestKey, args)
}

func (t *SimpleChaincode) getTransportRequestList(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
    return t.getContractList(stub, transportRequestKey, args)
}

func (t *SimpleChaincode) getGasRequestList(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
    return t.getContractList(stub, gasRequestKey, args)
}

func (t *SimpleChaincode) addIOTData (stub shim.ChaincodeStubInterface, args[] string ) ([]byte, error) {
    fmt.Println("Adding new IOT Data: "+ args[0])
    
	
    return nil, nil
}
                                                                                         
func (t *SimpleChaincode) readAssetSchemas (stub shim.ChaincodeStubInterface, args[] string ) ([]byte, error) {
    fmt.Println("readAssetSchemas: Adding new IOT Data Obj: "+ args)
    fmt.Println("readAssetSchemas: Adding new IOT Data: "+ args[0]])
	
    return nil, nil
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
	} else if function == "createTransportRequest" {
		return t.createTransportRequest(stub, args)
	} else if function == "createGasRequest" {
		return t.createGasRequest(stub, args)
	} else if function == "changePassword" {
		return t.changePassword(stub, args)
	} else if function == "updateContractStatus" {
		return t.updateContractStatus(stub, args)
	} else if function == "updateBusinessPlan" {
		return t.updateBusinessPlan(stub, args)
	} else if function == "topupBankBalance" {
		return t.topupBankBalance(stub, args)
	} else if function == "addIOTData" {
		return t.addIOTData(stub, args)
	}  else if function == "readAssetSchemas" {
		return t.readAssetSchemas(stub, args)
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
	} else if function == "getUserInfo" {
		return t.getUserInfo(stub, args)
	} else if function == "getCompanyList" {
		return t.getCompanyList(stub, args)
	} else if function == "getTradeRequestList" {
		return t.getTradeRequestList(stub, args)
    } else if function == "getTransportRequestList" {
		return t.getTransportRequestList(stub, args)
    } else if function == "getGasRequestList" {
		return t.getGasRequestList(stub, args)
    } else if function == "getBusinessPlanList" {
		return t.getBusinessPlanList(stub, args)
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
