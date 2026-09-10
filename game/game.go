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
	"github.com/filipstrom/goPet/pet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets/*
var assetFS embed.FS

const (
	speed                    = 10
	acc                      = 0.5
	screenW                  = 780
	screenH                  = 300
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
}

func (g *game) reset() {
	g.ai.ResetBody()
}

func (g *game) initialize() {
	g.backgroundColor = color.RGBA{0, 181, 226, 255}
	startPosition := linalg.Vec2{X: 0, Y: 0}
	body := pet.NewBody(acc, speed, splayerSpeed, 5, 100, 100, startPosition)
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
		g.ai.Move(pet.DirectionUp)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.ai.Move(pet.DirectionDown)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.ai.Move(pet.DirectionLeft)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.ai.Move(pet.DirectionRight)
	}

	g.ai.Update()

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
