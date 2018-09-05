package structs

type Bot struct {
	GridEntity
	DebugMode bool
	Claims    []string
	Score     int
}

type BotStatus struct {
	Claims   []string
	Id       string
	Location GridLocation
	Score    int
}

// GetStatus returns some basic information about this bot, including Location, Claims, Actions, and Score
// TODO: Actions need to be a real thing
func (bot Bot) GetStatus() BotStatus {
	return BotStatus{
		Claims:   bot.Claims,
		Id:       bot.Id,
		Location: bot.Location,
		Score:    bot.Score,
	}
}
