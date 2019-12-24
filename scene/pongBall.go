package scene

type PongBall struct {
	x		int
	y		int
}

func (ball *PongBall) SetX(newX int){
	ball.x = newX
}

func (ball *PongBall) SetY(newY int){
	ball.y = newY 
}

func (ball *PongBall) X() int {
	return ball.x
}

func (ball *PongBall) Y() int {
	return ball.y
}




