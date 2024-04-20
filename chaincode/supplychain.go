package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//  ---------------------------- data ------------------------------------------

type SupplyChain struct {
	contractapi.SystemContract
}

type CounterNO struct {
	Counter int `json:"Counter"`
}

type User struct {
	Name     string `json:"Name"`
	UserID   string `json:"UserID"`
	UserType string `json:"UserType"`
	Email    string `json:"Email"`
	Address  string `json:"Address"`
	Password string `json:"Password"`
}

type UserInfo struct {
	Name     string `json:"Name"`
	UserID   string `json:"UserID"`
	UserType string `json:"UserType"`
	Email    string `json:"Email"`
	Address  string `json:"Address"`
}

type ProductPos struct {
	Date      string `json:"Date"`
	Latitude  string `json:"Latitude"`
	Longitude string `json:"Longitude"`
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
	Position       []ProductPos `json:"Position"`
}

func (t *SupplyChain) Invoke(ctx contractapi.TransactionContextInterface) error {
	function, args := ctx.GetStub().GetFunctionAndParameters()
	fmt.Println("invoke is running" + function)

	switch function {
	case "InitLedger":
		return t.InitLedger(ctx)
	case "signIn":
		return t.signIn(ctx, args)
	case "createUser":
		return t.createUser(ctx, args)
	case "createProduct":
		if len(args) != 5 {
			return fmt.Errorf("insufficient arguments expected 5 for create-product")
		}
		name, userID, longitude, latitude, price := args[0], args[1], args[2], args[3], args[4]
		return t.createProduct(ctx, name, userID, longitude, latitude, price)
	case "updateProduct":
		if len(args) != 4 {
			return fmt.Errorf("insufficient arguments expected 4 for update-product")
		}
		userID, productID, name, price := args[0], args[1], args[2], args[3]
		return t.updateProduct(ctx, userID, productID, name, price)
	case "toSupplier":
		if len(args) != 4 {
			return fmt.Errorf("insufficient arguments expected 4 for to-supplier")
		}
		productID, supplierID, longitude, latitude := args[0], args[1], args[2], args[3]
		return t.toSupplier(ctx, productID, supplierID, longitude, latitude)
	case "toTransporter":
		if len(args) != 4 {
			return fmt.Errorf("insufficient arguments expected 4 for to-transporter")
		}
		productID, transporterID, longitude, latitude := args[0], args[1], args[2], args[3]
		return t.toTransporter(ctx, productID, transporterID, longitude, latitude)
	case "sellToCustomer":
		if len(args) != 3 {
			return fmt.Errorf("insufficient arguments expected 3 for sell-to-customer")
		}
		productID, customerID, latitude, longitude := args[0], args[1], args[2], args[3]
		return t.sellToCustomer(ctx, productID, customerID, latitude, longitude)
	// Add more functions here...
	default:
		return fmt.Errorf("invalid function name: %s", function)
	}
}

// //  ---------------------------- functions ------------------------------------------

func getCounter(ctx contractapi.TransactionContextInterface, AssetType string) int {
	counterAsBytes, _ := ctx.GetStub().GetState(AssetType)
	counter := CounterNO{}

	json.Unmarshal(counterAsBytes, &counter)
	fmt.Printf("Counter Current Value %d of  Asset Type %s  ", counter.Counter, AssetType)

	return counter.Counter
}

func incrementCounter(ctx contractapi.TransactionContextInterface, AssetType string) int {
	counterAsBytes, _ := ctx.GetStub().GetState(AssetType)
	counter := CounterNO{}

	json.Unmarshal(counterAsBytes, &counter)
	counter.Counter++
	counterAsBytes, _ = json.Marshal(counter)

	err := ctx.GetStub().PutState(AssetType, counterAsBytes)
	if err != nil {
		fmt.Printf("Failed to Increment Counter")
	}
	return counter.Counter
}

// Get the TimeStamp of transaction when chaicode was executed
func (t *SupplyChain) GetTxTimestamp(ctx contractapi.TransactionContextInterface) (string, error) {
	txTimeAsPtr, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "Error", err
	}
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()
	return timeStr, nil
}

func (t *SupplyChain) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Init Manufacturer admin
	userManufacturer := User{
		Name:     "Manufacturer_Admin",
		UserID:   "manufacturer-admin",
		UserType: "admin",
		Email:    "mfg.admin@scm.com",
		Address:  "fabric",
		Password: "admin@123",
	}

	userManufacturerBytes, err := json.Marshal(userManufacturer)
	if err != nil {
		return fmt.Errorf("marshal error: %s", err.Error())
	}

	err = ctx.GetStub().PutState(userManufacturer.UserID, userManufacturerBytes)
	if err != nil {
		return fmt.Errorf("failed to put manufacturer to world state: %s", err.Error())
	}

	// Init Supplier admin
	userSupplier := User{
		Name:     "Supplier_Admin",
		UserID:   "supplier-admin",
		UserType: "admin",
		Email:    "supplier.admin@scm.com",
		Address:  "fabric",
		Password: "admin@123",
	}

	userSupplierBytes, err := json.Marshal(userSupplier)
	if err != nil {
		return fmt.Errorf("marshal error: %s", err.Error())
	}

	err = ctx.GetStub().PutState(userSupplier.UserID, userSupplierBytes)
	if err != nil {
		return fmt.Errorf("failed to put supplier to world state: %s", err.Error())
	}

	// Init Transporter admin
	userTransporter := User{
		Name:     "Transporter_Admin",
		UserID:   "transporter-admin",
		UserType: "admin",
		Email:    "transporter.admin@scm.com",
		Address:  "fabric",
		Password: "admin@123",
	}

	userTransporterBytes, err := json.Marshal(userTransporter)
	if err != nil {
		return fmt.Errorf("marshal error: %s", err.Error())
	}

	err = ctx.GetStub().PutState(userTransporter.UserID, userTransporterBytes)
	if err != nil {
		return fmt.Errorf("failed to put supplier to world state: %s", err.Error())
	}

	return nil

}

func (t *SupplyChain) signIn(ctx contractapi.TransactionContextInterface, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("insufficient arguments. expected 2")
	}

	if len(args[0]) == 0 {
		return fmt.Errorf("user id must be provided")
	}
	if len(args[1]) == 0 {
		return fmt.Errorf("password must be provided")
	}

	userID := args[0]
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("error retrieving user data: %w", err)
	}
	if userBytes == nil {
		return fmt.Errorf("User not found: %s", userID)
	}

	user := User{}
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return fmt.Errorf("unmarshalling error: %w", err)
	}

	// Verify password
	if user.Password != args[1] {
		return fmt.Errorf("incorred Userid or passsword")
	}

	// No data returned, only error handling (success implied by lack of error)
	return nil
}

func (t *SupplyChain) createUser(ctx contractapi.TransactionContextInterface, args []string) error {
	if len(args) != 5 {
		return fmt.Errorf("insufficient arguments, expected 5")
	}

	if len(args[0]) == 0 {
		return fmt.Errorf("provide name for User")
	}

	if len(args[1]) == 0 {
		return fmt.Errorf("provide Email")
	}

	if len(args[2]) == 0 {
		return fmt.Errorf("please specify type of uer")
	}

	if len(args[3]) == 0 {
		return fmt.Errorf("please provide non-empty address")
	}

	if len(args[4]) == 0 {
		return fmt.Errorf("please enter valid non-empty password")
	}

	userCounter := getCounter(ctx, "UserCounterNO")
	userCounter++

	user := User{
		Name:     args[0],
		UserID:   "User" + strconv.Itoa(userCounter),
		Email:    args[1],
		UserType: args[2],
		Address:  args[3],
		Password: args[4],
	}

	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal error: %s", err.Error())
	}

	errPut := ctx.GetStub().PutState(user.UserID, userAsBytes)
	if errPut != nil {
		return fmt.Errorf("error storing date: %w", errPut)
	}

	incrementCounter(ctx, "UserCounterNO")
	fmt.Println("Successfully created user")

	return nil

}

func (t *SupplyChain) createProduct(ctx contractapi.TransactionContextInterface, name string, userId string, longitude string, latitude string, price string) error {

	userBytes, _ := ctx.GetStub().GetState(userId)
	if userBytes == nil {
		return fmt.Errorf("can not find user: %v", userId)
	}

	user := User{}
	json.Unmarshal(userBytes, &user)

	if user.UserType != "manufacturer" {
		return fmt.Errorf("only manufacturer can create product")
	}

	priceAsFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return fmt.Errorf("error converting price: %s", err.Error())
	}

	productCounter := getCounter(ctx, "ProductCounterNO")
	productCounter++

	txTimeAsPtr, err := t.GetTxTimestamp(ctx)
	if err != nil {
		return fmt.Errorf("error in transaction timestamp")
	}

	position := ProductPos{}
	position.Date = txTimeAsPtr
	position.Latitude = latitude
	position.Longitude = longitude

	product := Product{
		ProductID:      "Product" + strconv.Itoa(productCounter),
		Name:           name,
		ManufacturerID: user.UserID,
		SupplierID:     "",
		TransporterID:  "",
		CustomerID:     "",
		Status:         "Available",
		Position:       []ProductPos{position},
		Price:          priceAsFloat,
	}

	productAsBytes, err := json.Marshal(product)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(product.ProductID, productAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put to world state %v", err.Error())
	}

	incrementCounter(ctx, "ProductCounterNO")

	return nil
}

func (t *SupplyChain) updateProduct(ctx contractapi.TransactionContextInterface, userID string, productID string, name string, price string) error {

	userBytes, _ := ctx.GetStub().GetState(userID)
	if userBytes == nil {
		return fmt.Errorf("can not find user")
	}

	user := User{}
	json.Unmarshal(userBytes, &user)
	if user.UserType == "customer" {
		return fmt.Errorf("customer can not update product")
	}

	productBytes, err := ctx.GetStub().GetState(productID)

	if err != nil {
		return fmt.Errorf("failed to read product from world-state: %s", err.Error())
	}

	if productBytes == nil {
		return fmt.Errorf("can not find the product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	if product.TransporterID != "" {
		return fmt.Errorf("product sent to transporter. can not update price")
	}

	priceAsFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return fmt.Errorf("failed to convert price to float: %s", err.Error())
	}

	product.Name = name
	product.Price = priceAsFloat

	updateProductAsBytes, err := json.Marshal(product)
	if err != nil {
		return err
	}

	ctx.GetStub().PutState(product.ProductID, updateProductAsBytes)
	return nil

}

func (t *SupplyChain) toSupplier(ctx contractapi.TransactionContextInterface, productID string, supplierID string, longitude string, latitude string) error {

	userBytes, _ := ctx.GetStub().GetState(supplierID)

	if userBytes == nil {
		return fmt.Errorf("can not find Supplier")
	}

	user := User{}
	json.Unmarshal(userBytes, &user)

	if user.UserType != "supplier" {
		return fmt.Errorf("User must be a Supplier")
	}

	productBytes, err := ctx.GetStub().GetState(productID)

	if err != nil {
		return fmt.Errorf("failed to get product from world state: %s", productID)
	}

	if productBytes == nil {
		return fmt.Errorf("can not find the product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	if product.SupplierID != "" {
		return fmt.Errorf("Product is sent to Supplier already")
	}

	// Trnasaction Timestamp
	txTimeAsPtr, errTx := t.GetTxTimestamp(ctx)
	if errTx != nil {
		return fmt.Errorf("error getting transaction timestamp")
	}

	product.SupplierID = user.UserID
	product.Position = append(product.Position, ProductPos{Date: txTimeAsPtr, Latitude: latitude, Longitude: longitude})
	product.Status = "At warehouse"

	updateProductAsBytes, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Marshal Error: %s", err))
	}

	ctx.GetStub().PutState(product.ProductID, updateProductAsBytes)
	return nil
}

func (t *SupplyChain) toTransporter(ctx contractapi.TransactionContextInterface, productID string, transporterID string, longitude string, latitude string) error {

	userBytes, _ := ctx.GetStub().GetState(transporterID)

	if userBytes == nil {
		return fmt.Errorf("can not find user")
	}

	user := User{}
	json.Unmarshal(userBytes, &user)

	if user.UserType != "transporter" {
		return fmt.Errorf("User must be a Transporter")
	}

	productBytes, _ := ctx.GetStub().GetState(productID)
	if productBytes == nil {
		return fmt.Errorf("can not find the product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	if product.SupplierID == "" {
		return fmt.Errorf("product not sent to supplier yet")
	}

	if product.TransporterID != "" {
		return fmt.Errorf("product is sent to transporter already")
	}

	// Trnasaction Timestamp
	txTimeAsPtr, errTx := t.GetTxTimestamp(ctx)
	if errTx != nil {
		return fmt.Errorf("error getting transaction timeStamp")
	}

	product.TransporterID = user.UserID
	product.Position = append(product.Position, ProductPos{Date: txTimeAsPtr, Latitude: latitude, Longitude: longitude})
	product.Status = "In transit"

	updateProductAsBytes, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		return fmt.Errorf("marshal error in user %s", errMarshal)
	}

	errPut := ctx.GetStub().PutState(product.ProductID, updateProductAsBytes)
	if errPut != nil {
		return fmt.Errorf("failed to send to transporter: %s", product.ProductID)
	}

	fmt.Println("Product successfully sent for Transporting")
	return nil

}

func (t *SupplyChain) sellToCustomer(ctx contractapi.TransactionContextInterface, productID string, customerID string, longitude string, latitude string) error {
	productBytes, _ := ctx.GetStub().GetState(productID)
	if productBytes == nil {
		return fmt.Errorf("can not find the product")
	}

	product := Product{}
	json.Unmarshal(productBytes, &product)

	if product.TransporterID == "" {
		return fmt.Errorf("Product not sent to transporter yet")
	}
	if product.CustomerID != "" {
		return fmt.Errorf("Product already sold")
	}

	// Transaction Timestamp
	txTimeAsPtr, errTx := t.GetTxTimestamp(ctx)
	if errTx != nil {
		return fmt.Errorf("error in timestamp")
	}

	product.CustomerID = customerID
	product.Position = append(product.Position, ProductPos{Date: txTimeAsPtr, Latitude: latitude, Longitude: longitude})
	product.Status = "Sold"

	updateProductAsBytes, errMarshal := json.Marshal(product)
	if errMarshal != nil {
		return fmt.Errorf("marshal error in user %s", errMarshal)
	}

	ctx.GetStub().PutState(product.ProductID, updateProductAsBytes)
	return nil
}

func (t *SupplyChain) QueryProduct(ctx contractapi.TransactionContextInterface, productId string) (*Product, error) {
	productAsBytes, err := ctx.GetStub().GetState(productId)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %s", err.Error())
	}
	if productAsBytes == nil {
		return nil, fmt.Errorf("Product %s does not exist", productId)
	}
	product := new(Product)
	err = json.Unmarshal(productAsBytes, &product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (t *SupplyChain) QueryAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	startKey := "Product1"
	endKey := "Product999"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []*Product{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		product := new(Product)
		_ = json.Unmarshal(queryResponse.Value, product)
		results = append(results, product)
	}

	return results, nil
}

//  ---------------------------- main ------------------------------------------

func main() {
	chaincode, err := contractapi.NewChaincode(new(SupplyChain))
	if err != nil {
		fmt.Printf("Error creating chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting supply chain chaincode: %s", err.Error())
	}

}
