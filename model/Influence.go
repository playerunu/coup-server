package models

type Influence int

const (
	Duke Influence = iota
	Captain
	Assassin
	Contessa
	Ambassador
	length
)

func InfluenceToStr(influence Influence) string {
	var influenceStr string

	switch influence {
	case Duke:
		influenceStr = "Duke"
	case Captain:
		influenceStr = "Captain"
	case Assassin:
		influenceStr = "Assassin"
	case Contessa:
		influenceStr = "Contessa"
	case Ambassador:
		influenceStr = "Ambassador"
	}

	return influenceStr
}

func StrToInfluence(influenceStr string) Influence {
	var influence Influence

	switch influenceStr {
	case "Duke":
		influence = Duke
	case "Captain":
		influence = Captain
	case "Assassin":
		influence = Assassin
	case "Contessa":
		influence = Contessa
	case "Ambassador":
		influence = Ambassador
	}

	return influence
}
