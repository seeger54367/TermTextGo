package main

import (
	"log"

	. "github.com/rthornton128/goncurses"
)

func main() {
	// If you want to limit the number of messages, change maxMessageNum
	// maxMessageNum = ""
	var active int
	stdscr, err := Init()
	if err != nil {
		log.Fatal(err)
	}
	defer End()

	data := map[string]string{"to": "", "from": "", "message": ""}
	//Add numbers for names to appear in interface
	knownNumbers := map[string]string{"+19999999999": "Me"}

	//Change this to your twilio number
	data["from"] = "+19999999999"

	client := setupClient()
	contactMenu := getContactList(client)
	convoMenu := getContactMessages(client, contactMenu[0], knownNumbers)

	Raw(true)
	Echo(false)
	Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	printControls(stdscr)
	my, mx := stdscr.MaxYX()
	CONTACTHEIGHT := (my - 1)
	CONTACTWIDTH := 30
	CONVOHEIGHT := (my - 6)
	CONVOWIDTH := (mx - CONTACTWIDTH - 1)
	CONTROLHEIGHT := (my - CONVOHEIGHT)
	CONTROLWIDTH := CONVOWIDTH

	//Contact window starting coordinates
	conty, contx := 1, 0

	//Convo window starting coordinates
	convoy, convox := 1, (CONTACTWIDTH)

	// Control Window starting coordinates
	controly, controlx := CONVOHEIGHT+1, (CONTACTWIDTH + 1)

	contactWin, _ := NewWindow(CONTACTHEIGHT, CONTACTWIDTH, conty, contx)
	contactWin.Keypad(true)

	convoWin, _ := NewWindow(CONVOHEIGHT, CONVOWIDTH, convoy, convox)
	convoWin.Keypad(true)

	controlWin, _ := NewWindow(CONTROLHEIGHT, CONTROLWIDTH, controly, controlx)
	controlWin.Keypad(true)

	stdscr.Refresh()

	//map messages to contact
	contactMessages := map[string][]string{}
	for _, v := range contactMenu[:] {
		contactMessages[v] = getContactMessages(client, v, knownNumbers)
	}

	availWindows := []*Window{contactWin, convoWin}
	activeWindowIndex := 0
	activeWindow := availWindows[activeWindowIndex]
	activeItem := map[*Window]int{}
	availMenu := map[*Window][]string{contactWin: contactMenu, convoWin: convoMenu}
	activeMenu := availMenu[activeWindow]
	for _, w := range availWindows[:] {
		activeItem[w] = 0
	}

	printContacts(contactWin, contactMenu, active, knownNumbers)
	printMessages(convoWin, contactMessages[contactMenu[activeItem[activeWindow]]], convoy, 1)
	data["to"] = activeMenu[activeItem[activeWindow]]

	for {
		ch := stdscr.GetChar()
		switch Key(ch) {
		case 'q':
			return
		case KEY_UP, 'k':
			if activeItem[activeWindow] == 0 {
				activeItem[activeWindow] = len(activeMenu) - 1
				data["to"] = activeMenu[activeItem[activeWindow]]
			} else {
				activeItem[activeWindow] -= 1
				data["to"] = activeMenu[activeItem[activeWindow]]
			}
		case KEY_DOWN, 'j':
			if activeItem[activeWindow] == len(activeMenu)-1 {
				activeItem[activeWindow] = 0
				data["to"] = activeMenu[activeItem[activeWindow]]
			} else {
				activeItem[activeWindow] += 1
				data["to"] = activeMenu[activeItem[activeWindow]]
			}
		case 'r':
			for _, v := range contactMenu[:] {
				contactMessages[v] = getContactMessages(client, v, knownNumbers)
			}

		case KEY_RETURN, KEY_ENTER, Key('\r'):
			stdscr.Clear()
			stdscr.Refresh()
			data["message"] = composeMessage(stdscr, activeMenu[activeItem[activeWindow]], contactMenu[activeItem[activeWindow]])
			sendMessage(client, data)
			contactMessages[contactMenu[activeItem[activeWindow]]] = getContactMessages(client, activeMenu[activeItem[activeWindow]], knownNumbers)
			printContacts(contactWin, contactMenu, active, knownNumbers)
			printContacts(convoWin, convoMenu, active, knownNumbers)
		default:
			stdscr.MovePrintf(0, 0, "Character pressed = %3d/%c",
				ch, ch)
			stdscr.ClearToEOL()
			stdscr.Refresh()
		}

		printContacts(activeWindow, activeMenu, activeItem[activeWindow], knownNumbers)
		printMessages(convoWin, contactMessages[contactMenu[activeItem[activeWindow]]], convoy, 1)
		printControls(stdscr)
	}
}

func printContacts(w *Window, menu []string, active int, kn map[string]string) {
	y, x := 2, 2
	w.Box(0, 0)
	for i, s := range menu {
		if val, ok := kn[s]; ok {
			s = val
		}
		if i == active {
			w.AttrOn(A_REVERSE)
			w.MovePrint(y+i, x, s)
			w.AttrOff(A_REVERSE)
		} else {
			w.MovePrint(y+i, x, s)
		}
	}
	w.Refresh()
}

func printControls(w *Window) {
	msg := "'k' for up. 'j' for down, 'ENTR' to send message. 'r' to reload messages for all contacts (may seem unresponsive for a couple seconds)"
	w.MovePrint(0, 0, msg)
}

func printMessages(w *Window, menu []string, y int, x int) {
	w.Clear()
	w.Refresh()
	w.Box(0, 0)
	for i, s := range menu {
		w.MovePrint(y+i, x, s)
	}
	w.Refresh()
}
func composeMessage(s *Window, c string, n string) string {
	s.Clear()
	s.Refresh()
	msg := "Enter your Message for " + c + " (Enter to send): "
	row, col := 2, 2
	s.MovePrint(row, col, msg)

	/* GetString will only retrieve the specified number of characters. Any
	attempts by the user to enter more characters will elicit an audible
	beep */
	var str string
	var err error
	Echo(true)
	Cursor(1)
	str, err = s.GetString(500)
	Raw(true)
	Echo(false)
	Cursor(0)
	if err != nil {
		Echo(false)
		Cursor(0)
		s.MovePrint(row+1, col, "GetString Error:", err)
	}
	return str

}
