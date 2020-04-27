package types

const (
	ModuleName   = "iscn"
	StoreKey     = ModuleName
	QuerierRoute = ModuleName
	RouterKey    = ModuleName
)

var (
	IscnRecordKey = []byte{0x01}
	IscnCountKey  = []byte{0x02}
	EntityKey     = []byte{0x03}
	RightTermsKey = []byte{0x04}
	CidBlockKey   = []byte{0x05}
)

func GetIscnRecordKey(iscnId []byte) []byte {
	return append(IscnRecordKey, iscnId...)
}

func GetEntityKey(entityCid []byte) []byte {
	return append(EntityKey, entityCid...)
}

func GetRightTermsKey(rightTermsHash []byte) []byte {
	return append(RightTermsKey, rightTermsHash...)
}

func GetCidBlockKey(cid []byte) []byte {
	return append(CidBlockKey, cid...)
}
