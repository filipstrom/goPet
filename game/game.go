package game

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/filipstrom/goPet/object"
	"github.com/filipstrom/goPet/pet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/jakecoffman/cp/v2"
)

////go:embed assets/*
//// var assetFS embed.FS

const (
	speed                    = 10
	acc                      = 0.5
	screenW                  = 500
	screenH                  = 500
	playrWidth, playerHeight = 100, 100
	splayerSpeed             = 2
)

/*
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
*/
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
	body := cp.NewBody(1, 10)

	object := object.Object{Appearance: []int{1, 2, 3}, Texture: []int{1, 2, 3}, Shape: cp.NewBox(body, 50, 30, 0)}

	g.ai = pet.NewAI("NN", object, "This place")
	//g.playerImg = loadImage("assets/pet.png")

	playerImage := ebiten.NewImage(
		int(g.ai.GetWidth()),
		int(g.ai.GetHeigt()),
	)

	playerImage.Fill(color.White)
	g.playerImg = playerImage
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
		g.ai.Control(pet.DirectionUp, g.walls)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.ai.Control(pet.DirectionDown, g.walls)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.ai.Control(pet.DirectionLeft, g.walls)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.ai.Control(pet.DirectionRight, g.walls)
	}

	g.ai.Update(g.walls)

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backgroundColor)
	body := o.Shape.Body()
	pos := body.Position()

	vector.FillCircle(
		screen,
		float32(pos.X),
		float32(pos.Y),
		20,
		color.White,
		true,
	)

}

func Start() {
	ebiten.SetWindowSize(900, 500)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
