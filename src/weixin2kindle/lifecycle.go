package main

func init() {
	initC()
	initDb()
}

func destroy() {
	CloseDb()
}
