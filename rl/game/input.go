package game

func GetPlayerAction(key string, player *Actor) Action {
	switch key {
	case ".":
		return WaitAction{}
	/* movement */
	case "up", "k", "down", "j", "left", "h", "right", "l", "y", "u", "b", "n":
		var pos Vec2
		switch key {
		case "up", "k":
			pos = Vec2{0, -1}
		case "down", "j":
			pos = Vec2{0, 1}
		case "left", "h":
			pos = Vec2{-1, 0}
		case "right", "l":
			pos = Vec2{1, 0}
		case "y":
			pos = Vec2{-1, -1}
		case "u":
			pos = Vec2{1, -1}
		case "b":
			pos = Vec2{-1, 1}
		case "n":
			pos = Vec2{1, 1}
		}
		return MoveAction{Actor: player, Target: pos}

	default:
		return InvalidKeyAction{} // do not process turn
	}
}
