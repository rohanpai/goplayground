package main

import "github.com/tapir/glfw3-go"

import gl "github.com/chsc/gogl/gl21"
import "runtime"

func main() {
	runtime.LockOSThread()
	if !glfw.Init() {
		println("glfw init failure")
	}
	defer glfw.Terminate()
	//glfw.SetErrorCallback(func(err int, desc string) { panic(desc) })
	window := glfw.CreateWindow(500, 500, "egles", nil, nil)
	defer window.Destroy()

	window.MakeContextCurrent()

	for {
		if window.ShouldClose() {
			break
		}
		width, height := window.GetSize()
		ratio := float32(width) / float32(height)
		println(window.GetParameter(glfw.ClientApi))

		gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(gl.Double(-ratio), gl.Double(ratio), -1.0, 1.0, 1.0, -1.0)
		gl.MatrixMode(gl.MODELVIEW)

		gl.LoadIdentity()
		gl.Rotatef(gl.Float(glfw.GetTime()*50), 0, 0, 1.0)

		gl.Begin(gl.TRIANGLES)
		gl.Color3f(1.0, 0.0, 0.0)
		gl.Vertex3f(-0.6, -0.4, 0.0)
		gl.Color3f(0.0, 1.0, 0.0)
		gl.Vertex3f(0.6, -0.4, 0.0)
		gl.Color3f(0.0, 0.0, 1.0)
		gl.Vertex3f(0.0, 0.6, 0.0)
		gl.End()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
