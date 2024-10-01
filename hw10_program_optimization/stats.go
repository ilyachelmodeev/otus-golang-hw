package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader, domain string) (result users, err error) {
	s := bufio.NewReader(r)
	domainRegexp := regexp.MustCompile("\\." + domain)
	var line []byte

	for i := 0; ; i++ {
		line, _, err = s.ReadLine()
		if errors.Is(err, io.EOF) {
			err = nil
			break
		}
		if err != nil {
			return result, err
		}
		if !domainRegexp.Match(line) {
			continue
		}
		var user User
		if err = json.Unmarshal(line, &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainRegexp := regexp.MustCompile("\\." + domain)
	for _, user := range u {
		matched := domainRegexp.MatchString(user.Email)
		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
