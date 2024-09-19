package keeper_test

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	store "github.com/vitwit/avail-da-module/keeper"
)

func (s *TestSuite) TestPrepareProposal() {
	s.ctx = s.ctx.WithBlockHeight(int64(10))

	testCases := []struct {
		name          string
		req           *abci.RequestPrepareProposal
		voteEndHeight uint64
		status        uint32
		expectErr     bool
	}{
		{
			"preapre proposal txs",
			&abci.RequestPrepareProposal{
				ProposerAddress: s.addrs[1],
			},
			uint64(10),
			uint32(2),
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			res, err := s.proofofBlobProposerHandler.PrepareProposal(s.ctx, tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			}
		})
	}
}

func (s *TestSuite) TestProcessProposal() {
	testCases := []struct {
		name       string
		addmockTxs bool
		req        *abci.RequestProcessProposal
	}{
		{
			"process proposal txs",
			true,
			&abci.RequestProcessProposal{
				ProposerAddress: s.addrs[1].Bytes(),
				Height:          10,
				Txs:             [][]byte{s.getMockTx(), s.getMockTx()},
			},
		},
		{
			"txs are nil",
			false,
			&abci.RequestProcessProposal{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			res, err := s.proofofBlobProposerHandler.ProcessProposal(s.ctx, tc.req)
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *TestSuite) getMockTx() []byte {
	msg := banktypes.NewMsgSend(s.addrs[0], s.addrs[1], sdk.NewCoins(sdk.NewInt64Coin("stake", 100)))

	txBuilder := s.encCfg.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msg)
	s.Require().NoError(err)

	txBytes, err := s.encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	s.Require().NoError(err)

	return txBytes
}

func (s *TestSuite) TestPreBlocker() {
	testCases := []struct {
		name          string
		req           *abci.RequestFinalizeBlock
		blobStatus    uint32
		voteEndHeight uint64
		provenHeight  uint64
		currentHeight int64
	}{
		{
			"process preblocker",
			&abci.RequestFinalizeBlock{
				ProposerAddress: s.addrs[1],
				Txs:             [][]byte{s.getMockTx(), s.getMockTx()},
			},
			2,
			20,
			10,
			20,
		},
		{
			"process preblocker",
			&abci.RequestFinalizeBlock{
				ProposerAddress: s.addrs[1],
				Txs:             [][]byte{s.getMockTx(), s.getMockTx()},
			},
			0,
			20,
			10,
			6,
		},
		{
			"process preblocker",
			&abci.RequestFinalizeBlock{
				ProposerAddress: s.addrs[1],
				Txs:             [][]byte{s.getMockTx(), s.getMockTx()},
			},
			0,
			20,
			5,
			6,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ctx = s.ctx.WithBlockHeight(tc.currentHeight)
			s.keeper.ProposerAddress = s.addrs[1]

			fmt.Printf("s.relayer.AvailConfig: %v\n", s.relayer.AvailConfig)
			//s.relayer.AvailConfig.MaxBlobBlocks = 5

			err := store.UpdateBlobStatus(s.ctx, s.store, tc.blobStatus)
			s.Require().NoError(err)

			err = store.UpdateVotingEndHeight(s.ctx, s.store, tc.voteEndHeight, true)
			s.Require().NoError(err)

			err = store.UpdateProvenHeight(s.ctx, s.store, tc.provenHeight)
			s.Require().NoError(err)

			err = s.keeper.PreBlocker(s.ctx, tc.req)
			s.Require().NoError(err)
		})
	}
}
