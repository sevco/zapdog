package zapdog

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ZapDog")
}

var _ = BeforeSuite(func() {
	// block all HTTP requests
	httpmock.Activate()
})

var _ = BeforeEach(func() {
	// remove any mocks
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

var _ = Describe("ZapDog", func() {
	It("should build url with properties", func() {
		u, _ := ddURL("https://base.url/", Options{
			Source:   "zadog",
			Service:  "unittest",
			Hostname: "unittest-hostname",
			Tags:     []string{"tag1:one", "tag2:two"},
		})
		Expect(u).To(Equal("https://base.url/?ddsource=zadog&ddtags=tag1%3Aone%2Ctag2%3Atwo&hostname=unittest-hostname&service=unittest"))
	})

	It("should buffer on write", func() {
		l, _ := NewDataDogLogger(context.TODO(), "", Options{})
		_, err := l.Write([]byte("This is a message"))
		if err != nil {
			Fail(err.Error())
		}
		Expect(l.Lines).To(Equal([]DataDogLog{{Message: "This is a message"}}))
	})

	It("should post on sync", func() {
		var lines []DataDogLog
		httpmock.RegisterResponder("POST", "https://base.url",
			func(req *http.Request) (*http.Response, error) {
				var requestLines []DataDogLog
				if err := json.NewDecoder(req.Body).Decode(&requestLines); err != nil {
					return httpmock.NewStringResponse(400, ""), nil
				}
				lines = append(lines, requestLines...)
				return httpmock.NewStringResponse(200, ""), nil
			})

		l, _ := NewDataDogLogger(context.TODO(), "", Options{Host: "https://base.url"})
		if _, err := l.Write([]byte("Message one")); err != nil {
			Fail(err.Error())
		}
		if err := l.Sync(); err != nil {
			Fail(err.Error())
		}
		Expect(lines).To(Equal([]DataDogLog{{Message: "Message one"}}))
		Expect(len(l.Lines), 0)
	})
})
