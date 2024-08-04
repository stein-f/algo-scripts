package main

import (
	"errors"
	"fmt"
	"github.com/multiformats/go-multicodec"
	"regexp"
	"strings"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

var (
	ErrUnknownSpec      = errors.New("unsupported template-ipfs spec")
	ErrUnsupportedField = errors.New("unsupported ipfscid field, only reserve is currently supported")
	ErrUnsupportedCodec = errors.New("unknown multicodec type in ipfscid spec")
	ErrUnsupportedHash  = errors.New("unknown hash type in ipfscid spec")
	ErrInvalidV0        = errors.New("cid v0 must always be dag-pb and sha2-256 codec/hash type")
	ErrHashEncoding     = errors.New("error encoding new hash")
	TemplateIPFSRegexp  = regexp.MustCompile(`template-ipfs://{ipfscid:(?P<version>[01]):(?P<codec>[a-z0-9\-]+):(?P<field>[a-z0-9\-]+):(?P<hash>[a-z0-9\-]+)}`)
)

func main() {
	templateImageURL := "template-ipfs://{ipfscid:1:raw:reserve:sha2-256}"
	reserveAddress := "CM5VWFSBRTNAHPRMH6YAO2SN4OKD4TXBZHW3CES4QUHSCXRUH44CMV7NNU"
	fmt.Println(RenderArc19URL(templateImageURL, reserveAddress))
}

func RenderArc19URL(templateImageURL string, reserveAddress string) string {
	address, err := types.DecodeAddress(reserveAddress)
	if err != nil {
		panic(err)
	}
	metadataUrl, err := renderArc19Template(templateImageURL, address)
	if err != nil {
		panic(err)
	}
	return metadataUrl
}

func renderArc19Template(template string, reserveAddress types.Address) (string, error) {
	matches := TemplateIPFSRegexp.FindStringSubmatch(template)
	if matches == nil {
		if strings.HasPrefix(template, "template-ipfs://") {
			return "", ErrUnknownSpec
		}
		return template, nil
	}
	if matches[TemplateIPFSRegexp.SubexpIndex("field")] != "reserve" {
		return "", ErrUnsupportedField
	}
	var (
		codec         multicodec.Code
		multihashType uint64
		hash          []byte
		err           error
		cidResult     cid.Cid
	)
	if err = codec.Set(matches[TemplateIPFSRegexp.SubexpIndex("codec")]); err != nil {
		return "", ErrUnsupportedCodec
	}
	multihashType = multihash.Names[matches[TemplateIPFSRegexp.SubexpIndex("hash")]]
	if multihashType == 0 {
		return "", ErrUnsupportedHash
	}
	hash, err = multihash.Encode(reserveAddress[:], multihashType)
	if err != nil {
		return "", ErrHashEncoding
	}
	if matches[TemplateIPFSRegexp.SubexpIndex("version")] == "1" {
		if multihashType != multihash.SHA2_256 {
			return "", ErrInvalidV0
		}
		cidResult = cid.NewCidV1(uint64(codec), hash)
	} else {
		cidResult = cid.NewCidV0(hash)
	}
	ipfsCID := strings.ReplaceAll(template, matches[0], cidResult.String())

	cidCleaned := strings.ReplaceAll(ipfsCID, "#arc69", "")
	return fmt.Sprintf("ipfs://%s", cidCleaned), nil
}
