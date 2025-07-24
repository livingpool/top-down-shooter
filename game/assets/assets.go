package assets

import (
	"embed"
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed PNG
var assets embed.FS

func MustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		log.Fatalf("error loading asset: %v", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("error decoding image: %v", err)
	}

	return ebiten.NewImageFromImage(img)
}

func MustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		log.Fatalf("fs.Glob failed: %v", err)
	}
	slog.Info("found matching files", "count", len(matches), "path", path)

	re := regexp.MustCompile(`tile_(\d+)\.png`)
	slices.SortFunc(matches, func(a, b string) int {
		matchesA := re.FindStringSubmatch(a)
		matchesB := re.FindStringSubmatch(b)

		idxA, _ := strconv.Atoi(matchesA[1])
		idxB, _ := strconv.Atoi(matchesB[1])

		if idxA < idxB {
			return -1
		} else if idxA > idxB {
			return 1
		} else {
			return 0
		}
	})

	// some tile numbers are missing, e.g. tile_27.png,
	// so for those, i filled in nil.
	lastImgNum, _ := strconv.Atoi(re.FindStringSubmatch(matches[len(matches)-1])[1])

	// images start from index 1, i.e., tile_01.png
	images := make([]*ebiten.Image, lastImgNum)
	cursor := 0
	for i := range lastImgNum {
		idx := strconv.Itoa(i + 1)
		if len(idx) < 2 {
			idx = "0" + idx
		}

		if strings.Contains(matches[cursor], idx) {
			images[i] = MustLoadImage(matches[cursor])
			cursor++
		} else {
			images[i] = nil
			slog.Info("nil image (0-based)", "idx", i)
		}
	}

	slog.Info("loaded images", "count", len(images))
	return images
}

func MustLoadFont(name string) font.Face {
	f, err := assets.ReadFile(name)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	tt, err := opentype.Parse(f)
	if err != nil {
		log.Fatalf("error parsing font: %v", err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingVertical,
	})

	if err != nil {
		log.Fatalf("error generating font face: %v", err)
	}

	return face
}

/* Bullet */

var Bullet = MustLoadImage(`PNG/Tiles/tile_187.png`)

/* Background */

var Tiles = MustLoadImages(`PNG/Tiles/*`)

/* Humanoid */

var ManBlueGunSprite = MustLoadImage(`PNG/Man Blue/manBlue_gun.png`)
var ManBlueHoldSprite = MustLoadImage(`PNG/Man Blue/manBlue_hold.png`)
var ManBlueMachineSprite = MustLoadImage(`PNG/Man Blue/manBlue_machine.png`)
var ManBlueReloadSprite = MustLoadImage(`PNG/Man Blue/manBlue_reload.png`)
var ManBlueSilencerSprite = MustLoadImage(`PNG/Man Blue/manBlue_silencer.png`)
var ManBlueStandSprite = MustLoadImage(`PNG/Man Blue/manBlue_stand.png`)

var Zombie1GunSprite = MustLoadImage(`PNG/Zombie 1/zoimbie1_gun.png`)
var Zombie1HoldSprite = MustLoadImage(`PNG/Zombie 1/zoimbie1_hold.png`)
var Zombie1MachineSprite = MustLoadImage(`PNG/Zombie 1/zoimbie1_machine.png`)
var Zombie1ReloadSprite = MustLoadImage(`PNG/Zombie 1/zoimbie1_reload.png`)
var Zombie1SilencerSprite = MustLoadImage(`PNG/Zombie 1/zoimbie1_silencer.png`)
var Zombie1StandSprite = MustLoadImage(`PNG/Zombie 1/zoimbie1_stand.png`)
