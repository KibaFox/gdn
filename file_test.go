package gdn_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/kibafox/gdn"
)

var _ = Describe("File", func() {
	Context("TypeByExtension", func() {
		It("correctly gets the Markdown type", func() {
			Expect(gdn.TypeByExtension(".md")).Should(Equal(gdn.Markdown))
			Expect(gdn.TypeByExtension(".mkd")).Should(Equal(gdn.Markdown))
			Expect(gdn.TypeByExtension(".markdown")).Should(Equal(gdn.Markdown))
		})

		It("correctly returns the Unknown type", func() {
			Expect(gdn.TypeByExtension(".jpeg")).Should(Equal(gdn.Unknown))
			Expect(gdn.TypeByExtension(".txt")).Should(Equal(gdn.Unknown))
			Expect(gdn.TypeByExtension(".unknown")).Should(Equal(gdn.Unknown))
		})
	})

	Context("CopyFile", func() {
		It("copies a file", func() {
			tmp := tmpDir()
			defer os.RemoveAll(tmp)

			src := filepath.Join("testdata", "src", "example", "mytext.txt")
			dst := filepath.Join(tmp, "mytext.txt")

			Expect(gdn.CopyFile(src, dst)).To(Succeed())
			Expect(dst).Should(BeARegularFile())

			srcByt, err := ioutil.ReadFile(src)
			Expect(err).ShouldNot(HaveOccurred())

			dstByt, err := ioutil.ReadFile(dst)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(dstByt).Should(Equal(srcByt))
		})
	})

	Context("ReplaceExt", func() {
		Expect(gdn.ChExt("some/file", ".html")).
			Should(Equal("some/file.html"))
		Expect(gdn.ChExt("/asdf/qwer/markdown.md", ".html")).
			Should(Equal("/asdf/qwer/markdown.html"))
		Expect(gdn.ChExt("some/image.png", ".jpeg")).
			Should(Equal("some/image.jpeg"))
	})
})
