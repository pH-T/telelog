package telelog

import (
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

var (
	TELEGRAM_BOT_API              string = "https://api.telegram.org/bot"
	TELEGRAM_BOT_API_SEND_MESSAGE string = "/sendMessage"
	TELELOG_LEVEL_ERROR           string = "\u274C [ERROR]"
	TELELOG_LEVEL_ALERT           string = "\u2757 [ALERT]"
	TELELOG_LEVEL_INFO            string = "\u2139 [INFO]"
	TELELOG_LEVEL_SUCCESS         string = "\u2705 [SUCCESS]"
)

type Telelog struct {
	token  string
	chatID string
}

func NewTelelog(token string, chatID string) *Telelog {
	return &Telelog{token, chatID}
}

func sendMessage(msg string, token string, chatID string) error {

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", msg)

	resp, err := http.PostForm(TELEGRAM_BOT_API+token+TELEGRAM_BOT_API_SEND_MESSAGE, data)
	if err != nil {
		return errors.Wrap(err, "Could not POST data")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(err, "Could not POST data: Status was not 200")
	}

	return nil

}

func (telelog *Telelog) Log(msg string, level string) error {

	err := sendMessage(level+": "+time.Now().Format(time.RFC3339)+":\n "+msg, telelog.token, telelog.chatID)
	if err != nil {
		return errors.Wrap(err, "log failed")
	}

	return nil

}

func (telelog *Telelog) logError(msg string) error {

	return telelog.log(msg, TELELOG_LEVEL_ERROR)

}

func (telelog *Telelog) logAlert(msg string) error {

	return telelog.log(msg, TELELOG_LEVEL_ALERT)

}

func (telelog *Telelog) logInfo(msg string) error {

	return telelog.log(msg, TELELOG_LEVEL_INFO)

}

func (telelog *Telelog) logSuccess(msg string) error {

	return telelog.log(msg, TELELOG_LEVEL_SUCCESS)

}
