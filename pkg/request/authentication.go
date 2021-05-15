package request

import (
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/inheaden/simple-mail-api/pkg/utils"
)

// RatelimitIP limites the requests from an ip to RATE_LIMIT_SECONDS
// will return false if limit has been reached
func RatelimitIP(r *http.Request) (bool, error) {
	ip := getIP(r)

	ok := true
	ipModel, err := database.GetIP(ip)
	if err != nil {
		return false, err
	}

	rateLimitNano := generalConfig.RateLimitSeconds * int(time.Second)
	if ipModel != nil && ipModel.LastCallTime.After(time.Now().Add(-time.Duration(rateLimitNano))) {
		log.Debugf("Ratelimited %s", ip)
		ok = false
	}

	err = database.SaveIP(ip)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func getIP(r *http.Request) string {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if generalConfig.IPHeader != "" {
		ip = r.Header.Get(generalConfig.IPHeader)
	}
	return ip
}

// VerifyNonce verifies if the hash given is valid
// it will retrieve the nonce for this ip and verify the hash
func VerifyNonce(r *http.Request, nonce string, hash string) (bool, error) {
	ip := getIP(r)

	nonceModel, err := database.GetNonce(nonce)
	if nonceModel == nil || err != nil || nonceModel.IP != ip {
		return false, err
	}

	correctHash := utils.ValidateNonce(nonceModel.Nonce, hash, generalConfig.Secret)

	return correctHash, database.DeleteNonceByNonce(nonce)
}

// GetNonce creates and saves a nonce for the given ip
// it will override the old one
func GetNonce(r *http.Request) (string, error) {
	ip := getIP(r)
	nonce := utils.GetRandomString(10)

	err := database.SaveNonce(nonce, ip)
	return nonce, err
}
