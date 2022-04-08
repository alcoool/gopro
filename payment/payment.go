package payment

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jarcoal/httpmock"
)

func Do(amount int) error {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"http://payment-api",
		func(request *http.Request) (*http.Response, error) {
			time.Sleep(2 * time.Second)

			now := time.Now()

			if now.UnixNano()%2 == 0 {
				return httpmock.NewStringResponse(http.StatusOK, `{"success": true}`), nil
			} else {
				return httpmock.NewStringResponse(http.StatusInternalServerError, `{"success": false}`), nil
			}
		},
	)

	resp, err := http.Get(fmt.Sprintf("http://payment-api?ampunt=%d", amount))
	if err != nil {
		log.Fatalln(err)
	}

	if resp.Status != strconv.Itoa(http.StatusOK) {
		return errors.New("unable to process payment")
	}

	return nil
}
