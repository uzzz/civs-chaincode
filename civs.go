package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("civs")

type CivsChaincode struct {
}

func (t *CivsChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *CivsChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	logger.Info("Invoking function '%s'", function)

	switch function {
	case "PutElection":
		return t.putElection(stub, args)
	case "StartElection":
		return t.startElection(stub, args)
	case "StopElection":
		return t.stopElection(stub, args)
	case "PutVotes":
		return t.putVotes(stub, args)
	default:
		return nil, errors.New("Unsupported operation")
	}

	return nil, nil
}

func (t *CivsChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	logger.Info("Querying function '%s'", function)

	switch function {
	case "GetElection":
		return t.getElection(stub, args)
	case "GetVotes":
		return t.getVotes(stub, args)
	default:
		return nil, errors.New("Unsupported operation")
	}

	return nil, nil
}

func (t *CivsChaincode) putElection(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("CreateElection operation must include two arguments")
	}
	key := args[0]
	value := args[1]

	var i interface{}
	var err error

	err = json.Unmarshal([]byte(value), &i)
	if err != nil {
		return nil, fmt.Errorf("CreateElection operation failed. Invalid json: %s", err)
	}

	err = stub.PutState(key, []byte(value))
	if err != nil {
		return nil, fmt.Errorf("put operation failed. Error updating state: %s", err)
	}

	return nil, nil
}

func (t *CivsChaincode) getElection(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("getElection operation must include one argument, a key")
	}

	return t.fetchElection(stub, args[0])
}

func (t *CivsChaincode) startElection(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("StartElection operation must include one argument, a key")
	}

	return t.setElectionState(stub, args[0], "started")
}

func (t *CivsChaincode) stopElection(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("StopElection operation must include one argument, a key")
	}

	return t.setElectionState(stub, args[0], "stopped")
}

func (t *CivsChaincode) setElectionState(stub *shim.ChaincodeStub, key, state string) ([]byte, error) {
	electionData, err := t.fetchElection(stub, key)
	if err != nil {
		return nil, err
	}

	var election map[string]interface{}
	err = json.Unmarshal(electionData, &election)
	if err != nil {
		return nil, fmt.Errorf("StartElection operation failed. Invalid json: %s", err)
	}

	election["state"] = state

	value, err := json.Marshal(election)
	if err != nil {
		return nil, fmt.Errorf("StartElection operation failed. Failed to serialize: %s", err)
	}

	err = stub.PutState(key, value)
	if err != nil {
		return nil, fmt.Errorf("StartElection operation failed. Error updating state: %s", err)
	}

	return nil, nil
}

func (t *CivsChaincode) putVotes(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("PutVotes operation must include two arguments")
	}
	electionKey := args[0]
	key := electionKey + "_votes"
	value := args[1]

	var i interface{}
	var err error

	err = json.Unmarshal([]byte(value), &i)
	if err != nil {
		return nil, fmt.Errorf("PutVotes operation failed. Invalid json: %s", err)
	}

	err = stub.PutState(key, []byte(value))
	if err != nil {
		return nil, fmt.Errorf("put operation failed. Error updating state: %s", err)
	}

	return nil, nil
}

func (t *CivsChaincode) getVotes(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("GetVotes operation must include one argument, an election key")
	}
	electionKey := args[0]

	data, err := stub.GetState(electionKey + "_votes")
	if err != nil {
		logger.Error("Error getting state: %s", err)
		return nil, fmt.Errorf("fetch election operation failed: %s", err)
	}

	if len(data) == 0 {
		return []byte("{}"), nil
	}

	return data, nil
}

func (t *CivsChaincode) fetchElection(stub *shim.ChaincodeStub, key string) ([]byte, error) {
	data, err := stub.GetState(key)
	if err != nil {
		logger.Error("Error getting state: %s", err)
		return nil, fmt.Errorf("fetch election operation failed: %s", err)
	}

	return data, nil
}

func main() {
	err := shim.Start(new(CivsChaincode))
	if err != nil {
		fmt.Printf("Error starting CIVS chaincode: %s\n", err)
	}
}
