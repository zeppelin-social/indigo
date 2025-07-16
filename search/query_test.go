//go:build localsearch

package search

import (
	"context"
	"crypto/tls"
	"io"
	"log/slog"
	"net/http"
	"testing"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"

	"github.com/ipfs/go-cid"
	es "github.com/opensearch-project/opensearch-go/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testPostIndex    = "palomar_test_post"
	testProfileIndex = "palomar_test_profile"
)

func testEsClient(t *testing.T) *es.Client {
	cfg := es.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "admin",
		Password:  "0penSearch-Pal0mar",
		CACert:    nil,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 5,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	escli, err := es.NewClient(cfg)
	if err != nil {
		t.Fatal(err)
	}
	info, err := escli.Info()
	if err != nil {
		t.Fatal(err)
	}
	info.Body.Close()
	return escli

}

func testServer(ctx context.Context, t *testing.T, escli *es.Client, dir identity.Directory) *Server {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	srv, err := NewServer(
		escli,
		dir,
		ServerConfig{
			PostIndex:    testPostIndex,
			ProfileIndex: testProfileIndex,
			Logger:       slog.Default(),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	idx, err := NewIndexer(db, escli, dir, IndexerConfig{
		RelayHost:           "wss://relay.invalid",
		RelaySyncRateLimit:  1,
		IndexMaxConcurrency: 1,
		PostIndex:           testPostIndex,
		ProfileIndex:        testProfileIndex,
	})
	if err != nil {
		t.Fatal(err)
	}

	srv.Indexer = idx

	// NOTE: skipping errors
	resp, _ := srv.escli.Indices.Delete([]string{testPostIndex, testProfileIndex})
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	if err := srv.EnsureIndices(ctx); err != nil {
		t.Fatal(err)
	}

	return srv
}

func TestJapaneseRegressions(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	escli := testEsClient(t)
	dir := identity.NewMockDirectory()
	srv := testServer(ctx, t, escli, &dir)
	ident := identity.Identity{
		DID:    syntax.DID("did:plc:abc111"),
		Handle: syntax.Handle("handle.example.com"),
	}

	res, err := DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "english",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(0, len(res.Hits.Hits))

	p1 := appbsky.FeedPost{Text: "basic english post", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p1, rkey: "3kpnillluoh2y", rcid: cid.Undef},
	}))

	// https://github.com/bluesky-social/indigo/issues/302
	p2 := appbsky.FeedPost{Text: "学校から帰って熱いお風呂に入ったら力一杯がんばる", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p2, rkey: "3kpnillluo222", rcid: cid.Undef},
	}))
	p3 := appbsky.FeedPost{Text: "熱力学", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p3, rkey: "3kpnillluo333", rcid: cid.Undef},
	}))
	p4 := appbsky.FeedPost{Text: "東京都", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p4, rkey: "3kpnillluo444", rcid: cid.Undef},
	}))
	p5 := appbsky.FeedPost{Text: "京都", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p5, rkey: "3kpnillluo555", rcid: cid.Undef},
	}))
	p6 := appbsky.FeedPost{Text: "パリ", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p6, rkey: "3kpnillluo666", rcid: cid.Undef},
	}))
	p7 := appbsky.FeedPost{Text: "ハリー・ポッター", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p7, rkey: "3kpnillluo777", rcid: cid.Undef},
	}))
	p8 := appbsky.FeedPost{Text: "ハリ", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p8, rkey: "3kpnillluo223", rcid: cid.Undef},
	}))
	p9 := appbsky.FeedPost{Text: "multilingual 多言語", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p9, rkey: "3kpnillluo224", rcid: cid.Undef},
	}))

	_, err = srv.escli.Indices.Refresh()
	assert.NoError(err)

	// expect all to be indexed
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "*",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(9, len(res.Hits.Hits))

	// check that english matches (single post)
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "english",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// "thermodynamics"; should return only one match
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "熱力学",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// "Kyoto"; should return only one match
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "京都",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// "Paris"; should return only one match
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "パリ",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// should return only one match
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "ハリー",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// part of a word; should match none
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "ハ",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(0, len(res.Hits.Hits))

	// should match both ways, and together
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "multilingual",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "多言語",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "multilingual 多言語",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "multilingual 多言語",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "multilingual 多言語",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "\"multilingual 多言語\"",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
}

func TestParsedQuery(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	escli := testEsClient(t)
	dir := identity.NewMockDirectory()
	srv := testServer(ctx, t, escli, &dir)
	ident := identity.Identity{
		DID:    syntax.DID("did:plc:abc111"),
		Handle: syntax.Handle("handle.example.com"),
	}
	other := identity.Identity{
		DID:    syntax.DID("did:plc:abc222"),
		Handle: syntax.Handle("other.example.com"),
	}
	dir.Insert(ident)
	dir.Insert(other)

	res, err := DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "english",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(0, len(res.Hits.Hits))

	p1 := appbsky.FeedPost{Text: "basic english post", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p1, rkey: "3kpnillluoh2y", rcid: cid.Undef},
	}))
	p2 := appbsky.FeedPost{Text: "another english post", CreatedAt: "2024-01-02T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p2, rkey: "3kpnilllu2222", rcid: cid.Undef},
	}))
	p3 := appbsky.FeedPost{
		Text:      "#cat post with hashtag",
		CreatedAt: "2024-01-02T03:04:05.006Z",
		Facets: []*appbsky.RichtextFacet{
			&appbsky.RichtextFacet{
				Features: []*appbsky.RichtextFacet_Features_Elem{
					&appbsky.RichtextFacet_Features_Elem{
						RichtextFacet_Tag: &appbsky.RichtextFacet_Tag{
							Tag: "trick",
						},
					},
				},
				Index: &appbsky.RichtextFacet_ByteSlice{
					ByteStart: 0,
					ByteEnd:   4,
				},
			},
		},
	}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p3, rkey: "3kpnilllu3333", rcid: cid.Undef},
	}))
	p4 := appbsky.FeedPost{
		Text:      "@other.example.com post with mention",
		CreatedAt: "2024-01-02T03:04:05.006Z",
		Facets: []*appbsky.RichtextFacet{
			&appbsky.RichtextFacet{
				Features: []*appbsky.RichtextFacet_Features_Elem{
					&appbsky.RichtextFacet_Features_Elem{
						RichtextFacet_Mention: &appbsky.RichtextFacet_Mention{
							Did: "did:plc:abc222",
						},
					},
				},
				Index: &appbsky.RichtextFacet_ByteSlice{
					ByteStart: 0,
					ByteEnd:   18,
				},
			},
		},
	}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p4, rkey: "3kpnilllu4444", rcid: cid.Undef},
	}))
	p5 := appbsky.FeedPost{
		Text:      "https://bsky.app... post with hashtag #cat",
		CreatedAt: "2024-01-02T03:04:05.006Z",
		Facets: []*appbsky.RichtextFacet{
			&appbsky.RichtextFacet{
				Features: []*appbsky.RichtextFacet_Features_Elem{
					&appbsky.RichtextFacet_Features_Elem{
						RichtextFacet_Link: &appbsky.RichtextFacet_Link{
							Uri: "htTPS://www.en.wikipedia.org/wiki/CBOR?q=3&a=1&utm_campaign=123",
						},
					},
				},
				Index: &appbsky.RichtextFacet_ByteSlice{
					ByteStart: 0,
					ByteEnd:   19,
				},
			},
		},
	}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p5, rkey: "3kpnilllu5555", rcid: cid.Undef},
	}))
	p6 := appbsky.FeedPost{
		Text:      "post with lang (deutsch)",
		CreatedAt: "2024-01-02T03:04:05.006Z",
		Langs:     []string{"ja", "de-DE"},
	}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p6, rkey: "3kpnilllu6666", rcid: cid.Undef},
	}))
	p7 := appbsky.FeedPost{Text: "post with old date", CreatedAt: "2020-05-03T03:04:05.006Z"}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p7, rkey: "3kpnilllu7777", rcid: cid.Undef},
	}))
	p8 := appbsky.FeedPost{
		Text:      "post with parent",
		CreatedAt: "2024-01-02T03:04:05.006Z",
		Reply: &appbsky.FeedPost_ReplyRef{
			Parent: &comatproto.RepoStrongRef{
				Uri: "at://did:plc:abc111/app.bsky.feed.post/3kpnilllu4444",
			},
			Root: &comatproto.RepoStrongRef{
				Uri: "at://did:plc:abc111/app.bsky.feed.post/3kpnilllu4444",
			},
		},
	}
	assert.NoError(srv.Indexer.indexPosts(ctx, []*PostIndexJob{
		&PostIndexJob{did: ident.DID, record: &p8, rkey: "3kpnilllu8888", rcid: cid.Undef},
	}))

	_, err = srv.escli.Indices.Refresh()
	assert.NoError(err)

	// expect all to be indexed
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "*",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(8, len(res.Hits.Hits))

	// check that english matches both
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "english",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(2, len(res.Hits.Hits))

	// phrase only matches one
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "\"basic english\"",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// posts-by
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "from:handle.example.com",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(8, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "from:@handle.example.com",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(8, len(res.Hits.Hits))

	// hashtag query
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "post #trick",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "post #Trick",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "post #trick #allMustMatch",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(0, len(res.Hits.Hits))

	// mention query
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "@other.example.com",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// to parent query
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "to:handle.example.com",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// URL and domain queries
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "https://en.wikipedia.org/wiki/CBOR?a=1&q=3",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "\"https://en.wikipedia.org/wiki/CBOR?a=1&q=3\"",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(0, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "https://en.wikipedia.org/wiki/CBOR",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(0, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "domain:en.wikipedia.org",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// lang filter
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "lang:de",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))

	// date range filters
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "since:2023-01-01T00:00:00Z",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(7, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "since:2023-01-01",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(7, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "until:2023-01-01",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(1, len(res.Hits.Hits))
	res, err = DoSearchPosts(ctx, &dir, escli, testPostIndex, &PostSearchParams{
		Query:  "until:asdf",
		Offset: 0,
		Size:   20,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(8, len(res.Hits.Hits))
}
