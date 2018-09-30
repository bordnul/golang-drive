package authorize

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var tokenPath string
var clientPath string

//Takes config and opens a link to get authorization code and returns token
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Token file missing. Please open this link for authorization code:"+
		" \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

//Saves given token information to file
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving token to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()

	if err != nil {
		return err
	}
	json.NewEncoder(f).Encode(token)

	return nil
}

//Reads token from file and returns the data
func tokenFromFile(path string) (*oauth2.Token, error) {

	//Opens token from file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	//Makes Token struct and fills it with data from the file
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err

}

//Takes client config and grabs token from file or web if token file is missing
//Returns configured client ready for use
func GetClient(config *oauth2.Config, tokenPath string, clientPath string) (*http.Client, error) {

	tok, err := tokenFromFile(tokenPath)

	if err != nil {

		tok, err = getTokenFromWeb(config)
		if err != nil {

			return nil, err
		}

		err = saveToken(tokenPath, tok)

		if err != nil {

			return nil, err
		}

	}

	return config.Client(context.Background(), tok), nil

}
