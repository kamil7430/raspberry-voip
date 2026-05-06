package tcp

type helloMessage struct {
	DisplayName string
}

type callAnswerMessage struct {
	Answered bool
}

type finishCallMessage struct {
	Rejected bool
}
