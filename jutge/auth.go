package jutge

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const authURL = "https://jutge.org/services/authentication"
const authMsg = `---
email: %s
password: %s
`

// ValidateCredentials tells whether a user is valid in the Jutge.org platform
func ValidateCredentials(email string, password string) bool {
	c := &http.Client{}
	rb := fmt.Sprintf(authMsg, email, password)
	req, err := http.NewRequest("PUT", authURL, strings.NewReader(rb))
	res, err := c.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	return strings.Contains(string(body), "YES")
}
