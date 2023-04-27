package main

import (
	"fmt"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

func main() {
	reserve, err := getReserveAddressFromCID("bafkreiag3vu2ncvfq4mhpfezfxdpvhyc4e6jfrwvzvjqd2s6wvssoakjzm")
	if err != nil {
		panic(err)
	}
	fmt.Println(reserve)
}

func getReserveAddressFromCID(cidToDecode string) (string, error) {
	decodedCID, err := cid.Decode(cidToDecode)
	if err != nil {
		return "", fmt.Errorf("failed to decode cid. %w", err)
	}
	reserve, err := reserveAddressFromCID(decodedCID)
	if err != nil {
		return "", fmt.Errorf("failed to generate reserve from cid. %w", err)
	}
	return reserve, nil
}

func reserveAddressFromCID(cidToEncode cid.Cid) (string, error) {
	decodedMultiHash, err := multihash.Decode(cidToEncode.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to decode ipfs cid: %w", err)
	}
	return types.EncodeAddress(decodedMultiHash.Digest)
}
