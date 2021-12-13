package Api

import (
	"encoding/json"
	"io/ioutil"
	"koinfolio/Config"
	"koinfolio/Models"
	"net/http"
	"net/url"
	"strconv"
)

func CoinMarketCapAPI(amount int, code string) (response *Models.APIResponse, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", Config.CoinMarketCapApiUrl, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("amount", strconv.Itoa(amount))
	q.Add("symbol", code)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", Config.CoinMarketCapApiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
