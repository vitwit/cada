package avail

// implement this interface to interact with avail chain
type DA interface {
	Submit(data []byte) (BlockMetaData, error)
	IsDataAvailable(data []byte, availBlockHeight int) (bool, error)
	GetBlock(availBlockHeight int) (Block, error)
}
