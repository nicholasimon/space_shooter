package main

import (
	"math/rand"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var ( // MARK: var
	// imgs
	backfadeon bool
	backfade   = float32(0.4)
	backimg    rl.Texture2D
	playerimg  = rl.NewRectangle(0, 0, 16, 16)
	imgs       rl.Texture2D
	// enemies
	enemiesmap = make([]string, levela)
	// extras
	bulletmap = make([]string, levela)
	// player
	player, playerh, playerv int
	// level
	minh                     = 2
	maxh                     = (gridh / 2)
	drawblock, drawblocknext int
	levelw, levelh, levela   int
	levelmap                 = make([]string, levela)
	// core
	cameracount         int
	gridw, gridh, grida int
	monh32, monw32      int32
	monitorh, monitorw  int
	debugon             bool
	framecount          int
	camera              rl.Camera2D
)

func raylib() { // MARK: raylib
	rl.InitWindow(monw32, monh32, "galaxtica")
	rl.SetExitKey(rl.KeyEnd)          // key to end the game and close window
	imgs = rl.LoadTexture("imgs.png") // load images
	backimggen := rl.GenImageCellular(monitorw, monitorh, 32)
	backimg = rl.LoadTextureFromImage(backimggen)
	rl.UnloadImage(backimggen)
	rl.SetTargetFPS(24)
	for !rl.WindowShouldClose() { // MARK: WindowShouldClose

		framecount++
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.DrawTexture(backimg, 0, 0, rl.Fade(rl.Red, backfade))
		rl.BeginMode2D(camera)
		// MARK: draw map layer 1 / left up

		drawx := int32(0)
		drawy := int32(0)
		drawblock = drawblocknext
		count := 0

		for a := 0; a < grida; a++ {

			checklevel := levelmap[drawblock]
			checkbullet := bulletmap[drawblock]
			checkenemy := enemiesmap[drawblock]

			switch checklevel {
			case ".":
				rl.DrawRectangle(drawx, drawy, 15, 15, rl.Black)
				//	v2 := rl.NewVector2(float32(drawx), float32(drawy))
			//	rl.DrawTextureRec(imgs, tile1, v2, rl.White)
			case " ":
				//	rl.DrawRectangleLines(drawx, drawy, 15, 15, rl.Fade(rl.White, 0.2))
			}
			switch checkbullet {
			case "b":
				rl.DrawRectangle(drawx, drawy, 15, 15, rl.Orange)
			}
			switch checkenemy {
			case "e1":
				rl.DrawRectangle(drawx, drawy, 15, 15, rl.Green)
			}

			// draw player
			if player == drawblock {
				//	rl.DrawRectangle(drawx, drawy, 15, 15, rl.Red)
				v2 := rl.NewVector2(float32(drawx), float32(drawy))
				rl.DrawTextureRec(imgs, playerimg, v2, rl.White)
			}

			count++
			drawx += 16
			drawblock++
			if count == gridw {
				count = 0
				drawx = 0
				drawy += 16
				drawblock += levelw - gridw
			}

		}

		rl.EndMode2D() // MARK: draw no camera
		fx()
		if debugon {
			debug()
		}

		rl.EndDrawing()

		input()
		updateall()

	}
	rl.CloseWindow()
}
func createlevel() { // MARK: createlevel
	maxh = (gridh / 2) + 4
	levelmap = make([]string, levela)
	for a := 0; a < levela; a++ {
		levelmap[a] = " "
	}

	for a := 0; a < levelw*2; a++ {
		levelmap[a] = "."
	}

	block := levelw * 2
	length := rInt(minh, maxh)
	levelmap[block] = "."
	for a := 0; a < levelw; a++ {
		for b := 0; b < length; b++ {
			levelmap[block+(b*levelw)] = "."
		}
		block++
		length += rInt(-minh, minh)
		if length <= 2 {
			length += (minh * 2)
		}
		if length > maxh {
			length -= (minh)
		}
	}

	block = levelw * (gridh - 2)

	for a := 0; a < levelw*2; a++ {
		levelmap[block] = "."
		block++
	}
	block = levelw * (gridh - 2)
	length = rInt(minh, maxh)
	levelmap[block] = "."
	for a := 0; a < levelw; a++ {
		for b := 0; b < length; b++ {
			levelmap[block-(b*levelw)] = "."
		}
		block++
		length += rInt(-minh, minh)
		if length <= 2 {
			length += (minh * 2)
		}
		if length > maxh {
			length -= (minh)
		}
	}

}
func main() { // MARK: main
	rand.Seed(time.Now().UnixNano()) // random numbers
	rl.SetTraceLog(rl.LogError)      // hides info window
	initialize()
	raylib()
}
func horizvert() { // MARK: horizvert
	playerh, playerv = player/levelw, player%levelw

}
func createmaps() {
	bulletmap = make([]string, levela)
	enemiesmap = make([]string, levela)
}
func initialize() {
	rl.InitWindow(monw32, monh32, "galaxtica")
	setscreen()
	rl.CloseWindow()
	createlevel()
	startsettings()
	createmaps()
}
func startsettings() { // MARK:startsettings
	player = 10
	player += levelw * (gridh / 2)
}
func setscreen() { // MARK: setscreen
	monitorh = rl.GetScreenHeight()
	monitorw = rl.GetScreenWidth()
	monh32 = int32(monitorh)
	monw32 = int32(monitorw)
	rl.SetWindowSize(monitorw, monitorh)

	camera.Zoom = 2.0
	camera.Target.X = 0
	camera.Target.Y = 0

	gridw = (monitorw/16 + 1) / 2
	gridh = (monitorh/16 + 1) / 2
	grida = gridw * gridh

	levelw = gridw * 20
	levelh = gridh
	levela = levelw * levelh

}
func updateall() { // MARK: updateall

	if drawblocknext < levelw-gridw {
		drawblocknext++
		player++
	}

	playercollisions()
	enemies()
	animations()

}
func animations() { // MARK: animations

	if backfadeon {
		backfade += 0.02
		if backfade >= 0.4 {
			backfadeon = false
		}
	} else {
		backfade -= 0.02
		if backfade <= 0.2 {
			backfadeon = true
		}
	}

}
func fx() { // MARK: fx

	scany := int32(0)
	for a := 0; a < monitorh; a++ {
		rl.DrawLine(0, scany, monw32, scany, rl.Fade(rl.Black, 0.8))
		scany += 2
		a++
	}

}
func enemies() { // MARK: enemies

	if framecount%60 == 0 {
		block := drawblocknext + gridw
		block += (gridh / 2) * levelw
		enemiesmap[block] = "e1"
		enemiesmap[block+1] = "e1"
		enemiesmap[block+levelw] = "e1"
		enemiesmap[(block+1)+levelw] = "e1"
	}

}
func playercollisions() { // MARK: playercollisions

	if levelmap[player+1] == "." {
		horizvert()
		if playerh < gridh/2 {
			player += levelw
		} else if playerh > gridh/2 {
			player -= levelw
		}
	}

}
func checkcollision(dir int) bool { // MARK: checkcollision
	collision := false
	switch dir {
	case 1:
		if levelmap[player-levelw] == "." {
			collision = true
		}
	case 2:
		if levelmap[player+1] == "." {
			collision = true
		}
	case 3:
		if levelmap[player+levelw] == "." {
			collision = true
		}
	case 4:
		if levelmap[player-1] == "." {
			collision = true
		}
	}
	return collision
}
func bullet() { // MARK: bullet
	bulletmap[player] = "b"
}
func input() { // MARK: input

	if rl.IsKeyPressed(rl.KeySpace) {
		bullet()
	}

	if rl.IsKeyDown(rl.KeyLeft) {
		horizvert()
		if checkcollision(4) {

		} else {
			if playerv > drawblocknext+2 {
				player--
			}
		}
	}
	if rl.IsKeyDown(rl.KeyRight) {
		horizvert()
		if checkcollision(2) {

		} else {
			if playerv < drawblocknext+(gridw-3) {
				player++
			}
		}
	}
	if rl.IsKeyDown(rl.KeyUp) {
		horizvert()
		if checkcollision(1) {

		} else {
			if playerh > 2 {
				player -= levelw
			}
		}
	}
	if rl.IsKeyDown(rl.KeyDown) {
		if checkcollision(3) {

		} else {
			horizvert()
			if playerh < gridh-2 {
				player += levelw
			}
		}
	}

	if rl.IsKeyPressed(rl.KeyKpDecimal) {
		if debugon {
			debugon = false
		} else {
			debugon = true
		}
	}
}
func debug() { // MARK: debug
	rl.DrawRectangle(monw32-300, 0, 500, monw32, rl.Fade(rl.Blue, 0.4))
	rl.DrawFPS(monw32-290, monh32-100)

	gridwTEXT := strconv.Itoa(gridw)
	gridhTEXT := strconv.Itoa(gridh)
	gridaTEXT := strconv.Itoa(grida)
	playerTEXT := strconv.Itoa(player)
	playerhTEXT := strconv.Itoa(playerv)
	playervTEXT := strconv.Itoa(playerh)
	drawblocknextTEXT := strconv.Itoa(drawblocknext)

	rl.DrawText(gridwTEXT, monw32-290, 10, 10, rl.White)
	rl.DrawText("gridw", monw32-200, 10, 10, rl.White)
	rl.DrawText(gridhTEXT, monw32-290, 20, 10, rl.White)
	rl.DrawText("gridh", monw32-200, 20, 10, rl.White)
	rl.DrawText(gridaTEXT, monw32-290, 30, 10, rl.White)
	rl.DrawText("grida", monw32-200, 30, 10, rl.White)
	rl.DrawText(playerTEXT, monw32-290, 40, 10, rl.White)
	rl.DrawText("player", monw32-200, 40, 10, rl.White)
	rl.DrawText(playerhTEXT, monw32-290, 50, 10, rl.White)
	rl.DrawText("playerv", monw32-200, 50, 10, rl.White)
	rl.DrawText(playervTEXT, monw32-290, 60, 10, rl.White)
	rl.DrawText("playerh", monw32-200, 60, 10, rl.White)
	rl.DrawText(drawblocknextTEXT, monw32-290, 70, 10, rl.White)
	rl.DrawText("drawblocknext", monw32-200, 70, 10, rl.White)

}

// random numbers
func rInt(min, max int) int {
	return rand.Intn(max-min) + min
}
func rInt32(min, max int) int32 {
	a := int32(rand.Intn(max-min) + min)
	return a
}
func rFloat32(min, max int) float32 {
	a := float32(rand.Intn(max-min) + min)
	return a
}
func flipcoin() bool {
	var b bool
	a := rInt(0, 10001)
	if a < 5000 {
		b = true
	}
	return b
}
func rolldice() int {
	a := rInt(1, 7)
	return a
}
