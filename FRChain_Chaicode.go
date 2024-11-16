package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type FRChainChaincode struct {
	contractapi.Contract
}

type FlowRule struct {
	Src       string `json:"src"`
	Dst       string `json:"dst"`
	Controller string `json:"controller"`
	Status    string `json:"status"` // Valid or Invalid
}

// Initialize Ledger
func (f *FRChainChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
	initialFlows := []FlowRule{}

	for _, flow := range initialFlows {
		flowJSON, err := json.Marshal(flow)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(flow.Src+flow.Dst, flowJSON)
		if err != nil {
			return err
		}
	}

	return nil
}

// Validate Flow Rule
func (f *FRChainChaincode) ValidateFlowRule(ctx contractapi.TransactionContextInterface, src string, dst string, controller string) (bool, error) {
	// Example validation rule: Disallow flows between certain addresses
	if src == "malicious_src" || dst == "malicious_dst" {
		return false, nil
	}

	flow := FlowRule{
		Src:       src,
		Dst:       dst,
		Controller: controller,
		Status:    "Valid",
	}

	flowJSON, err := json.Marshal(flow)
	if err != nil {
		return false, err
	}

	err = ctx.GetStub().PutState(src+dst, flowJSON)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Query Flow Rule
func (f *FRChainChaincode) QueryFlowRule(ctx contractapi.TransactionContextInterface, src string, dst string) (*FlowRule, error) {
	flowJSON, err := ctx.GetStub().GetState(src + dst)
	if err != nil {
		return nil, fmt.Errorf("failed to get flow: %v", err)
	}
	if flowJSON == nil {
		return nil, fmt.Errorf("flow not found")
	}

	var flow FlowRule
	err = json.Unmarshal(flowJSON, &flow)
	if err != nil {
		return nil, err
	}

	return &flow, nil
}
