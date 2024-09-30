package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	sendMessageMethod = "sendMessage"
	getUpdatesMethod  = "getUpdates"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

// NewClient creates a new user client
func NewClient(host string, token string) Client {
	return Client{
		host:     host,
		basePath: "bot" + token,
		client:   http.Client{},
	}
}

// Updates is use to retrieve updates. offset is used to retrieve later updates, see more info in telegram API getUpdates chapter "getting updates"
func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	query := url.Values{}
	// strconv to int => ASCII
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))
	data, err := c.doRequest(getUpdatesMethod, query)
	if err != nil {
		return nil, err
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// SendMessage is used to send a message, receiving predefined args. See more in chapter "sendMessage" telegram API
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("can't send this message %w", err)
	}

	return nil
}

// doRequest is used to send a request
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	// prepare a new request
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't prepare a new request %w", err)
	}

	req.URL.RawQuery = query.Encode()

	// send prepared request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't send a prepared request %w", err)
	}

	return body, nil
}
