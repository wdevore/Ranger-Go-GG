package tests

import (
	"math"
	"testing"

	"github.com/wdevore/GameEngine/engine"
)

func Test_Translate(t *testing.T) {
	v := engine.NewVector3()
	at := engine.NewAffineTransform()
	t.Logf("v: %s\n", v)
	at.SetToTranslate(2.0, 0.0)
	at.ApplyToVector(v)
	t.Logf("v: %s\n", v)
	if v.X != 2.0 {
		t.Error("Expected v.X == 2")
	}
}

func Test_Rotate(t *testing.T) {
	v := engine.NewVector3()
	v.Set2Components(1.0, 0.0)
	at := engine.NewAffineTransform()
	t.Logf("v: %s\n", v)
	angle := 45.0 * (math.Pi / 180.0)
	t.Logf("Angle: %f\n", angle)
	at.SetToRotate(angle)
	t.Logf("at: \n%s\n", at)
	at.ApplyToVector(v)
	t.Logf("v: %s\n", v)

	// Assuming +Y is downward then:
	if v.X != math.Cos(angle) {
		t.Error("Expected v.X == 0.707107")
	}
	if v.Y != -math.Sin(angle) {
		t.Error("Expected v.Y == -0.707107")
	}
}
