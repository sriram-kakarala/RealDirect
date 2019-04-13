/*
	The example is inspired by the marbles example of Fabric-samples
	To start with it does not have any validations as such but eventually the idea is to
	1. Create Regulatory Body Profile
	2. Create Trader Profile
	3. Create User Profile
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type asset struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	Quantity   int    `json:"quantity"`
	Owner      string `json:"owner"`
	Price      int    `json:"price"`
}

type assetPrivateDetails struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	Price      int    `json:"price"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	switch function {
	case "initasset":
		//create a new asset
		return t.initasset(stub, args)
	case "readasset":
		//read a asset
		return t.readasset(stub, args)
	case "readassetPrivateDetails":
		//read a asset private details
		return t.readassetPrivateDetails(stub, args)
	case "transferasset":
		//change owner of a specific asset
		return t.transferasset(stub, args)
	case "delete":
		//delete a asset
		return t.delete(stub, args)
	case "queryassetsByOwner":
		//find assets for owner X using rich query
		return t.queryassetsByOwner(stub, args)
	case "queryassets":
		//find assets based on an ad hoc rich query
		return t.queryassets(stub, args)
	case "getassetsByRange":
		//get assets based on range query
		return t.getassetsByRange(stub, args)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

// ============================================================
// initasset - create a new asset, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initasset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	type assetTransientInput struct {
		Name     string `json:"name"` //the fieldtags are needed to keep case from bouncing around
		Quantity int    `json:"quantity"`
		Owner    string `json:"owner"`
		Price    int    `json:"price"`
	}

	// ==== Input sanitation ====
	fmt.Println("- start init asset")

	if len(args) == 0 {
		return shim.Error("Incorrect number of arguments. Private asset data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	if _, ok := transMap["name"]; !ok {
		return shim.Error("asset must be a key in the transient map")
	}

	if len(transMap["name"]) == 0 {
		return shim.Error("asset value in the transient map must be a non-empty JSON string")
	}

	var assetInput assetTransientInput
	err = json.Unmarshal(transMap["asset"], &assetInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["asset"]))
	}

	if len(assetInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}

	if assetInput.Quantity <= 0 {
		return shim.Error("quanity field must be a positive integer")
	}
	if len(assetInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if assetInput.Price <= 0 {
		return shim.Error("price field must be a positive integer")
	}

	// ==== Check if asset already exists ====
	assetAsBytes, err := stub.GetPrivateData("collectionassets", assetInput.Name)
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if assetAsBytes != nil {
		fmt.Println("This asset already exists: " + assetInput.Name)
		return shim.Error("This asset already exists: " + assetInput.Name)
	}

	// ==== Create asset object, marshal to JSON, and save to state ====
	asset := &asset{
		ObjectType: "asset",
		Name:       assetInput.Name,
		Quantity:   assetInput.Quantity,
		Owner:      assetInput.Owner,
	}
	assetJSONasBytes, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save asset to state ===
	err = stub.PutPrivateData("collectionassets", assetInput.Name, assetJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Create asset private details object with price, marshal to JSON, and save to state ====
	assetPrivateDetails := &assetPrivateDetails{
		ObjectType: "assetPrivateDetails",
		Name:       assetInput.Name,
		Price:      assetInput.Price,
	}
	assetPrivateDetailsBytes, err := json.Marshal(assetPrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData("collectionassetPrivateDetails", assetInput.Name, assetPrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the asset to enable owner-based range queries, e.g. return all Sriram's assets ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~owner~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~owner~*
	indexName := "owner~name"
	ownerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{asset.Owner, asset.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the asset.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutPrivateData("collectionassets", ownerNameIndexKey, value)

	// ==== asset saved and indexed. Return success ====
	fmt.Println("- end init asset")
	return shim.Success(nil)
}

// ===============================================
// readasset - read a asset from chaincode state
// ===============================================
func (t *SimpleChaincode) readasset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the asset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionassets", name) //get the asset from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"asset does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ===============================================
// readassetreadassetPrivateDetails - read a asset private details from chaincode state
// ===============================================
func (t *SimpleChaincode) readassetPrivateDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the asset to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionassetPrivateDetails", name) //get the asset private details from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get private details for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"asset private details does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ==================================================
// delete - remove a asset key/value pair from state
// ==================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("- start delete asset")

	type assetDeleteTransientInput struct {
		Name string `json:"name"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private asset name must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	if _, ok := transMap["asset_delete"]; !ok {
		return shim.Error("asset_delete must be a key in the transient map")
	}

	if len(transMap["asset_delete"]) == 0 {
		return shim.Error("asset_delete value in the transient map must be a non-empty JSON string")
	}

	var assetDeleteInput assetDeleteTransientInput
	err = json.Unmarshal(transMap["asset_delete"], &assetDeleteInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["asset_delete"]))
	}

	if len(assetDeleteInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}

	// to maintain the owner~name index, we need to read the asset first and get its owner
	valAsbytes, err := stub.GetPrivateData("collectionassets", assetDeleteInput.Name) //get the asset from chaincode state
	if err != nil {
		return shim.Error("Failed to get state for " + assetDeleteInput.Name)
	} else if valAsbytes == nil {
		return shim.Error("asset does not exist: " + assetDeleteInput.Name)
	}

	var assetToDelete asset
	err = json.Unmarshal([]byte(valAsbytes), &assetToDelete)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(valAsbytes))
	}

	// delete the asset from state
	err = stub.DelPrivateData("collectionassets", assetDeleteInput.Name)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// Also delete the asset from the owner~name index
	indexName := "owner~name"
	ownerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{assetToDelete.Owner, assetToDelete.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.DelPrivateData("collectionassets", ownerNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// Finally, delete private details of asset
	err = stub.DelPrivateData("collectionassetPrivateDetails", assetDeleteInput.Name)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// ===========================================================
// transfer a asset by setting a new owner name on the asset
// ===========================================================
func (t *SimpleChaincode) transferasset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("- start transfer asset")

	type assetTransferTransientInput struct {
		Name  string `json:"name"`
		Owner string `json:"owner"`
	}

	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Private asset data must be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	if _, ok := transMap["asset_owner"]; !ok {
		return shim.Error("asset_owner must be a key in the transient map")
	}

	if len(transMap["asset_owner"]) == 0 {
		return shim.Error("asset_owner value in the transient map must be a non-empty JSON string")
	}

	var assetTransferInput assetTransferTransientInput
	err = json.Unmarshal(transMap["asset_owner"], &assetTransferInput)
	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["asset_owner"]))
	}

	if len(assetTransferInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(assetTransferInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}

	assetAsBytes, err := stub.GetPrivateData("collectionassets", assetTransferInput.Name)
	if err != nil {
		return shim.Error("Failed to get asset:" + err.Error())
	} else if assetAsBytes == nil {
		return shim.Error("asset does not exist: " + assetTransferInput.Name)
	}

	assetToTransfer := asset{}
	err = json.Unmarshal(assetAsBytes, &assetToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	assetToTransfer.Owner = assetTransferInput.Owner //change the owner

	assetJSONasBytes, _ := json.Marshal(assetToTransfer)
	err = stub.PutPrivateData("collectionassets", assetToTransfer.Name, assetJSONasBytes) //rewrite the asset
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferasset (success)")
	return shim.Success(nil)
}

// ===========================================================================================
// getassetsByRange performs a range query based on the start and end keys provided.

// Read-only function results are not typically submitted to ordering. If the read-only
// results are submitted to ordering, or if the query is used in an update transaction
// and submitted to ordering, then the committing peers will re-execute to guarantee that
// result sets are stable between endorsement time and commit time. The transaction is
// invalidated by the committing peers if the result set has changed between endorsement
// time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// ===========================================================================================
func (t *SimpleChaincode) getassetsByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetPrivateDataByRange("collectionassets", startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getassetsByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// =======Rich queries =========================================================================
// Two examples of rich queries are provided below (parameterized query and ad hoc query).
// Rich queries pass a query string to the state database.
// Rich queries are only supported by state database implementations
//  that support rich query (e.g. CouchDB).
// The query string is in the syntax of the underlying state database.
// With rich queries there is no guarantee that the result set hasn't changed between
//  endorsement time and commit time, aka 'phantom reads'.
// Therefore, rich queries should not be used in update transactions, unless the
// application handles the possibility of result set changes between endorsement and commit time.
// Rich queries can be used for point-in-time queries against a peer.
// ============================================================================================

// ===== Example: Parameterized rich query =================================================
// queryassetsByOwner queries for assets based on a passed in owner.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *SimpleChaincode) queryassetsByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	owner := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"asset\",\"owner\":\"%s\"}}", owner)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ===== Example: Ad hoc rich query ========================================================
// queryassets uses a query string to perform a query for assets.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the queryassetsForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *SimpleChaincode) queryassets(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString-Sriram queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetPrivateDataQueryResult("collectionassets", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}
