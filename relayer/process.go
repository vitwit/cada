package relayer

import (
	"context"
	"time"
)

// Start begins the relayer process
func (r *Relayer) Start() error {
	ctx := context.Background()

	timer := time.NewTimer(r.pollInterval)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case height := <-r.commitHeights:
			r.latestCommitHeight = height
		case height := <-r.provenHeights:
			r.updateHeight(height)
		case <-timer.C:
			// TODO: client update
		}
	}
}

// NotifyCommitHeight is called by the app to notify the relayer of the latest commit height
func (r *Relayer) NotifyCommitHeight(height int64) {
	r.commitHeights <- height
}

// NotifyProvenHeight is called by the app to notify the relayer of the latest proven height
// i.e. the height of the highest incremental block that was proven to be posted to Celestia.
func (r *Relayer) NotifyProvenHeight(height int64) {
	r.provenHeights <- height
}

// updateHeight is called when the provenHeight has changed
func (r *Relayer) updateHeight(height int64) {
	if height > r.latestProvenHeight {
		// fmt.Println("Latest proven height:", height) // TODO: remove, debug only
		r.latestProvenHeight = height
		r.pruneCache(height)
	}
}

// pruneCache will delete any headers or proofs that are no longer needed
func (r *Relayer) pruneCache(provenHeight int64) {
	r.mu.Lock()
	// TODO: proofs deletions after completion
	r.mu.Unlock()
}

// // only validators need to check for new block proofs and update clients
// // if a validator, they should not be querying celestia until they finished catching up
// func (r *Relayer) isValidatorAndCaughtUp(ctx context.Context) bool {
// 	status, err := r.localProvider.Status(ctx)
// 	if err != nil {
// 		return false
// 	}
// 	if status.ValidatorInfo.VotingPower > 0 && !status.SyncInfo.CatchingUp {
// 		return true
// 	}
// 	return false
// }

// func (r *Relayer) getClientState(ctx context.Context) *celestia.ClientState {

// 	// TODO: client
// 	// clientState, err := r.localProvider.QueryCelestiaClientState(ctx)
// 	// if err != nil {
// 	// 	return nil
// 	// }

// 	// if clientState.LatestHeight.RevisionHeight == 0 {
// 	// 	return nil
// 	// }

// 	// // update celestia's last queried height to avoid unnecessary queries
// 	// if r.celestiaLastQueriedHeight < int64(clientState.LatestHeight.RevisionHeight) {
// 	// 	r.celestiaLastQueriedHeight = int64(clientState.LatestHeight.RevisionHeight)
// 	// }

// 	// return clientState
// 	return nil
// }
