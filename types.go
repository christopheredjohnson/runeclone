package main

type Point struct {
	X, Y int
}

type Recipe struct {
	Name   string
	Inputs []ItemSlot
	Output ItemSlot
}
