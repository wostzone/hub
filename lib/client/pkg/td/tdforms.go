// Package td with TD form creation
package td

import "github.com/wostzone/hub/lib/client/pkg/vocab"

// CreateTDForm creates a form object description how to connect to a Thing via the Hub
//
// NOTE: In WoST actions are always routed via the Hub using the Hub's protocol binding.
// Under normal circumstances forms only apply to the top level describing the Hub's protocols
// Returns a form object with operations
func CreateTDForm(op string, href string, contentType string, httpMethodName string) map[string]interface{} {
	form := make(map[string]interface{})
	form[vocab.WoTOperation] = op
	form[vocab.WoTHref] = href
	form["contentType"] = contentType
	if httpMethodName != "" {
		form["htv:methodName"] = httpMethodName
	}
	return form
}
