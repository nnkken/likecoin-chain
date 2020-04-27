package types

import (
	"encoding/json"
)

// These are for internal usage, not for transactions.
// Transactions will contain records which are more human readable, and could be transformed into structures.

// TODO: better name?
type Entity struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	// TODO: custom fields
}

type CID = []byte

func (entity Entity) String() string {
	bz, err := json.Marshal(entity)
	if err != nil {
		panic(err)
	}
	return string(bz)
}

type Stakeholder struct {
	Type   string `json:"type" yaml:"type"`
	Entity CID    `json:"entity" yaml:"entity"`
	Stake  uint32 `json:"stake" yaml:"stake"`
}

type Period struct {
	From int64 `json:"from" yaml:"from"`
	To   int64 `json:"to" yaml:"to"`
}

func (p Period) String() string {
	return ""
}

type Right struct {
	Type      string `json:"type" yaml:"type"`
	Holder    CID    `json:"holder" yaml:"holder"`
	Terms     CID    `json:"terms" yaml:"terms"`
	Period    Period `json:"period" yaml:"period"`
	Territory string `json:"territory" yaml:"territory"`
}

type RightTerms struct {
	Content string `json:"content" yaml:"content"`
}

type IscnContent struct {
	Type        string   `json:"type" yaml:"type"`
	Source      string   `json:"source" yaml:"source"`
	Fingerprint string   `json:"fingerprint" yaml:"fingerprint"`
	Feature     string   `json:"feature" yaml:"feature"`
	Edition     string   `json:"edition" yaml:"edition"`
	Title       string   `json:"title" yaml:"title"`
	Description string   `json:"description" yaml:"description"`
	Tags        []string `json:"tags" yaml:"tags"`
	// TODO: allow custom fields
}

type IscnRecord struct {
	Stakeholders []Stakeholder `json:"stakeholders" yaml:"stakeholders"`
	Timestamp    int64         `json:"timestamp" yaml:"timestamp"`
	Parent       CID           `json:"parent" yaml:"parent"`
	Version      uint32        `json:"version" yaml:"version"`
	Rights       []Right       `json:"rights" yaml:"rights"`
	Content      IscnContent   `json:"content" yaml:"content"`
}

func (iscnRecord IscnRecord) String() string {
	// TODO: timestamp should be ISO-8601
	bz, err := json.Marshal(iscnRecord)
	if err != nil {
		panic(err)
	}
	return string(bz)
}
