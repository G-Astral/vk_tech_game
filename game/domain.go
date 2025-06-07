package main

type player struct {
	Location 		string 				`json:"location"`
}

type item struct {
	Name 			string 				`json:"name"`
	Location 		string 				`json:"location"`
}

type backpack struct {
	item
	IsEquipped 		bool 				`json:"isEquipped"`
	ItemsInside 	[]string 			`json:"itemsInside"`
}

type room struct {
	Name 			string 				`json:"name"`
	EnterCondition 	string 				`json:"enterCondition"`
	Items 			[]string 			`json:"items"`
	Exits 			map[string]string 	`json:"exits"`
}

type gameData struct {
	Player 			player 				`json:"player"`
	Items 			[]item 				`json:"items"`
	Backpack 		backpack 			`json:"backpack"`
	Rooms 			[]room 				`json:"rooms"`
}