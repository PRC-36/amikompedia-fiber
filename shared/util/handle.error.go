package util

func ReturnIfError(err error) {
	if err != nil {
		panic(err)
	}
}
