package main

import (
	"unsafe"

	"github.com/bbredesen/go-vk"
	"github.com/bbredesen/vkm"
)

type Vertex struct {
	Pos   vkm.Pt2
	Color vkm.Vec3
}

var verts = []Vertex{
	{vkm.NewPt2(-0.5, -0.5), vkm.Vec3{1.0, 0.0, 0.0}},
	{vkm.NewPt2(0.5, -0.5), vkm.Vec3{0.0, 1.0, 0.0}},
	{vkm.NewPt2(0.5, 0.5), vkm.Vec3{0.0, 0.0, 1.0}},
	{vkm.NewPt2(-0.5, 0.5), vkm.Vec3{1.0, 1.0, 1.0}},
}

var indices = []uint16{0, 1, 2, 2, 3, 0}

func getVertexBindingDescription() vk.VertexInputBindingDescription {
	return vk.VertexInputBindingDescription{
		Binding:   0,
		Stride:    uint32(unsafe.Sizeof(Vertex{})),
		InputRate: vk.VERTEX_INPUT_RATE_VERTEX,
	}
}

func getVertexAttributeDescriptions() []vk.VertexInputAttributeDescription {
	rval := make([]vk.VertexInputAttributeDescription, 2)

	// Attributes of the vertex points; format R32G32 indicates that we have two 32 bit signed floating point components
	// If, for example, Vertex.Pos was declared as vkm.Pt (four components under the hood), then we would use vk.FORMAT_R32G32B32A32_SFLOAT
	rval[0] = vk.VertexInputAttributeDescription{
		Location: 0,
		Binding:  0,
		Format:   vk.FORMAT_R32G32_SFLOAT,
		Offset:   uint32(unsafe.Offsetof(Vertex{}.Pos)),
	}
	rval[1] = vk.VertexInputAttributeDescription{
		Location: 1,
		Binding:  0,
		Format:   vk.FORMAT_R32G32B32_SFLOAT,
		Offset:   uint32(unsafe.Offsetof(Vertex{}.Color)),
	}

	return rval
}

func (app *App_03) createVertexBuffer() {

	size := vk.DeviceSize(unsafe.Sizeof(Vertex{})) * vk.DeviceSize(len(verts))

	stagingBuffer, stagingBufferMemory := app.createBuffer(vk.BUFFER_USAGE_TRANSFER_SRC_BIT, size, vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)
	r, ptr := vk.MapMemory(app.device, stagingBufferMemory, 0, size, 0)
	if r != vk.SUCCESS {
		panic("Could not map memory for vertex buffer: " + r.String())
	}

	vk.MemCopySlice(unsafe.Pointer(ptr), verts)

	vk.UnmapMemory(app.device, stagingBufferMemory)

	app.vertexBuffer, app.vertexBufferMemory = app.createBuffer(vk.BUFFER_USAGE_VERTEX_BUFFER_BIT|vk.BUFFER_USAGE_TRANSFER_DST_BIT, size, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)

	app.copyBuffer(stagingBuffer, app.vertexBuffer, size)

	vk.DestroyBuffer(app.device, stagingBuffer, nil)
	vk.FreeMemory(app.device, stagingBufferMemory, nil)

}

func (app *App_03) createIndexBuffer() {

	size := vk.DeviceSize(unsafe.Sizeof(uint16(0))) * vk.DeviceSize(len(indices))

	stagingBuffer, stagingBufferMemory := app.createBuffer(vk.BUFFER_USAGE_TRANSFER_SRC_BIT, size, vk.MEMORY_PROPERTY_HOST_VISIBLE_BIT|vk.MEMORY_PROPERTY_HOST_COHERENT_BIT)
	r, ptr := vk.MapMemory(app.device, stagingBufferMemory, 0, size, 0)
	if r != vk.SUCCESS {
		panic("Could not map memory for vertex buffer: " + r.String())
	}

	vk.MemCopySlice(unsafe.Pointer(ptr), indices)

	// // MapMemory needs to return a []byte maybe...shouldn't have to do this:
	// var sl = struct {
	// 	addr uintptr
	// 	len  int
	// 	cap  int
	// }{uintptr(unsafe.Pointer(ptr)), int(size), int(size)}
	// bytes := *(*[]byte)(unsafe.Pointer(&sl))

	// vb := AnySliceToBytes(indices)
	// copy(bytes, vb)

	vk.UnmapMemory(app.device, stagingBufferMemory)

	app.indexBuffer, app.indexBufferMemory = app.createBuffer(vk.BUFFER_USAGE_INDEX_BUFFER_BIT|vk.BUFFER_USAGE_TRANSFER_DST_BIT, size, vk.MEMORY_PROPERTY_DEVICE_LOCAL_BIT)

	app.copyBuffer(stagingBuffer, app.indexBuffer, size)

	vk.DestroyBuffer(app.device, stagingBuffer, nil)
	vk.FreeMemory(app.device, stagingBufferMemory, nil)

}

func (app *App_03) findMemoryType(typeFilter uint32, flags vk.MemoryPropertyFlags) uint32 {
	memProps := vk.GetPhysicalDeviceMemoryProperties(app.physicalDevice)
	var i uint32
	for i = 0; i < memProps.MemoryTypeCount; i++ {
		if typeFilter&1<<i != 0 && memProps.MemoryTypes[i].PropertyFlags&flags == flags {
			return i
		}
	}
	panic("Could not find appropriate memory type.")
}

// // VERY MUCH TODO: This or something similar should be a convenience function provided with go-vk
// func AnySliceToBytes[T any](input []T) []byte {
// 	type sl struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}

// 	if len(input) == 0 {
// 		return []byte{}
// 	}

// 	inputLen := len(input) * int(unsafe.Sizeof(input[0]))
// 	sl_input := sl{uintptr(unsafe.Pointer(&input[0])), inputLen, inputLen}

// 	return *(*[]byte)(unsafe.Pointer(&sl_input))
// }

// func AnyTypeToBytes[T any](input T) []byte {
// 	type sl struct {
// 		addr uintptr
// 		len  int
// 		cap  int
// 	}
// 	sl_input := sl{uintptr(unsafe.Pointer(&input)), int(unsafe.Sizeof(input)), int(unsafe.Sizeof(input))}
// 	return *(*[]byte)(unsafe.Pointer(&sl_input))

// }
