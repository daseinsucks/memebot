package main

import (
	//	"bytes"

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/StarkBotsIndustries/telegraph/v2"
	pogreb "github.com/akrylysov/pogreb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"strconv"
	//	"strings"
)

var tgApiKey, err = ioutil.ReadFile(".secret")
var bot, _ = tgbotapi.NewBotAPI(string(tgApiKey))

const telegrap_base_url = "https://telegra.ph/"

func main() {
	tgbd, _ := pogreb.Open("telegramdb", nil)

	fmt.Println(tgApiKey)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			if update.Message.Photo != nil {

				n, _ := ioutil.ReadFile("amount.txt")
				m := string(n)

				amount, _ := strconv.ParseInt(m, 10, 64)

				amount++
				bts := []byte(fmt.Sprint(amount))
				ioutil.WriteFile("amount.txt", bts, 0644)

				length := len(update.Message.Photo)
				id := update.Message.Photo[length-1].FileID
				fileurl, _ := bot.GetFileDirectURL(id)
				fmt.Println(fileurl)

				file_name := fmt.Sprint(amount)

				file := createFile(file_name)
				GetFile(file, httpClient(), fileurl, file_name)
				telegraph_link := UploadFileToTelegraph(file_name)
				os.Remove(file_name)

				file_name = "0"

				//tgbd.Put([]byte(file_name), []byte(telegraph_link))

				randomNumber := "0"

				telegraph_link = getfromdb(tgbd, randomNumber)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, telegrap_base_url+telegraph_link)
				bot.Send(msg)

			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Memes only please.....")
				bot.Send(msg)
			}
		}
	}
}

func UploadFileToTelegraph(file_name string) string {
	file, _ := os.Open(file_name)
	// os.File is io.Reader so just pass it.
	link, _ := telegraph.Upload(file, "photo")
	log.Println(link)
	return link
}

func createFile(file_name string) *os.File {
	file, err := os.Create(file_name)

	checkError(err)
	return file
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetFile(file *os.File, client *http.Client, url string, file_name string) {
	resp, err := client.Get(url)
	checkError(err)
	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)
	defer file.Close()
	checkError(err)

	fmt.Println("Just Downloaded a file %s with size %d", file_name, size)
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func getfromdb(database *pogreb.DB, key string) string {
	defer handlePanic()

	thing, _ := database.Get([]byte(key))
	thingstring := string(thing)

	return thingstring
}

func handlePanic() {

	// detect if panic occurs or not
	a := recover()

	if a != nil {
		fmt.Println("RECOVER", a)
	}

}
