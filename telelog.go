package telelog

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", msg)

	req, err := http.NewRequest("POST", TELEGRAM_BOT_API+token+TELEGRAM_BOT_API_SEND_MESSAGE, bytes.NewBuffer([]byte(data.Encode())))
	if err != nil {
		return errors.Wrap(err, "Could not create request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Could not POST data")
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Could not read response")
	}

	type telegramResponse struct {
		Ok          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}
	res := telegramResponse{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return errors.Wrap(err, "Could not read json response into struct")
	}

	if resp.StatusCode != http.StatusOK || !res.Ok {
		return errors.Wrap(err, "Could not POST data: HTTP-status was not 200 or status was not ok")
	}

	return nil
}

// Log logs the given msg with the given level (e.g. see 'TELELOG_LEVEL_ERROR' above) to telegram
func (telelog *Telelog) Log(msg string, level string) error {
	err := sendMessage(level+":\n"+msg, telelog.token, telelog.chatID)
	if err != nil {
		return errors.Wrap(err, "log failed")
	}

	return nil
}

func (telelog *Telelog) LogError(msg string) error {
	return telelog.Log(msg, TELELOG_LEVEL_ERROR)
}

func (telelog *Telelog) LogAlert(msg string) error {
	return telelog.Log(msg, TELELOG_LEVEL_ALERT)
}

func (telelog *Telelog) LogInfo(msg string) error {
	return telelog.Log(msg, TELELOG_LEVEL_INFO)
}

func (telelog *Telelog) LogSuccess(msg string) error {
	return telelog.Log(msg, TELELOG_LEVEL_SUCCESS)
}
