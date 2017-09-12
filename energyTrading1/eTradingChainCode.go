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
var tradeRequestKey = "TRADEREQUESTIDLIST"
var transportRequestKey = "TRANSPORTREQUESTIDLIST"
var gasRequestKey = "GASREQUESTIDLIST"
var planKey = "PLANIDLIST"
var allKeys = "ALLKEYS"

var userIDAffix = "_USERLIST"    //<CompanyType>_USERLIST
var planIDAffix = "_PLAN"      // <CompanyID>_PLAN
var iotKeyAffix = "_IOTDATA"    //<CompanyID>_IOTDATA
var invoiceAffix = "_INVOICELIST"   //<ContractID>_INVOICELIST
var incidentAffix = "_INCIDENTLIST" //<ContractID>_INCIDENTLIST

type SimpleChaincode struct {

}

type company struct {
	CompanyID 		string 	`json:"company_id"`
	CompanyType		string 	`json:"company_type"`
	CompanyName 	string	`json:"company_name"`
	CompanyLocation	string	`json:"company_location"`
	BankBalance		float64	`json:"bank_balance"`
    BalanceUpdatedDateMS		int	`json:"bank_balance_date_ms"`
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
}

type contractInfo struct {
    Contract     contract       `json:"contract"`
    InitiatorCompany company    `json:"initiator_company"`
    ReceiverCompany  company    `json:"receiver_company"`
    BusinessPlan     businessPlan `json:"business_plan"`
    InvoiceList     []invoice   `json:"invoice_list"`
    IncidentList    []incident  `json:"incident_list"`
}

type flowMeterData struct {
	DeviceID           string  `json:"device_id"`
	DeviceLocation     string  `json:"device_location"`
	CompanyID          string  `json:"company_id"`
	PressureKPA        int     `json:"pressure_kpa"`
    TemperatureC       int     `json:"temperature_c"`
	SpecificGravity    float64 `json:"specific_gravity"`
	EnergyMWH          float64 `json:"energy_mwh"`
	TimestampMS        int     `json:"timestamp_ms"`
}

type invoice struct {
	InvoiceID          int     `json:"invoice_id"`
	InvoiceDateMS      int     `json:"invoice_date_ms"`
	PaymentStatus      string  `json:"payment_status"`
	PaymentDateMS      int     `json:"payment_date_ms"`
    ContractID         int     `json:"contract_id"`
}

type incident struct {
	IncidentID          int     `json:"incident_id"`
	IncidentDateMS      int     `json:"incident_date_ms"`
    IncidentStatus      string  `json:"incident_status"`
	ExpectedEnergyMWH   float64 `json:"expected_energy_mwh"`
    ActualEnergyMWH     float64 `json:"actual_energy_mwh"`
    ContractID          int     `json:"contract_id"`
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

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) ([]byte, error) {
    //Initialize master keys list
    masterKeyList := []string{companyKey, tradeRequestKey, transportRequestKey, gasRequestKey, planKey, 
                              "buyer"+ userIDAffix, "shipper"+ userIDAffix, "producer"+ userIDAffix, "transporter"+ userIDAffix,
                             "BUYER1"+ iotKeyAffix, "BUYER2"+ iotKeyAffix, "PRODUCER1"+ iotKeyAffix, "PRODUCER2"+ iotKeyAffix, "TRANSPORTER1"+ iotKeyAffix, "TRANSPORTER2" + iotKeyAffix, "TRANSPORTER3" + iotKeyAffix}
    
    var currentDate int
    currentDate = 0
    
    var currentDateStr string
    year, month, day := time.Now().Date()
    var monthInNumber int = int(month) //convert time.Month to integer
    currentDateStr = strconv.Itoa(day) + "/" + strconv.Itoa(monthInNumber) + "/" + strconv.Itoa(year)
    
    //Create default companies
    var compIDArr CompanyIDList
    t.addCompany (stub, compIDArr, "BUYER1", "Buyer", "EnBW", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "BUYER1")
    t.addCompany (stub, compIDArr, "BUYER2", "Buyer", "Vattenfall", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "BUYER2")
    
    t.addCompany (stub, compIDArr, "SHIPPER1", "Shipper", "RWE Supply and Trading", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "SHIPPER1")
    t.addCompany (stub, compIDArr, "SHIPPER2", "Shipper", "UNIPER Energy Trading", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "SHIPPER2")
    
    t.addCompany (stub, compIDArr, "PRODUCER1", "Producer", "Dong Energy", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "PRODUCER1")
    t.addCompany (stub, compIDArr, "PRODUCER2", "Producer", "Gaz Promp", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "PRODUCER2")
    
    t.addCompany (stub, compIDArr, "TRANSPORTER1", "Transporter", "Open Grid Europe", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "TRANSPORTER1")
    t.addCompany (stub, compIDArr, "TRANSPORTER2", "Transporter", "ONTRAS GMBH", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "TRANSPORTER2")
    t.addCompany (stub, compIDArr, "TRANSPORTER3", "Transporter", "Gasunie DTS", "Europe", 100000, currentDate)
    compIDArr = append(compIDArr, "TRANSPORTER3")
    
    masterKeyList = append(masterKeyList, compIDArr...)
    
	//create default users
	var userIDArr UserIDList
    t.addUser(stub, userIDArr, "buyer1", "buyer1", "BUYER1", "Buyer")	
    userIDArr = append(userIDArr, "buyer1")
    t.addUser(stub, userIDArr, "buyer2", "buyer2", "BUYER2", "Buyer")	
    userIDArr = append(userIDArr, "buyer2")
    
    t.addUser(stub, userIDArr, "shipper1", "shipper1", "SHIPPER1", "Shipper")	
    userIDArr = append(userIDArr, "shipper1")
    t.addUser(stub, userIDArr, "shipper2", "shipper2", "SHIPPER2", "Shipper")	
    userIDArr = append(userIDArr, "shipper2")
    
    t.addUser(stub, userIDArr, "producer1", "producer1", "PRODUCER1", "Producer")
    userIDArr = append(userIDArr, "producer1")
    t.addUser(stub, userIDArr, "producer2", "producer2", "PRODUCER2", "Producer")     
	userIDArr = append(userIDArr, "producer2")
    
    t.addUser(stub, userIDArr, "transporter1", "transporter1", "TRANSPORTER1", "Transporter")
    userIDArr = append(userIDArr, "transporter1")
    t.addUser(stub, userIDArr, "transporter2", "transporter2", "TRANSPORTER2", "Transporter")
    userIDArr = append(userIDArr, "transporter2")
    t.addUser(stub, userIDArr, "transporter3", "transporter3", "TRANSPORTER3", "Transporter")
    userIDArr = append(userIDArr, "transporter3")
	
    masterKeyList = append(masterKeyList, userIDArr...)
    
    //Create business plans
    var planID string
    var bpIDList BusinessPlanIDList
    
    //Create business plan for shippers
    planID = "SHIPPER1" + planIDAffix
    t.createBusinessPlan(stub, bpIDList, planID, currentDateStr, 14.0, "Europe", 0, "Bunder-Tief, Steinbrink", 0, "SHIPPER1") 
    bpIDList = append(bpIDList, planID)
    
    planID = "SHIPPER2" + planIDAffix
    t.createBusinessPlan(stub, bpIDList, planID, currentDateStr, 15.0, "Steinitz", 0, "Steinitz", 0, "SHIPPER2")  
    bpIDList = append(bpIDList, planID)
    
    //Create business plan for producers
    planID = "PRODUCER1" + planIDAffix
    t.createBusinessPlan(stub, bpIDList, planID, currentDateStr, 12.0, "Wardenburg", 200, "Wardenburg", 200, "PRODUCER1")     
    bpIDList = append(bpIDList, planID)
    
    planID = "PRODUCER2" + planIDAffix
    t.createBusinessPlan(stub, bpIDList, planID, currentDateStr, 10.0, "Ellund", 300, "Ellund", 300, "PRODUCER2")
    bpIDList = append(bpIDList, planID)
    
    //Create business plan for trasporters
    planID = "TRANSPORTER1" + planIDAffix
    t.createBusinessPlan(stub, bpIDList, planID, currentDateStr, 11.0, "Wardenburg", 200, "Bunder-Tief", 100, "TRANSPORTER1")  
    bpIDList = append(bpIDList, planID)
    
    planID = "TRANSPORTER2" + planIDAffix
    t.createBusinessPlan(stub, bpIDList, planID, currentDateStr, 9.0, "Ellund", 300, "Steinbrink", 150, "TRANSPORTER2")
    bpIDList = append(bpIDList, planID)
    
    planID = "TRANSPORTER3" + planIDAffix
    t.createBusinessPlan(stub, bpIDList, planID, currentDateStr, 8.0, "Ellund", 350, "Steinitz", 175, "TRANSPORTER3")
    bpIDList = append(bpIDList, planID)
    
    masterKeyList = append(masterKeyList, bpIDList...)
    
    updateMasterKeyList (stub, masterKeyList)
	return nil, nil
}

func (t *SimpleChaincode) updateMasterKeyList(stub shim.ChaincodeStubInterface, keys []string) ([]byte, error) {
    var masterKeyList []string
    
    //Get the existing master array of keys
    keyListBytes, _ := stub.GetState(allKeys)
    if keyListBytes != nil {
	   _ = json.Unmarshal(keyListBytes, &masterKeyList)
    }
    
    //If the key already exists in the master list, then return to avoid duplicates
    if(len(keys) == 1 && contains(masterKeyList, keys[0]) ) {
        return nil, nil
    }
    
    //Append the new key to the master array
    masterKeyList = append(masterKeyList, keys...)
    keyListBytes, err := json.Marshal(masterKeyList)
    _ = stub.PutState(allKeys, keyListBytes)  
    
    return nil, nil
}

func (t *SimpleChaincode) getMasterKeyList(stub shim.ChaincodeStubInterface) ([]byte, error) {
    var masterKeyList []string
    
    //Get the existing master array of keys
    keyListBytes, _ := stub.GetState(allKeys)
    if keyListBytes != nil
	   _ = json.Unmarshal(keyListBytes, &masterKeyList)
    
    
    return []byte(masterKeyList), nil
}


func (t *SimpleChaincode) addCompany (stub shim.ChaincodeStubInterface, compIDArr CompanyIDList, compID string, 
				       compType string, compName string, compLoc string, bankBalance float64,  balanceDate int) bool {
    fmt.Println("Adding new company:"+ compName)
   
	var newCompany company
    
	newCompany = company{CompanyID: compID, CompanyType: compType, CompanyName: compName, 
                         CompanyLocation: compLoc, BankBalance: bankBalance, BalanceUpdatedDateMS: balanceDate}
    
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
        }
        
        lenArr = lenArr - 1 
        if (lenArr != 0) {
            returnMessage = returnMessage + ","
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
        
    //Add new username to master key list
    uList := []string{userName}
    t.updateMasterKeyList(stub, uList)
    
	return nil, nil

}

func (t *SimpleChaincode) getUserInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var userName, compID, returnMessage string
    var compStruct company
    var busPlanStruct businessPlan
    var userInfoObj userInfo
    
    if len (args) < 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2 (userName, compID)")
	}
    
    userName = args[0]
    compID = args[1]
    fmt.Println("Getting user info for: "+userName)		
        
    //Get Company details
    compInfo, _ := stub.GetState(compID)	
    _ = json.Unmarshal(compInfo, &compStruct)

    userInfoObj.UserID = userName
    userInfoObj.Company = compStruct
    fmt.Println(userInfoObj)

    //Get Business Plan info
    if compStruct.CompanyType == "Producer" || compStruct.CompanyType == "Transporter"  || compStruct.CompanyType == "Shipper" {
        bpInfo, _ := stub.GetState(compID + planIDAffix)	
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
}

func (t *SimpleChaincode) validateUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var userName, returnMessage string
    
	if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting 2 (userName and password).")
	}

    //Requires 2 arguments
	userName = args[0]
    //password = args[1]
    
    fmt.Println("Validating and getting user info: "+userName)
    
	validUser, err1, compID := t.verifyUser(stub, args)
	if err1 != nil {
        fmt.Println(err1)
    }
    fmt.Println(compID)
    
	if validUser == true {
        tArgs := []string{userName, compID}
        return t.getUserInfo(stub, tArgs)
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
	var compID string
    var topupDate int
	var topupAmount float64
	var companyObj company
    
    fmt.Println("Entered function topupBankBalance()")
    
    if len(args) < 3 {
        return nil, errors.New("Incorrect number of arguments. Expecting 3 arguments (CompanyID, top-up amount, top-up date).")
	}

	compID = args[0]
	topupAmount, _ = strconv.ParseFloat(args[1], 64)
    topupDate, _ = strconv.Atoi(args[2])
    
    //Get the company object from DB
    compObjBytes, _ := stub.GetState(compID)
    _ = json.Unmarshal(compObjBytes, &companyObj)
    fmt.Println(companyObj)
    
    //Topup the amount   
    companyObj.BankBalance = companyObj.BankBalance + topupAmount   
    companyObj.BalanceUpdatedDateMS = topupDate
        
    companyObjBytes, err := json.Marshal(&companyObj)
    if err != nil {
        fmt.Println("Failed to marshal company info.")
        fmt.Println(err)
    }
    err3 := stub.PutState(compID, companyObjBytes)
    if err3 != nil {
        fmt.Println(err3)
        return nil, errors.New("Failed to save Company info")
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

func (t *SimpleChaincode) getBusinessPlanList(stub shim.ChaincodeStubInterface) ([]byte, error) {
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
    
	var initiatorID, contractIDString, receiverID, contractStartDate, contractEndDate, contractStatus, entryLocation string
	var contractID int
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
    entryLocation = "Europe";
    
    if(len(args) == 7) { // Buyer adds location for gas request
        entryLocation = args[6];
    } 
    
	contractObj = contract{ContractID: contractID, InitiatorID: initiatorID, ReceiverID: receiverID,
                           EnergyMWH: energyMWH, EntryLocation: entryLocation, ContractStartDate: contractStartDate, ContractEndDate: contractEndDate, ContractStatus: contractStatus }

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
    
    //Add new contractID to master key list
    idList := []string{contractIDString, contractIDString + invoiceAffix, contractIDString + incidentAffix}
    t.updateMasterKeyList(stub, idList)
    
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
	
    fmt.Println("Updating contract ID " + args[0] + " to status " + args[1])
    
	contractIDString = args[0]
	contractObjBytes, _ := stub.GetState(contractIDString)
	err1 := json.Unmarshal(contractObjBytes, &contractObj)
	if err1 != nil {
        fmt.Println(err1)
		return nil, err1
	}
	
	//Update the status
	contractObj.ContractStatus = args[1]
	
	//Save the updated trade request
	contractBytes, err2 := json.Marshal(&contractObj)
	if err2 != nil {
        fmt.Println(err2)
		return nil, err2
	}
	err3 := stub.PutState(contractIDString, contractBytes)
	if err3 != nil {
        fmt.Println(err3)
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
	var companyID, contractIDStr, returnMessage, key string
	var lenMap int	
    var contractIDList []string
    var contractObj contract
    var contractFullObj contractInfo
    var businessPlanObj businessPlan
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
            
            //Add Business plan that is linked to this contract
            key = receiverCompany.CompanyID + planIDAffix
            planObjBytes, _ := stub.GetState(key)
            _ = json.Unmarshal(planObjBytes, &businessPlanObj)
            contractFullObj.BusinessPlan = businessPlanObj
            
            //Add invoices and incidents related to the contracts
            contractIDStr = strconv.Itoa(contractObj.ContractID)
            invoiceList, incidentList := t.getInvoiceIncidentList(stub, contractIDStr)
            
            contractFullObj.InvoiceList = invoiceList
            contractFullObj.IncidentList = incidentList
            
            fmt.Println(contractFullObj)
            
            contractFullObjBytes, err1 := json.Marshal(contractFullObj)
            if err1 != nil {
              return nil, err1
            }
            
            returnMessage = returnMessage + string(contractFullObjBytes)
        }
        
        lenMap = lenMap - 1
        if (lenMap!= 0) {
            returnMessage = returnMessage + ","
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

func (t *SimpleChaincode) getContractObjList(stub shim.ChaincodeStubInterface, idArrKey string, companyID string) ([]contract) {
    var contractIDList []string
    var contractObj contract
    var contractObjList []contract
    
    fmt.Println("Getting Contract Objects for company: "+ companyID)
	
	contractListObjBytes, _ := stub.GetState(idArrKey)
	_ = json.Unmarshal(contractListObjBytes, &contractIDList)
    
	for _, k := range contractIDList {
        
		contractObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(contractObjBytes, &contractObj)
        
        if(contractObj.InitiatorID == companyID || contractObj.ReceiverID == companyID) {
            if(contractObj.ContractStatus == "Accepted") {
                fmt.Println(contractObj)
                
                contractObjList = append(contractObjList, contractObj)
            }
        }            
	}
	
	return contractObjList
}

func (t *SimpleChaincode) getContractArrKey (stub shim.ChaincodeStubInterface, companyID string ) (string) {
    var companyType, companyTypeKey string
    var companyObj company
    
    compObjBytes, _ := stub.GetState(companyID)   
    _ = json.Unmarshal(compObjBytes, &companyObj)
    companyType = companyObj.CompanyType
    
    if(companyType == "Producer") {
        companyTypeKey = tradeRequestKey
    } else if(companyType == "Transporter") {
        companyTypeKey = transportRequestKey
    } else if(companyType == "Buyer") {
        companyTypeKey = gasRequestKey
    }
    
    return companyTypeKey
}

func (t *SimpleChaincode) addIOTData (stub shim.ChaincodeStubInterface, args[] string ) ([]byte, error) {
    //args[0] = {"device_id": "GasFlowMeter_1", "device_location": "Location 1", "company_id": "TRANSPORTER1", "pressure_kpa": 100, "temperature_c": 20, "specific_gravity": 0.65, "energy_mwh": 100,"timestamp_ms":1503416349302}
    
    fmt.Println("Adding new IOT Data: "+ args[0])
    
    var flowMeter flowMeterData
    var flowMeterList []flowMeterData
    var contractObjList []contract
    var contractArrKey string
    
    //Convert json string to json object
    _ = json.Unmarshal([]byte(args[0]), &flowMeter)
    fmt.Println(flowMeter)
    
    //Get the flow meter data list for this company
    var arrKey = flowMeter.CompanyID + iotKeyAffix
    flowMeterObjBytes, _ := stub.GetState(arrKey)   
	
    if flowMeterObjBytes != nil {
		 _ = json.Unmarshal(flowMeterObjBytes, &flowMeterList)
	}   
    
    //Add flow meter data to the list and save it
    flowMeterList = append(flowMeterList, flowMeter)	
    flowMeterObjBytes, _ = json.Marshal(&flowMeterList)
    _ = stub.PutState(arrKey, flowMeterObjBytes)
    
    fmt.Println(flowMeterList)
    
    //Get the key where the contracts are stored for the peer who owns IOT data
    contractArrKey = t.getContractArrKey(stub, flowMeter.CompanyID)
    
    //Check for invoice or incident to be created
    
    //Get all the contracts with this producer/transporter/buyer
    contractObjList = t.getContractObjList(stub, contractArrKey, flowMeter.CompanyID)
    
    for _, contractObj := range contractObjList {
        
        fmt.Println(contractObj)
        // If the energy from flow meter is higher or equal to the energy set in the contract, then create an invoice
        // Else create an incident
        if(flowMeter.EnergyMWH >= contractObj.EnergyMWH){
            //Create invoice
            t.createInvoice(stub, flowMeter.TimestampMS, contractObj.ContractID)
        } else {
            //Create incident
            t.createIncident(stub, flowMeter.TimestampMS, contractObj.EnergyMWH, flowMeter.EnergyMWH, contractObj.ContractID)
        }
    }
    return nil, nil
}

func (t *SimpleChaincode) createInvoice (stub shim.ChaincodeStubInterface, invoiceID int,  contractID int ) ([]byte, error) {
	var contractIDStr, invoiceIDStr, paymentStatus string 
    var invoiceDateMS int
    var invoiceObj invoice
    var invoiceIDArr []string
    
    fmt.Println("Creating new invoice...")
    
    invoiceIDStr = strconv.Itoa(invoiceID)
    invoiceDateMS = invoiceID
    paymentStatus = "Pending"
    contractIDStr = strconv.Itoa(contractID)
    
    //Create invoice and store in database
    invoiceObj = invoice {InvoiceID: invoiceID, InvoiceDateMS: invoiceDateMS, PaymentStatus: paymentStatus, ContractID: contractID}
    invoiceObjBytes, err1 := json.Marshal(invoiceObj)
	if err1 != nil {
		return nil, err1
	}
	_ = stub.PutState(invoiceIDStr, invoiceObjBytes)
    
    //Add new invoiceID to master key list
    idList := []string{invoiceIDStr}
    t.updateMasterKeyList(stub, idList)
    
    //Add invoice id into contract's invoice list
    var arrKey = contractIDStr+invoiceAffix
    invoiceIDListObjBytes, err2 := stub.GetState(arrKey)
	if err2 != nil {
		return nil, err2
	}
	if invoiceIDListObjBytes != nil {
		_ = json.Unmarshal(invoiceIDListObjBytes, &invoiceIDArr)
	}    
    invoiceIDArr = append(invoiceIDArr, invoiceIDStr)	
    invoiceIDListObjBytes, _ = json.Marshal(&invoiceIDArr)
    _ = stub.PutState(arrKey, invoiceIDListObjBytes)
    
    return nil, nil
}

func (t *SimpleChaincode) createIncident (stub shim.ChaincodeStubInterface, incidentID int,  expectedEnergyMWH float64, actualEnergyMWH float64, contractID int  ) ([]byte, error) {
	var contractIDStr, incidentIDStr, incidentStatus string 
    var incidentDateMS int
    var incidentObj incident
    var incidentIDArr []string
    
    fmt.Println("Creating new incident...")
    
    incidentIDStr = strconv.Itoa(incidentID)
    incidentDateMS = incidentID
    incidentStatus = "New"
    contractIDStr = strconv.Itoa(contractID)
    
    //Create incident and store in database
    incidentObj = incident {IncidentID: incidentID, IncidentDateMS: incidentDateMS, IncidentStatus: incidentStatus, ExpectedEnergyMWH: expectedEnergyMWH, ActualEnergyMWH: actualEnergyMWH, ContractID: contractID}
    incidentObjBytes, err1 := json.Marshal(incidentObj)
	if err1 != nil {
		return nil, err1
	}
	_ = stub.PutState(incidentIDStr, incidentObjBytes)
    
    //Add new incidentID to master key list
    idList := []string{incidentIDStr}
    t.updateMasterKeyList(stub, idList)
    
    //Add incident id into contract's incident list
    var arrKey = contractIDStr+incidentAffix
    incidentIDListObjBytes, err2 := stub.GetState(arrKey)
	if err2 != nil {
		return nil, err2
	}
	if incidentIDListObjBytes != nil {
		_ = json.Unmarshal(incidentIDListObjBytes, &incidentIDArr)
	}    
    incidentIDArr = append(incidentIDArr, incidentIDStr)	
    incidentIDListObjBytes, _ = json.Marshal(&incidentIDArr)
    _ = stub.PutState(arrKey, incidentIDListObjBytes)
    
    return nil, nil
}
                                                                                         
func (t *SimpleChaincode) getIOTData (stub shim.ChaincodeStubInterface, args[] string ) ([]byte, error) {
    fmt.Println("getIOTData for company ID: "+ args[0])
    var companyID, returnMessage string
	var flowMeterList []flowMeterData
    
    companyID = args[0]
    
    //Get the flow meter data list for this company
    var arrKey = companyID + iotKeyAffix
    flowMeterObjBytes, _ := stub.GetState(arrKey) 
    _ = json.Unmarshal(flowMeterObjBytes, &flowMeterList)
    
    returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : " + string(flowMeterObjBytes) + "}"
    
    fmt.Println(flowMeterList)
    
    return []byte(returnMessage), nil
}

func (t *SimpleChaincode) getIOTDataForShipper (stub shim.ChaincodeStubInterface, args[] string ) ([]byte, error) {
    fmt.Println("getIOTDataForShipper company ID: "+ args[0])
    var companyID, returnMessage, arrKey string
	var flowMeterList []flowMeterData
    var flowMeterFullList []flowMeterData
    var contractObjList []contract
    
    companyID = args[0]
    
    returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : "
    
    //Get the IOT data from producers
    contractObjList = t.getContractObjList(stub, tradeRequestKey, companyID)
    for _, contractObj := range contractObjList {
        //Get the flow meter data list for this company (Contract.Receiver ID)
        arrKey = contractObj.ReceiverID + iotKeyAffix
        flowMeterObjBytes, _ := stub.GetState(arrKey) 
        if flowMeterObjBytes != nil {
            _ = json.Unmarshal(flowMeterObjBytes, &flowMeterList)
            flowMeterFullList = append(flowMeterFullList, flowMeterList...)
        }
    }
    
    //Get the IOT data from transporters
    contractObjList = t.getContractObjList(stub, transportRequestKey, companyID)
    for _, contractObj := range contractObjList {
        //Get the flow meter data list for this company (Contract.Receiver ID)
        arrKey = contractObj.ReceiverID + iotKeyAffix
        flowMeterObjBytes, _ := stub.GetState(arrKey) 
        if flowMeterObjBytes != nil {
            _ = json.Unmarshal(flowMeterObjBytes, &flowMeterList)

            flowMeterFullList = append(flowMeterFullList, flowMeterList...)
        }
    }
    
    //Get the IOT data from buyers
    contractObjList = t.getContractObjList(stub, gasRequestKey, companyID)
    for _, contractObj := range contractObjList {
        //Get the flow meter data list for this company (Contract.Initiator ID)
        arrKey = contractObj.InitiatorID + iotKeyAffix
        flowMeterObjBytes, _ := stub.GetState(arrKey) 
        if flowMeterObjBytes != nil {
            _ = json.Unmarshal(flowMeterObjBytes, &flowMeterList)

            flowMeterFullList = append(flowMeterFullList, flowMeterList...)
        }
    }
    
    flowMeterFullListBytes, _ := json.Marshal(flowMeterFullList)
    returnMessage = returnMessage + string(flowMeterFullListBytes)
    returnMessage = returnMessage + "}"
    
    fmt.Println(flowMeterFullList)
    fmt.Println(returnMessage)
    
    return []byte(returnMessage), nil
}
    
func (t *SimpleChaincode) getInvoiceList (stub shim.ChaincodeStubInterface, args[] string ) ([]byte, error) {
    fmt.Println("Getting list of invoices for company ID: " + args[0])
    var lenArr int
	var invoiceIDArr []string
	var contractID, returnMessage string
    var invoiceObj invoice
    
    if len(args) < 1 {
        return nil, errors.New("Incorrect number of arguments. 1 expected (Contract ID)")
	}
    
    contractID = args[0]
        
    var arrKey = contractID + invoiceAffix
    invoiceIDArrBytes, _ := stub.GetState(arrKey)
	_ = json.Unmarshal(invoiceIDArrBytes, &invoiceIDArr)
    
    fmt.Println(invoiceIDArr)
    
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["
	lenArr = len(invoiceIDArr)
	for _, k := range invoiceIDArr {
		invoiceObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(invoiceObjBytes, &invoiceObj)
        fmt.Println(invoiceObj)
        
        returnMessage = returnMessage + string(invoiceObjBytes) 
        
        lenArr = lenArr - 1 
        if (lenArr != 0) {
            returnMessage = returnMessage + ","
        } 
	} 
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil
} 

func (t *SimpleChaincode) getIncidentList (stub shim.ChaincodeStubInterface, args[] string ) ([]byte, error) {
    fmt.Println("Getting list of incidents for company ID: " + args[0])
    var lenArr int
	var incidentIDArr []string
	var contractID, returnMessage string
    var incidentObj incident
    
    
    if len(args) < 1 {
        return nil, errors.New("Incorrect number of arguments. 1 expected (Contract ID)")
	}
    
    contractID = args[0]
        
    var arrKey = contractID + incidentAffix
    incidentIDArrBytes, _ := stub.GetState(arrKey)
	_ = json.Unmarshal(incidentIDArrBytes, &incidentIDArr)
    
    fmt.Println(incidentIDArr)
    
	returnMessage = "{\"statusCode\" : \"SUCCESS\", \"body\" : ["
	lenArr = len(incidentIDArr)
	for _, k := range incidentIDArr {
		incidentObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(incidentObjBytes, &incidentObj)
        fmt.Println(incidentObj)
        
        returnMessage = returnMessage + string(incidentObjBytes) 
        
        lenArr = lenArr - 1 
        if (lenArr != 0) {
            returnMessage = returnMessage + ","
        } 
	} 
	returnMessage = returnMessage + "]}"
	return []byte(returnMessage), nil
} 

func (t *SimpleChaincode) getInvoiceIncidentList(stub shim.ChaincodeStubInterface, contractID string) ([]invoice, []incident) {
    var invoiceIDList, incidentIDList []string
    var invoiceObj invoice
    var invoiceObjList []invoice
    var incidentObj incident
    var incidentObjList []incident
    
    fmt.Println("Getting Invoice and Incident Objects for contract: "+ contractID)
	
    var idArrKey = contractID + invoiceAffix
	invoiceIDListObjBytes, _ := stub.GetState(idArrKey)
	_ = json.Unmarshal(invoiceIDListObjBytes, &invoiceIDList)
    
	for _, k := range invoiceIDList {
        
		invoiceObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(invoiceObjBytes, &invoiceObj)
        
        invoiceObjList = append(invoiceObjList, invoiceObj)
	}
	
    fmt.Println(invoiceObjList)
    
    idArrKey = contractID + incidentAffix
	incidentIDListObjBytes, _ := stub.GetState(idArrKey)
	_ = json.Unmarshal(incidentIDListObjBytes, &incidentIDList)
    
	for _, k := range incidentIDList {
        
		incidentObjBytes, _ := stub.GetState(k)
        _ = json.Unmarshal(incidentObjBytes, &incidentObj)
        
        incidentObjList = append(incidentObjList, incidentObj)
	}
    
    fmt.Println(incidentObjList)
    
	return invoiceObjList, incidentObjList
}

func (t *SimpleChaincode) makePayment (stub shim.ChaincodeStubInterface, args[] string ) ([]byte, error) {
    var returnMessage, invoiceIDStr, contractIDStr, planIDKey, totalCostStr, bankBalStr string
    var contractObj contract
    var planObj businessPlan
    var totalCost float64
    var initiatorCompany, receiverCompany company
    var invoiceObj invoice
    var currentDate int
    
    fmt.Println("Pay for the contract (Invoice ID: "+ args[0] + ")")
    if len(args) < 3 {
        return nil, errors.New("Incorrect number of arguments. 3 expected (Invoice ID, Contract ID, Current Date in MilliSecs)")
	}
    
    invoiceIDStr = args[0]
    contractIDStr = args[1]
    currentDate, _ = strconv.Atoi(args[2])
    
    contractObjBytes, _ := stub.GetState(contractIDStr)
    _ = json.Unmarshal(contractObjBytes, &contractObj)
        
    //Fetch gas price from the Business Plan
    planIDKey = contractObj.ReceiverID + planIDAffix 
    planObjBytes, _ := stub.GetState(planIDKey)
    _ = json.Unmarshal(planObjBytes, &planObj)
    
    //Energy consumed * gas price per mwh
    totalCost = contractObj.EnergyMWH * planObj.GasPrice
    
    //Fetch Initiator company
    initiatorCompanyObjBytes, _ := stub.GetState(contractObj.InitiatorID)
    _ = json.Unmarshal(initiatorCompanyObjBytes, &initiatorCompany)
    
    //Subtract amount from initiator company
    if (initiatorCompany.BankBalance < totalCost) {
        totalCostStr = strconv.FormatFloat(totalCost, 'E', -1, 64)
        bankBalStr = strconv.FormatFloat(initiatorCompany.BankBalance, 'E', -1, 64)
        returnMessage = "{\"statusCode\" : \"FAIL\", \"body\" : \"Transaction FAILED: Insufficient funds (Bank Balance: "+ bankBalStr +", Invoice payment amount: "+totalCostStr+")\"}"
        
        return []byte(returnMessage), nil
    } else {
        initiatorCompany.BankBalance = initiatorCompany.BankBalance - totalCost
        initiatorCompany.BalanceUpdatedDateMS = currentDate
        initiatorCompanyObjBytes, _ = json.Marshal(&initiatorCompany)
        _ = stub.PutState(initiatorCompany.CompanyID, initiatorCompanyObjBytes)
        
        fmt.Println(initiatorCompany)
    }
    
    
    //Add the amount to Receiver company
    receiverCompanyObjBytes, _ := stub.GetState(contractObj.ReceiverID)
    _ = json.Unmarshal(receiverCompanyObjBytes, &receiverCompany)
    
    receiverCompany.BankBalance = receiverCompany.BankBalance + totalCost
    receiverCompany.BalanceUpdatedDateMS = currentDate
    
    receiverCompanyObjBytes, _ = json.Marshal(&receiverCompany)
    _ = stub.PutState(receiverCompany.CompanyID, receiverCompanyObjBytes)
    
    fmt.Println(receiverCompany)
    
    
    //Update the invoice payment status and date
    invoiceObjBytes, _ := stub.GetState(invoiceIDStr)
    _ = json.Unmarshal(invoiceObjBytes, &invoiceObj)
    
    invoiceObj.PaymentDateMS = currentDate
    invoiceObj.PaymentStatus = "Paid"
    
    invoiceObjBytes, _ = json.Marshal(&invoiceObj)
    _ = stub.PutState(invoiceIDStr, invoiceObjBytes)
    
    fmt.Println(invoiceObj)
    
    return nil, nil
}

func (t *SimpleChaincode) Reset(stub shim.ChaincodeStubInterface) ([]byte, error) {
    var masterKeyList []string
    
    //Get the existing master array of keys
    keyListBytes, _ := stub.GetState(allKeys)
    if keyListBytes != nil
	   _ = json.Unmarshal(keyListBytes, &masterKeyList)
    
	for _, key := range masterKeyList {
        fmt.Println("Deleting data with key: " + key)
        err := stub.DelState(key)
        if err != nil {
            fmt.Println(err)
        }
    }
    
    t.Init(stub)

	return nil, nil
}
    
func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Running Invoke function")

	if function == "init" {
		return t.Init(stub)
	} else if function == "delete" {
		return t.deleteData(stub, args)
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
	} else if function == "makePayment" {
		return t.makePayment(stub, args)
	} else if function == "reset" {
		return t.Reset(stub)
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
	} else if function == "validateUser" {
		return t.validateUser(stub, args)
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
		return t.getBusinessPlanList(stub)
    } else if function == "getIOTData" {
		return t.getIOTData(stub, args)
    } else if function == "getInvoiceList" {
		return t.getInvoiceList(stub, args)
    } else if function == "getIncidentList" {
		return t.getIncidentList(stub, args)
    } else if function == "getIOTDataForShipper" {
		return t.getIOTDataForShipper(stub, args)
    } else if function == "getMasterKeyList" {
		return t.getMasterKeyList(stub)
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

func (t *SimpleChaincode) deleteData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key string
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
    
    key = args[0]
    
    fmt.Println("Deleting data with key:" + key)

	// Delete the key from the state in ledger
	err := stub.DelState(key)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}
