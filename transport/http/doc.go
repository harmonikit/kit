// Package http provides HTTP transport bindings for harmoni endpoints.
//
// It implements transport.Server and provides client-side endpoint wrappers
// for calling remote HTTP services.
//
// Example:
//
//	server := http.NewServer(myEndpoint, decodeRequest, encodeResponse)
//	log.Fatal(server.ListenAndServe(":8080"))
package http
