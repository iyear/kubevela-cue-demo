package test

type InlineStruct struct {
	// Field1 doc
	Field1 string `json:"field1"` // Field1 comment
	Field2 string `json:"field2"`
	// Field3 doc
	Field3 struct {
		// Field3.Field1 doc
		Field1 string `json:"field1"` // Field3.Field1 comment
		Field2 string `json:"field2"` // Field3.Field2 comment
	} `json:"field3"`
}

type InlineStruct2 struct {
	// Field1 doc
	Field1       string           `json:"field1"`
	InlineStruct `json:",inline"` // Field1 comment
}

type Inline struct {
	Field1        string `json:"field1"`
	InlineStruct  `json:",inline"`
	InlineStruct2 `json:",inline"`
}
