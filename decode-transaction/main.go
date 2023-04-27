package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

func main() {
	encodedTxn := "iqRhcmN2xCDo8Iyy/28kAPOHXWNDLwF8m8ox481nf9f7AiLEwBtapqNmZWXNA+iiZnbOAbGKSKNnZW6sbWFpbm5ldC12MS4womdoxCDAYcTY/B293tLXYEvkVo4/bQQZh6w3veS2ILWrOSSK36NncnDEIL/7xcRNMEZVOKbBdwW17ilTcCFKeVhHLygvgsdMpa4Gomx2zgGxjjCjc25kxCDo8Iyy/28kAPOHXWNDLwF8m8ox481nf9f7AiLEwBtapqR0eXBlpWF4ZmVypHhhaWTOFXV0sg=="
	decodedFromB64, err := base64.StdEncoding.DecodeString(encodedTxn)
	if err != nil {
		panic(err)
	}
	var decodedTxn types.Transaction
	if err := msgpack.Decode(decodedFromB64, &decodedTxn); err != nil {
		panic(err)
	}
	jsonTxn, err := json.Marshal(decodedTxn)
	if err != nil {
		panic(err)
	}
	println(string(jsonTxn))
}
