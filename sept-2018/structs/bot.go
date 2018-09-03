package structs

type Bot struct {
	GridEntity
	Claims []string
}

type BotStatus struct {
	Actions  []string
	Claims   []string
	Id       string
	Location GridLocation
	Score    uint64
}

// GetStatus returns some basic information about this bot, including Location, Claims, Actions, and Score
// TODO: Actions need to be a real thing
// TODO: Score needs to be a real thing
func (bot Bot) GetStatus() BotStatus {
	return BotStatus{
		Actions:  []string{},
		Claims:   bot.Claims,
		Id:       bot.Id,
		Location: bot.Location,
		Score:    0,
	}
}
