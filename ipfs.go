package main

import (
	"strings"
)

func uploadJson2ipfs(jsonStringReader *strings.Reader) (jsonCID string, err error) {
	jsonCID, err = IPFS_SHELL.Add(jsonStringReader)
	return
}
