package request

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Data struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	City  string `json:"city"`
}

// Fetch data from API in batches
func FetchData(apiURL string, batchSize, page int, ch chan<- []Data) {
	client := resty.New()

	resp, err := client.R().
		SetResult(&[]Data{}).
		Get(fmt.Sprintf("%s?page=%d&size=%d", apiURL, page, batchSize))

	if err != nil {
		fmt.Println("API Error:", err)
		close(ch)
		return
	}

	data := *resp.Result().(*[]Data)
	ch <- data
}
