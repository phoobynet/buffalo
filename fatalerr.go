package main

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
