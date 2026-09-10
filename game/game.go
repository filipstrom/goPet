package game

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/filipstrom/goPet/linalg"
	"github.com/filipstrom/goPet/object"
	"github.com/filipstrom/goPet/pet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

//go:embed assets/*
var assetFS embed.FS

const (
	speed                    = 10
	acc                      = 0.5
	screenW                  = 500
	screenH                  = 500
	playrWidth, playerHeight = 100, 100
	splayerSpeed             = 2
)

func loadImage(assetPath string) *ebiten.Image {
	f, err := assetFS.Open(assetPath)
	if err != nil {
		log.Panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

type game struct {
	backgroundColor color.RGBA
	playerImg       *ebiten.Image
	ai              pet.AI
	playerVelocity  float64
	fullscreen      bool
	initialized     bool
	walls           []object.Object
}

func (g *game) reset() {
	g.ai.ResetBody()
}

func (g *game) initialize() {
	g.backgroundColor = color.RGBA{0, 181, 226, 255}
	startPosition := linalg.Vec2{X: 250, Y: 250}
	body := pet.NewBody(acc, speed, splayerSpeed, 5, 100, 100, startPosition)

	g.walls = append(g.walls, object.Object{
		Appearance: []int{0, 0, 0},
		Texture:    []int{0, 0, 0},
		Position:   linalg.Vec2{X: 0, Y: 0},
		Width:      screenW,
		Height:     2,
	})

	g.walls = append(g.walls, object.Object{
		Appearance: []int{0, 0, 0},
		Texture:    []int{0, 0, 0},
		Position:   linalg.Vec2{X: 0, Y: screenH - 2},
		Width:      screenW,
		Height:     2,
	})

	g.walls = append(g.walls, object.Object{
		Appearance: []int{0, 0, 0},
		Texture:    []int{0, 0, 0},
		Position:   linalg.Vec2{X: screenW - 2, Y: 0},
		Width:      2,
		Height:     screenH,
	})
	g.walls = append(g.walls, object.Object{
		Appearance: []int{0, 0, 0},
		Texture:    []int{0, 0, 0},
		Position:   linalg.Vec2{X: 0, Y: 0},
		Width:      2,
		Height:     screenH,
	})

	g.ai = pet.NewAI("NN", body, "This place")
	g.playerImg = loadImage("assets/pet.png")
	g.fullscreen = true
	g.initialized = true
}

func newGame() *game {
	return &game{
		backgroundColor: color.RGBA{0, 0, 255, 255},
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}

func (g *game) Update() error {

	if !g.initialized {
		g.initialize()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt) {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.ai.ResetBody()
		fmt.Println("Hej")
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.ai.Control(pet.DirectionUp)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.ai.Control(pet.DirectionDown)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.ai.Control(pet.DirectionLeft)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.ai.Control(pet.DirectionRight)
	}

	g.ai.Update(g.walls)

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backgroundColor)

	op := &ebiten.DrawImageOptions{}

	bounds := g.playerImg.Bounds()
	imgW := float64(bounds.Dx())
	imgH := float64(bounds.Dy())

	scaleX := float64(playrWidth) / (imgW)
	scaleY := float64(playerHeight) / (imgH)

	for _, ob := range g.walls {
		vector.FillRect(
			screen,
			float32(ob.Position.X),
			float32(ob.Position.Y),
			float32(ob.Width),
			float32(ob.Height),
			color.White,
			false,
		)

	}

	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(-g.ai.GetWidth()/2, -g.ai.GetHeigt()/2)
	op.GeoM.Rotate(g.ai.GetRotation() * math.Pi / 180)
	op.GeoM.Translate(g.ai.Position())

	screen.DrawImage(g.playerImg, op)
}

func Start() {
	ebiten.SetWindowSize(900, 500)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
