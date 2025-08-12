package model 

type Product struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Bytes []byte `json:"bytes"`
}