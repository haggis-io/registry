package model

type DocumentQuery struct {
	Name    string
	Pattern bool
	Version string
	Limit   int32
	Author  string
	Status  string
}
