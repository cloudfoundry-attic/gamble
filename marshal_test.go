package gamble_test

import (
	. "github.com/cloudfoundry/gamble"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("marshaling", func() {
	It("marshals strings", func() {
		document, err := Marshal("some-string")
		Expect(document).To(Equal(
			`--- some-string
...
`,
		))
		Expect(err).NotTo(HaveOccurred())
	})

	It("turns maps into YAML", func() {
		map1 := map[string]interface{}{
			"key1": 5,
			"key2": 10.0,
			"key3": "cool",
			"key4": nil,
		}
		document, err := Marshal(map1)
		Expect(document).To(Equal(
			`---
key1: 5
key2: 10.00
key3: cool
key4: null
...
`,
		))
		Expect(err).NotTo(HaveOccurred())
	})

	It("turns slices into YAML", func() {
		map1 := map[string]interface{}{
			"key1": []interface{}{
				"my-string",
				5.0,
				nil,
			},
		}
		document, err := Marshal(map1)
		Expect(err).NotTo(HaveOccurred())
		Expect(document).To(Equal(
			`---
key1:
  - my-string
  - 5.00
  - null
...
`,
		))
	})

	It("gracefully handles unknown types", func() {
		_, err := Marshal(struct{}{})
		Expect(err).To(HaveOccurred())
	})
})
