package main

import (
	"github.com/valentinpelus/go-package/notification"

	"github.com/rs/zerolog/log"
)

func postMessageSlack(alertName string, namespace string, confPath *string) {
	notification.LoadConfSlack(*confPath)

	log.Info().Msgf("Alert notif %s %s", alertName, namespace)

	url := notification.ConfigurationSlack.WebhookUrl
	username := notification.ConfigurationSlack.SlackClient.UserName
	channel := notification.ConfigurationSlack.SlackClient.Channel
	clusterName := notification.ConfigurationSlack.ClusterName

	log.Info().Msgf("Slack url %s", url)
	log.Info().Msgf("Slack clusterName %s", clusterName)

	// Loading slack
	sc := notification.SlackClient{
		WebHookUrl: url,
		UserName:   username,
		Channel:    channel,
	}

	sr := notification.SlackJobNotification{
		Title: "Remediate - The auto-remediation has been triggered",
		Text:  "Remediate - The auto-remediation has been triggered",
		Details: "*Alert remediated*: " + alertName +
			"\r\n *Cluster*: " + clusterName +
			"\r\n *Namespace*: " + namespace,
		Color:     "#5581d9",
		IconEmoji: "necron",
	}

	err := sc.SendJobNotification(sr)
	if err != nil {
		log.Fatal().Err(err)
	}
}
