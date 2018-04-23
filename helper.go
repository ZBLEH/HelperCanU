package gogl

import (
	"errors"
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
	//"image/jpeg"
	"image/png"
	"os"
)

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

func CreateShader(shadersource string, shaderType uint32) (ShaderID, error) {
	shadersource = shadersource + "\x00"
	shaderId := gl.CreateShader(shaderType)
	csource, free := gl.Strs(shadersource)
	gl.ShaderSource(shaderId, 1, csource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLenght int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLenght)
		log := strings.Repeat("\x00", int(logLenght+1))
		gl.GetShaderInfoLog(shaderId, logLenght, nil, gl.Str(log))
		return 0, errors.New("Failed to compil shader\n" + log)
	}
	return ShaderID(shaderId), nil
}

func CreateProgram(vertPath string, fragPath string) (ProgramID, error) {
	var success int32

	vert, err := LoadShader(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	frag, err := LoadShader(fragPath, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.LinkProgram(shaderProgram)
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLenght int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLenght)
		log := strings.Repeat("\x00", int(logLenght+1))
		gl.GetProgramInfoLog(shaderProgram, logLenght, nil, gl.Str(log))
		return 0, errors.New("Failed to compil\n" + log)
	}

	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))
	return ProgramID(shaderProgram), nil
}

func GenBindBuffer(target uint32) BufferID {
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(target, VBO)
	return BufferID(VBO)
}

func GenBindVertexArray() BufferID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return BufferID(VAO)
}

func GenEBO() uint32 {
	var EBO uint32
	gl.GenBuffers(1, &EBO)
	return EBO
}

func BindVertexArray(vao BufferID) {
	gl.BindVertexArray(uint32(vao))
}

func BufferDataFloat(target uint32, data []float32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}
func BufferDataInt(target uint32, data []uint32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

func UseProgram(prog ProgramID) {
	gl.UseProgram(uint32(prog))
}

func UnbindVertex() {
	gl.BindVertexArray(0)
}

func (shader *Shader) SetFloat(name string, f float32) {
	name_cstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(uint32(shader.id), name_cstr)
	gl.Uniform1f(location, f)
}


func GenBindText() TextureID {
	var textID uint32
	gl.GenTextures(1, &textID)
	gl.BindTexture(gl.TEXTURE_2D, textID)
	return TextureID(textID)
}

func BindTexture(id TextureID) {
	gl.BindTexture(gl.TEXTURE_2D, uint32(id))
}
func LoadTexture(filename string) TextureID {
	infile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer infile.Close()
	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r/256)
				bIndex++
			pixels[bIndex] = byte(g/256)
				bIndex++
			pixels[bIndex] = byte(b/256)
				bIndex++
			pixels[bIndex] = byte(a/256)

		}
	}
	texture := GenBindText()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))

	gl.GenerateMipmap(gl.TEXTURE_2D)
	return TextureID(texture)
}