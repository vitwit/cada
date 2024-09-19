package avail

// implement this interface to interact with avail chain
type AvailDA interface {
	Submit(data []byte) (AvailBlockMetaData, error)
	IsDataAvailable(data []byte, availBlockHeight int) (bool, error)
	GetBlock(availBlockHeight int) (AvailBlock, error)
}
