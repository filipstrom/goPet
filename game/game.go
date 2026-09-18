package game

import (
	"encoding/json"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"os/exec"

	"github.com/filipstrom/goPet/object"
	"github.com/filipstrom/goPet/pet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	space           *cp.Space
	circles         []object.Object
	decoder         *json.Decoder
	encoder         *json.Encoder
}

func (g *game) reset() {
	// g.ai.ResetBody()
}

func (g *game) initialize() {
	g.backgroundColor = color.RGBA{0, 181, 226, 255}
	body := cp.NewBody(1, 10)
	body.SetPosition(cp.Vector{X: 100, Y: 100})
	shape := cp.NewCircle(body, 20, cp.Vector{})
	shape.SetFriction(0.8)

	ob := object.Object{Appearance: []int{1, 2, 3}, Texture: []int{1, 2, 3}, Shape: shape}

	eatSensor := cp.NewCircle(
		body,
		8,                      // radius
		cp.Vector{X: 25, Y: 0}, // framför kroppen
	)

	//g.playerImg = loadImage("assets/pet.png")
	g.space = cp.NewSpace()
	g.space.SetDamping(0.05)
	g.space.AddBody(body)
	g.space.AddShape(shape)
	g.ai = pet.NewAI("NN", ob, g.space, eatSensor)

	eatSensor.SetSensor(true)

	g.space.AddShape(eatSensor)

	// making walls

	// Left wall
	leftBody := cp.NewStaticBody()
	leftBody.SetPosition(cp.Vector{X: 2.5, Y: 250})

	leftWall := cp.NewBox(leftBody, 5, 500, 0)

	g.space.AddBody(leftBody)
	g.space.AddShape(leftWall)

	// Right wall
	rightBody := cp.NewStaticBody()
	rightBody.SetPosition(cp.Vector{X: 500 - 2.5, Y: 250})

	rightWall := cp.NewBox(rightBody, 5, 500, 0)

	g.space.AddBody(rightBody)
	g.space.AddShape(rightWall)

	// Top wall
	topBody := cp.NewStaticBody()
	topBody.SetPosition(cp.Vector{X: 250, Y: 2.5})

	topWall := cp.NewBox(topBody, 500, 5, 0)

	g.space.AddBody(topBody)
	g.space.AddShape(topWall)

	// Bottom wall
	bottomBody := cp.NewStaticBody()
	bottomBody.SetPosition(cp.Vector{X: 250, Y: 500 - 2.5})

	bottomWall := cp.NewBox(bottomBody, 500, 5, 0)

	g.space.AddBody(bottomBody)
	g.space.AddShape(bottomWall)

	// Adding wall
	bottomWall.UserData = object.ShapeData{Appearance: []float32{0, 0, 1}, IsFood: false}
	topWall.UserData = object.ShapeData{Appearance: []float32{0, 0, 1}, IsFood: false}
	leftWall.UserData = object.ShapeData{Appearance: []float32{0, 0, 1}, IsFood: false}
	rightWall.UserData = object.ShapeData{Appearance: []float32{0, 0, 1}, IsFood: false}

	// Making food
	spawnFood(g)
	g.fullscreen = true
	g.initialized = true

	// Making brain
	cmd := exec.Command("python", "brain.py")
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(stdout)
	g.decoder = decoder
	encoder := json.NewEncoder(stdin)
	g.encoder = encoder

}

func spawnFood(g *game) {
	foodBody := cp.NewBody(1, 10)

	x := 30 + rand.Float64()*(500-60)
	y := 30 + rand.Float64()*(500-60)

	foodBody.SetPosition(cp.Vector{
		X: x,
		Y: y,
	})
	foodShape := cp.NewCircle(foodBody, 20, cp.Vector{})
	g.space.AddBody(foodBody)
	g.space.AddShape(foodShape)
	foodShape.UserData = object.ShapeData{Appearance: []float32{1, 0, 0}, IsFood: true}

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
		eated := g.ai.Eat(g.space)
		if eated {
			spawnFood(g)
		}

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

	// p := g.ai.Look(&g.space)
	// fmt.Println(p)

	inputs := g.ai.GetInputs()

	if err := g.encoder.Encode(inputs); err != nil {
		panic(err)
	}

	outputs := []float32{}

	if err := g.decoder.Decode(&outputs); err != nil {
		panic(err)
	}
	fmt.Println(outputs)

	g.ai.Move(outputs)

	g.space.Step(1.0 / 60.0)

	// g.ai.Update(g.walls)

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(g.backgroundColor)
	body := g.ai.GetBody()
	pos := body.Shape.Body().Position()

	vector.FillCircle(
		screen,
		float32(pos.X),
		float32(pos.Y),
		20,
		color.White,
		true,
	)
	leftEye := cp.Vector{X: 15, Y: -8}
	rightEye := cp.Vector{X: 15, Y: 8}
	leftWorld := g.ai.GetBody().Shape.Body().LocalToWorld(leftEye)
	rightWorld := g.ai.GetBody().Shape.Body().LocalToWorld(rightEye)

	g.space.EachShape(func(shape *cp.Shape) {
		switch s := shape.Class.(type) {

		case *cp.Circle:
			pos := s.TransformC()

			vector.FillCircle(
				screen,
				float32(pos.X),
				float32(pos.Y),
				float32(s.Radius()),
				color.White,
				true,
			)

		case *cp.PolyShape:
			// Walls/Rectangles
			if s.Count() == 4 {

				bb := s.BB()

				width := bb.R - bb.L
				height := bb.T - bb.B
				vector.FillRect(
					screen,
					float32(bb.L),
					float32(bb.B),
					float32(width),
					float32(height),
					color.White,
					false,
				)

			}

		case *cp.Segment:
			//fmt.Println("segment")
		}
	})

	ebitenutil.DebugPrintAt(screen, g.ai.GetMood().String(), 20, 20)
	// Eyes

	vector.FillCircle(
		screen,
		float32(leftWorld.X),
		float32(leftWorld.Y),
		3,
		color.Black,
		true,
	)

	vector.FillCircle(
		screen,
		float32(rightWorld.X),
		float32(rightWorld.Y),
		3,
		color.Black,
		true,
	)

	// g.ai.Look()

}

func Start() {
	ebiten.SetWindowSize(900, 500)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
