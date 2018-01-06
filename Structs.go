package main

type Data struct {
	ItemList []Item `json:"items"`
}

type Item struct {
	RegName   string `json:"registryName"`
	StackSize int    `json:"stackSize"`
	MetaData  int    `json:"meta"`
	Nbt       string `json:"nbt"`
	ModID     string `json:"modid"`
	//Each item is given a uuid when it is uploaded
	UUID string `json:"UUID"`
}

type RemoveJson struct {
	//The uuid of the item to remove
	UUID string `json:"UUID"`
}

type RemoveResponse struct {
	Success bool `json:"success"`
	Item    Item `json:"item"`
}


type User struct {
	Coins int `json:"coins"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Transaction struct {
	TransactionType string  `json:"type"`
	UserName string  `json:"username"`
	UUID string  `json:"uuid"`
	ItemName string  `json:"itemName"`
	Cost int  `json:"cost"`
}