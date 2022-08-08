package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Container struct {
	Name string
	List []string
}

type Player struct {
	Backpack bool
	Inv      map[string]bool
	Place    *Location
}

type Location struct {
	Name   string
	Conts  []*Container
	Neighs []*Location
}

var player = new(Player)
var Door bool

func (player *Player) Inspect() string {
	var ans string
	empty := true
	switch player.Place.Name {
	case "кухня":
		ans += "ты находишься на кухне, "
		for _, cont := range player.Place.Conts {
			if len(cont.List) != 0 {
				empty = false
				ans += "на " + cont.Name + "е: "
				for _, item := range cont.List {
					ans += item + ", "
				}
			}
		}
		if empty {
			ans += "здесь пусто, "
		}
		ans += "надо "
		if !(player.Inv["ключи"] && player.Inv["конспекты"]) {
			ans += "собрать рюкзак и "
		}
		ans += "идти в универ. "
	case "комната":
		for i, cont := range player.Place.Conts {
			if len(cont.List) != 0 {
				empty = false
				if i > 0 {
					ans += ", "
				}
				ans += "на " + cont.Name + "е: "
				for idx, item := range cont.List {
					if idx > 0 {
						ans += ", " + item
					} else {
						ans += item
					}
				}
			}
		}
		if empty {
			ans += "пустая комната. "
		} else {
			ans += ". "
		}
	}
	ans += "можно пройти - "
	for idx, neigh := range player.Place.Neighs {
		ans += neigh.Name
		if idx < len(player.Place.Neighs)-1 {
			ans += ", "
		}
	}
	return ans
}

func (player *Player) Move(loc string) string {
	var ans string
	path := false
	var dest *Location
	for _, neigh := range player.Place.Neighs {
		if loc == neigh.Name {
			dest = neigh
			path = true
			break
		}
	}
	if !path {
		return "нет пути в " + loc
	}
	switch loc {
	case "улица":
		if Door {
			return "на улице весна. можно пройти - домой"
		} else {
			return "дверь закрыта"
		}
	case "комната":
		ans += "ты в своей комнате. "
	case "кухня":
		ans += "кухня, ничего интересного. "
	case "коридор":
		ans += "ничего интересного. "
	}
	player.Place = dest
	ans += "можно пройти - "
	for idx, neigh := range player.Place.Neighs {
		ans += neigh.Name
		if idx < len(player.Place.Neighs)-1 {
			ans += ", "
		}
	}
	return ans
}

func (player *Player) Take(takeable string) string {
	ans := "нет такого"
	if !player.Backpack {
		return "некуда класть"
	}
	for idx, cont := range player.Place.Conts {
		var list []string
		here := false
		for _, item := range cont.List {
			if item == takeable {
				here = true
			} else {
				list = append(list, item)
			}
		}
		player.Place.Conts[idx].List = list
		if !player.Inv[takeable] && here {
			player.Inv[takeable] = true
			return "предмет добавлен в инвентарь: " + takeable
		}
	}
	return ans
}

func (player *Player) Use(item, object string) string {
	if !player.Inv[item] {
		return "нет предмета в инвентаре - " + item
	}
	if item == "ключи" && object == "дверь" && player.Place.Name == "коридор" {
		Door = true
		return "дверь открыта"
	}
	return "не к чему применить"
}

func (player *Player) PutOn(bp string) string {
	ans := "нет такого"
	for idx, cont := range player.Place.Conts {
		var list []string
		for _, item := range cont.List {
			if item == bp {
				player.Backpack = true
			} else {
				list = append(list, item)
			}
		}
		player.Place.Conts[idx].List = list
		if player.Backpack {
			return "вы надели: " + bp
		}
	}
	return ans
}

func initGame() {
	var roomTable = new(Container)
	var kitchenTable = new(Container)
	var chair = new(Container)
	var kitchen = new(Location)
	var room = new(Location)
	var hall = new(Location)
	var street = new(Location)
	roomTable.Name = "стол"
	kitchenTable.Name = "стол"
	chair.Name = "стул"
	kitchenTable.List = append(kitchenTable.List, "чай")
	roomTable.List = append(roomTable.List, "ключи", "конспекты")
	chair.List = append(chair.List, "рюкзак")
	room.Name = "комната"
	kitchen.Name = "кухня"
	hall.Name = "коридор"
	street.Name = "улица"
	kitchen.Conts = append(kitchen.Conts, kitchenTable)
	room.Conts = append(room.Conts, roomTable, chair)
	kitchen.Neighs = append(kitchen.Neighs, hall)
	room.Neighs = append(room.Neighs, hall)
	hall.Neighs = append(hall.Neighs, kitchen, room, street)
	street.Neighs = append(street.Neighs, hall)
	player.Backpack = false
	player.Inv = map[string]bool{
		"ключи":     false,
		"конспекты": false,
		"чай":       false,
	}
	player.Place = kitchen
	Door = false
}

func handleCommand(command string) string {
	var orders = strings.Split(command, " ")
	switch orders[0] {
	case "осмотреться":
		return player.Inspect()
	case "идти":
		return player.Move(orders[1])
	case "взять":
		return player.Take(orders[1])
	case "надеть":
		return player.PutOn(orders[1])
	case "применить":
		return player.Use(orders[1], orders[2])
	default:
		return "неизвестная команда"
	}
}
func main() {
	initGame()
	var command string
	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		command = scanner.Text()
		if command == "Stop" {
			break
		}
		answer := handleCommand(command)
		fmt.Println(answer)
	}
}
