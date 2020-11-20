package gdn_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"gitlab.com/kibafox/gdn"
)

var _ = Describe("Garden", func() {
	Context("Branch NewTree", func() {
		It("can create a root branch", func() {
			root := gdn.NewTree("mygarden/src", "mygarden/dest")
			Expect(root.Src).Should(Equal("mygarden/src"))
			Expect(root.Dst).Should(Equal("mygarden/dest"))
			Expect(root.Path).Should(Equal("/"))
			Expect(root.Branches).Should(BeEmpty())
			Expect(root.Leaves).Should(BeEmpty())
		})
	})

	Context("Branch Scan", func() {
		It("scans input paths", func() {
			root := gdn.NewTree("testdata/src", "tmp")

			Expect(root.Scan()).Should(Succeed())

			bi := func(element interface{}) string {
				return element.(*gdn.Branch).Path
			}

			li := func(element interface{}) string {
				return element.(*gdn.Leaf).Path
			}

			Expect(root.Branches).Should(MatchAllElements(bi, Elements{
				"/example": PointTo(MatchAllFields(Fields{
					"Src":      Equal("testdata/src/example"),
					"Dst":      Equal("tmp/example"),
					"Path":     Equal("/example"),
					"Branches": BeNil(),
					"Leaves": MatchAllElements(li, Elements{
						"/example/mytext.txt": PointTo(MatchAllFields(Fields{
							"Src":    Equal("testdata/src/example/mytext.txt"),
							"DstDir": Equal("tmp/example"),
							"Path":   Equal("/example/mytext.txt"),
							"Typ":    Equal(gdn.Unknown),
						})),

						"/example/mydoc.md": PointTo(MatchAllFields(Fields{
							"Src":    Equal("testdata/src/example/mydoc.md"),
							"DstDir": Equal("tmp/example"),
							"Path":   Equal("/example/mydoc.md"),
							"Typ":    Equal(gdn.Markdown),
						})),
					}),
				})),
			}))
		})

		It("gives error when the source path is not set", func() {
			root := gdn.NewTree("", "tmp")
			Expect(root.Scan()).Should(BeAssignableToTypeOf(gdn.ErrSrcNotSet))
		})

		It("gives error when the destination path is not set", func() {
			root := gdn.NewTree(filepath.Join("testdata", "src"), "")
			Expect(root.Scan()).Should(BeAssignableToTypeOf(gdn.ErrDstNotSet))
		})

		It("gives error when the scanning an empty directory", func() {
			tmp := tmpDir()
			defer os.RemoveAll(tmp)

			root := gdn.NewTree(tmp, "tmp")
			Expect(root.Grow()).Should(BeAssignableToTypeOf(gdn.ErrEmptyTree))
		})
	})

	Context("Branch Grow", func() {
		It("grows testdata/src to match testdata/expected", func() {
			tmp := tmpDir()
			defer os.RemoveAll(tmp)

			root := gdn.NewTree(filepath.Join("testdata", "src"), tmp)

			Expect(root.Scan()).To(Succeed())

			Expect(root.Grow()).To(Succeed())

			Expect(tmp).Should(MatchDir(filepath.Join("testdata", "expected")))
		})

		It("gives error when the source path is not set", func() {
			root := gdn.NewTree("", "tmp")
			Expect(root.Grow()).Should(BeAssignableToTypeOf(gdn.ErrSrcNotSet))
		})

		It("gives error when the destination path is not set", func() {
			root := gdn.NewTree(filepath.Join("testdata", "src"), "")
			Expect(root.Grow()).Should(BeAssignableToTypeOf(gdn.ErrDstNotSet))
		})

		It("warns when the branch is empty and not scanned", func() {
			root := gdn.NewTree(filepath.Join("testdata", "src"), "tmp")
			Expect(root.Grow()).Should(BeAssignableToTypeOf(gdn.ErrNotScanned))
		})
	})

	Context("Leaf", func() {
		It("sets destination extension to .html for markdown files", func() {
			leaf := gdn.Leaf{
				Src:    "asdf/my.md",
				DstDir: "qwer",
				Path:   "/my.md",
				Typ:    gdn.Markdown,
			}

			Expect(leaf.Dst()).Should(Equal("qwer/my.html"))
		})

		It("keeps the same destination extension for unknown files", func() {
			leaf := gdn.Leaf{
				Src:    "asdf/my.txt",
				DstDir: "qwer",
				Path:   "/my.txt",
				Typ:    gdn.Unknown,
			}

			Expect(leaf.Dst()).Should(Equal("qwer/my.txt"))
		})
	})

	Context("Leaf Grow", func() {
		It("copies unknown files", func() {
			tmp := tmpDir()
			defer os.RemoveAll(tmp)

			src := filepath.Join("testdata", "src", "example", "mytext.txt")
			dst := filepath.Join(tmp, "mytext.txt")

			leaf := gdn.Leaf{
				Src:    src,
				DstDir: tmp,
				Path:   "/example/mytext.txt",
				Typ:    gdn.Unknown,
			}

			Expect(leaf.Grow()).Should(Succeed())
			Expect(dst).Should(BeARegularFile())

			srcByt, err := ioutil.ReadFile(src)
			Expect(err).ShouldNot(HaveOccurred())

			dstByt, err := ioutil.ReadFile(dst)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(dstByt).Should(Equal(srcByt))
		})

		It("renders markdown files", func() {
			tmp := tmpDir()
			defer os.RemoveAll(tmp)

			src := filepath.Join("testdata", "src", "example", "mydoc.md")
			dst := filepath.Join(tmp, "mydoc.html")

			leaf := gdn.Leaf{
				Src:    src,
				DstDir: tmp,
				Path:   "/example/mytext.md",
				Typ:    gdn.Markdown,
			}

			Expect(leaf.Grow()).Should(Succeed())
			Expect(dst).Should(BeARegularFile())

			dstByt, err := ioutil.ReadFile(dst)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(string(dstByt)).Should(Equal(`<h1>My Document</h1>

<p>This is <em>my</em> document with some <strong>Markdown</strong>.</p>
`))
		})

		It("gives error when the source path is not set", func() {
			leaf := gdn.Leaf{Src: "", DstDir: "tmp"}
			Expect(leaf.Grow()).Should(BeAssignableToTypeOf(gdn.ErrSrcNotSet))
		})

		It("gives error when the destination path is not set", func() {
			leaf := gdn.Leaf{Src: "tmp", DstDir: ""}
			Expect(leaf.Grow()).Should(BeAssignableToTypeOf(gdn.ErrDstNotSet))
		})
	})
})
