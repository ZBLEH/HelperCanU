package gogl

import (
	"errors"
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
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
