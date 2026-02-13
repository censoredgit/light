package controller

import "net/http"

type BackRedirectResponse struct {
	*RedirectResponse
}

func newBackRedirectResponse(req *http.Request, flashStorage ContextFlashStorage) *BackRedirectResponse {
	backUrl := req.Referer()
	if backUrl == "" {
		backUrl = *rootAction.uri
	}

	return &BackRedirectResponse{
		RedirectResponse: newRedirect(backUrl, flashStorage),
	}
}

func (b *BackRedirectResponse) With(it func(response ResponseExtendData)) Response {
	it(b)
	return b
}

func (b *BackRedirectResponse) Url() string {
	return b.url
}
