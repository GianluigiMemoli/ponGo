package scene

type Racket struct {
	x		int
	y		int
}

func (racket *Racket) SetX(newX int){
	racket.x = newX
}

func (racket *Racket) SetY(newY int){
	racket.y = newY
}

func (racket *Racket) X() int{
	return racket.x
}

func (racket *Racket) Y() int{
	return racket.y
}




