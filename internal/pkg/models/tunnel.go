package models

import (
	"fmt"

	"github.com/breathbath/go_utils/utils/testing"
)

type OperationStatus struct {
	Status string `json:"status"`
}

func (os *OperationStatus) KeyValues() []testing.KeyValueStr {
	return []testing.KeyValueStr{
		{
			Key:   "Status",
			Value: os.Status,
		},
	}
}

type Tunnel struct {
	ID          string `json:"id"`
	Client      string `json:"client"`
	Lhost       string `json:"lhost"`
	Lport       string `json:"lport"`
	Rhost       string `json:"rhost"`
	Rport       string `json:"rport"`
	LportRandom bool   `json:"lport_random"`
	Scheme      string `json:"scheme"`
	ACL         string `json:"acl"`
}

func (t *Tunnel) Headers() []string {
	return []string{
		"ID",
		"CLIENT",
		"LHOST",
		"LPORT",
		"RHOST",
		"RPORT",
		"LPORTRAND",
		"SCHEME",
		"ACL",
	}
}

func (t *Tunnel) Row() []string {
	return []string{
		t.ID,
		t.Client,
		t.Lhost,
		t.Lport,
		t.Rhost,
		t.Rport,
		fmt.Sprint(t.LportRandom),
		t.Scheme,
		t.ACL,
	}
}

func (t *Tunnel) KeyValues() []testing.KeyValueStr {
	return []testing.KeyValueStr{
		{
			Key:   "ID",
			Value: t.ID,
		},
		{
			Key:   "CLIENT",
			Value: t.Client,
		},
		{
			Key:   "LHOST",
			Value: t.Lhost,
		},
		{
			Key:   "LPORT",
			Value: t.Lport,
		},
		{
			Key:   "RHOST",
			Value: t.Rhost,
		},
		{
			Key:   "RPORT",
			Value: t.Rport,
		},
		{
			Key:   "LPORT RANDOM",
			Value: fmt.Sprintf("%v", t.LportRandom),
		},
		{
			Key:   "SCHEME",
			Value: t.Scheme,
		},
		{
			Key:   "ACL",
			Value: t.ACL,
		},
	}
}
