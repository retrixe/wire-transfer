package main

import "encoding/json"

const idChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func errorJson(error string) string {
	json, err := json.Marshal(map[string]string{"error": error})
	if err != nil {
		panic(err)
	}
	return string(json)
}
