# dike

dike is a simple library for generating golang structure verification, it does not rely on references when verifying

# condition or tags

| name | description                | support type                                                                                                                             |
|------|----------------------------|------------------------------------------------------------------------------------------------------------------------------------------|
| re   | required                   | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string<br/>slice,array,struct,pointer |
| op   | optional                   | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string<br/>slice,array,struct,pointer |
| lt   | less than                  | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string                                |
| lte  | less than or equal to      | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string                                |
| gt   | great than                 | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string                                |
| gte  | great than or equal to     | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string                                |
| eq   | equal                      | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string,<br/>array                     |
| neq  | not equal to               | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string<br/>array                      |
| bet  | between the two            | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string                                |
| in   | in the collection          | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string                                |
| ni   | not in the collection      | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string                                |
| len  | specific length limit      | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string<br/>slice,array                |
| size | interval length limit      | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string<br/>slice,array                |
| reg  | regular matching condition | string                                                                                                                                   |
| dc   | field prompt description   | int8,uint8<br/>int16,uint16<br/>int32,uint32<br/>int,uint<br/>int64,uint64<br/>float32,float64<br/>string<br/>slice,array,struct,pointer |

# usage

```
At present, it is necessary to generate code with the help of go test, and it should be noted that when the struct structure is similar to type A struct {B *B}, the parameter &A{B:&B} needs to be passed when calling dike.GenerateCheck!
for example:
func Test_GenCheckA(t *testing.T) {
     dike.GenerateCheck(&A{B:&B{})
}
```

# current defect

The code repetition rate is high, and some codes have not been optimized
