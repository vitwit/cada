package avail

type BlockMetaData struct {
	BlockNumber int    `json:"block_number"`
	BlockHash   string `json:"block_hash"`
	Hash        string `json:"hash"`
	Index       int    `json:"index"`
}

type Block struct {
	Block      int64        `json:"block_number"`
	Extrinsics []Extrinsics `json:"data_transactions"`
}

type Extrinsics struct {
	Data string `json:"data"`
}
