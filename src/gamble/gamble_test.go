package gamble

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gamble", func() {
	var (
		document string
		node Node
		err error
	)

	Context("parsing a single should value", func() {
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

	Context("parsing a sequence of strings", func() {
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

	Context("parsing a map of string to strings", func() {
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

	Context("parsing a nested map", func() {
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

	Context("parsing an invalid document", func() {
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
