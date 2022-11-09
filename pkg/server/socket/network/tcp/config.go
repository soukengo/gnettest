package tcp

type Config struct {
	Address   string
	Multicore bool
	ReadBuf   int
	SendBuf   int
}
