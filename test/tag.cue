package test

InlineStruct: {
	// Field1 comment
	// Field1 doc
	field1: string
	field2: string
	// Field3 doc
	field3: {
		// Field3.Field1 comment
		// Field3.Field1 doc
		field1: string
		// Field3.Field2 comment
		field2: string
	}
}
InlineStruct2: {
	// Field1 doc
	field1: string
	{
		// Field1 comment
		// Field1 doc
		field1: string
		field2: string
		// Field3 doc
		field3: {
			// Field3.Field1 comment
			// Field3.Field1 doc
			field1: string
			// Field3.Field2 comment
			field2: string
		}
	}
}
Inline: {
	field1: string
	{
		// Field1 comment
		// Field1 doc
		field1: string
		field2: string
		// Field3 doc
		field3: {
			// Field3.Field1 comment
			// Field3.Field1 doc
			field1: string
			// Field3.Field2 comment
			field2: string
		}
	}
	{
		// Field1 doc
		field1: string
		{
			// Field1 comment
			// Field1 doc
			field1: string
			field2: string
			// Field3 doc
			field3: {
				// Field3.Field1 comment
				// Field3.Field1 doc
				field1: string
				// Field3.Field2 comment
				field2: string
			}
		}
	}
}

output: InlineStruct2 & {field2: "ssss"}
