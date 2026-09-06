package main

import "fmt"

type vehicle interface {
	drive()
	startEngine()
}

type engine struct {
}

type bikes struct {
	engine
	modeltype string
}

type car struct {
	e         engine
	modeltype string
}

func (e engine) startEngine() {
	fmt.Println("Engine started")
}

func (b *bikes) startEngine() {
	fmt.Println("riding bike")
}

func (b bikes) drive() {
	fmt.Println("Riding bike:", b.modeltype)
}

func (b *car) startEngine() {
	fmt.Println("riding car")
}

func (b car) drive() {
	fmt.Println("Riding car:", b.modeltype)
}

func main1() {
	/*bike := &bikes{engine{}, "yamaha"}
	car := &car{engine{}, "mercedes"}

	var v vehicle
	v = bike

	bike.startEngine()
	v.startEngine()

	v = car
	v.startEngine()*/

	// bikes embeds engine -> promoted
	bike := &bikes{engine{}, "yamaha"}
	bike.startEngine()        // calls bikes.startEngine (shadowing) if bikes defines it
	bike.engine.startEngine() // calls engine.startEngine explicitly

	// car has named field e -> no promotion
	car := &car{engine{}, "merc"}
	car.startEngine()   // calls car's startEngine (you have it)
	car.e.startEngine() // calls engine.startEngine (must go through e)
}
