package main

import (
	"net/http"
	"testing"
)

func BenchmarkAPI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://localhost:60001/authenticate?system=cowshed&uid=uid_admin&permission=read")
		if err != nil {
			b.Error(err)
			return
		}

		if resp.StatusCode != 200 {
			b.Error("authenticate failed, %d", resp.StatusCode)
		}
	}
}
