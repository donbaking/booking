package repository

//可以將功能寫在這可以讓其他程式使用函式，擴充功能很方便

type DatabaseRepo interface {
	ALLUsers() bool
}