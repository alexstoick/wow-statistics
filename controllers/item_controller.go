package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var REDIS_SET_KEY = os.Getenv("REDIS_SET_KEY")

func NewData(c *gin.Context) {
	redis := FetchRedisFromContext(c)
	c.JSON(200, redis)

	item_list := makeRequest("https://item-list.wow.stoica.xyz/v1")
	var arr []string
	json.Unmarshal(item_list, &arr)

	for _, value := range arr {
		res := makeRequest("https://api.wow.stoica.xyz/v1/items/" + value)
		var lolmap map[string]interface{}
		json.Unmarshal(res, &lolmap)
		lol := lolmap["Item"].(map[string]interface{})
		itemID := lol["ItemID"].(float64)
		itemName := lol["ItemName"].(string)
		new_price := lolmap["BuyPrice"].(float64)
		itemID_string := strconv.FormatFloat(itemID, 'f', 0, 64)
		old_price, _ := redis.Get(itemID_string).Result()

		var text string
		if int(new_price) > len(old_price) {
			text = fmt.Sprintf("Price increase for %s from %.0f to %s", itemName, new_price, old_price)
		} else {
			text = fmt.Sprintf("Price decrease for %s from %s to %s", itemName, new_price, old_price)
		}

		sendNotification(text)

		redis.Set(itemID_string, int(new_price), 0).Err()
	}
}

func makeRequest(url string) []byte {
	response, err := http.Get(url)
	var contents []byte
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer response.Body.Close()
		contents, _ = ioutil.ReadAll(response.Body)
	}
	return contents
}

func sendNotification(text string) {
	url := "http://bot.management.stoica.xyz/unsafe_endpoint"
	body := map[string]string{"text": text}
	bytesarr, _ := json.Marshal(body)
	fmt.Println(string(bytesarr))

	client := &http.Client{}
	client.Post(url, "application/json", bytes.NewBuffer(bytesarr))

}
