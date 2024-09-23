package keeper_test

import (
	abci "github.com/cometbft/cometbft/abci/types"
	store "github.com/vitwit/avail-da-module/keeper"
)

func (s *TestSuite) TestExtendVoteHandler() {
	testCases := []struct {
		name          string
		startHeight   uint64
		endHeight     uint64
		availHeight   uint64
		blobStatus    uint32
		currentHeight int64
		voteEndHeight uint64
		isLastVoting  bool
		expectErr     bool
	}{
		{
			"blob status is not in voting state",
			uint64(1),
			uint64(20),
			uint64(30),
			uint32(3),
			int64(22),
			uint64(20),
			false,
			false,
		},
		{
			"blob status is in voting state",
			uint64(1),
			uint64(20),
			uint64(30),
			uint32(2),
			int64(30),
			uint64(31),
			true,
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.ctx = s.ctx.WithBlockHeight(tc.currentHeight)

			err := store.UpdateStartHeight(s.ctx, s.store, tc.startHeight)
			s.Require().NoError(err)

			err = store.UpdateEndHeight(s.ctx, s.store, tc.endHeight)
			s.Require().NoError(err)

			err = store.UpdateAvailHeight(s.ctx, s.store, tc.availHeight)
			s.Require().NoError(err)

			err = store.UpdateBlobStatus(s.ctx, s.store, tc.blobStatus)
			s.Require().NoError(err)

			err = store.UpdateVotingEndHeight(s.ctx, s.store, tc.voteEndHeight, tc.isLastVoting)
			s.Require().NoError(err)

			extendVoteHandler := s.voteExtensionHandler.ExtendVoteHandler()

			req := &abci.RequestExtendVote{
				ProposerAddress: s.addrs[1],
			}

			res, err := extendVoteHandler(s.ctx, req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(res)
			}
		})
	}
}
