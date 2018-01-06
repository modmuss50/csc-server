package main

func GetUser(uuid string) User {
	user := User{}
	err := DataBase.Read("users", uuid, &user)
	if err != nil {
		return CreateUser(uuid)
	}
	return user
}

func CreateUser(uuid string) User {
	user := User{Coins:5}
	DataBase.Write("users", uuid, user)
	return user
}

func UpdateUser(uuid string, user User){
	DataBase.Write("users", uuid, user)
}

func GetCoins(uuid string) int{
	return GetUser(uuid).Coins
}

func AddCoin(uuid string){
	user := GetUser(uuid)
	user.Coins ++
	UpdateUser(uuid, user)
}

func RemoveCoin(uuid string){
	user := GetUser(uuid)
	user.Coins --
	UpdateUser(uuid, user)
}