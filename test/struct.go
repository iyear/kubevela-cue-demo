package test

import (
	"crypto"
	"net/http"
)

type BasicType struct {
	Field1  string  `cue:"field1"`
	Field2  int     `cue:"field2"`
	Field3  bool    `cue:"field3"`
	Field4  float32 `cue:"field4"`
	Field5  float64 `cue:"field5"`
	Field6  int8    `cue:"field6"`
	Field7  int16   `cue:"field7"`
	Field8  int32   `cue:"field8"`
	Field9  int64   `cue:"field9"`
	Field10 uint    `cue:"field10"`
	Field11 uint8   `cue:"field11"`
	Field12 uint16  `cue:"field12"`
	Field13 uint32  `cue:"field13"`
	Field14 uint64  `cue:"field14"`
	Field15 uintptr `cue:"field15"`
	Field16 byte    `cue:"field16"`
	Field17 rune    `cue:"field17"`
}

type TagName struct {
	Field1 string `cue:"f1"`
	Field2 string `cue:"f2"`
	Field3 string `cue:"f3"`
}

type SliceAndArray struct {
	Field1  []string   `cue:"field1"`
	Field2  [3]string  `cue:"field2"`
	Field3  []int      `cue:"field3"`
	Field4  [3]int     `cue:"field4"`
	Field5  []bool     `cue:"field5"`
	Field6  [3]bool    `cue:"field6"`
	Field7  []float32  `cue:"field7"`
	Field8  [3]float32 `cue:"field8"`
	Field9  []float64  `cue:"field9"`
	Field10 [3]float64 `cue:"field10"`
}

type SmallStruct struct {
	Field1 string `cue:"field1"`
	Field2 string `cue:"field2"`
}

type AnonymousField struct {
	SmallStruct
}

type ReferenceField struct {
	Field1 *SmallStruct `cue:"field1"`
}

type StructField struct {
	Field1 SmallStruct  `cue:"field1"`
	Field2 *SmallStruct `cue:"field2"`
}

type EmbedStruct struct {
	Field1 struct {
		Field1 string `cue:"field1"`
		Field2 string `cue:"field2"`
	} `cue:"field1"`
	Field2 struct {
		Field1 string `cue:"field1"`
		Field2 string `cue:"field2"`
		Field3 struct {
			Field1 string `cue:"field1"`
			Field2 string `cue:"field2"`
			Field3 struct {
				Field1 string `cue:"field1"`
				Field2 string `cue:"field2"`
				Field3 struct {
					Field1 string `cue:"field1"`
					Field2 string `cue:"field2"`
				} `cue:"field3"`
			} `cue:"field3"`
		} `cue:"field3"`
	} `cue:"field2"`
	Field3 http.Header `cue:"field3"`
	Field4 crypto.Hash `cue:"field4"`
	// Field5 *http.Request `cue:"field5"` // inf loop currently
}
