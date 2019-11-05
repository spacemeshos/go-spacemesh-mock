package utils

import (
	"bytes"
	"os/exec"
	"regexp"
)

const (
	numBlock = "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
	ipPattern = numBlock + "\\." + numBlock + "\\." + numBlock + "\\." + numBlock
	portPat = `\d{1,5}`
	fullAddressPat = ipPattern + ":" + portPat
	fullLocalhostAddPat = "localhost:" + portPat
)

// validate full address is of format "ip:port"
func ValidateFullAddress(address string) bool {
	fullAddReg := regexp.MustCompile(fullAddressPat)
	isMatch := fullAddReg.MatchString(address)

	fullLocalAddReg := regexp.MustCompile(fullLocalhostAddPat)
	isMatchLocal := fullLocalAddReg.MatchString(address)

	if isMatch || isMatchLocal {
		return true
	}

	return false
}

func FindExecPath(process string) (string, error) {
	cmd := exec.Command("which",  process)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

