package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/lmittmann/tint"
)

const (
	endpoint_schedules     = "https://app.datadoghq.eu/api/unstable/on-call/schedules/"
	datadog_schedules_link = "https://app.datadoghq.eu/on-call/schedules/"
)

func get_on_call_schedule(
	schedule_name string, dd_api_key string, dd_app_key string,
) (string, error) {
	slog.Info("getting schedule for", "name", schedule_name)

	headers := map[string]string{
		"DD-API-KEY":         dd_api_key,
		"DD-APPLICATION-KEY": dd_app_key,
		"Accept":             "application/json",
	}

	// make the request
	body, err := http_req("schedules", "GET", endpoint_schedules, headers, nil)
	if err != nil {
		slog.Error(err.Error())
		return "", fmt.Errorf("error getting schedules http request")
	}

	// body is json so we need to parse it to get the schedule id
	var schedules SchedulesBody
	err = json.Unmarshal(body, &schedules)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling schedules response body: %v", err)
	}

	// loop through the schedules to find the schedule id from the schedule name
	for _, schedule := range schedules.Data {
		if schedule.Type == "schedules" && schedule.Attributes.Name == schedule_name {
			slog.Info("schedule found:", "id", schedule.Id)
			return schedule.Id, nil
		}
	}
	return "", fmt.Errorf("schedule not found")
}

func get_on_call_engineer(
	schedule_id string, dd_api_key string, dd_app_key string,
) (string, error) {
	slog.Info("getting on-call engineer for", "schedule_id", schedule_id)

	headers := map[string]string{
		"DD-API-KEY":         dd_api_key,
		"DD-APPLICATION-KEY": dd_app_key,
		"Accept":             "application/json",
	}

	// make the request
	body, err := http_req("engineer", "GET", endpoint_schedules+schedule_id+"/on-call", headers, nil)
	if err != nil {
		slog.Error(err.Error())
		return "", fmt.Errorf("error getting on-call engineer http request")
	}

	// body is json so we need to parse it to get the on-call engineer
	var on_call OnCallBody
	err = json.Unmarshal(body, &on_call)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling engineer response body: %v", err)
	}

	// get the on-call engineer id
	user_id := on_call.Data.Relationships.User.Data.Id
	slog.Info("on-call engineer id found:", "id", user_id)

	// loop through the included to find the on-call engineer name from the on-call engineer id
	for _, included := range on_call.Included {
		if included.Id == user_id && included.Type == "users" {
			slog.Info("on-call engineer email found:", "email", included.Attributes.Email)
			return included.Attributes.Email, nil
		}
	}
	return "", fmt.Errorf("on-call engineer not found")
}

func send_slack_message(
	on_call_engineer_email string,
	schedule_id string,
	endpoint_slack string,
	schedule_name string,
) (string, error) {
	slog.Info("sending slack message:", "on_call_engineer_email", on_call_engineer_email)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// prepare body for the slack message
	messageBody, err := json.Marshal(map[string]string{
		"engineer_email": on_call_engineer_email,
		"schedule_name":  schedule_name,
		"schedule_link":  datadog_schedules_link + schedule_id,
	})
	if err != nil {
		return "", fmt.Errorf("error marshalling slack message: %v", err)
	}

	// make the request
	body, err := http_req("slack", "POST", endpoint_slack, headers, messageBody)
	if err != nil {
		slog.Error(err.Error())
		return "", fmt.Errorf("error sending slack message with http request")
	}

	slog.Info("slack message sent:", "response", string(body))
	return string(body), nil
}

func http_req(
	kind string,
	method string,
	endpoint string,
	headers map[string]string,
	payload []byte,
) ([]byte, error) {
	// create http request to the datadog to get the on-call schedule engineer
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating %s request: %v", kind, err)
	}
	// set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making %s request: %v", kind, err)
	}
	defer resp.Body.Close()

	// read the body into a byte slice before unmarshalling it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading %s response body: %v", kind, err)
	}

	// check if the response code is 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"received non-200 response code for %s request: code=%v body=%s",
			kind,
			resp.StatusCode,
			body,
		)
	}
	return body, nil
}

func get_vars() (string, string, string, string, error) {
	schedule_name := os.Getenv("SCHEDULE_NAME")
	dd_api_key := os.Getenv("DD_API_KEY")
	dd_app_key := os.Getenv("DD_APP_KEY")
	endpoint_slack := os.Getenv("ENDPOINT_SLACK")
	var err error

	switch {
	case schedule_name == "":
		err = fmt.Errorf("SCHEDULE_NAME environment variable not set")
	case dd_api_key == "":
		err = fmt.Errorf("DD_API_KEY environment variable not set")
	case dd_app_key == "":
		err = fmt.Errorf("DD_APP_KEY environment variable not set")
	case endpoint_slack == "":
		err = fmt.Errorf("ENDPOINT_SLACK environment variable not set")
	}

	return schedule_name, dd_api_key, dd_app_key, endpoint_slack, err
}

func main() {
	var schedule_name, dd_api_key, dd_app_key, endpoint_slack string
	var err error

	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, nil)))

	schedule_name, dd_api_key, dd_app_key, endpoint_slack, err = get_vars()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	schedule_id, err := get_on_call_schedule(schedule_name, dd_api_key, dd_app_key)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	on_call_engineer_email, err := get_on_call_engineer(schedule_id, dd_api_key, dd_app_key)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	_, err = send_slack_message(
		on_call_engineer_email, schedule_id, endpoint_slack, schedule_name,
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
