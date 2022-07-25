package infra

func ShowErrorAndExit(err error) {
	panic(err.Error())
}
