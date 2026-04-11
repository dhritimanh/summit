package ui

import (
	"fmt"
	"image/color"
	"math"
	"summit/sim"
	"summit/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 720
	
	// Layout Zones
	MountainWidth  = 640
	MountainHeight = 540
	CommandHeight  = 180
	SidebarWidth   = 320
)

type Button struct {
	Label  string
	X, Y   int
	W, H   int
	Action func()
}

type GameState int

const (
	StatePlaying GameState = iota
	StateGameOver
	StateVictory
)

type App struct {
	Sim             *sim.Sim
	SelectedClimber *world.Climber
	Buttons         []Button
	Counter         int 
	State           GameState
	EndMessage      string
}

func (a *App) Update() error {
	a.Counter++
	
	// Default selection
	if a.SelectedClimber == nil && len(a.Sim.Team) > 0 {
		a.SelectedClimber = a.Sim.Team[0]
	}

	// Handle Clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		a.handleClick(mx, my)
	}

	if a.State != StatePlaying {
		return nil // Pause interacting
	}

	// Tick Simulation
	delta := 1.0 / 60.0 
	a.Sim.Tick(delta)

	// Check Win/Loss conditions
	for _, c := range a.Sim.Team {
		done, outcome := a.Sim.CheckTerminal(c)
		if done {
			if outcome == "success" {
				a.State = StateVictory
				a.EndMessage = fmt.Sprintf("VICTORY!\n\n%s made it back to Base Camp alive.\nAn incredible ascent.", c.Name)
			} else {
				a.State = StateGameOver
				if outcome == "ams_death" {
					a.EndMessage = fmt.Sprintf("TRAGEDY\n\n%s has collapsed from severe AMS (HACE).\nThere is no recovery.", c.Name)
				} else {
					a.EndMessage = fmt.Sprintf("TRAGEDY\n\n%s has succumbed to exhaustion.", c.Name)
				}
			}
		}
	}
	
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	// 1. Background
	screen.Fill(color.RGBA{10, 10, 20, 255})

	// 2. Draw Sectors
	a.drawMountainSector(screen)
	a.drawCommandBar(screen)
	a.drawSidebar(screen)

	// 3. Draw Game Over Overlay
	if a.State != StatePlaying {
		vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, 220}, true)
		ebitenutil.DebugPrintAt(screen, a.EndMessage, ScreenWidth/2-100, ScreenHeight/2-20)
	}
}

func (a *App) drawMountainSector(screen *ebiten.Image) {
	// Day/Night Background in Mountain Zone
	bg := a.getBackground()
	vector.DrawFilledRect(screen, 0, 0, MountainWidth, MountainHeight, bg, true)

	// Mountain Triangle
	var path vector.Path
	path.MoveTo(50, MountainHeight-50)
	path.LineTo(MountainWidth/2, 100)
	path.LineTo(MountainWidth-50, MountainHeight-50)
	path.Close()
	
	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vertices {
		vertices[i].SrcX, vertices[i].SrcY = 1, 1
		vertices[i].ColorR, vertices[i].ColorG, vertices[i].ColorB, vertices[i].ColorA = 0.3, 0.3, 0.4, 1
	}
	whiteImage := ebiten.NewImage(3, 3)
	whiteImage.Fill(color.White)
	screen.DrawTriangles(vertices, indices, whiteImage, &ebiten.DrawTrianglesOptions{})

	// Climbers
	for _, c := range a.Sim.Team {
		cx, cy := a.getClimberPos(c)
		
		// Panic Pulse
		pulse := math.Sin(float64(a.Counter) * 0.1) * 2
		radius := 7.0 + pulse
		if c.AMS >= 80 {
			radius = 7.0 + (math.Sin(float64(a.Counter)*0.3) * 5)
		}

		dotColor := color.RGBA{255, 255, 255, 255}
		if c.AMS >= 55 && c.AMS < 80 {
			dotColor = color.RGBA{255, 165, 0, 255} 
		} else if c.AMS >= 80 {
			dotColor = color.RGBA{255, 50, 50, 255}
		}

		vector.DrawFilledCircle(screen, float32(cx), float32(cy), float32(radius), dotColor, true)
		if c == a.SelectedClimber {
			vector.StrokeCircle(screen, float32(cx), float32(cy), float32(radius+4), 2, color.White, true)
		}
		ebitenutil.DebugPrintAt(screen, c.Name, cx+15, cy-8)
	}

	// Weather Info (Top Left of mountain)
	msg := fmt.Sprintf("Day %d | %02d:00\nWind: %d km/h", a.Sim.State.Day, a.Sim.State.Hour, a.Sim.State.WindSpeed)
	ebitenutil.DebugPrintAt(screen, msg, 20, 20)
}

func (a *App) drawCommandBar(screen *ebiten.Image) {
	yStart := MountainHeight
	vector.DrawFilledRect(screen, 0, float32(yStart), MountainWidth, CommandHeight, color.RGBA{20, 20, 30, 255}, true)
	
	// Draw Speed Buttons
	speedLabel := fmt.Sprintf("TIME: %dx", a.Sim.Speed)
	if a.Sim.Speed == 0 { speedLabel = "TIME: PAUSED" }
	ebitenutil.DebugPrintAt(screen, speedLabel, 20, yStart+20)

	a.Buttons = []Button{
		{"PAUSE", 20, yStart + 50, 60, 30, func() { a.Sim.Speed = 0 }},
		{"1x", 90, yStart + 50, 40, 30, func() { a.Sim.Speed = 1 }},
		{"4x", 140, yStart + 50, 40, 30, func() { a.Sim.Speed = 4 }},
		{"12x", 190, yStart + 50, 40, 30, func() { a.Sim.Speed = 12 }},
	}

	// Draw Climber Orders if selected
	if c := a.SelectedClimber; c != nil {
		ebitenutil.DebugPrintAt(screen, "ORDERS FOR "+c.Name, 260, yStart+20)
		a.Buttons = append(a.Buttons, []Button{
			{"CLIMB", 260, yStart + 50, 70, 40, func() { c.Order = world.OrderClimb }},
			{"REST", 340, yStart + 50, 70, 40, func() { c.Order = world.OrderRest }},
			{"HOLD", 420, yStart + 50, 70, 40, func() { c.Order = world.OrderHold }},
			{"DESCEND", 500, yStart + 50, 70, 40, func() { c.Order = world.OrderDescend }},
			{"O2 ON/OFF", 580, yStart + 50, 70, 40, func() { c.O2Active = !c.O2Active }},
		}...)
	}

	for _, b := range a.Buttons {
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), color.RGBA{40, 45, 60, 255}, true)
		ebitenutil.DebugPrintAt(screen, b.Label, b.X+5, b.Y+12)
	}
}

func (a *App) drawSidebar(screen *ebiten.Image) {
	xStart := MountainWidth
	vector.DrawFilledRect(screen, float32(xStart), 0, SidebarWidth, ScreenHeight, color.RGBA{15, 15, 25, 255}, true)

	if c := a.SelectedClimber; c != nil {
		ebitenutil.DebugPrintAt(screen, "CLIMER INTEL", xStart+20, 20)
		ebitenutil.DebugPrintAt(screen, "NAME: "+c.Name, xStart+20, 50)
		
		// Fitness Bar
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FITNESS: %d%%", c.Fitness), xStart+20, 80)
		a.drawBar(screen, xStart+20, 100, 200, 15, float32(c.Fitness)/100.0, color.RGBA{100, 255, 100, 255})

		// AMS Bar
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("AMS RISK: %d%%", c.AMS), xStart+20, 130)
		amsColor := color.RGBA{255, 200, 0, 255}
		if c.AMS > 70 { amsColor = color.RGBA{255, 50, 50, 255} }
		a.drawBar(screen, xStart+20, 150, 200, 15, float32(c.AMS)/100.0, amsColor)

		// O2 Tanks
		ebitenutil.DebugPrintAt(screen, "OXYGEN SUPPLY:", xStart+20, 190)
		for i := 0; i < 3; i++ {
			o2Col := color.RGBA{50, 50, 50, 255}
			if i < c.O2Charges { o2Col = color.RGBA{100, 200, 255, 255} }
			vector.DrawFilledRect(screen, float32(xStart+20+(i*40)), 215, 30, 50, o2Col, true)
		}

		status := "STATUS: " + string(c.Order)
		if c.O2Active { status += " | O2 ON" }
		ebitenutil.DebugPrintAt(screen, status, xStart+20, 280)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("ALTITUDE: %dm", c.Altitude), xStart+20, 305)
	}

	// Radio Log at bottom of sidebar
	ebitenutil.DebugPrintAt(screen, "RADIO LOG", xStart+20, 480)
	for i, line := range a.Sim.Log {
		ebitenutil.DebugPrintAt(screen, "> "+line, xStart+10, 500+(i*15))
	}
}

func (a *App) drawBar(screen *ebiten.Image, x, y, w, h int, val float32, col color.Color) {
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{40, 40, 40, 255}, true)
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w)*val, float32(h), col, true)
}

func (a *App) handleClick(x, y int) {
	// 1. Check Buttons
	for _, b := range a.Buttons {
		if x >= b.X && x <= b.X+b.W && y >= b.Y && y <= b.Y+b.H {
			b.Action()
			return
		}
	}

	// 2. Check Climbers
	for _, c := range a.Sim.Team {
		cx, cy := a.getClimberPos(c)
		dist := math.Sqrt(math.Pow(float64(x)-float64(cx), 2) + math.Pow(float64(y)-float64(cy), 2))
		if dist < 20 {
			a.SelectedClimber = c
			return
		}
	}
}

func (a *App) getClimberPos(c *world.Climber) (int, int) {
	altPercent := float32(c.Altitude-a.Sim.Mountain.CampAltitudes[0]) / 
		float32(a.Sim.Mountain.CampAltitudes[4]-a.Sim.Mountain.CampAltitudes[0])
	mapY := (MountainHeight - 50) - (altPercent * (MountainHeight - 150))
	mapX := MountainWidth / 2
	return int(mapX), int(mapY)
}

func (a *App) Layout(w, h int) (int, int) { return ScreenWidth, ScreenHeight }

func (a *App) getBackground() color.Color {
	h := a.Sim.State.Hour
	if h >= 6 && h < 17 { return color.RGBA{135, 206, 235, 255} }
	if h >= 17 && h < 20 { return color.RGBA{200, 100, 50, 255} }
	return color.RGBA{20, 20, 40, 255}
}
