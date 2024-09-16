package keeper_test

func (s *TestSuite) TestAddAvailblobDataToTxs() {
	// Define test cases
	testCases := []struct {
		name         string
		injectDataBz []byte
		maxTxBytes   int64
		txs          [][]byte
		expectedTxs  [][]byte
	}{
		{
			name:         "No injected data, returns original transactions",
			injectDataBz: nil,
			maxTxBytes:   200,
			txs:          [][]byte{[]byte("tx1"), []byte("tx2")},
			expectedTxs:  [][]byte{[]byte("tx1"), []byte("tx2")},
		},
		{
			name:         "Injected data is empty, returns original transactions",
			injectDataBz: []byte{},
			maxTxBytes:   200,
			txs:          [][]byte{[]byte("tx1"), []byte("tx2")},
			expectedTxs:  [][]byte{[]byte("tx1"), []byte("tx2")},
		},
		{
			name:         "Injected data fits within maxTxBytes limit",
			injectDataBz: []byte("injectedData"),
			maxTxBytes:   2000,
			txs:          [][]byte{[]byte("tx1"), []byte("tx2")},
			expectedTxs:  [][]byte{[]byte("injectedData"), []byte("tx1"), []byte("tx2")},
		},
		// {
		// 	name:         "Injected data does not allow all transactions to fit",
		// 	injectDataBz: []byte("injectedData"),
		// 	maxTxBytes:   20,
		// 	txs:          [][]byte{[]byte("tx1"), []byte("tx2")},
		// 	expectedTxs:  [][]byte{[]byte("injectedData"), []byte("tx1")},
		// },
		{
			name:         "All transactions fit exactly into maxTxBytes",
			injectDataBz: []byte("injectedData"),
			maxTxBytes:   27,
			txs:          [][]byte{[]byte("tx1")},
			expectedTxs:  [][]byte{[]byte("injectedData"), []byte("tx1")},
		},
	}

	// Loop through each test case
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resultTxs := s.keeper.AddAvailblobDataToTxs(tc.injectDataBz, tc.maxTxBytes, tc.txs)

			s.Require().Equal(len(tc.expectedTxs), len(resultTxs), "The length of the result transactions is not as expected")

			for i := range tc.expectedTxs {
				s.Require().Equal(tc.expectedTxs[i], resultTxs[i], "Transaction at index %d does not match the expected output", i)
			}
		})
	}
}