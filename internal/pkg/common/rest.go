package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

/*
/ The callPost function makes a REST API POST call.
/ The first parameter is the URL of the API.
/ The second parameter is the payload that needs to be sent with the POST request.
/ The third parameter is the authorization token that needs to be sent along with the request.
/ The function returns the response body, response status code and any error occurred during the request.
*/
func CallPost(url string, payload interface{}, user string, pass string) ([]byte, int, error) {
	jsonValue, _ := json.Marshal(payload)
	Debug("HTTP JSON payload:%s", jsonValue)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	CheckIfError(err)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user, pass)

	Debug("Sending POST request to: %s", req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Error(err.Error())
		return nil, 500, err
	}
	defer resp.Body.Close()

	Debug("HTTP Status: %s", resp.Status)
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		Error(err.Error())
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func ResponseJsonParser(resp []byte) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
