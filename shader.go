package gogl

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"io/ioutil"
	"os"
	"time"
)

type ShaderID uint32
type ProgramID uint32
type BufferID uint32

type Shader struct {
	id               ProgramID
	vertexPath       string
	fragmentPath     string
	vertexModified   time.Time
	fragmentModified time.Time
}

func NewShader(vertexPath string, fragmentPath string) (*Shader, error) {
	program, err := CreateProgram(vertexPath, fragmentPath)
	if err != nil {
		return nil, err
	}
	result := &Shader{program, vertexPath, fragmentPath, getModTime(vertexPath), getModTime(fragmentPath)}
	return result, nil
}

func LoadShader(path string, shaderType uint32) (ShaderID, error) {
	shaderFile, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	shaderFileStr := string(shaderFile)

	shader, err := CreateShader(shaderFileStr, shaderType)
	if err != nil {
		return 0, err
	}
	return shader, nil
}

func getModTime(filePath string) time.Time {
	file, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}
	return file.ModTime()
}
func (shader *Shader) Use() {
	UseProgram(shader.id)
}
func (shader *Shader) CheckShaderForChanges() {

	vertexModTime := getModTime(shader.vertexPath)
	fragmentModTime := getModTime(shader.fragmentPath)
	if !vertexModTime.Equal(shader.vertexModified) || !fragmentModTime.Equal(shader.fragmentModified) {
		program, err := CreateProgram(shader.vertexPath, shader.fragmentPath)
		if err != nil {
			fmt.Print("error create program\n")
		} else {
			gl.DeleteProgram(uint32(shader.id))
			shader.id = program

		}
	}
}
