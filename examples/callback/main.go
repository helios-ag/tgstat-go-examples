package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/callback"
	"net/http"
	"os"
)

var qs = []*survey.Question{
	{
		Name:     "Token",
		Prompt:   &survey.Input{Message: "Enter your token"},
		Validate: survey.Required,
	},
	{
		Name:   "CallbackURL",
		Prompt: &survey.Input{Message: "Enter callback url"},
	},
	{
		Name:   "ChannelId",
		Prompt: &survey.Input{Message: "Enter ChannelId"},
	},
}

// Simple example that can be used with ngrok for testing purposes
func main() {
	answers := struct {
		Token       string
		CallbackURL string
		ChannelId   string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tgstat.Token = answers.Token

	callbackReq := ""

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		cbReq, res, setErr := callback.SetCallback(context.Background(), answers.CallbackURL)
		callbackReq = cbReq.VerifyCode
		fmt.Fprint(w, callbackReq)
		if setErr != nil {
			fmt.Printf("error setting callBack: %v\n", setErr)
			fmt.Printf("status: %v\n", res.Status)
			fmt.Printf("status: %d\n", res.StatusCode)
			fmt.Printf("status: %v\n", res.Body)
			os.Exit(1)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			fmt.Fprint(w, callbackReq)
		}
	})

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		resp, _, errInfo := callback.GetCallbackInfo(context.Background())
		if errInfo == nil {
			fmt.Fprintf(w, resp.Status)
			fmt.Fprintf(w, resp.Response.Url)
			fmt.Fprintf(w, resp.Response.LastErrorMessage)
			fmt.Fprint(w, resp.Response.LastErrorDate)
			fmt.Fprint(w, resp.Response.PendingUpdateCount)

		}
		fmt.Fprintf(w, "error getting callBack info: %v\n", errInfo)
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		request := callback.SubscribeChannelRequest{
			ChannelId:  answers.ChannelId,
			EventTypes: "new_post",
		}
		resp, _, errInfo := callback.SubscribeChannel(context.Background(), request)
		if errInfo != nil {
			fmt.Fprintf(w, "error subscribing: %v\n", errInfo)
		} else {
			fmt.Fprint(w, fmt.Sprint(resp.SubscriptionId))
		}
	})

	http.ListenAndServe(":8081", nil)
}
