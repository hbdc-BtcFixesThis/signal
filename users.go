package main

type User struct {
	Name     string `json:"name"`
	PassHash string `json:"pass_hash"`
	ActiveDB string `json:"db"`
}
