package controllers

import (
	"app/base/core"
	"app/base/utils"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdvisoriesExportJSON(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequest("GET", "/", nil, "application/json", AdvisoriesExportHandler)

	var output []AdvisoriesDBLookupV3
	CheckResponse(t, w, http.StatusOK, &output)

	assert.Equal(t, 12, len(output))
	assert.Equal(t, output[2].ID, "RH-1")
	assert.Equal(t, output[2].Description, "adv-1-des")
	assert.Equal(t, output[2].Synopsis, "adv-1-syn")
	assert.Equal(t, output[2].AdvisoryTypeName, "enhancement")
	assert.Equal(t, output[2].CveCount, 0)
	assert.Equal(t, output[2].RebootRequired, false)
	assert.Equal(t, output[2].ReleaseVersions, RelList{"7.0", "7Server"})
	assert.Equal(t, output[2].InstallableSystems, 4)
	assert.Equal(t, output[2].ApplicableSystems, 2+output[2].InstallableSystems)
}

func TestAdvisoriesExportCSV(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequest("GET", "/", nil, "text/csv", AdvisoriesExportHandler)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	lines := strings.Split(body, "\r\n")

	assert.Equal(t, 14, len(lines))
	assert.Equal(t, "RH-1,adv-1-des,2016-09-22T16:00:00Z,adv-1-syn,enhancement,,0,false,\"7.0,7Server\",4,6", lines[3])
}

func TestAdvisoriesExportWrongFormat(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequest("GET", "/", nil, "test-format", AdvisoriesExportHandler)

	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
	body := w.Body.String()
	assert.Equal(t, InvalidContentTypeErr, body)
}

func TestAdvisoriesExportCSVFilter(t *testing.T) {
	core.SetupTest(t)

	for _, URL := range []string{"/?filter[id]=RH-1", "/?filter[synopsis]=adv-1-syn"} {
		w := CreateRequest("GET", URL, nil, "text/csv", AdvisoriesExportHandler)

		assert.Equal(t, http.StatusOK, w.Code)
		body := w.Body.String()
		lines := strings.Split(body, "\r\n")

		assert.Equal(t, 3, len(lines))
		assert.Equal(t, "RH-1,adv-1-des,2016-09-22T16:00:00Z,adv-1-syn,enhancement,,0,false,\"7.0,7Server\",4,6", lines[1])
		assert.Equal(t, "", lines[2])
	}
}

func TestAdvisoriesExportTagsInvalid(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequestRouterWithPath("GET", "/?tags=ns1/k3=val4&tags=invalidTag", nil, "", AdvisoriesExportHandler, "/")

	var errResp utils.ErrorResponse
	CheckResponse(t, w, http.StatusBadRequest, &errResp)
	assert.Equal(t, fmt.Sprintf(InvalidTagMsg, "invalidTag"), errResp.Error)
}

func TestAdvisoriesExportIncorrectFilter(t *testing.T) {
	core.SetupTest(t)
	w := CreateRequestRouterWithPath("GET", "/?filter[filteriamnotexitst]=abcd", nil, "text/csv",
		AdvisoriesExportHandler, "/")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
