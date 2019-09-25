package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"time"

	"github.com/pelletier/go-toml"
)

var (
	netTransport *http.Transport
	netClient    *http.Client
)

type (
	// Config of the plugin
	Config struct {
		Token string

		Timeout string
	}
	// Plugin defines the sonar-scaner plugin parameters.
	Plugin struct {
		Config Config
	}
	// SonarReport it is the representation of .scannerwork/report-task.txt
	SonarReport struct {
		ProjectKey   string `toml:"projectKey"`
		ServerURL    string `toml:"serverUrl"`
		DashboardURL string `toml:"dashboardUrl"`
		CeTaskID     string `toml:"ceTaskId"`
		CeTaskURL    string `toml:"ceTaskUrl"`
	}

	// TaskResponse Give Compute Engine task details such as type, status, duration and associated component.
	TaskResponse struct {
		Task struct {
			ID            string `json:"id"`
			Type          string `json:"type"`
			ComponentID   string `json:"componentId"`
			ComponentKey  string `json:"componentKey"`
			ComponentName string `json:"componentName"`
			AnalysisID    string `json:"analysisId"`
			Status        string `json:"status"`
		} `json:"task"`
	}

	// ProjectStatusResponse Get the quality gate status of a project or a Compute Engine task
	ProjectStatusResponse struct {
		ProjectStatus struct {
			Status string `json:"status"`
		} `json:"projectStatus"`
	}
)

func init() {
	netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
}

// Exec executes the plugin step
func (p Plugin) Exec() error {
	cmd := exec.Command("sed", "-e", "s/=/=\"/", "-e", "s/$/\"/", ".scannerwork/report-task.txt")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// log.Printf("%s\n",output)

	report := SonarReport{}
	err = toml.Unmarshal(output, &report)
	if err != nil {
		return err
	}
	// log.Printf("%#v\n", report)

	taskRequest, err := http.NewRequest("GET", report.CeTaskURL, nil)
	taskRequest.Header.Add("Authorization", "Basic "+p.Config.Token)
	taskResponse, err := netClient.Do(taskRequest)
	if err != nil {
		return err
	}
	buf, _ := ioutil.ReadAll(taskResponse.Body)
	task := TaskResponse{}
	json.Unmarshal(buf, &task)
	// log.Printf("%#v\n",task)

	reportRequest := url.Values{
		"analysisId": {task.Task.AnalysisID},
	}
	projectRequest, err := http.NewRequest("GET", report.ServerURL+"/api/qualitygates/project_status?"+reportRequest.Encode(), nil)
	projectRequest.Header.Add("Authorization", "Basic "+p.Config.Token)
	projectResponse, err := netClient.Do(projectRequest)
	if err != nil {
		return err
	}
	buf, _ = ioutil.ReadAll(projectResponse.Body)
	project := ProjectStatusResponse{}
	json.Unmarshal(buf, &project)
	// log.Printf("%#v\n", project)
	if "OK" != project.ProjectStatus.Status {
		return errors.New("Quality Gate Failed")
	}
	fmt.Println("Quality Gate successfully passed")
	return nil
}
