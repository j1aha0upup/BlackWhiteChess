package main

import (
	"os"
	"strconv"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

const (
	no_chess = iota
	blue_chess
	red_chess
)

const (
	play_with_bot = iota
	play_with_player
)

type chessGame struct {
	width         int //chess game window width
	height        int //chess game window height
	X             int //move window start x
	Y             int //move window start x
	rootX         int //chess on table positionX
	rootY         int //chess on table positionY
	chessWidth    int //every chess width
	chessHeight   int //every chess height
	botId         int //means when bot play chess need time to think
	bot_or_player int

	chessWindow       *gtk.Window //chess game window widget
	buttonClose       *gtk.Button
	buttonMin         *gtk.Button
	start_with_bot    *gtk.Button
	start_with_player *gtk.Button
	red_chess_count   *gtk.Label
	blue_chess_count  *gtk.Label
	countdown         *gtk.Label
	red_chess_tip     *gtk.Image
	blue_chess_tip    *gtk.Image
	current_chess     int
	countdown_id      int
	chess_tip_id      int
	countdown_num     int
	chessTable        [8][8]int
}

func judgeResult(ch *chessGame) bool {
	for tempi := 0; tempi < 8; tempi++ {
		for tempj := 0; tempj < 8; tempj++ {
			i, j := tempi, tempj

			if ch.chessTable[i][j] == no_chess {
				j--
				for j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					j--
				}

				if j != -1 && ch.chessTable[i][j] == ch.current_chess && j != tempj-1 {
					return true
				}
				i = tempi
				j = tempj
			}

			if ch.chessTable[i][j] == no_chess {
				j++
				for j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					j++
				}

				if j != 8 && ch.chessTable[i][j] == ch.current_chess && j != tempj+1 {
					return true
				}
				i = tempi
				j = tempj
			}

			if ch.chessTable[i][j] == no_chess {
				i--
				for i != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i--
				}

				if i != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 {
					return true
				}
				i = tempi
				j = tempj
			}
			if ch.chessTable[i][j] == no_chess {
				i++
				for i != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i++
				}

				if i != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 {
					return true
				}
				i = tempi
				j = tempj
			}
			if ch.chessTable[i][j] == no_chess {
				i++
				j++
				for i != 8 && j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i++
					j++
				}

				if i != 8 && j != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 && j != tempj+1 {
					return true
				}
				i = tempi
				j = tempj
			}

			if ch.chessTable[i][j] == no_chess {
				i++
				j--
				for i != 8 && j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i++
					j--
				}

				if i != 8 && j != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 && j != tempj-1 {
					return true
				}
				i = tempi
				j = tempj
			}
			if ch.chessTable[i][j] == no_chess {
				i--
				j++
				for i != -1 && j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i--
					j++
				}

				if i != -1 && j != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 && j != tempj+1 {
					return true
				}
				i = tempi
				j = tempj

			}
			if ch.chessTable[i][j] == no_chess {
				i--
				j--
				for i != -1 && j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i--
					j--
				}

				if i != -1 && j != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 && j != tempj-1 {
					return true
				}
			}
		}
	}
	return false
}

func imagewidget_set_image(img *gtk.Image, imageFile string) {
	w, h := img.GetSizeRequest()
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale(imageFile, w, h, false)
	img.SetFromPixbuf(pixbuf)
	pixbuf.Unref()
}

func button_set_image(btn *gtk.Button, imageFile string) {
	w, h := btn.GetSizeRequest()
	pixbuf, _ := gdkpixbuf.NewPixbufFromFileAtScale(imageFile, w, h, false)
	image := gtk.NewImageFromPixbuf(pixbuf)
	btn.SetImage(image)
	pixbuf.Unref()
}

func move_window(ctx *glib.CallbackContext) {
	//获取鼠标属性结构体变量，系统内部的变量，不是用户传参变量
	arg := ctx.Args(0)
	//还是EventButton
	event := *(**gdk.EventButton)(unsafe.Pointer(&arg))
	temp := ctx.Data()
	ch, ok := temp.(*chessGame)

	if ok && ch.Y < 30 {
		ch.chessWindow.Move(int(event.XRoot)-int(ch.X), int(event.YRoot)-int(ch.Y))
	}
}

// from 8 dirctor eat chess
func chesseat(i int, j int, ch *chessGame) bool {
	tempi, tempj, result := i, j, false

	if ch.chessTable[i][j] == no_chess {
		j--
		for j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
			j--
		}

		if j != -1 && ch.chessTable[i][j] == ch.current_chess && j != tempj-1 {
			result = true
			j++
			for ; j < tempj; j++ {
				ch.chessTable[i][j] = ch.current_chess
			}
			ch.chessWindow.QueueDraw()
		}
		i = tempi
		j = tempj
	}

	if ch.chessTable[i][j] == no_chess {
		j++
		for j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
			j++
		}

		if j != 8 && ch.chessTable[i][j] == ch.current_chess && j != tempj+1 {
			result = true
			j--
			for ; j > tempj; j-- {
				ch.chessTable[i][j] = ch.current_chess
			}
			ch.chessWindow.QueueDraw()
		}
		i = tempi
		j = tempj
	}

	if ch.chessTable[i][j] == no_chess {
		i--
		for i != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
			i--
		}

		if i != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 {
			result = true
			i++
			for ; i < tempi; i++ {
				ch.chessTable[i][j] = ch.current_chess
			}
			ch.chessWindow.QueueDraw()
		}
		i = tempi
		j = tempj
	}

	if ch.chessTable[i][j] == no_chess {
		i++
		for i != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
			i++
		}

		if i != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 {
			result = true
			i--
			for ; i > tempi; i-- {
				ch.chessTable[i][j] = ch.current_chess
			}
			ch.chessWindow.QueueDraw()
		}
		i = tempi
		j = tempj
	}

	if ch.chessTable[i][j] == no_chess {
		i++
		j++
		for i != 8 && j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
			i++
			j++
		}

		if i != 8 && j != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 && j != tempj+1 {
			result = true
			i--
			j--
			for i > tempi && j > tempj {
				ch.chessTable[i][j] = ch.current_chess
				i--
				j--
			}
			ch.chessWindow.QueueDraw()
		}
		i = tempi
		j = tempj
	}

	if ch.chessTable[i][j] == no_chess {
		i++
		j--
		for i != 8 && j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
			i++
			j--
		}

		if i != 8 && j != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 && j != tempj-1 {
			result = true
			i--
			j++
			for i > tempi && j < tempj {
				ch.chessTable[i][j] = ch.current_chess
				i--
				j++
			}
			ch.chessWindow.QueueDraw()
		}
		i = tempi
		j = tempj
	}

	if ch.chessTable[i][j] == no_chess {
		i--
		j++
		for i != -1 && j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
			i--
			j++
		}

		if i != -1 && j != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 && j != tempj+1 {
			result = true
			i++
			j--
			for i < tempi && j > tempj {
				ch.chessTable[i][j] = ch.current_chess
				i++
				j--
			}
			ch.chessWindow.QueueDraw()
		}
		i = tempi
		j = tempj
	}

	if ch.chessTable[i][j] == no_chess {
		i--
		j--
		for i != -1 && j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
			i--
			j--
		}

		if i != -1 && j != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 && j != tempj-1 {
			result = true
			i++
			j++
			for i < tempi && j < tempj {
				ch.chessTable[i][j] = ch.current_chess
				i++
				j++
			}
			ch.chessWindow.QueueDraw()
		}
	}

	if result {
		ch.chessTable[tempi][tempj] = ch.current_chess
	}

	return result

}

func keep_move_window_start_coordinate(ctx *glib.CallbackContext) {
	//获取鼠标属性结构体变量，系统内部的变量，不是用户传参变量
	arg := ctx.Args(0)
	//还是EventButton
	event := *(**gdk.EventButton)(unsafe.Pointer(&arg))
	temp := ctx.Data()
	ch, ok := temp.(*chessGame)
	if ok && event.Button == 1 {
		ch.X = int(event.X)
		ch.Y = int(event.Y)

		if int(event.X) > ch.rootX && int(event.Y) > ch.rootY && int(event.X) < ch.rootX+8*ch.chessWidth && int(event.Y) < ch.rootY+8*ch.chessHeight {
			i := (int(event.X) - ch.rootX) / ch.chessWidth
			j := (int(event.Y) - ch.rootY) / ch.chessHeight

			if !(ch.bot_or_player == play_with_bot && ch.current_chess == blue_chess) {
				if chesseat(i, j, ch) {
					change_role(ch)
				}
			}

		}
	}

}

func close_window(ctx *glib.CallbackContext) {
	temp := ctx.Data()
	ch, ok := temp.(*chessGame)
	if ok {
		glib.TimeoutRemove(ch.countdown_id)
		glib.TimeoutRemove(ch.chess_tip_id)
		gtk.MainQuit()
	}
}

func (ch *chessGame) createWindow() {
	//get window widget and turn chessWindow type
	builder := gtk.NewBuilder()
	builder.AddFromFile("chess_game.glade")
	ch.chessWindow = &gtk.Window{Bin: gtk.Bin{Container: gtk.Container{Widget: *gtk.WidgetFromObject(builder.GetObject("chessWindow"))}}}

	//set about window attributes
	ch.width = 800
	ch.height = 480
	ch.chessWindow.SetSizeRequest(ch.width, ch.height) //setting size
	ch.chessWindow.SetPosition(gtk.WIN_POS_CENTER)     //settting position
	ch.chessWindow.SetAppPaintable(true)               //set allow paint
	ch.chessWindow.SetDecorated(false)                 //set no border
	ch.chessWindow.SetIconFromFile("images/face.png")  //set application icon

	//set move window
	ch.chessWindow.SetEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))
	ch.chessWindow.Connect("button-press-event", keep_move_window_start_coordinate, ch)
	ch.chessWindow.Connect("motion-notify-event", move_window, ch)

	//set background photo
	ch.chessWindow.Connect("expose-event", func() {
		//指定窗口为绘图区域，在窗口上绘图
		painter := ch.chessWindow.GetWindow().GetDrawable()
		gc := gdk.NewGC(painter)
		bk, _ := gdkpixbuf.NewPixbufFromFileAtScale("images/bg.jpg", ch.width, ch.height, false)
		painter.DrawPixbuf(gc, bk, 0, 0, 0, 0, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)

		blueChess, _ := gdkpixbuf.NewPixbufFromFileAtScale("images/blue.png", ch.chessWidth, ch.chessHeight, false)
		redChess, _ := gdkpixbuf.NewPixbufFromFileAtScale("images/red.png", ch.chessWidth, ch.chessHeight, false)

		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				if ch.chessTable[i][j] == blue_chess {
					painter.DrawPixbuf(gc, blueChess, 0, 0, ch.rootX+i*ch.chessWidth, ch.rootY+j*ch.chessHeight, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
				} else if ch.chessTable[i][j] == red_chess {
					painter.DrawPixbuf(gc, redChess, 0, 0, ch.rootX+i*ch.chessWidth, ch.rootY+j*ch.chessHeight, -1, -1, gdk.RGB_DITHER_NONE, 0, 0)
				}
			}

		}

		bk.Unref()
	})

	//get and change button type than set button connect
	ch.buttonClose = &gtk.Button{
		Bin: gtk.Bin{
			Container: gtk.Container{
				Widget: *gtk.WidgetFromObject(builder.GetObject("buttonClose")),
			},
		},
		Activatable: gtk.Activatable(*gtk.WidgetFromObject(builder.GetObject("buttonClose"))),
	}
	ch.buttonMin = &gtk.Button{
		Bin: gtk.Bin{
			Container: gtk.Container{
				Widget: *gtk.WidgetFromObject(builder.GetObject("buttonMin")),
			},
		},
		Activatable: gtk.Activatable(*gtk.WidgetFromObject(builder.GetObject("buttonMin"))),
	}
	ch.buttonClose.Connect("clicked", close_window, ch)     //set buttonclose connect
	ch.buttonMin.Connect("clicked", ch.chessWindow.Iconify) //set buttonmin connect
	ch.buttonClose.SetCanFocus(false)                       //remove button foucsborder
	ch.buttonMin.SetCanFocus(false)                         //remove button foucsborder
	button_set_image(ch.buttonClose, "images/close.png")
	button_set_image(ch.buttonMin, "images/min.png")

	//get label and set label default
	ch.red_chess_count = &gtk.Label{Misc: gtk.Misc{Widget: *gtk.WidgetFromObject(builder.GetObject("red_chess_count"))}}
	ch.blue_chess_count = &gtk.Label{Misc: gtk.Misc{Widget: *gtk.WidgetFromObject(builder.GetObject("blue_chess_count"))}}
	ch.countdown = &gtk.Label{Misc: gtk.Misc{Widget: *gtk.WidgetFromObject(builder.GetObject("countdown"))}}
	ch.red_chess_count.SetText("2")
	ch.blue_chess_count.SetText("2")
	ch.countdown.SetText("20")
	ch.blue_chess_count.ModifyFontEasy("50")
	ch.red_chess_count.ModifyFontEasy("50")
	ch.countdown.ModifyFontEasy("30")
	ch.blue_chess_count.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	ch.red_chess_count.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))
	ch.countdown.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	//get image widget and set image
	ch.red_chess_tip = &gtk.Image{Misc: gtk.Misc{Widget: *gtk.WidgetFromObject(builder.GetObject("red_chess_tip"))}}
	ch.blue_chess_tip = &gtk.Image{Misc: gtk.Misc{Widget: *gtk.WidgetFromObject(builder.GetObject("blue_chess_tip"))}}
	imagewidget_set_image(ch.blue_chess_tip, "images/blue.png")
	imagewidget_set_image(ch.red_chess_tip, "images/red.png")

	//empty chessTable
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			ch.chessTable[i][j] = no_chess
		}
	}

	//other attributes
	ch.red_chess_tip.Hide()
	ch.blue_chess_tip.Hide()
	ch.current_chess = red_chess
	ch.countdown_num = 21
	ch.rootX = 205
	ch.rootY = 55
	ch.chessWidth = 49
	ch.chessHeight = 41

	//new button select play with bot or player
	ch.start_with_bot = &gtk.Button{
		Bin: gtk.Bin{
			Container: gtk.Container{
				Widget: *gtk.WidgetFromObject(builder.GetObject("start_with_bot")),
			},
		},
		Activatable: gtk.Activatable(*gtk.WidgetFromObject(builder.GetObject("start_with_bot"))),
	}
	ch.start_with_player = &gtk.Button{
		Bin: gtk.Bin{
			Container: gtk.Container{
				Widget: *gtk.WidgetFromObject(builder.GetObject("start_with_player")),
			},
		},
		Activatable: gtk.Activatable(*gtk.WidgetFromObject(builder.GetObject("start_with_player"))),
	}
	ch.start_with_bot.SetLabel("start with bot")
	ch.start_with_player.SetLabel("player pk")
	ch.start_with_bot.Connect("clicked", func() {
		ch.start_with_bot.Hide()
		ch.start_with_player.Hide()
		ch.chessWindow.QueueDraw()
		ch.bot_or_player = play_with_bot
		ch.chessInit()

	})
	ch.start_with_player.Connect("clicked", func() {
		ch.start_with_bot.Hide()
		ch.start_with_player.Hide()
		ch.chessWindow.QueueDraw()
		ch.chessInit()
		ch.bot_or_player = play_with_player
	})

	ch.chessWindow.ShowAll()
}

func (ch *chessGame) chessInit() {
	ch.chessTable[3][3] = blue_chess
	ch.chessTable[4][4] = blue_chess
	ch.chessTable[3][4] = red_chess
	ch.chessTable[4][3] = red_chess

	ch.countdown_id = glib.TimeoutAdd(1000, func() bool {
		ch.countdown_num--
		ch.countdown.SetText(strconv.Itoa(ch.countdown_num))

		if ch.countdown_num == 0 {
			change_role(ch)
		}

		return true
	})

	ch.chess_tip_id = glib.TimeoutAdd(1000, func() bool {
		change_hideordisplay(ch)
		return true
	})

}
func change_hideordisplay(ch *chessGame) {
	if ch.current_chess == red_chess {
		ch.blue_chess_tip.Hide()
		if ch.red_chess_tip.GetVisible() {
			ch.red_chess_tip.Hide()
		} else {
			ch.red_chess_tip.Show()
		}
	} else {
		ch.red_chess_tip.Hide()
		if ch.blue_chess_tip.GetVisible() {
			ch.blue_chess_tip.Hide()
		} else {
			ch.blue_chess_tip.Show()
		}
	}
}

func botPlay(ch *chessGame) {
	eats, max, resulti, resultj := 0, 0, 0, 0

	for tempi := 0; tempi < 8; tempi++ {
		for tempj := 0; tempj < 8; tempj++ {
			i, j := tempi, tempj

			if ch.chessTable[i][j] == no_chess {
				j--
				for j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					j--
				}

				if j != -1 && ch.chessTable[i][j] == ch.current_chess && j != tempj-1 {
					eats++
				}
				i = tempi
				j = tempj
			}

			if ch.chessTable[i][j] == no_chess {
				j++
				for j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					j++
				}

				if j != 8 && ch.chessTable[i][j] == ch.current_chess && j != tempj+1 {
					eats++
				}
				i = tempi
				j = tempj
			}

			if ch.chessTable[i][j] == no_chess {
				i--
				for i != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i--
				}

				if i != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 {
					eats++
				}
				i = tempi
				j = tempj
			}
			if ch.chessTable[i][j] == no_chess {
				i++
				for i != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i++
				}

				if i != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 {
					eats++
				}
				i = tempi
				j = tempj
			}
			if ch.chessTable[i][j] == no_chess {
				i++
				j++
				for i != 8 && j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i++
					j++
				}

				if i != 8 && j != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 && j != tempj+1 {
					eats++
				}
				i = tempi
				j = tempj
			}

			if ch.chessTable[i][j] == no_chess {
				i++
				j--
				for i != 8 && j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i++
					j--
				}

				if i != 8 && j != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi+1 && j != tempj-1 {
					eats++
				}
				i = tempi
				j = tempj
			}
			if ch.chessTable[i][j] == no_chess {
				i--
				j++
				for i != -1 && j != 8 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i--
					j++
				}

				if i != -1 && j != 8 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 && j != tempj+1 {
					eats++
				}
				i = tempi
				j = tempj

			}
			if ch.chessTable[i][j] == no_chess {
				i--
				j--
				for i != -1 && j != -1 && ch.chessTable[i][j] != ch.current_chess && ch.chessTable[i][j] != no_chess {
					i--
					j--
				}

				if i != -1 && j != -1 && ch.chessTable[i][j] == ch.current_chess && i != tempi-1 && j != tempj-1 {
					eats++
				}
			}

			if eats > max {
				max = eats
				resulti = tempi
				resultj = tempj
			}
			eats = 0
		}
	}
	if chesseat(resulti, resultj, ch) {
		change_role(ch)
	}

}

func change_role(ch *chessGame) {
	ch.countdown_num = 21
	redChessCount, blueChessCount := 0, 0

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if ch.chessTable[i][j] == red_chess {
				redChessCount++
			} else if ch.chessTable[i][j] == blue_chess {
				blueChessCount++
			}
		}
	}

	ch.blue_chess_count.SetText(strconv.Itoa(blueChessCount))
	ch.red_chess_count.SetText(strconv.Itoa(redChessCount))

	if ch.current_chess == blue_chess {
		ch.current_chess = red_chess
	} else {
		ch.current_chess = blue_chess
	}

	if !judgeResult(ch) {
		glib.TimeoutRemove(ch.chess_tip_id)
		glib.TimeoutRemove(ch.countdown_id)
		if redChessCount > blueChessCount {
			dialog1 := gtk.NewMessageDialog(
				ch.chessWindow.GetTopLevelAsWindow(), //指定父窗口
				gtk.DIALOG_MODAL,                     //模态对话框
				gtk.MESSAGE_QUESTION,                 //指定对话框类型
				gtk.BUTTONS_YES_NO,                   //默认按钮
				"red win,Again?",                     //设置内容
			)
			dialog1.SetTitle("Red Win!")

			flag := dialog1.Run()

			if flag == gtk.RESPONSE_YES {
				dialog1.Destroy() //销毁对话框
				ch.createWindow()
			} else if flag == gtk.RESPONSE_NO {
				dialog1.Destroy() //销毁对话框
				gtk.MainQuit()
			}
		} else if redChessCount < blueChessCount {
			dialog1 := gtk.NewMessageDialog(
				ch.chessWindow.GetTopLevelAsWindow(), //指定父窗口
				gtk.DIALOG_MODAL,                     //模态对话框
				gtk.MESSAGE_QUESTION,                 //指定对话框类型
				gtk.BUTTONS_YES_NO,                   //默认按钮
				"blue win,Again?",                    //设置内容
			)
			dialog1.SetTitle("Blue Win!")

			flag := dialog1.Run()

			if flag == gtk.RESPONSE_YES {
				dialog1.Destroy() //销毁对话框
				ch.createWindow()
			} else if flag == gtk.RESPONSE_NO {
				dialog1.Destroy() //销毁对话框
				gtk.MainQuit()
			}

		} else {
			dialog1 := gtk.NewMessageDialog(
				ch.chessWindow.GetTopLevelAsWindow(), //指定父窗口
				gtk.DIALOG_MODAL,                     //模态对话框
				gtk.MESSAGE_QUESTION,                 //指定对话框类型
				gtk.BUTTONS_YES_NO,                   //默认按钮
				"Draw,Again?",                        //设置内容
			)
			dialog1.SetTitle("Draw!")

			flag := dialog1.Run()

			if flag == gtk.RESPONSE_YES {
				dialog1.Destroy() //销毁对话框
				ch.createWindow()
			} else if flag == gtk.RESPONSE_NO {
				dialog1.Destroy() //销毁对话框
				gtk.MainQuit()
			}
		}
	}
	if ch.current_chess == blue_chess && ch.bot_or_player == play_with_bot {
		ch.botId = glib.TimeoutAdd(1500, func() bool {
			botPlay(ch)
			return false
		})
	}
}

func main() {
	gtk.Init(&os.Args)

	var chessGamer chessGame
	chessGamer.createWindow()

	gtk.Main()
}
