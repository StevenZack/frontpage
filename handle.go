package frontpage

func handleMsg(b []byte) {
	if strHandler != nil {
		strHandler(string(b))
		return
	}

}
