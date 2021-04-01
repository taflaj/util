// ipinfo.go

package ipinfo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var headers map[string]string = make(map[string]string)

func init() {
	headers["Accept"] = "application/json"
}

func SetToken(token string) {
	headers["Authorization"] = "Bearer " + token
}

func SetAgent(agent string) {
	headers["User-Agent"] = agent
}

type Info struct {
	IP       string
	City     string
	Region   string
	Country  string
	Postal   string
	HostName string
	Org      string
	Bogon    bool
	Error    string
}

func GetInfo(address string) (*Info, error) {
	success := make(chan *Info)
	failure := make(chan error)
	go func() {
		client := &http.Client{}
		url := "http://ipinfo.io/" + address
		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			failure <- err
		}
		for k, v := range headers {
			request.Header.Set(k, v)
		}
		response, err := client.Do(request)
		if err != nil {
			failure <- err
		}
		defer response.Body.Close()
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			failure <- err
		}
		var info Info
		if err = json.Unmarshal(data, &info); err != nil {
			failure <- err
		}
		success <- &info
	}()
	select {
	case info := <-success:
		return info, nil
	case err := <-failure:
		return nil, err
	}
}
