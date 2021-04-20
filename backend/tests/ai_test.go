package tests

import "testing"

func TestGenericAi(t *testing.T) {
	var u1 = []byte("{\"name\":\"Lemi1\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":5}]}}")
	var u2 = []byte("{\"name\":\"Lemi2\",\"units\":{\"archer\":[],\"knight\":[{\"x\":5,\"y\":8}]}}")

	world := CreateTestWorld(u1, u2, t)

	world.Update()
}
