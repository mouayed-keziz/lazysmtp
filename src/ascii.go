package main

func GetColoredASCIIArt() string {
	cyan := "\x1b[0;36m"
	green := "\x1b[0;32m"
	yellow := "\x1b[0;33m"
	magenta := "\x1b[0;35m"
	blue := "\x1b[0;34m"
	reset := "\x1b[0m"

	art := ` ` + cyan + ` ___      _______  _______  __   __  _______  __   __  _______  _______ 
 ` + green + `|   |    |   _   ||       ||  | |  ||       ||  |_|  ||       ||       |
 ` + yellow + `|   |    |  |_|  ||____   ||  |_|  ||  _____||       ||_     _||    _  |
 ` + magenta + `|   |    |       | ____|  ||       || |_____ |       |  |   |  |   |_| |
 ` + blue + `|   |___ |       || ______||_     _||_____  ||       |  |   |  |    ___|
 ` + cyan + `|       ||   _   || |_____   |   |   _____| || ||_|| |  |   |  |   |    
 ` + green + `|_______||__| |__||_______|  |___|  |_______||_|   |_|  |___|  |___|    
 ` + yellow + `                                         ` + reset + `

` + cyan + `    SMTP Testing Tool for Developers
`
	return art
}
