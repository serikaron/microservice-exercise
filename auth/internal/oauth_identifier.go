package internal

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"mse/pkg"

	"golang.org/x/oauth2"
)

func IdentifyWithOAuth(code string) (*pkg.Identity, error) {
	conf := &oauth2.Config{
		RedirectURL:  "http://localhost:38080/get_code",
		ClientID:     "hblZUghd1RBvAnvjciTMP-NkVVP1SHIzYu2NA4esL-8",
		ClientSecret: "DBqIdw6pDU5DoyImNsJX7PMTLgXRWtnuQJiWcsP03yk",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://unsplash.com/oauth/token",
		},
	}

	t, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	c := conf.Client(context.Background(), t)

	rsp, err := c.Get("https://api.unsplash.com/me")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rsp.Body.Close()
	}()

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	id := pkg.Identity{}
	err = json.Unmarshal(data, &id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
