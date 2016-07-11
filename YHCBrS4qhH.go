package main

import (
	&#34;fmt&#34;
	&#34;local/gl&#34;
	&#34;local/glfw&#34;
	&#34;runtime&#34;
	&#34;time&#34;
)

func main() {
	runtime.LockOSThread()
	glfw.SetErrorCallback(errorCallback)
	if !glfw.Init() {
		panic(&#34;INIT FAIL&#34;)
	}
	defer glfw.Terminate()
	window, err := glfw.CreateWindow(640, 480, &#34;Testing&#34;, nil, nil)
	if err != nil {
		panic(err)
	}
	if err := gl.Init(); err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// OPENGL STUFF
	gl.ClearColor(0.0, 0.0, 0.4, 1.0)

	// VERTEX SHADER
	VertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	vssrc := gl.GLString(vertexShaderSource)
	defer gl.GLStringFree(vssrc)
	gl.ShaderSource(VertexShader, 1, &amp;vssrc, nil)
	gl.CompileShader(VertexShader)
	var vstatus gl.Int
	gl.GetShaderiv(VertexShader, gl.COMPILE_STATUS, &amp;vstatus)
	fmt.Printf(&#34;Compiled Vertex Shader: %v\n&#34;, vstatus)

	// FRAGMENT SHADER
	FragShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fssrc := gl.GLString(fragmentShaderSource)
	defer gl.GLStringFree(fssrc)
	gl.ShaderSource(FragShader, 1, &amp;fssrc, nil)
	gl.CompileShader(FragShader)
	var fstatus gl.Int
	gl.GetShaderiv(FragShader, gl.COMPILE_STATUS, &amp;fstatus)
	fmt.Printf(&#34;Compiled Fragment Shader: %v\n&#34;, fstatus)

	// CREATE PROGRAM
	shaderprogram := gl.CreateProgram()
	gl.AttachShader(shaderprogram, VertexShader)
	gl.AttachShader(shaderprogram, FragShader)
	fragoutstring := gl.GLString(&#34;outColor&#34;)
	defer gl.GLStringFree(fragoutstring)
	gl.BindFragDataLocation(shaderprogram, gl.Uint(0), fragoutstring)

	gl.LinkProgram(shaderprogram)
	var linkstatus gl.Int
	gl.GetProgramiv(shaderprogram, gl.LINK_STATUS, &amp;linkstatus)
	fmt.Printf(&#34;Program Link: %v\n&#34;, linkstatus)

	/////////////////// WORKS FOR SOME REASON
	gl.VertexAttribPointer(gl.Uint(0), gl.Int(2), gl.FLOAT, gl.FALSE, gl.Sizei(0), gl.Pointer(&amp;verts[0]))
	gl.EnableVertexAttribArray(0)

	////////////////DOES NOT WORK
	// VERTEX BUFFER
	/*var vbo gl.Uint
	gl.GenBuffers(1, &amp;vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(len(verts)*4), gl.Pointer(&amp;verts), gl.STATIC_DRAW)

	// VERTEX ARRAY
	var vao gl.Uint
	gl.GenVertexArrays(1, &amp;vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, gl.FALSE, 0, nil)*/

	// USE PARTICULAR PROGRAM
	gl.UseProgram(shaderprogram)

	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, gl.Int(0), gl.Sizei(len(verts)*4))

		window.SwapBuffers()

		time.Sleep(100 * time.Millisecond)
		glfw.PollEvents()
	}
}

var verts = []gl.Float{
	0.0, 1.0,
	1.0, -1.0,
	-1.0, -1.0}

const (
	Title              = &#34;OpenGL Shader&#34;
	Width              = 640
	Height             = 480
	vertexShaderSource = `
#version 140
in vec2 position;
void main()
{
	gl_Position = vec4(position, 0.0, 1.0);
}
`
	fragmentShaderSource = `
#version 140
out vec4 outColor;
void main()
{
	outColor = vec4(1.0, 1.0, 0.0, 1.0);
}
`
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf(&#34;%v: %v\n&#34;, err, desc)
}
