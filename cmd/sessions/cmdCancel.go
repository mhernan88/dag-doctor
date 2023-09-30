package sessions

func(sm SessionManager) cancelSession(ID string) error{
	err := sm.UpdateSessionStatus(ID, "cancelled")
	if err != nil {
		return err
	}
	return nil
}
