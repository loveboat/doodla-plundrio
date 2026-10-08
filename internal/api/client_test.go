package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/elsbrock/go-putio"
)

func TestGetAllTransferFilesKeepsSubdirectoryPaths(t *testing.T) {
	folder := func(id int, name string) string {
		return fmt.Sprintf(`{"id":%d,"name":%q,"content_type":"application/x-directory","file_type":"FOLDER"}`, id, name)
	}
	file := func(id int, name string) string {
		return fmt.Sprintf(`{"id":%d,"name":%q,"size":10,"content_type":"video/mp4","file_type":"VIDEO"}`, id, name)
	}
	listings := map[string]string{
		"100": "[" + file(1, "Show.S01E01.mp4") + "," + folder(200, "Subs") + "]",
		"200": "[" + folder(300, "Show.S01E01") + "," + folder(301, "Show.S01E02") + "]",
		"300": "[" + file(3, "2_eng.srt") + "]",
		"301": "[" + file(4, "2_eng.srt") + "]",
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v2/files/100" {
			fmt.Fprintf(w, `{"status":"OK","file":%s}`, folder(100, "Show.S01"))
			return
		}
		files, ok := listings[r.URL.Query().Get("parent_id")]
		if r.URL.Path != "/v2/files/list" || !ok {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"status":"ERROR","error_type":"NotFound","error_message":"missing"}`)
			return
		}
		fmt.Fprintf(w, `{"status":"OK","files":%s,"parent":%s,"cursor":""}`, files, folder(0, "parent"))
	}))
	defer srv.Close()
	client := putio.NewClient(srv.Client())
	client.BaseURL, _ = url.Parse(srv.URL)

	files, err := (&Client{client: client}).GetAllTransferFiles(context.Background(), 100)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, len(files))
	for i, f := range files {
		names[i] = f.Name
	}

	want := []string{"Show.S01E01.mp4", "Subs/Show.S01E01/2_eng.srt", "Subs/Show.S01E02/2_eng.srt"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("names = %q, want %q", names, want)
	}
}

func TestTransferRootNotFoundIsDistinctFromMissingChild(t *testing.T) {
	for _, missingRoot := range []bool{true, false} {
		t.Run(fmt.Sprint(missingRoot), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if !missingRoot && r.URL.Path == "/v2/files/500" {
					fmt.Fprint(w, `{"status":"OK","file":{"id":500,"name":"Book","content_type":"application/x-directory","file_type":"FOLDER"}}`)
					return
				}
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"status":"ERROR","error_type":"NotFound","error_message":"missing"}`)
			}))
			defer srv.Close()
			client := putio.NewClient(srv.Client())
			client.BaseURL, _ = url.Parse(srv.URL)
			_, err := (&Client{client: client}).GetAllTransferFiles(context.Background(), 500)
			var rootErr *TransferSourceNotFoundError
			if err == nil || errors.As(err, &rootErr) != missingRoot {
				t.Fatalf("root missing=%v, got %v", missingRoot, err)
			}
			var response *putio.ErrorResponse
			if !errors.As(err, &response) || response.Type != "NotFound" {
				t.Fatalf("must retain underlying API error: %v", err)
			}
		})
	}
}
