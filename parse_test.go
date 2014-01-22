package gamble_test

import (
	. "github.com/cloudfoundry/gamble"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parsing", func() {
	var (
		document string
		node Node
		err error
	)

	Describe("parsing a single should value", func() {
		BeforeEach(func() {
			document = "the_string"
			node, err = Parse(document)
		})

		It("should return a single string value", func() {
			Expect(node).To(Equal("the_string"))
		})

		It("should not have returned an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("parsing a sequence of strings", func() {
		BeforeEach(func() {
			document = `
- foo
- bar
- baz
`
			node, err = Parse(document)
		})

		It("should return a slice of strings", func() {
			Expect(node).To(Equal([]interface{}{"foo", "bar", "baz"}))
		})

		It("should not have returned an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("parsing a map of string to strings", func() {
		BeforeEach(func() {
			document = `
key1: value1
key2: value2
`
			node, err = Parse(document)
		})

		It("should return a map of string to interface{}", func() {
			Expect(node).To(Equal(map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			}))
		})

		It("should not have returned an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("parsing a nested map", func() {
		BeforeEach(func() {
			document = `
---
globals:
- taco
- burrito
- kimchi
collections:
- name: oceans
  locals:
    foo: bar
    bar: baz
  sequences:
  - one
  - two
  - three
- name: seas
  age: 55
`
			node, err = Parse(document)
		})

		It("should return a nested map", func() {
			Expect(node).To(Equal(map[string]interface{}{
				"globals": []interface{}{
					"taco",
					"burrito",
					"kimchi",
				},
				"collections": []interface{}{
					map[string]interface{}{
						"name": "oceans",
						"locals": map[string]interface{}{
							"foo": "bar",
							"bar": "baz",
						},
						"sequences": []interface{}{
							"one",
							"two",
							"three",
						},
					},
					map[string]interface{}{
						"name": "seas",
						"age":  "55",
					},
				},
			}))
		})

		It("should not have returned an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("parsing nulls", func() {
		It("returns nils", func() {
			node, err = Parse("some_key: null")
			Expect(err).NotTo(HaveOccurred())
			Expect(node).To(Equal(map[string]interface{}{ "some_key": nil }))
		})

		It("parses the string 'null' correctly when quoted", func() {
			node, err = Parse("some_key: \"null\"")
			Expect(node).To(Equal(map[string]interface{}{ "some_key": "null" }))

			node, err = Parse("some_key: 'null'")
			Expect(node).To(Equal(map[string]interface{}{ "some_key": "null" }))
		})
	})

	Describe("parsing empty values in a map", func() {
		It("returns nil", func() {
			node, err = Parse("some_key:")
			Expect(err).NotTo(HaveOccurred())
			Expect(node).To(Equal(map[string]interface{}{ "some_key": nil }))
		})
	})

	Describe("parsing an invalid document", func() {
		BeforeEach(func() {
			node, err = Parse(`
---
-
  -
		-
`)
		})

		It("should have returned an error", func() {
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Error parsing YAML."))
		})
	})
})

type MyStruct struct {
	Age int
	Name string
	Address Address
}

type Address struct {
	Street1 string
}

var _ = Describe("unmarshaling into a struct", func() {
	It("populates the fields of a struct", func() {
		document :=
`
---
age: 54
name: john
`
		var result MyStruct
		err := Unmarshal(document, &result)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Age).To(Equal(54))
		Expect(result.Name).To(Equal("john"))
	})

	It("populates nested structs", func() {
		document :=
`
---
age: 54
name: john
address:
  street1: 123 Fake St
`
		var result MyStruct
		err := Unmarshal(document, &result)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Address.Street1).To(Equal("123 Fake St"))
	})

	It("populates slices", func() {
		type structWithSlice struct {
			Things []string
		}

		document :=
`
---
things:
- thing1
- thing2
`

		var result structWithSlice
		err := Unmarshal(document, &result)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Things).To(Equal([]string{ "thing1", "thing2" }))
	})

	It("populates maps with string keys", func() {
		var myMap map[string]int
		document :=
`
---
thing1: 5
thing2: 6
`

		err := Unmarshal(document, &myMap)
		Expect(err).NotTo(HaveOccurred())
		Expect(myMap).To(Equal(map[string]int {
			"thing1": 5,
			"thing2": 6,
		}))
	})

	It("populates slices of structs", func() {
		type thing struct {
			Name string
			Age int
		}

		type structWithSlice struct {
			Things []thing
		}

		document :=
`
---
things:
- name: thing1
  age: 12
- name: thing2
  age: 54
`

		var result structWithSlice
		err := Unmarshal(document, &result)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Things).To(Equal([]thing{
			thing{ Name: "thing1", Age: 12 },
			thing{ Name: "thing2", Age: 54 },
	    }))
	})
})

