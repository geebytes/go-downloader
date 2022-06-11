package godownloader

import (
	"testing"

	"github.com/geebytes/go-downloader/pb"
)

func TestDownload(t *testing.T) {
	downloader := NewDownloader("/data/work/geebytes/go-downloader")
	// req := NewRequest("https://d-e01.winudf.com/b/XAPK/Y29tLmFjdGl2aXNpb24uY2FsbG9mZHV0eS5zaG9vdGVyXzE0MTYyXzcwMDE3MGY4?_fn=Q2FsbCBvZiBEdXR5IE1vYmlsZSBTZWFzb24gNV92MS4wLjMzX2Fwa3B1cmUuY29tLnhhcGs&_p=Y29tLmFjdGl2aXNpb24uY2FsbG9mZHV0eS5zaG9vdGVy&am=l9vAeeKgFjE8VG2ObBc-SQ&at=1654742655&download_id=otr_1546801073013106&k=a69f740b6565596bf16b122ecb08ea9862a2b000&r=https%3A%2F%2Fapkpure.com%2F&uu=https%3A%2F%2Fd-34.winudf.com%2Fb%2FXAPK%2FY29tLmFjdGl2aXNpb24uY2FsbG9mZHV0eS5zaG9vdGVyXzE0MTYyXzcwMDE3MGY4%3Fk%3D74f99af48548cd7ed04f3a7e0a8207a062a2b000", "GET", RequestWithRequestProxy(Proxy{ProxyUrl: "http://127.0.0.1:1081"}))
	local, err := downloader.Download(&pb.Request{
		Url:    "https://lh3.googleusercontent.com/i5dYZRkVCUK97bfprQ3WXyrT9BnLSZtVKGJlKQ919uaUB0sxbngVCioaiyu9r6snqfi2aaTyIvv6DHm4m2R3y7hMajbsv14pSZK8mhs=s10000?fit=max&h=2500&w=2500&auto=format&s=09f8240e4e777c5a9aa20f32bd6e8148",
		Method: "GET",
		Proxy:  &pb.Proxy{ProxyUrl: "http://127.0.0.1:1081"},
	}, "test03")
	if err != nil {
		t.Error(err)
	}
	t.Log(local)
}
