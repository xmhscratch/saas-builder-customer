package models

import (
	"github.com/vanng822/go-solr/solr"
)

// SearchResultParser comment
type SearchResultParser struct {
	RawResponse map[string]interface{}
}

// Parse comment
func (ctx *SearchResultParser) Parse(resp *solr.SolrResponse) (*solr.SolrResult, error) {
	sr := &solr.SolrResult{}
	sr.Results = new(solr.Collection)
	sr.Status = resp.Status

	if resp.Status == 0 {
		ctx.ParseResponse(resp, sr)
		ctx.ParseFacetCounts(resp, sr)
		ctx.ParseHighlighting(resp, sr)
	} else {
		ctx.ParseError(resp, sr)
	}
	ctx.RawResponse = resp.Response

	return sr, nil
}

// ParseError comment
func (ctx *SearchResultParser) ParseError(resp *solr.SolrResponse, sr *solr.SolrResult) {
	if error, ok := resp.Response["error"]; ok {
		sr.Error = error.(map[string]interface{})
	}
}

// ParseResponse comment
func (ctx *SearchResultParser) ParseResponse(resp *solr.SolrResponse, sr *solr.SolrResult) {
	if resp, ok := resp.Response["resp"].(map[string]interface{}); ok {
		sr.Results.NumFound = int(resp["numFound"].(float64))
		sr.Results.Start = int(resp["start"].(float64))
		if docs, ok := resp["docs"].([]interface{}); ok {
			sr.Results.Docs = make([]solr.Document, len(docs))
			// remove version
			for i, v := range docs {
				d := solr.Document{}
				for k, v := range v.(map[string]interface{}) {
					if k != "_version_" {
						d.Set(k, v)
					}
				}
				sr.Results.Docs[i] = d
			}
		}
	} else {
		panic(`Standard ctx can only parse solr resp with resp object,
					ie resp.resp and resp.resp.docs.
					Please use other ctx or implement your own ctx`)
	}
}

// ParseFacetCounts comment
func (ctx *SearchResultParser) ParseFacetCounts(resp *solr.SolrResponse, sr *solr.SolrResult) {
	if facetCounts, ok := resp.Response["facet_counts"]; ok {
		sr.FacetCounts = facetCounts.(map[string]interface{})
	}
}

// ParseHighlighting comment
func (ctx *SearchResultParser) ParseHighlighting(resp *solr.SolrResponse, sr *solr.SolrResult) {
	if highlighting, ok := resp.Response["highlighting"]; ok {
		sr.Highlighting = highlighting.(map[string]interface{})
	}
}
