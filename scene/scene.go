package scene

import (
	"fmt"
	FileLogger "github.com/GianluigiMemoli/FileLog"
	"github.com/gdamore/tcell"
	"sync"
)

const (
	ballRune = '\u25CF'
	racketSize 	= 3
	toLeft 		= 'l'
	toRight 	= 'r'
	toTop		= 't'
	toBottom	= 'b'
	straight	= 's'
)
var log *FileLogger.FileLogger

type Component interface {
	SetX(newX int)
	SetY(newY int)
	X()		int
	Y()		int
}

type Scene struct {
	leftRacket 			Component
	rightRacket			Component
	ball 				Component
	myScreen			tcell.Screen
	width				int
	height				int
	ballDirX			rune
	ballDirY			rune
	ballXincrementer 	float32
	ballVelocity		float32
	player1Score		int
	player2Score		int
}

func NewScene() *Scene{
	log = FileLogger.NewLogger()
	leftRacket := &Racket{
		x: 0,
		y: 0,
	}

	rightRacket := &Racket{
		x: 0,
		y: 0,
	}

	ball := &PongBall{
		x: 0,
		y: 0,
	}

	newScr, err := tcell.NewScreen()

	if err != nil{
		fmt.Println("err:", err)
		panic(err)
	}

	return &Scene{
		leftRacket:  		leftRacket,
		rightRacket: 		rightRacket,
		ball:        		ball,
		myScreen:    		newScr,
		ballDirX:	 		toLeft,
		ballDirY:	 		straight,
		ballXincrementer: 	0,
		ballVelocity: 		0.001,
		player1Score: 		0,
		player2Score:		0,
	}
}

func (scene *Scene) initialPosition() {
	//Setting rackets' position on middle of screen height
	racketHeight := scene.height/2 - racketSize/2
	scene.leftRacket.SetY(racketHeight)
	scene.rightRacket.SetY(racketHeight)
	scene.rightRacket.SetX(scene.width-1)
	//Setting ball's position on center of screen @ToReview
	scene.ball.SetX(scene.width/2)
	scene.ball.SetY(scene.height/2)
	scene.ballDirY = straight
	scene.ballDirX = toLeft
}

func (scene *Scene) SetupScene(){
	err := scene.myScreen.Init()
	//Showing screen for getting Size
	scene.myScreen.Show()
	if err != nil {
		panic(err)
	}
	//Getting screen dimension to set Components' positions
	w, h := scene.myScreen.Size()
	scene.height = h
	scene.width = w
	scene.initialPosition()
}
func intToRune(score int) rune {
	switch score {
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	case 3:
		return '3'
	case 4:
		return '4'
	case 5:
		return '5'
	}
	return '0'
}
func (scene *Scene) DrawScene() {
	/*
		This method must be called periodically after every ball move
	*/
	scene.myScreen.Clear()
	//Drawing ball
	scene.myScreen.SetContent(scene.ball.X(), scene.ball.Y(), ballRune, nil, tcell.StyleDefault)
	//Drawing rackets
	for i := 0; i < racketSize; i++ {
		scene.myScreen.SetContent(scene.leftRacket.X(), scene.leftRacket.Y() + i , tcell.RuneBlock, nil,  tcell.StyleDefault)
		scene.myScreen.SetContent(scene.rightRacket.X(), scene.rightRacket.Y() + i, tcell.RuneBlock, nil,  tcell.StyleDefault)
	}
	//Drawing middle line
	for i := 0; i < scene.height; i++ {
		scene.myScreen.SetContent(scene.width/2, i , tcell.RuneVLine, nil,  tcell.StyleDefault.Dim(true))
	}
	//Players' Score
	halfScreen := scene.width/2
	halfCenter  := halfScreen - halfScreen/2
	scene.myScreen.SetContent(halfCenter, 0, intToRune(scene.player1Score), nil, tcell.StyleDefault)
	scene.myScreen.SetContent(halfCenter + halfScreen, 0, intToRune(scene.player2Score), nil, tcell.StyleDefault)
	//Showing screen modifications
	scene.myScreen.Show()
}


func (scene *Scene) ballMover(){
	scene.ballXincrementer += scene.ballVelocity
	if int(scene.ballXincrementer) >= 1 {
		//Horizontal ball behaviour
		switch scene.ballDirX {
		case toLeft:
			scene.ball.SetX(scene.ball.X() - 1)
		case toRight:
			scene.ball.SetX(scene.ball.X() + 1)
		}
		//Vertical ball behaviour
		switch scene.ballDirY {
		case toBottom:
			scene.ball.SetY(scene.ball.Y() + 1)
		case toTop:
			scene.ball.SetY(scene.ball.Y() - 1)
		case straight:
			break
		}

		//Checking if ball is colliding with the left or right racket
		if scene.ball.X() == scene.leftRacket.X()+1 {
			if scene.ball.Y() >= scene.leftRacket.Y() && scene.ball.Y() < scene.leftRacket.Y()+racketSize {
				scene.ballDirX = toRight
				scene.ballDirY = topOrBtm(scene.ball.Y(), scene.leftRacket.Y())
			}
		} else if scene.ball.X() == scene.rightRacket.X()-1 {
			if scene.ball.Y() >= scene.rightRacket.Y() && scene.ball.Y() < scene.rightRacket.Y()+racketSize {
				scene.ballDirX = toLeft
				scene.ballDirY = topOrBtm(scene.ball.Y(), scene.rightRacket.Y())
			}
		}
		//Checking if ball is colliding with top or bottom wall
		if scene.ball.X() > scene.leftRacket.X() && scene.ball.X() < scene.rightRacket.X() {
			if scene.ball.Y() == 0 {
				scene.ballDirY = toBottom
			} else if scene.ball.Y() == scene.height-1 {
				scene.ballDirY = toTop
			}
		}

		if scene.ball.X() == scene.width {
			scene.player1Score ++
			scene.initialPosition()
		} else if scene.ball.X() == 0 {
			scene.player2Score ++
			scene.initialPosition()
		}
		scene.ballXincrementer = 0
	}
}



func topOrBtm(ballY, racketY int) rune{
	racketCenter := racketSize / 2
	collisionPoint := ballY - racketY
	if collisionPoint % racketSize < racketCenter {
		return toTop
	} else if collisionPoint % racketSize == racketCenter {
		return straight
	} else {
		return toBottom
	}
}

func (scene *Scene) eventDispatcher(group *sync.WaitGroup, exiter chan bool){
	out := false
	exiter <- out
	scr := scene.myScreen
	for ; !out; {
		event := scr.PollEvent()
		switch typedEv := event.(type) {
		case *tcell.EventKey:
			if typedEv.Key() == tcell.KeyEsc {
				out = true
			} else {
				scene.keyboardController(typedEv)
			}
		}
		exiter <- out
	}
	group.Done()
}

func (scene *Scene) AI(){
	if scene.ball.Y() > scene.rightRacket.Y(){
		scene.keyboardController(makeEventKey('D'))
	} else if scene.ball.Y() < scene.rightRacket.Y(){
		scene.keyboardController(makeEventKey('U'))
	}
}

func makeEventKey(ch rune) *tcell.EventKey{
	return tcell.NewEventKey(tcell.KeyRune, ch, tcell.ModNone)
}

func (scene *Scene) keyboardController(kbdEv *tcell.EventKey){
	if kbdEv.Rune() == 'w' || kbdEv.Rune() == 'W' {
		if scene.leftRacket.Y() > 0 {
			scene.leftRacket.SetY(scene.leftRacket.Y() - 1)
		}
	} else if kbdEv.Rune() == 's' || kbdEv.Rune() == 'S' {
		if scene.leftRacket.Y() + racketSize < scene.height {
				scene.leftRacket.SetY(scene.leftRacket.Y() + 1)
		}
	} else if kbdEv.Rune() == 'U' {
		if scene.rightRacket.Y() > 0 {
			scene.rightRacket.SetY(scene.rightRacket.Y() - 1)
		}
	} else if kbdEv.Rune() == 'D' {
		if scene.rightRacket.Y() + racketSize  < scene.height {
			scene.rightRacket.SetY(scene.rightRacket.Y() + 1)
		}
	}
}


func (scene *Scene) Animate() {
	exiter := make(chan bool)
	var group sync.WaitGroup
	group.Add(1)
	go scene.eventDispatcher(&group, exiter)
	out := <- exiter
	for ;!out; {
		scene.ballMover()
		scene.AI()
		scene.DrawScene()
		select {
		case new := <-exiter:
			out = new
		default:
			out = false
		}
	}
	group.Wait()
}

func (scene *Scene) Shutdown(){
	scene.myScreen.Fini()
}
