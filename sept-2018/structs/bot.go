package structs

type Bot struct {
	GridEntity
	DebugMode bool
	Claims    []string
	Score     uint64
}

type BotStatus struct {
	Actions   []string
	Claims    []string
	DebugMode bool
	Id        string
	Location  GridLocation
	Score     uint64
}

// GetStatus returns some basic information about this bot, including Location, Claims, Actions, and Score
// TODO: Actions need to be a real thing
func (bot Bot) GetStatus() BotStatus {
	return BotStatus{
		Actions:   []string{},
		Claims:    bot.Claims,
		DebugMode: bot.DebugMode,
		Id:        bot.Id,
		Location:  bot.Location,
		Score:     bot.Score,
	}
}
