package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type GetHeightArguments struct {
	ChainName string `json:"chain" jsonschema:"required,description=The name of the blockchain"`
}

func getHeight(rpc string) (uint64, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return 0, err
	}
	return client.BlockNumber(context.Background())
}

func getRpc(chainName string) (string, error) {
	switch chainName {
	case "ethereum":
		return "https://1rpc.io/eth", nil
	case "bsc":
		return "https://bsc-dataseed.bnbchain.org", nil
	}
	return "", errors.New("does not supported chain")
}

func main() {
	done := make(chan struct{})
	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
	err := server.RegisterTool("GetHeight", "Get the latest block height of blockchain", func(arguments GetHeightArguments) (*mcp_golang.ToolResponse, error) {
		rpc, err := getRpc(arguments.ChainName)
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(err.Error())), nil
		}
		height, err := getHeight(rpc)
		if err != nil {
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(err.Error())), nil
		}
		resText := fmt.Sprintf("The latest block height of %s is %v", arguments.ChainName, height)
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(resText)), nil
	})

	if err != nil {
		panic(err)
	}
	err = server.Serve()
	if err != nil {
		panic(err)
	}
	<-done
}
