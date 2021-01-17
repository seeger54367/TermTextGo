package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"

	. "github.com/kevinburke/twilio-go"
)

func getContactMessages(c *Client, sc string, kn map[string]string) []string {
	contactMessages := []string{}
	allMessages := []*Message{}
	data := url.Values{}
	data2 := url.Values{}
	data.Set("From", sc)
	data.Set("PageSize", "100")
	data2.Set("To", sc)
	data2.Set("PageSize", "100")
	msgStr := ""
	msg, _ := c.Messages.GetPage(context.TODO(), data)
	for _, m := range msg.Messages {
		// Create message strings
		allMessages = append(allMessages, m)
	}
	msg, _ = c.Messages.GetPage(context.TODO(), data2)
	for _, m := range msg.Messages {
		// Create message strings
		allMessages = append(allMessages, m)
	}

	for _, m := range allMessages {
		// See if sender is a known contact
		newFrom := ""
		if val, ok := kn[string(m.From)]; ok {
			newFrom = val
		} else {
			newFrom = string(m.From)
		}
		msgStr = m.DateSent.Time.Format("2006-01-02 15:04:05") + " |  " + newFrom + ": " + m.Body
		contactMessages = append(contactMessages, msgStr)
	}
	sort.Strings(contactMessages)
	return contactMessages
}

func sendMessage(c *Client, m map[string]string) {
	c.Messages.SendMessage(m["from"], m["to"], m["message"], nil)
}

func setupClient() *Client {
	//Gets Twilio AccoundID from a plaintext file
	sidFile := "/path/to/file"
	//Gets Twilio Auth Token from a plaintext file
	authFile := "/path/to/file"

	sid, err := ioutil.ReadFile(sidFile)
	if err != nil {
		fmt.Println("An Error Has Occurred")
	}
	sid = sid[:len(string(sid))-1]

	token, err := ioutil.ReadFile(authFile)
	if err != nil {
		fmt.Println("An Error Has Occurred")
	}
	token = token[:len(string(token))-1]

	client := NewClient(string(sid), string(token), nil)

	return client
}

func getContactList(c *Client) []string {
	var allConvoNumbers = []string{}
	data := url.Values{}
	data2 := url.Values{}
	//Replace "myNum" value with your Twilio number
	myNum := "+19999999999"
	data.Set("From", myNum)
	data.Set("PageSize", "100")
	msg, _ := c.Messages.GetPage(context.TODO(), data)
	for _, m := range msg.Messages {
		allConvoNumbers = append(allConvoNumbers, string(m.To))
	}
	data2.Set("To", myNum)
	data2.Set("PageSize", "100")
	msg, _ = c.Messages.GetPage(context.TODO(), data2)
	for _, m := range msg.Messages {
		allConvoNumbers = append(allConvoNumbers, string(m.From))
	}
	uniqueConvoNumbers := unique(allConvoNumbers)
	return uniqueConvoNumbers
}

func unique(pn []string) []string {
	// From
	// https://www.geeksforgeeks.org/how-to-remove-duplicate-values-from-slice-in-golang/
	pm := make(map[string]bool)
	ul := []string{}

	for _, entry := range pn {
		if _, value := pm[entry]; !value {
			pm[entry] = true
			ul = append(ul, entry)
		}
	}
	return ul
}
