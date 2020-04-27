package types

const (
	QueryIscnRecord      = "records"
	QueryEntity          = "entity"
	QueryParams          = "params"
	QueryCidBlockGet     = "cid_get"
	QueryCidBlockGetSize = "cid_get_size"
	QueryCidBlockHas     = "cid_has"
)

type QueryEntityParams struct {
	Cid []byte
}

type QueryRightTermsParams struct {
	Cid []byte
}

type QueryRecordParams struct {
	Id []byte // TODO: string?
}

type QueryCidParams struct {
	Cid []byte
}
