package internal

const (
	Alphabet           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers            = "1234567890"
	AlphaNumaric       = Alphabet + Numbers
	WordCharacters     = AlphaNumaric + "_"
	FileNameCharacters = WordCharacters + " -."
)
