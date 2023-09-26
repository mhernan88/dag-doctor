package sessions

import (
	"fmt"

	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/mhernan88/dag-bisect/db/models"
)

func (sm SessionManager) cleanup(ID string, sess *models.Session, ui *cmd.UI) error {
	var err error
	if (len(ui.DAG.Nodes) == 0) || (len(ui.DAG.Roots) == 0) {
		ui.Terminate()
		sm.UpdateErrNode(ID, ui.LastFailedNode)
		if ui.LastFailedNode == "" {
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
		err = cmd.SaveState(sess.State, *ui)
		if err != nil {
			return fmt.Errorf("failed to save state | %v", err)
		}
		fmt.Printf("successfully iterated session %s\n", ID)
	}
	return nil
}
