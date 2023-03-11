package main

import gonanoid "github.com/matoous/go-nanoid/v2"

const idChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateFileId() string {
	return gonanoid.MustGenerate(idChars, 8)
}
