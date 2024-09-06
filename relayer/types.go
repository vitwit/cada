package relayer

type BlockData struct {
	Block      int64           `json:"block_number"`
	Extrinsics []ExtrinsicData `json:"data_transactions"`
}

type ExtrinsicData struct {
	Data string `json:"data"`
}
