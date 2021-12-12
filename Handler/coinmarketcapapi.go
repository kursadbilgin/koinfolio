package Handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"koinfolio/Models"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func CoinMarketCapAPI(amount int, code string) (response Models.APIResponse) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/tools/price-conversion", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("amount", strconv.Itoa(amount))
	q.Add("symbol", code)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "30593477-b629-4f2b-bbb5-6a95d8e88211")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	return response
}
