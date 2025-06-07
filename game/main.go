package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func (p *player) moveFunction(direction string) string{
	var curRoom string
	var currentRoom room
	var res string

	getRoomByName(&curRoom, &currentRoom)

	direction, ok := currentRoom.Exits[direction]
	if !ok {
		return "Нет пути в эту комнату"
	}

	for _, room := range game.Rooms {
		if room.Name == direction {
			if room.EnterCondition != "" {
				res =  "Дверь закрыта!"
			} else {
				p.Location = direction
				res = fmt.Sprintf("Был совершен переход между комнатами: %s -> %s", curRoom, direction)
				if direction == "домой" {
					p.Location = "коридор"
				}
			}
		}
	}

	return res
}

func (b *backpack) putBackpackOn(arg string) string {
	var res string
	
	if arg != b.Name {
		res = "Неизвестный предмет или неподходящая команда!"
		return res
	}
	
	var curRoom string
	var currentRoom room

	getRoomByName(&curRoom, &currentRoom)

	if b.Location != curRoom || b.IsEquipped {
		res = "В этой комнате нет рюкзака или он уже надет!"
	} else {
		res = "Вы надели рюкзак!"
		b.Location = ""
		b.IsEquipped = true
	}

	return res
}

func (b *backpack) takeItem(arg string) string {
	var res string

	for i := range game.Items {
		if arg != game.Items[i].Name {
			res = "Такой предмет не существует или неподходящая команда!"
		} else {
			var curRoom string
			var currentRoom room

			getRoomByName(&curRoom, &currentRoom)

			if game.Items[i].Location != curRoom {
				res = "Предмет в другой комнате либо в рюкзаке!"
			} else {
				if b.IsEquipped {
					b.ItemsInside = append(b.ItemsInside, game.Items[i].Name)
					game.Items[i].Location = ""
					res = fmt.Sprintf("Добавлено в рюкзак: %s", game.Items[i].Name)
				} else {
					res = "Сначала надо надеть рюкзак!"
				}
			}
			break
		}
	}

	return res
}

func (p *player) openDoor(arg1 string, arg2 string) string {
	var res string
	flag := false

	if game.Player.Location != "коридор" {
		res = "Открыть дверь на улицу можно только из коридора!"
		return res
	}
	
	if arg1 != "ключи" {
		res = "Этим предметом нельзя открыть дверь!"
		return res
	}
	
	if arg2 != "дверь" {
		res = "К этому предмету нельзя использовать ключи!"
		return res
	}

	for i := range game.Backpack.ItemsInside {
		if arg1 == game.Backpack.ItemsInside[i] {
			flag = true
			break
		} 
	}
		
	if flag {
		for i := range game.Rooms {
			if game.Rooms[i].Name == "улица" {
				game.Rooms[i].EnterCondition = ""
			}
		}
		res = "Дверь на улицу открыта!"
	} else {
		res = "Вы не взяли ключи!"
	}

	return res
}

func (g *gameData) checkState() string {
	var curRoom string
	var curItems []string
	var curExits []string
	var currentRoom room

	getRoomByName(&curRoom, &currentRoom)

	for _, it := range g.Items {
		if it.Location == g.Player.Location{
			curItems = append(curItems, it.Name)
		}
	}
	if g.Backpack.Location == g.Player.Location {
		curItems = append(curItems, g.Backpack.Name)
	}
	curItemsStr := strings.Join(curItems, ", ")

	for dir, dest := range currentRoom.Exits {
		for _, room := range g.Rooms {
			if room.Name == dest {
				curExits = append(curExits, dir)
			}
		}
	}
	curExitsStr := strings.Join(curExits, ", ")

	return fmt.Sprintf("Вы находитесь в комнате: %s.\nВ этой комнате есть: %s.\nИз этой комнаты можно перейти в: %s.", curRoom, curItemsStr, curExitsStr)
}

func loadGameData(filename string) (*gameData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var g gameData
	err = json.Unmarshal(data, &g)
	if err != nil {
		return nil, err
	}

	return &g, nil
}

var game *gameData

func main() {
	var err error
	game, err = loadGameData("data/game_data.json")
	if err != nil {
		panic(err)
	}

	initGame()
}

func initGame() {
	var line string
	scanner := bufio.NewScanner(os.Stdin)

	for line != "выход"{
		fmt.Print("\nВведите комманду: ")
		
		if scanner.Scan() {
			line = scanner.Text()
			fmt.Println(handleCommand(line))
		} else {
			if errLine := scanner.Err(); errLine != nil {
				fmt.Fprintln(os.Stderr, "Ошибка чтения ввода:", errLine)
			}
			break
		}
	}
}

func handleCommand(inputLine string) string {
	var res string

	parts := strings.Split(inputLine, " ")
	command := parts[0]
	arg1 := ""
	if len(parts) == 2 {
		arg1 = parts[1]
	}
	arg2 := ""
	if len(parts) == 3 {
		arg1 = parts[1]
		arg2 = parts[2]
	}

	switch command {
	case "осмотреться":
		res = game.checkState()
	case "выход":
		res = "Выход из игры..."
	case "идти":
		res = game.Player.moveFunction(arg1)
	case "надеть":
		res = game.Backpack.putBackpackOn(arg1)//рюкзак надет
	case "взять":
		res = game.Backpack.takeItem(arg1)//вещь в рюкзаке
	case "применить":
		res = game.Player.openDoor(arg1, arg2)//дверь открыта
	default:
		res = "Команда не существует!"
	}

	return res
}