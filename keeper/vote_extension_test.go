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

func (s *TestSuite) TestVerifyVoteExtensionHandler() {
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
		req           *abci.RequestVerifyVoteExtension
		expectedError string
	}{
		{
			name:          "nil request",
			startHeight:   1,
			endHeight:     20,
			availHeight:   30,
			blobStatus:    3,
			currentHeight: 22,
			voteEndHeight: 20,
			isLastVoting:  false,
			expectErr:     true,
			req:           nil,
			expectedError: "request is nil",
		},
		{
			name:          "invalid vote height",
			startHeight:   1,
			endHeight:     20,
			availHeight:   30,
			blobStatus:    2,
			currentHeight: 30,
			voteEndHeight: 31,
			isLastVoting:  true,
			expectErr:     true,
			req: &abci.RequestVerifyVoteExtension{
				Height:        -1,
				VoteExtension: []byte("valid extension"),
			},
			expectedError: "invalid vote height: -1",
		},
		{
			name:          "empty vote extension",
			startHeight:   1,
			endHeight:     20,
			availHeight:   30,
			blobStatus:    2,
			currentHeight: 30,
			voteEndHeight: 31,
			isLastVoting:  true,
			expectErr:     true,
			req: &abci.RequestVerifyVoteExtension{
				Height:        1,
				VoteExtension: []byte{},
			},
			expectedError: "vote extension data is empty",
		},
		{
			name:          "valid request",
			startHeight:   1,
			endHeight:     20,
			availHeight:   30,
			blobStatus:    2,
			currentHeight: 30,
			voteEndHeight: 31,
			isLastVoting:  true,
			expectErr:     false,
			req: &abci.RequestVerifyVoteExtension{
				Height:        1,
				VoteExtension: []byte("valid extension"),
			},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Set up block height and blob status
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

			verifyVoteExtensionHandler := s.voteExtensionHandler.VerifyVoteExtensionHandler()

			// Handle nil request
			var res *abci.ResponseVerifyVoteExtension
			if tc.req == nil {
				res, err = verifyVoteExtensionHandler(s.ctx, nil)
			} else {
				res, err = verifyVoteExtensionHandler(s.ctx, tc.req)
			}

			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedError)
				s.Require().Equal(abci.ResponseVerifyVoteExtension_REJECT, res.Status)
			} else {
				s.Require().NoError(err)
				s.Require().NotNil(res)
				s.Require().Equal(abci.ResponseVerifyVoteExtension_ACCEPT, res.Status)
			}
		})
	}
}
