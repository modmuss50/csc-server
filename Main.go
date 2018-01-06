package main

import (
	"io"
	"net/http"
	"github.com/modmuss50/goutils"
	"github.com/thoas/stats"
	"encoding/json"
	"io/ioutil"
	"time"
	"math/rand"
	"fmt"
	"errors"
)

var (
	Info = Data{}
)

func main() {

	fmt.Println("Loading Cross Server Storage - Server")

	Load()

	//Info.ItemList = append(Info.ItemList, Item{RegName:"minecraft:stone", StackSize:32, ModID:"minecraft", UUID:"123"})
	//Save()

	middleware := stats.New()
	mux := http.NewServeMux()

	mux.HandleFunc("/list", listItems)
	mux.HandleFunc("/addItem", addItem)
	mux.HandleFunc("/removeItem", removeItem)

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(middleware.Data())
		w.Write(b)
	})
	http.ListenAndServe(":8000", middleware.Handler(mux))

}

func listItems(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, goutils.ToJson(Info))
}

func addItem(w http.ResponseWriter, r *http.Request) {

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

	//Generates a random string for the item to aid with removing
	item.UUID = RandString(16, int64(len(Info.ItemList)))

	Info.ItemList = append(Info.ItemList, item)

	uuid := r.Header.Get("uuid")
	username := r.Header.Get("username")

	Log(uuid + "(" + username + ") added " + item.UUID + " to the list")

	io.WriteString(w, goutils.ToJson(item))

	//Humm, might not want to do this after every request
	Save()
}

func removeItem(w http.ResponseWriter, r *http.Request) {

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

	removedItem, err := DeleteItem(remove.UUID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	Log(uuid + "(" + username + ") removed " + remove.UUID + " from the list")

	io.WriteString(w, goutils.ToJson(RemoveResponse{Success:true,Item:removedItem}))

	//Humm, might not want to do this after every request
	Save()
}

type Data struct {
	ItemList []Item `json:"items"`
}

type Item struct {
	RegName string `json:"registryName"`
	StackSize int `json:"stackSize"`
	Nbt string `json:"nbt"`
	ModID string `json:"modid"`
	//Each item is given a uuid when it is uploaded
	UUID string `json:"UUID"`
}

type RemoveJson struct {
	//The uuid of the item to remove
	UUID string `json:"UUID"`
}

type RemoveResponse struct {
	Success bool `json:"success"`
	Item Item `json:"item"`
}

func DeleteItem(uuid string) (Item, error){
	for index, item := range Info.ItemList {
		if item.UUID == uuid {
			Info.ItemList = RemoveItemFromList(Info.ItemList, index)
			return item, nil
		}
	}
	return Item{}, errors.New("Failed to remove item from list")
}

func RemoveItemFromList(s []Item, i int) []Item {
	fmt.Println("removing ", i , " from slice")
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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

func Log(str string){
	goutils.AppendStringToFile(str, "log.txt")
}

func Save(){
	json := goutils.ToJson(Info)
	goutils.WriteStringToFile(json, "data.json")
}

func Load(){
	if !goutils.FileExists("data.json") {
		return
	}

	jsonStr := goutils.ReadStringFromFile("data.json")
	var readData Data

	err := json.Unmarshal([]byte(jsonStr), &readData)
	if err != nil {
		fmt.Println(err)
		return
	}

	Info = readData

	fmt.Println("Loaded ", len(Info.ItemList), " items from data.json")
}