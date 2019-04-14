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
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type user struct {
	ObjectType  string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name        string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	DisplayName string `json:"displayname"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}

type asset struct {
	ObjectType  string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name        string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	DisplayName string `json:"displayname"`
	Quantity    int    `json:"quantity"`
	Owner       string `json:"owner"`
	Price       int    `json:"price"`
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
	case "inituser":
		return t.inituser(stub, args)
	case "signinuser":
		return t.signinuser(stub, args)
	case "initasset":
		//create a new asset
		return t.initasset(stub, args)
	case "readasset":
		//read a asset
		return t.readasset(stub, args)
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
	case "getHistoryForAsset":
		return t.getHistoryForAsset(stub, args)
	case "queryMarblesWithPagination":
		return t.queryAssetsWithPagination(stub, args)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

func (t *SimpleChaincode) inituser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1
	// "Name", "Password"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init asset")
	if len(args[0]) <= 0 {
		return shim.Error("Username cannot be empty")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Password cannot be empty")
	}

	userName := strings.ToLower(args[0])
	userDisplayName := args[0]
	userEmail := strings.ToLower(args[1])
	userPassword := args[2]

	// ==== Check if asset already exists ====
	userAsBytes, err := stub.GetState(userEmail)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userAsBytes != nil {
		fmt.Println("This user already exists: " + userName)
		return shim.Error("300")
	}

	// ==== Create asset object and marshal to JSON ====
	objectType := "user"
	user := &user{objectType, userName, userDisplayName, userPassword, userEmail}
	userJSONasBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error("301")
	}

	// === Save asset to state ===
	err = stub.PutState(userEmail, userJSONasBytes)
	if err != nil {
		return shim.Error("302")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) signinuser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1
	// "email", "Password"
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init asset")
	if len(args[0]) <= 0 {
		return shim.Error("Username cannot be empty")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Password cannot be empty")
	}

	userEmail := strings.ToLower(args[0])
	userPassword := args[1]

	// ==== Check if asset already exists ====
	userAsBytes, err := stub.GetState(userEmail)
	if err != nil {
		return shim.Error("Failed to get user: " + err.Error())
	} else if userAsBytes == nil {
		fmt.Println("User does not exist: " + userEmail)
		return shim.Error("300")
	}

	var userJSON user
	err = json.Unmarshal(userAsBytes, &userJSON)
	if err != nil {
		return shim.Error("301")
	}

	if userJSON.Password == userPassword {
		return shim.Success(nil)
	}

	return shim.Error("302")
}

// ============================================================
// initasset - create a new asset, store into chaincode state
// ============================================================
func (t *SimpleChaincode) initasset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1       2     3
	// "Name", "Quantity", "Owner", "Price"
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init asset")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	assetName := strings.ToLower(args[0])
	assetDisplayName := args[0]
	assetQuantity, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Quanity must be a numeric string")
	}
	assetOwner := strings.ToLower(args[2])
	assetPrice, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Price argument must be a numeric string")
	}

	// ==== Check if asset already exists ====
	assetAsBytes, err := stub.GetState(assetName)
	if err != nil {
		return shim.Error("Failed to get asset: " + err.Error())
	} else if assetAsBytes != nil {
		fmt.Println("This asset already exists: " + assetName)
		return shim.Error("This asset already exists: " + assetName)
	}

	// ==== Create asset object and marshal to JSON ====
	objectType := "asset"
	asset := &asset{objectType, assetName, assetDisplayName, assetQuantity, assetOwner, assetPrice}
	assetJSONasBytes, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save asset to state ===
	err = stub.PutState(assetName, assetJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the asset to enable color-based range queries, e.g. return all assets belonging to Sriram ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	indexName := "owner~name"
	ownerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{asset.Owner, asset.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the asset.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(ownerNameIndexKey, value)

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

	name = strings.ToLower(args[0])
	valAsbytes, err := stub.GetState(name) //get the asset from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Asset does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ==================================================
// delete - remove a asset key/value pair from state
// ==================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var assetJSON asset
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	assetName := strings.ToLower(args[0])

	// to maintain the owner~name index, we need to read the asset first and get its owner
	valAsbytes, err := stub.GetState(assetName) //get the asset from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + assetName + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Asset does not exist: " + assetName + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &assetJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + assetName + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(assetName) //remove the asset from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// maintain the index
	indexName := "owner~name"
	ownerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{assetJSON.Owner, assetJSON.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(ownerNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

// ===========================================================
// transfer a asset by setting a new owner name on the asset
// ===========================================================
func (t *SimpleChaincode) transferasset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0       1
	// "name", "newowner"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	assetName := strings.ToLower(args[0])
	newOwner := strings.ToLower(args[1])
	fmt.Println("- start transferAsset ", assetName, newOwner)

	assetAsBytes, err := stub.GetState(assetName)
	if err != nil {
		return shim.Error("Failed to get asset:" + err.Error())
	} else if assetAsBytes == nil {
		return shim.Error("Asset does not exist")
	}

	assetToTransfer := asset{}
	err = json.Unmarshal(assetAsBytes, &assetToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	indexName := "owner~name"
	ownerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{assetToTransfer.Owner, assetToTransfer.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(ownerNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	//change the owner
	assetToTransfer.Owner = newOwner

	assetJSONasBytes, _ := json.Marshal(assetToTransfer)
	err = stub.PutState(assetName, assetJSONasBytes) //rewrite the asset
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Update the index
	newOwnerNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{assetToTransfer.Owner, assetToTransfer.Name})
	if err != nil {
		return shim.Error(err.Error())
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the asset.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	stub.PutState(newOwnerNameIndexKey, value)

	fmt.Println("- end transferAsset (success)")
	return shim.Success(nil)
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

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
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

	return &buffer, nil
}

func (t *SimpleChaincode) queryAssetsWithPagination(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "queryString"
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	queryString := args[0]
	//return type of ParseInt is int64
	pageSize, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	bookmark := args[2]

	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, int32(pageSize), bookmark)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func getQueryResultForQueryStringWithPagination(stub shim.ChaincodeStubInterface, queryString string, pageSize int32, bookmark string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", bufferWithPaginationInfo.String())

	return buffer.Bytes(), nil
}

// ===========================================================================================
// addPaginationMetadataToQueryResults adds QueryResponseMetadata, which contains pagination
// info, to the constructed query results
// ===========================================================================================
func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *pb.QueryResponseMetadata) *bytes.Buffer {

	buffer.WriteString("[{\"ResponseMetadata\":{\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}]")

	return buffer
}

func (t *SimpleChaincode) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	assetName := strings.ToLower(args[0])

	fmt.Printf("- start getHistoryForAsset: %s\n", assetName)

	resultsIterator, err := stub.GetHistoryForKey(assetName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the asset
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON asset)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
