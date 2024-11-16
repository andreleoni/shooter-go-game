package assets

import (
	"encoding/json"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Frame struct {
	Filename         string `json:"filename"`
	Frame            Rect   `json:"frame"`
	Rotated          bool   `json:"rotated"`
	Trimmed          bool   `json:"trimmed"`
	SpriteSourceSize Rect   `json:"spriteSourceSize"`
	SourceSize       Size   `json:"sourceSize"`
	Duration         int    `json:"duration"`
}

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type Size struct {
	W int `json:"w"`
	H int `json:"h"`
}

type AnimationData struct {
	Frames []Frame `json:"frames"`
}

type Animation struct {
	Frames       []*ebiten.Image
	CurrentFrame int
	FrameTimer   float64
	FrameDelay   float64
}

func NewAnimation(frameDelay float64) *Animation {
	return &Animation{
		FrameDelay: frameDelay,
	}
}

func (a *Animation) LoadFromJSON(jsonPath, tilesetPath string) error {
	// Carregar o tileset
	file, err := assets.Open(tilesetPath)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	tileset := ebiten.NewImageFromImage(img)

	// Carregar e parsear o arquivo JSON
	jsonFile, err := assets.Open(jsonPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var animationData AnimationData
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&animationData); err != nil {
		return err
	}

	// Extrair subimagens do tileset com base nos dados do JSON
	for _, frame := range animationData.Frames {
		subImage := tileset.SubImage(image.Rect(frame.Frame.X, frame.Frame.Y, frame.Frame.X+frame.Frame.W, frame.Frame.Y+frame.Frame.H)).(*ebiten.Image)
		a.Frames = append(a.Frames, subImage)
	}

	return nil
}

func (a *Animation) Update(deltaTime float64) {
	a.FrameTimer += deltaTime
	if a.FrameTimer >= a.FrameDelay {
		a.FrameTimer = 0
		a.CurrentFrame = (a.CurrentFrame + 1) % len(a.Frames)
	}
}

func (a *Animation) Draw(
	screen *ebiten.Image,
	x, y, width, height float64,
	invertHorizontal bool) {

	op := &ebiten.DrawImageOptions{}

	if invertHorizontal {
		op.GeoM.Scale(-1, 1)
	}

	op.GeoM.Scale(width/float64(a.Frames[a.CurrentFrame].Bounds().Dx()), height/float64(a.Frames[a.CurrentFrame].Bounds().Dy()))

	op.GeoM.Translate(x, y)

	screen.DrawImage(a.Frames[a.CurrentFrame], op)
}
