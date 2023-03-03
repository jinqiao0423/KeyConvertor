package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
)

const DELEGATED string = "delegated"

type KeyInfo struct {
	Type       string
	PrivateKey []byte
}

var lm = flag.String("lm", "", "Lotus private key string")
var ml = flag.String("ml", "", "Metamask private key string")

func main() {
	fmt.Println(parse())
}

func parse() string {
	flag.Parse()
	if len(*lm) != 0 && len(*ml) != 0 {
		panic("Cannot input -lm and -ml at the same command")
	}
	if len(*lm) == 0 && len(*ml) == 0 {
		panic("No valid private key input")
	}
	if len(*lm) != 0 {
		result, err := lotus2Meta(*lm)
		if err != nil {
			panic(err)
		}
		return result
	} else {
		result, err := meta2Lotus(*ml)
		if err != nil {
			panic(err)
		}
		return result
	}
}

func meta2Lotus(input string) (string, error) {
	privBytes, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	ki := KeyInfo{Type: "delegated", PrivateKey: privBytes}
	b, err := json.Marshal(ki)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func lotus2Meta(input string) (string, error) {
	b, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}
	var ki KeyInfo
	err = json.Unmarshal(b, &ki)
	if err != nil {
		return "", err
	}
	if ki.Type != DELEGATED {
		return "", errors.New("invalid key type")
	}
	return hex.EncodeToString(ki.PrivateKey), nil
}
