package keeper_test

import (
	store "github.com/vitwit/avail-da-module/keeper"
)

func (s *TestSuite) TestCanUpdateStatusToPending() {
	testCases := []struct {
		name         string
		updateStatus bool
		status       uint32
	}{
		{
			"status bytes are nil",
			false,
			0,
		},
		{
			"update status",
			true,
			2,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			res := store.CanUpdateStatusToPending(s.store)
			s.True(res)
			if tc.updateStatus {
				err := store.UpdateBlobStatus(s.ctx, s.store, tc.status)
				s.Require().NoError(err)

				res := store.CanUpdateStatusToPending(s.store)
				s.False(res)
			}
		})
	}
}

func (s *TestSuite) TestGetStatusFromStore() {
	testCases := []struct {
		name         string
		updateStatus bool
		status       uint32
		expectOutput uint32
	}{
		{
			"status bytes are nil",
			false,
			0,
			0,
		},
		{
			"update status",
			true,
			1,
			1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			res := store.GetStatusFromStore(s.store)
			s.Require().Equal(res, uint32(0))
			if tc.updateStatus {
				err := store.UpdateBlobStatus(s.ctx, s.store, tc.status)
				s.Require().NoError(err)

				res := store.GetStatusFromStore(s.store)
				s.Require().Equal(res, tc.expectOutput)
			}
		})
	}
}

func (s *TestSuite) TestUpdateStartHeight() {
	err := store.UpdateStartHeight(s.ctx, s.store, uint64(1))
	s.Require().NoError(err)

	height := s.keeper.GetStartHeightFromStore(s.ctx)
	s.Require().Equal(height, uint64(1))
}

func (s *TestSuite) TestUpdateEndHeight() {
	err := store.UpdateEndHeight(s.ctx, s.store, uint64(10))
	s.Require().NoError(err)

	height := s.keeper.GetEndHeightFromStore(s.ctx)
	s.Require().Equal(height, uint64(10))
}

func (s *TestSuite) TestUpdateProvenHeight() {
	err := store.UpdateProvenHeight(s.ctx, s.store, uint64(5))
	s.Require().NoError(err)

	height := s.keeper.GetProvenHeightFromStore(s.ctx)
	s.Require().Equal(height, uint64(5))
}

func (s *TestSuite) TestUpdateAvailHeight() {
	err := store.UpdateAvailHeight(s.ctx, s.store, uint64(20))
	s.Require().NoError(err)

	height := s.keeper.GetAvailHeightFromStore(s.ctx)
	s.Require().Equal(height, uint64(20))
}

func (s *TestSuite) TestUpdateVotingEndHeight() {
	err := store.UpdateVotingEndHeight(s.ctx, s.store, uint64(20))
	s.Require().NoError(err)

	height := s.keeper.GetVotingEndHeightFromStore(s.ctx)
	s.Require().Equal(height, uint64(20))
}
