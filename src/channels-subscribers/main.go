package main

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tgstat "github.com/helios-ag/tgstat-go"
	"github.com/helios-ag/tgstat-go/channels"
	"os"
	"strconv"
	"time"
)

var qs = []*survey.Question{
	{
		Name:     "Token",
		Prompt:   &survey.Input{Message: "Enter your token"},
		Validate: survey.Required,
	},
	{
		Name:     "ChannelId",
		Prompt:   &survey.Input{Message: "Enter channel id"},
		Validate: survey.Required,
	},
	{
		Name:   "Group",
		Prompt: &survey.Select{Message: "Choose grouping", Options: []string{"hour", "day", "week", "month"}},
	},
	{
		Name:   "StartTime",
		Prompt: &survey.Input{Message: "Start Time", Default: ""},
	},
	{
		Name:   "EndTime",
		Prompt: &survey.Input{Message: "End Time", Default: ""},
	},
}

func main() {
	answers := struct {
		Token     string
		ChannelId string
		Group     string
		StartTime string
		EndTime   string
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var startTime, endTime string
	if answers.StartTime != "" {
		startTime = strconv.FormatInt(time.Now().Unix()-86400, 10)
	}

	if answers.EndTime != "" {
		endTime = strconv.FormatInt(time.Now().Unix(), 10)
	}

	var group *string
	if answers.Group != "" {
		group = tgstat.String(answers.Group)
	}
	req := channels.ChannelSubscribersRequest{
		ChannelId: answers.ChannelId,
		StartDate: tgstat.String(startTime),
		EndDate:   tgstat.String(endTime),
		Group:     group,
	}

	tgstat.Token = answers.Token

	info, _, err := channels.Subscribers(context.Background(), req)

	if err != nil {
		fmt.Printf("error getting data: %v\n", err)
		os.Exit(1)
	}
	fmt.Print("Err values")
	for _, info := range info.Response {
		fmt.Printf("ParticipantsCount: %d\n", info.ParticipantsCount)
		fmt.Printf("Period: %s\n", info.Period)
	}

}
