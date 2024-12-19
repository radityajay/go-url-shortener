package logc

import (
	"fmt"
	"os"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
)

type FormatType string

const (
	UNHANDLED  FormatType = "Unhandled error on "
	ERR        FormatType = "Error occured while "
	SUSPICIOUS FormatType = "Suspicious behavior detected, "
	MYSTERIOUS FormatType = "Mysterious guy detected "
	LAZY       FormatType = "Lazy to handle, plz handle sar "
)

type SlackConfig struct {
	Username *string
	Channel  *string
}

func sendToSlack(payload slack.Payload) {
	err := slack.Send(golog.Slack.URL, "", payload)
	if err != nil {
		fmt.Println(err)
	}
}

func compose(s *SlackConfig, message string, messageType string, color string, emoji string, errors error) {
	if s.Username == nil {
		s.Username = &golog.Slack.Username
	}
	if s.Channel == nil {
		s.Channel = &golog.Slack.Channel
	}

	attachment := slack.Attachment{
		Color: &color,
	}
	attachment.AddField(slack.Field{
		Title: "Message",
		Value: message,
	}).AddField(slack.Field{
		Title: "Level",
		Value: messageType,
	})

	if errors != nil {
		attachment.AddField(slack.Field{
			Title: "Exception",
			Value: fmt.Sprintf("``` %s ```", errors.Error()),
		})
	}

	payload := slack.Payload{
		Username:    *s.Username,
		Channel:     *s.Channel,
		IconEmoji:   emoji,
		Text:        message,
		Attachments: []slack.Attachment{attachment},
	}
	sendToSlack(payload)
}

func composeWithData(s *SlackConfig, message string, messageType string, color string, emoji string, data []byte, e error) {
	if s.Username == nil {
		s.Username = &golog.Slack.Username
	}
	if s.Channel == nil {
		s.Channel = &golog.Slack.Channel
	}

	attachment := slack.Attachment{
		Color: &color,
	}
	attachment.AddField(slack.Field{
		Title: "Message",
		Value: message,
	}).AddField(slack.Field{
		Title: "Level",
		Value: messageType,
	}).AddField(slack.Field{
		Title: "Data",
		Value: fmt.Sprintf("``` %s ```", string(data)),
	})

	if e != nil {
		attachment.AddField(slack.Field{
			Title: "Exception",
			Value: fmt.Sprintf("``` %s ```", e.Error()),
		})
	}

	payload := slack.Payload{
		Username:    *s.Username,
		Channel:     *s.Channel,
		IconEmoji:   emoji,
		Text:        message,
		Attachments: []slack.Attachment{attachment},
	}
	sendToSlack(payload)
}

func GetRecoveryChannel() *string {
	slChannel := os.Getenv("SLACK_CHANNEL_RECOVERY")
	if slChannel == "" {
		slChannel = os.Getenv("SLACK_CHANNEL")
	}
	return &slChannel
}

func LogInfo(message string, data *[]byte, slackConfig *SlackConfig) {
	if data == nil {
		compose(slackConfig, message, "INFO", "#2eb886", ":ok_hand:", nil)
	} else {
		composeWithData(slackConfig, message, "INFO", "#2eb886", ":ok_hand:", *data, nil)
	}
}

func LogError(message string, data *[]byte, e error, slackConfig *SlackConfig) {
	if data == nil {
		compose(slackConfig, message, "ERROR", "#a30200", ":bomb:", e)
	} else {
		composeWithData(slackConfig, message, "ERROR", "#a30200", ":bomb:", *data, e)
	}
}

func LogWarning(message string, data *[]byte, e error, slackConfig *SlackConfig) {
	if data == nil {
		compose(slackConfig, message, "WARNING", "#ffc107", ":warning:", e)
	} else {
		composeWithData(slackConfig, message, "WARNING", "#ffc107", ":warning:", *data, e)
	}
}

func Error(msg string, err error, data *[]byte) {
	env := os.Getenv("ENV")
	if env != "test" {
		if data == nil {
			golog.Slack.Error(msg, err)
		} else {
			golog.Slack.ErrorWithData(msg, *data, err)
		}
	}
}

func ErrorWithFormat(title string, format FormatType, err error, data *[]byte) {
	env := os.Getenv("ENV")
	if env != "test" {
		if data == nil {
			golog.Slack.Error(string(format)+title, err)
		} else {
			golog.Slack.ErrorWithData(string(format)+title, *data, err)
		}
	}
}

func Info(msg string, data *[]byte) {
	env := os.Getenv("ENV")
	if env != "test" {
		if data != nil {
			golog.Slack.InfoWidthData(msg, *data)
		} else {
			golog.Slack.Info(msg)
		}
	}
}

func InfoWithFormat(title string, format FormatType, data *[]byte) {
	env := os.Getenv("ENV")
	if env != "test" {
		if data != nil {
			golog.Slack.InfoWidthData(string(format)+title, *data)
		} else {
			golog.Slack.Info(string(format) + title)
		}
	}
}

func Warn(msg string, err error, data *[]byte) {
	env := os.Getenv("ENV")
	if env != "test" {
		if data != nil {
			golog.Slack.WarningWithData(msg, *data, err)
		} else {
			golog.Slack.Warning(msg, err)
		}
	}
}

func WarnWithFormat(title string, format FormatType, err error, data *[]byte) {
	env := os.Getenv("ENV")
	if env != "test" {
		if data != nil {
			golog.Slack.WarningWithData(string(format)+title, *data, err)
		} else {
			golog.Slack.Warning(string(format)+title, err)
		}
	}
}

func PrintJson(data interface{}) {
	js, _ := sonic.Marshal(data)
	fmt.Println(string(js))
}
