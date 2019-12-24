package main

import (
	"github.com/GianluigiMemoli/ponGO/scene"
)

func main(){
	myScene := scene.NewScene()
	myScene.SetupScene()
	myScene.Animate()
	myScene.Shutdown()
}


