package dike

type TBuilder struct {
	Name string `dike:"re;size:1,32;bet:1,32;in:1,2,3,4;ni:5,6,7,8;dc:name"`
	Age  uint8  `dike:"op;bet:1,32;ni:1,2,3,4;ni:5,6,7,8;dc:age;neq:18;"`
	Like string `dike:"re;len:12;dc:like;size:1,32;bet:1,2;eq:1;reg:/d;in:唱,跳,rap,篮球"`
}
