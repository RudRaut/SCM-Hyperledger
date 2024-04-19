package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	protos "github.com/hyperledger/fabric-protos-go/peer"
)

//  ---------------------------- data ------------------------------------------

type SupplyChain struct {
	contractapi.SystemContract
}

type CounterNO struct {
	Counter int `json:"Conuter"`
}

type User struct {
	Name     string `json:"Name"`
	UserID   string `json:"UserID"`
	UserType string `json:"UserType"`
	Email    string `json:"Email"`
	Address  string `json:"Address"`
	Password string `json:"Password"`
}

type ProductDates struct {

	// Try Creating date struct {DD:MM:YYYY} and use that instead of string (later)

	ManufactureDate string `json:"ManufactureDate"` // Date of Manufacturing the product
	SupplyDate      string `json:"SupplyDate"`      // Supplier getting products
	OrderDate       string `json:"OrderDate"`       // Order placed by customer
	TransportDate   string `json:"TransportDate"`   // Product sent for transporting
	SoldDate        string `json:"SoldDate"`
	DeliveryDate    string `json:"DeliveryDate"` //  Product delivered to customer
}

type Product struct {
	// Product Data
	ProductID      string       `json:"ProductID"`
	OrderID        string       `json:"OrderID"`
	Name           string       `json:"Name"`
	CustomerID     string       `json:"CustomerID"`
	ManufacturerID string       `json:"ManufacturerID"`
	SupplierID     string       `json:"SupplierID"`
	TransporterID  string       `json:"TransporterID"`
	Status         string       `json:"Status"`
	Price          float64      `json:"Price"`
	Date           ProductDates `json:"Date"`
}

//  ---------------------------- main ------------------------------------------

func main() {
	err := shim.Start(new(SupplyChain))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SupplyChain) Init(APIstub shim.ChaincodeStubInterface) protos.Response {
	// Initializing Product Counter
	ProductCounterBytes, _ := APIstub.GetState("ProductCounterNO")
	if ProductCounterBytes == nil { // No value of ProductCounterbytes found, initialize a new one
		var ProductCounter = CounterNO{Counter: 0}
		ProductCounterBytes, _ := json.Marshal(ProductCounter)
		err := APIstub.PutState("ProductCounterNO", ProductCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to initiate Product Counter"))
		}
	}

	// Initializing Order Counter
	OrderCounterBytes, _ := APIstub.GetState("OrderCounterNO")
	if OrderCounterBytes == nil { // No value of ProductCounterbytes found, initialize a new one
		var OrderCounter = CounterNO{Counter: 0}
		OrderCounterBytes, _ := json.Marshal(OrderCounter)
		err := APIstub.PutState("OrderCounterNO", OrderCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to initiate Product Counter"))
		}
	}

	// Initializing User Counter
	UserCounterBytes, _ := APIstub.GetState("ProductCounterNO")
	if UserCounterBytes == nil { // No value of ProductCounterbytes found, initialize a new one
		var UserCounter = CounterNO{Counter: 0}
		UserCounterBytes, _ := json.Marshal(UserCounter)
		err := APIstub.PutState("UserCounterNO", UserCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to initiate Product Counter"))
		}
	}

	return shim.Success(nil) // Successful completion of function
}

// Invoke - All functions will be invoked by this function
func (t *SupplyChain) Invoke(APIstub shim.ChaincodeStubInterface) protos.Response {
	function, args := APIstub.GetFunctionAndParameters()
	fmt.Println("invoke is running" + function)

	// Handling different functionsd
	switch function {
	case "initLedger":
		//initLedger function
		return t.initLedger(APIstub, args)
	case "signIn":
		return t.signIn(APIstub, args)
	case "createUser":
		return t.createUser(APIstub, args)
	case "createProduct":
		return t.createProduct(APIstub, args)
	case "updateProduct":
		return t.updateProduct(APIstub, args)
	case "orderProduct":
		return t.orderProduct(APIstub, args)
	case "productDelivered":
		return t.productDelivered(APIstub, args)
	case "toSupplier":
		return t.toSupplier(APIstub, args)
	case "toTransporter":
		return t.toTransporter(APIstub, args)
	case "sellToCustomer":
		return t.sellToCustomer(APIstub, args)
	// case "queryAsset":
	// 	return t.queryAsset(APIstub, args);
	case "queryAll":
		return t.queryAll(APIstub, args)
	default:
		fmt.Println("Invalid function name:", function) // Use fmt.Println instead of console.error
		// return nil // Or throw an error
	}

	return shim.Success(nil)
}

func getCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	fmt.Sprintf("Counter Current Value %d of  Asset Type %s  ", counterAsset.Counter, AssetType)

	return counterAsset.Counter
}

func incrementCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.Counter++
	counterAsBytes, _ = json.Marshal(counterAsset)

	err := APIstub.PutState(AssetType, counterAsBytes)
	if err != nil {
		fmt.Sprintf("Failed to Increment Counter")
	}

	fmt.Println("Increment Successful for %v", counterAsset)

	return counterAsset.Counter
}

// Get the TimeStamp of transaction when chaicode was executed
func (t *SupplyChain) GetTxTime(APIstub shim.ChaincodeStubInterface) (string, error) {
	txTimeAsPtr, err := APIstub.GetTxTimestamp()
	if err != nil {
		fmt.Printf("Error returning TimeStamp \n")
		return "Error", err
	}
	fmt.Printf("\t returned value from APIstub: %v \n", txTimeAsPtr)
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()

	return timeStr, nil
}

func (t *SupplyChain) initLedger(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	entityUser := User{Name: "admin", UserID: "admin", Email: "email@scm.com", UserType: "admin", Address: "Pune", Password: "adminpw"}
	entityUserAsBytes, errMarshal := json.Marshal(entityUser)
	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in User %s", errMarshal))
	}

	errPut := APIstub.PutState(entityUser.UserID, entityUserAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create entity Asset: %s", entityUser.UserID))
	}

	fmt.Println("Added User", entityUser)
	return shim.Success(nil)
}

func (t *SupplyChain) signIn(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 2 {
		return shim.Error("Insufficeint Arguments, Expected 2")
	}

	if len(args[0]) == 0 {
		return shim.Error("User ID must be provided")
	}

	if len(args[1]) == 0 {
		return shim.Error("Password must be provided")
	}

	entityUserBytes, _ := APIstub.GetState(args[0])
	if entityUserBytes == nil {
		return shim.Error("Cannot Find Entity")
	}

	entityUser := User{}
	json.Unmarshal(entityUserBytes, &entityUser)

	if entityUser.Password != args[1] {
		return shim.Error("Either ID or password is wrong")
	}

	return shim.Success(entityUserBytes)
}

func (t *SupplyChain) createUser(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 5 {
		return shim.Error("Insufficeint Arguments, required 5")
	}

	if len(args[0]) == 0 {
		return shim.Error("Please provide a name for the User")
	}

	if len(args[1]) == 0 {
		return shim.Error("Please provide Email to create User")
	}

	if len(args[2]) == 0 {
		return shim.Error("Please specify Type of User")
	}

	if len(args[3]) == 0 {
		return shim.Error("Please provide non-Empty address")
	}

	if len(args[4]) == 0 {
		return shim.Error("please enter valid non-empty Password")
	}

	userCounter := getCounter(APIstub, "UserCounterNO")
	userCounter++

	var comAsset = User{Name: args[0], UserID: "User" + strconv.Itoa(userCounter), Email: args[1], UserType: args[2], Address: args[3], Password: args[4]}

	comAssetBytes, errMarshal := json.Marshal(comAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Failed to Register User: %s", comAsset.UserID))
	}

	errPut := APIstub.PutState(comAsset.UserID, comAssetBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to Register User: %s", comAsset.UserID))
	}

	incrementCounter(APIstub, "UserCounterNO")
	fmt.Println("User Registered Successfully %v", comAsset)

	return shim.Success(comAssetBytes)
}

func (t *SupplyChain) createProduct(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 3 {
		return shim.Error("Insufficient Arguments, required 3")
	}

	if len(args[0]) == 0 {
		return shim.Error("Please enter Product name to register")
	}

	if len(args[1]) == 0 {
		return shim.Error("UserID must be provided")
	}

	if len(args[2]) == 0 {
		return shim.Error("Please specify price for the Product")
	}

	// User details from blockchain using UserID
	userBytes, _ := APIstub.GetState(args[1])
	if userBytes == nil {
		return shim.Error("Can not find the User")
	}

	user := User{}
	json.Unmarshal(userBytes, &user)

	if user.UserType != "Manufacturer" || user.UserType != "manufacturer" {
		return shim.Error("You must be a manufacturer to CreateProduct")
	}

	// Price Conversion
	p1, errPrice := strconv.ParseFloat(args[2], 64)
	if errPrice != nil {
		return shim.Error(fmt.Sprintf("Failed to convert Price %s", errPrice))
	}

	productCounter := getCounter(APIstub, "ProductCounterNO")
	productCounter++

	//Creating transaction TimStamp
	txTimeAsPtr, errTx := t.GetTxTime(APIstub)
	if errTx != nil {
		return shim.Error(fmt.Sprintf("Error in TimeStamp"))
	}

	dates := ProductDates{}

	dates.ManufactureDate = txTimeAsPtr

	var comAsset = Product{ProductID: "Product" + strconv.Itoa(productCounter), OrderID: "", Name: args[0], CustomerID: "", ManufacturerID: args[1], SupplierID: "", TransporterID: "", Status: "Available", Date: dates, Price: p1}
	comAssetAsBytes, errMarshal := json.Marshal(comAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Product: %s", errMarshal))
	}

	errPut := APIstub.PutState(comAsset.ProductID, comAssetAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to Create Product: %v", comAsset))
	}
	incrementCounter(APIstub, "ProductCounterNO")

	fmt.Println("Successfully created Product: %v", comAsset)

	return shim.Success(comAssetAsBytes)
}

func (t *SupplyChain) updateProduct(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 4 {
		return shim.Error("Insufficient Arguments, required 4")
	}

	if len(args[0]) == 0 {
		return shim.Error("Provide ProductID")
	}

	if len(args[1]) == 0 {
		return shim.Error("Provide UserID")
	}

	if len(args[2]) == 0 {
		return shim.Error("Provide Product Name")
	}

	if len(args[3]) == 0 {
		return shim.Error("Provide updated proice for the Product")
	}

	userBytes, _ := APIstub.GetState(args[1])
	if userBytes == nil {
		shim.Error("Can not find the User")
	}

	user := User{}
	json.Unmarshal(userBytes, &user)

	if user.UserType == "Customer" || user.UserType == "customer" {
		return shim.Error("Customer cannot update Product")
	}

	productBytes, _ := APIstub.GetState(args[0])
	if productBytes == nil {
		return shim.Error("Can not find Product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	// Price Conversion
	p1, errPrice := strconv.ParseFloat(args[3], 64)
	if errPrice != nil {
		return shim.Error(fmt.Sprintf("Failed to convert Price: %s", errPrice))
	}
	product.Name = args[2]
	product.Price = p1

	updateProductAsBytes, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error: %s", errMarshal))
	}

	errPut := APIstub.PutState(product.ProductID, updateProductAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to sell to Customer: %s", product.ProductID))
	}

	fmt.Println("Successfully updated Product: %v", product.ProductID)
	return shim.Success(updateProductAsBytes)
}

func (t *SupplyChain) orderProduct(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 2 {
		return shim.Error("Insufficient arguments, required 2")
	}

	if len(args[0]) == 0 {
		return shim.Error("Please provide CustomerID")
	}

	if len(args[1]) == 0 {
		return shim.Error("Please provide ProductID")
	}
	userBytes, _ := APIstub.GetState(args[0])

	if userBytes == nil {
		return shim.Error("Can not find User")
	}

	user := User{}
	json.Unmarshal(userBytes, &user)
	if user.UserType != "Customer" || user.UserType != "customer" {
		return shim.Error("Only Customer can order product")
	}

	productBytes, _ := APIstub.GetState((args[1]))
	if productBytes == nil {
		return shim.Error("Can not find the Product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	orderCounter := getCounter(APIstub, "OrderCounterNO")
	orderCounter++

	// Transaction TimeStamp
	txTimeAsPtr, errTx := t.GetTxTime(APIstub)
	if errTx != nil {
		return shim.Error(fmt.Sprintf("Error in TimeStamp"))
	}

	product.OrderID = "Order" + strconv.Itoa(orderCounter)
	product.CustomerID = user.UserID
	product.Status = "Ordered"
	product.Date.OrderDate = txTimeAsPtr
	updateProductAsBytes, errMrashal := json.Marshal(product)
	if errMrashal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error: %s", errMrashal))
	}

	incrementCounter(APIstub, "OrderCounterNO")

	errPut := APIstub.PutState(product.ProductID, updateProductAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to place the Order: %s", product.ProductID))
	}
	fmt.Println("Order placed successfully: %v", product.ProductID)
	return shim.Success(updateProductAsBytes)
}

func (t *SupplyChain) productDelivered(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 1 {
		return shim.Error("Insufficient Arguments, required 1")
	}

	if len(args[0]) == 0 {
		return shim.Error("Provide ProductID")
	}

	productBytes, _ := APIstub.GetState((args[0]))
	if productBytes == nil {
		return shim.Error("Can not find the Product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	if product.Status != "sold" || product.Status != "Sold" {
		return shim.Error("Product is not delivered yet")
	}

	// Trnasaction Timestamp
	txTimeAsPtr, errTx := t.GetTxTime(APIstub)
	if errTx != nil {
		return shim.Error(fmt.Sprintf("Error in TimeStamp"))
	}

	product.Date.DeliveryDate = txTimeAsPtr
	product.Status = "Delivered"
	updateProductAsBytes, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Product: %s", errMarshal))
	}

	errPut := APIstub.PutState(product.ProductID, updateProductAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to update that Product is delivered: %s", product.ProductID))
	}

	fmt.Println("Successfully delivered Product: %v", product.ProductID)
	return shim.Success(updateProductAsBytes)

}

func (t *SupplyChain) toSupplier(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 2 {
		return shim.Error("Insufficient Arguments, required 2")
	}

	if len(args[0]) == 0 {
		return shim.Error("Please provide ProductID")
	}

	if len(args[1]) == 0 {
		return shim.Error("Please provide UserID")
	}

	userBytes, _ := APIstub.GetState(args[1])

	if userBytes == nil {
		return shim.Error("Can not find Supplier")
	}

	user := User{}
	json.Unmarshal(userBytes, &user)

	if user.UserType != "Supplier" || user.UserType != "Supplier" {
		return shim.Error("User must be a Supplier")
	}

	productBytes, _ := APIstub.GetState((args[0]))
	if productBytes == nil {
		return shim.Error("Can not find the Product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	if product.SupplierID != "" {
		return shim.Error("Product is sent to Supplier already")
	}

	// Trnasaction Timestamp
	txTimeAsPtr, errTx := t.GetTxTime(APIstub)
	if errTx != nil {
		return shim.Error(fmt.Sprintf("Error in TimeStamp"))
	}

	product.SupplierID = user.UserID
	product.Date.SupplyDate = txTimeAsPtr
	updateProductAsBytes, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error: %s", errMarshal))
	}

	errPut := APIstub.PutState(product.ProductID, updateProductAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to send to Supplier: %s", product.ProductID))
	}

	fmt.Println("Successfully sent Product for supply: %v", product.ProductID)
	return shim.Success(updateProductAsBytes)

}

func (t *SupplyChain) toTransporter(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 2 {
		return shim.Error("Insufficient Arguments, required 2")
	}

	if len(args[0]) == 0 {
		return shim.Error("Please provide ProductID")
	}

	if len(args[1]) == 0 {
		return shim.Error("Please provide UserID")
	}

	userBytes, _ := APIstub.GetState(args[1])

	if userBytes == nil {
		return shim.Error("Can not find Transporter")
	}

	user := User{}
	json.Unmarshal(userBytes, &user)

	if user.UserType != "Transporter" || user.UserType != "transporter" {
		return shim.Error("User must be a Transporter")
	}

	productBytes, _ := APIstub.GetState((args[0]))
	if productBytes == nil {
		return shim.Error("Can not find the Product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	if product.TransporterID != "" {
		return shim.Error("Product is sent to Transporter already")
	}

	// Trnasaction Timestamp
	txTimeAsPtr, errTx := t.GetTxTime(APIstub)
	if errTx != nil {
		return shim.Error(fmt.Sprintf("Error in TimeStamp"))
	}

	product.TransporterID = user.UserID
	product.Date.TransportDate = txTimeAsPtr

	updateProductAsBytes, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in User %s", errMarshal))
	}

	errPut := APIstub.PutState(product.ProductID, updateProductAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to send to Transporter: %s", product.ProductID))
	}

	fmt.Println("Product successfully sent for Transporting")
	return shim.Success(updateProductAsBytes)

}

func (t *SupplyChain) sellToCustomer(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 1 {
		return shim.Error("Insufficient Arguments, required 1")
	}

	if len(args[0]) == 0 {
		return shim.Error("Please provide ProductID")
	}

	productBytes, _ := APIstub.GetState((args[0]))
	if productBytes == nil {
		return shim.Error("Can not find the Product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	if product.OrderID == "" {
		return shim.Error("Product is not ordered yet")
	}

	if product.CustomerID == "" {
		return shim.Error("CustomerID should be set to sell to customers")
	}

	// Transaction Timestamp
	txTimeAsPtr, errTx := t.GetTxTime(APIstub)
	if errTx != nil {
		return shim.Error(fmt.Sprintf("Error in TimeStamp"))
	}

	product.Date.SoldDate = txTimeAsPtr
	product.Status = "Sold"

	updateProductAsBytes, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in User %s", errMarshal))
	}

	errPut := APIstub.PutState(product.ProductID, updateProductAsBytes)
	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to update that Product is delivered: %s", product.ProductID))
	}

	fmt.Println("Successfully sold Product to Customer: %v", product.ProductID)
	return shim.Success(updateProductAsBytes)
}

func (t *SupplyChain) queryAsset(APIstub shim.ChaincodeStub, args []string) protos.Response {
	if len(args) != 1 {
		return shim.Error("Insufficient Arguments, required 1")
	}

	productAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(productAsBytes)
}

func (t *SupplyChain) queryAll(APIstub shim.ChaincodeStubInterface, args []string) protos.Response {
	if len(args) != 1 {
		return shim.Error("Insufficient Arguments, required 1")
	}

	if len(args[0]) == 0 {
		return shim.Error("Please provide Asset Type")
	}

	assetType := args[0]
	assetCounter := getCounter(APIstub, assetType+"CounterNO")

	startKey := assetType + "1"
	endKey := assetType + strconv.Itoa(assetCounter+1)

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	fmt.Println("- queryAllAssets:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())

}
