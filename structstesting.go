package main

import "fmt"

type S1 struct {
	T1 []*S2
	T2 []byte
}

type S2 struct {
	T3 int
	T4 *S3
}

type S3 struct {
	T5 int
	T6 int
	T7 int
}

func main() {
	a := map[string]S1{}
	s1 := S1{}
	for i := 0; i < 3; i++ {
		s3 := S3{T5: i, T7: i, T6: i}
		s2 := S2{T4: &s3}
		s1.T1 = append(s1.T1, &s2)
	}
	a["test"] = s1
	for _, value := range a {
		for _, d := range value.T1 {
			d.T4.T5 = 222
		}
	}
	for _, value := range a["test"].T1 {
		fmt.Printf("%+v", value.T4)
	}

}
