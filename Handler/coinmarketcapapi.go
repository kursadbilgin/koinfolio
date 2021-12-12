package Handler

import (
	"encoding/json"
	"io/ioutil"
	"koinfolio/Models"
	"net/http"
	"net/url"
	"strconv"
)

func CoinMarketCapAPI(amount int, code string) (response *Models.APIResponse, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/tools/price-conversion", nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("amount", strconv.Itoa(amount))
	q.Add("symbol", code)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "30593477-b629-4f2b-bbb5-6a95d8e88211")
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
