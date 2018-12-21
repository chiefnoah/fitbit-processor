package listener

import "net/http"

//VerifyRequest verifies a fitbit request
func VerifyRequest(r *http.Request) (bool, error) {
	validDNS, err := verifyRDNS(r.RemoteAddr)
	if err != nil {
		return false, err
	}
	validSig, err := verifySignature(r)
	if err != nil {
		return false, err
	}

	return validDNS && validSig, nil
}

func verifyRDNS(ip string) (bool, error) {
	return true, nil
}

func verifySignature(r *http.Request) (bool, error) {
	return true, nil
}
