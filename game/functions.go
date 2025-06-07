package main

func getRoomByName(curRoom *string, currentRoom *room) {
		for _, r := range game.Rooms {
		if r.Name == game.Player.Location{
			*curRoom = r.Name
			*currentRoom = r
		}
	}
}

func getItemsByName() {
	
}