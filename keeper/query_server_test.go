package keeper_test

import (
	store "github.com/vitwit/avail-da-module/keeper"
	"github.com/vitwit/avail-da-module/types"
)

func (s *TestSuite) TestSubmitBlobStatus() {

	testCases := []struct {
		name string

		req          types.QuerySubmitBlobStatusRequest
		status       uint32
		expectOutput string
	}{
		{
			"get blobstatus",
			types.QuerySubmitBlobStatusRequest{},
			2,
			"IN_VOTING",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {

			err := store.UpdateBlobStatus(s.ctx, s.store, tc.status)
			s.Require().NoError(err)

			res, err := s.keeper.SubmitBlobStatus(s.ctx, &tc.req)
			s.Require().NoError(err)
			s.Require().NotNil(res)
			s.Require().Equal(res.Status, tc.expectOutput)
		})
	}
}
