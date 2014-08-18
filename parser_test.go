package vtex

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	spew.Config.SortKeys = true
}

func TestParser(t *testing.T) {
	Convey("parses vtex", t, func() {
		v := ParseFile("example.vtex")
		So(reflect.DeepEqual(v, exampleVTEX), ShouldBeTrue)
	})
}

var exampleVTEX = Element{
	Key: "CDmeVtex",
	Value: map[string]interface{}{
		"m_bNoLod": true,
		"m_inputTextureArray": []interface{}{
			Element{
				Key: "CDmeInputTexture",
				Value: map[string]interface{}{
					"m_colorSpace":          "linear",
					"m_fileName":            "materials\\\\particle\\\\particle_debris_burst\\\\particle_debris_burst_001.tga",
					"m_imageProcessorArray": []interface{}{},
					"m_name":                "SheetTexture",
					"m_typeString":          "2D",
				},
			},
		},
		"m_nOutputMaxDimension": int64(0),
		"m_nOutputMinDimension": int64(0),
		"m_outputClearColor":    Vector4{W: 0, X: 0, Y: 0, Z: 0},
		"m_outputFormat":        "DXT5",
		"m_outputTypeString":    "2D",
		"m_textureOutputChannelArray": []interface{}{
			Element{
				Key: "CDmeTextureOutputChannel",
				Value: map[string]interface{}{
					"CDmeImageProcessor": map[string]interface{}{
						"m_algorithm":  "",
						"m_stringArg":  "",
						"m_vFloat4Arg": Vector4{W: 0, X: 0, Y: 0, Z: 0},
					},
					"m_dstChannels": "rgba",
					"m_inputTextureArray": []string{
						"SheetTexture",
					},
					"m_outputColorSpace": "linear",
					"m_srcChannels":      "rgba",
				},
			},
		},
		"m_vClamp": Vector3{X: 0, Y: 0, Z: 0},
	},
}
