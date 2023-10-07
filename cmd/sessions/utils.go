package sessions

import (
	"fmt"

	"github.com/mhernan88/dag-bisect/resolver/interactive"
	"github.com/mhernan88/dag-bisect/models"
)

func (sm SessionManager) cleanup(ID string, sess *models.Session, state *models.State) error {
	var err error
	if (len(state.DAG.Nodes) == 0) || (len(state.DAG.Roots) == 0) {
		interactive.Terminate(state)
		sm.UpdateErrNode(ID, state.LastFailedNode)
		if state.LastFailedNode == "" {
			err = sm.UpdateSessionStatus(ID, "ok")
			if err != nil {
				return fmt.Errorf("failed to update session status | %v", err)
			}
		} else {
			err = sm.UpdateSessionStatus(ID, "err")
			if err != nil {
				return fmt.Errorf("failed to update session status | %v", err)
			}
		}
	} else {
		err = models.SaveState(sess.State, state)
		if err != nil {
			return fmt.Errorf("failed to save state | %v", err)
		}
		fmt.Printf("successfully iterated session %s\n", ID)
	}
	return nil
}

