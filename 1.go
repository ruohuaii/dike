package dike

type TBuilder struct {
	Name    string `dike:"re;size:1,32;bet:1,32;in:1,2,3,4;ni:5,6,7,8;dc:name"`
	Age     uint8
	Like    string    `dike:"re;len:12;dc:like;size:1,32;bet:1,2;eq:1;reg:/d;in:唱,跳,rap,篮球"`
	FavCol  *Favorite `dike:"op;dc:fav_col"`
	Class   Class     `dike:"op;dc:class"`
	Friends []string  `dike:"op;size:1,3;dc:friends"`
}

type Favorite struct {
	Color string `dike:"re;size:2,4;dc:color"`
	Ball  string `dike:"op;ni:足球,排球;dc:ball"`
}

type Class struct {
	Grade int `dike:"re;lt:4;dc:grade"`
}
