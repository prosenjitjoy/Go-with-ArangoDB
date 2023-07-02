package main

// A typical document type
type IntKeyValue struct {
	Key   string `json:"_key"`
	Value int    `json:"value"`
}

// A typical vertex type must
type MyVertexNode struct {
	Key string `json:"_key"`
	//....//
	Data   string  `json:"data"`
	Weight float64 `json:"weight"`
}

// A typical edge type must have fields mathing _from and _to
type MyEdgeLink struct {
	Key  string `json:"_key"`
	From string `json:"_from"`
	To   string `json:"_to"`
	//....//
	Weight float64 `json:"weight"`
}
