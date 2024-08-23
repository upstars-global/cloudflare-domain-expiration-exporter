package config

import (
	"fmt"
	"os"
	"strings"
)

type AuthInfo struct {
	Registrator string
	PananameAuth
	GodaddyAuth
	DomainsList
}

type PananameAuth struct {
	Token string
}

type GodaddyAuth struct {
	Secret string
	Token  string
}

type DomainsList struct {
	List []string
}

func New() (*AuthInfo, error) {
	info := new(AuthInfo)

	if val, ok, err := searchPananame(); ok && err == nil {
		info.PananameAuth = *val
		return info, nil
	} else if ok && err != nil {
		return nil, err
	}

	if val, ok, err := searchGodaddy(); ok && err == nil {
		info.GodaddyAuth = *val
		return info, nil
	} else if ok && err != nil {
		return nil, err
	}

	if val, ok, err := searchDomainsList(); ok && err == nil {
		info.DomainsList = *val
		return info, nil
	} else if ok && err != nil {
		return nil, err
	}

	return &AuthInfo{}, nil
}

func searchPananame() (*PananameAuth, bool, error) {
	if val, ok := os.LookupEnv("PANANAME_TOKEN"); ok {
		if len(val) == 0 {
			return &PananameAuth{Token: val}, true, nil
		} else {
			return nil, true, fmt.Errorf("PANANAME_TOKEN env var found but is empty")
		}
	}
	return nil, false, nil
}

func searchGodaddy() (*GodaddyAuth, bool, error) {
	var (
		secret string
		token  string

		secretOk bool
		tokenOk  bool

		info *GodaddyAuth
	)

	if secret, secretOk = os.LookupEnv("GODADDY_SECRET"); secretOk {
		if len(secret) == 0 {
			info.Secret = secret
		} else {
			return nil, true, fmt.Errorf("GODADDY_SECRET env var found but is empty")
		}
	}

	if token, tokenOk = os.LookupEnv("GODADDY_TOKEN"); tokenOk {
		if len(token) == 0 {
			info.Token = token
		} else {
			return nil, true, fmt.Errorf("GODADDY_TOKEN env var found but is empty")
		}
	}

	if secretOk && tokenOk {
		return info, true, nil
	}

	if secretOk && !tokenOk {
		return nil, true, fmt.Errorf("can't found GODADDY_TOKEN variable")
	}

	if !secretOk && tokenOk {
		return nil, true, fmt.Errorf("can't found GODADDY_SECRET variable")
	}

	return nil, false, nil
}

func searchDomainsList() (*DomainsList, bool, error) {
	dList := new(DomainsList)
	if val, ok := os.LookupEnv("DOMAINS_LIST"); ok && len(val) > 0 {
		lst := strings.Split(val, ",")
		if len(lst) >= 0 {
			dList.List = lst
			return dList, true, nil
		} else {
			return nil, true, fmt.Errorf("DOMAINS_LIST env var found but is empty")
		}
	}
	return nil, false, nil
}
