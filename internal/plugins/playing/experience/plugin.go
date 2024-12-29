package experience

import (
	"game/internal/assets"
	"game/internal/constants"
	"game/internal/core"
	"game/internal/plugins/playing/camera"
	"game/internal/plugins/playing/player"
	"image/color"
	"math"
	"math/rand"
	"time"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var crystalRadius = 10
var superCrystalRadius = 15

type ExperiencePlugin struct {
	kernel   *core.GameKernel
	crystals []*Crystal
	plugins  *core.PluginManager

	crystalAnimation      *assets.Animation
	superCrystalAnimation *assets.Animation
}

type Crystal struct {
	X, Y          float64
	Width, Height float64
	Active        bool
	Speed         float64
	Value         int
	animation     *assets.Animation
}

type SuperXPCrystal struct {
	X, Y          float64
	Width, Height float64
	Value         int
	Active        bool
	Speed         float64
	Size          float64
}

func NewExperiencePlugin(plugins *core.PluginManager) *ExperiencePlugin {
	return &ExperiencePlugin{
		crystals: []*Crystal{},
		plugins:  plugins,
	}
}

func (ep *ExperiencePlugin) ID() string {
	return "ExperienceSystem"
}

func (ep *ExperiencePlugin) Init(kernel *core.GameKernel) error {
	ep.kernel = kernel
	rand.Seed(time.Now().UnixNano())

	crystalAnimation := assets.NewAnimation(0.1)
	err := crystalAnimation.LoadFromJSON(
		"assets/images/experience/crystal/asset.json",
		"assets/images/experience/crystal/asset.png")
	if err != nil {
		log.Fatal("Failed to load crystal animation:", err)
	}

	superCrystalAnimation := assets.NewAnimation(0.1)
	err = superCrystalAnimation.LoadFromJSON(
		"assets/images/experience/supercrystal/asset.json",
		"assets/images/experience/supercrystal/asset.png")
	if err != nil {
		log.Fatal("Failed to load super crystal animation:", err)
	}

	ep.crystalAnimation = crystalAnimation
	ep.superCrystalAnimation = superCrystalAnimation

	return nil
}

func (ep *ExperiencePlugin) Update() error {
	playerPlugin := ep.plugins.GetPlugin("PlayerSystem").(*player.PlayerPlugin)

	playerX, playerY := playerPlugin.GetPosition()
	playerWidth, playerHeight := playerPlugin.GetSize()

	cameraPlugin := ep.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	// Group far crystals
	farCrystals := make([]*Crystal, 0)
	activeCrystals := make([]*Crystal, 0)

	for _, crystal := range ep.crystals {
		if !crystal.Active {
			continue
		}

		// Calculate screen position
		screenX := crystal.X - cameraX
		screenY := crystal.Y - cameraY

		// Check if crystal is far from screen
		if screenX < -500 || screenX > constants.ScreenWidth+500 ||
			screenY < -500 || screenY > constants.ScreenHeight+500 {
			farCrystals = append(farCrystals, crystal)
		} else {
			activeCrystals = append(activeCrystals, crystal)
		}
	}

	// Create super XP if enough far crystals
	if len(farCrystals) >= 5 {
		// Calculate total XP value
		totalXP := 0

		for _, fc := range farCrystals {
			totalXP += fc.Value
		}

		// Create super crystal in random position near screen
		margin := float64(0) // Margin from screen edge
		superX := cameraX + rand.Float64()*(constants.ScreenWidth-2*margin) + margin
		superY := cameraY + rand.Float64()*(constants.ScreenHeight-2*margin) + margin

		superCrystal := &Crystal{
			X:      superX,
			Y:      superY,
			Width:  float64(superCrystalRadius),
			Height: float64(superCrystalRadius),
			Active: true,
			Speed:  200,
			Value:  totalXP, // Bonus for collecting grouped XP
		}

		// Deactivate grouped crystals
		for _, crystal := range farCrystals {
			crystal.Active = false
		}

		// Add super crystal
		ep.crystals = append(activeCrystals, superCrystal)
	} else {
		ep.crystals = activeCrystals
	}

	// Update remaining crystals
	for _, crystal := range ep.crystals {
		if crystal.Active {
			if crystal.animation != nil {
				crystal.animation.Update(ep.kernel.DeltaTime)
			}

			// Move towards player if in range
			dx := (playerX + playerWidth/2) - (crystal.X + 5)
			dy := (playerY + playerHeight/2) - (crystal.Y + 5)
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance > 0 {
				dx /= distance
				dy /= distance
			}

			crystal.X += dx * crystal.Speed * ep.kernel.DeltaTime
			crystal.Y += dy * crystal.Speed * ep.kernel.DeltaTime

			if ep.inPlayerCollectionRadius(crystal, playerX, playerY, playerWidth, playerHeight) {
				crystal.Speed = 450
			}

			if ep.checkCollisionWithPlayer(crystal, playerX, playerY, playerWidth, playerHeight) {
				crystal.Active = false
				playerPlugin.AddExperience(crystal.Value)
			}
		}
	}

	return nil
}

func (ep *ExperiencePlugin) Draw(screen *ebiten.Image) {
	cameraPlugin := ep.plugins.GetPlugin("CameraSystem").(*camera.CameraPlugin)
	cameraX, cameraY := cameraPlugin.GetPosition()

	for _, crystal := range ep.crystals {
		if crystal.Active {
			screenX := crystal.X - cameraX
			screenY := crystal.Y - cameraY

			// Only draw if on screen (with margin)
			if screenX >= -crystal.Width && screenX <= constants.ScreenWidth+crystal.Width &&
				screenY >= -crystal.Height && screenY <= constants.ScreenHeight+crystal.Height {

				if ep.crystalAnimation != nil {
					ep.crystalAnimation.Draw(screen, assets.DrawInput{
						Width:  crystal.Width,
						Height: crystal.Height,
						X:      screenX,
						Y:      screenY,
					})
				} else {
					// Fallback to rectangle if animation fails
					vector.DrawFilledRect(screen,
						float32(screenX),
						float32(screenY),
						float32(crystal.Width),
						float32(crystal.Height),
						color.RGBA{0, 255, 255, 255},
						true)
				}
			}
		}
	}
}

func (ep *ExperiencePlugin) DropCrystal(x, y float64) {
	ep.crystals = append(ep.crystals, &Crystal{
		X:      x,
		Y:      y,
		Width:  float64(crystalRadius),
		Height: float64(crystalRadius),
		Active: true,
		Value:  1,
	})
}

func (ep *ExperiencePlugin) inPlayerCollectionRadius(crystal *Crystal, playerX, playerY, playerWidth, playerHeight float64) bool {
	collectionRadius := 50.0

	dx := (playerX + playerWidth/2) - (crystal.X + 5)
	dy := (playerY + playerHeight/2) - (crystal.Y + 5)

	distance := math.Sqrt(dx*dx + dy*dy)

	return distance <= collectionRadius
}

func (ep *ExperiencePlugin) checkCollisionWithPlayer(crystal *Crystal, playerX, playerY, playerWidth, playerHeight float64) bool {
	return crystal.X < playerX+playerWidth && crystal.X+10 > playerX &&
		crystal.Y < playerY+playerHeight && crystal.Y+10 > playerY
}
