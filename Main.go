package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/modmuss50/goutils"
	"github.com/nanobox-io/golang-scribble"
	"github.com/thoas/stats"
)

//Databse help: https://medium.com/@skdomino/scribble-a-tiny-json-database-in-golang-9817854deb05

var (
	DataBase *scribble.Driver
)

func main() {

	fmt.Println("Loading Cross Server Storage - Server")

	db, _ := scribble.New("./db", nil)
	DataBase = db

	//item := Item{RegName:"Test", UUID:"123"}
	//DataBase.Write("items", item.UUID, item)

	middleware := stats.New()
	mux := http.NewServeMux()

	mux.HandleFunc("/list", listItems)
	mux.HandleFunc("/addItem", addItem)
	mux.HandleFunc("/removeItem", removeItem)
	mux.HandleFunc("/coins", getCoins)

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(middleware.Data())
		w.Write(b)
	})
	http.ListenAndServe(":8000", middleware.Handler(mux))

}

//TODO merge list and coins?

func listItems(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, goutils.ToJson(ListItems()))
}

func getCoins(w http.ResponseWriter, r *http.Request) {
	uuid := r.Header.Get("uuid")
	io.WriteString(w, goutils.ToJson(GetUser(uuid)))
}

func addItem(w http.ResponseWriter, r *http.Request) {
	//Sets max size to 10KB
	r.Body = http.MaxBytesReader(w, r.Body, 10000)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var item Item
	err = json.Unmarshal(b, &item)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	uuid := r.Header.Get("uuid")
	username := r.Header.Get("username")

	//Generates a random string for the item to aid with removing
	items := ListItems()
	item.UUID = RandString(16, int64(len(items)))

	DataBase.Write("items", item.UUID, item)
	AddCoin(uuid)

	Log(uuid + "(" + username + ") added " + item.UUID + " to the list")
	io.WriteString(w, goutils.ToJson(item))
}

func removeItem(w http.ResponseWriter, r *http.Request) {
	//Sets max size to 10KB
	r.Body = http.MaxBytesReader(w, r.Body, 10000)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var remove RemoveJson
	err = json.Unmarshal(b, &remove)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	uuid := r.Header.Get("uuid")
	username := r.Header.Get("username")
	if GetCoins(uuid) == 0 {
		io.WriteString(w, goutils.ToJson(ErrorResponse{"Not enough coins"}))
		return
	}

	removedItem := Item{}
	err = DataBase.Read("items", remove.UUID, &removedItem)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	DataBase.Delete("items", remove.UUID)
	RemoveCoin(uuid)

	Log(uuid + "(" + username + ") removed " + remove.UUID + " from the list")

	io.WriteString(w, goutils.ToJson(RemoveResponse{Success: true, Item: removedItem}))

}

func ListItems() []Item {
	items, _ := DataBase.ReadAll("items")
	itemList := []Item{}
	for _, item := range items {
		f := Item{}
		json.Unmarshal([]byte(item), &f)
		itemList = append(itemList, f)
	}
	return itemList
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(size int, seed int64) string {
	rand.Seed(time.Now().UnixNano() + seed)
	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func Log(str string) {
	goutils.AppendStringToFile(str, "log.txt")
}
