package language

func Init() error {
	if err := initCpp(); err != nil {
		return err
	}
	if err := initPython(); err != nil {
		return err
	}
	return nil
}
