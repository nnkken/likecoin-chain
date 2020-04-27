package types

type IscnPair struct {
	Id     []byte     `json:"id" yaml:"id"`
	Record IscnRecord `json:"record" yaml:"record"`
}

type GenesisState struct {
	Params      Params     `json:"params" yaml:"params"`
	Entities    []Entity   `json:"entities" yaml:"entities"`
	RightTerms  []Right    `json:"rightTerms" yaml:"rightTerms"`
	IscnRecords []IscnPair `json:"iscnRecords" yaml:"iscnRecords"`
}

func DefaultGenesisState() GenesisState {
	return GenesisState{} // TODO: default param
}

func ValidateGenesis(data GenesisState) error {
	return nil // TODO
}
