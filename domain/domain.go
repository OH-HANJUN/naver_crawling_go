package domain

import(
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

// https://github.com/googleapis/google-cloud-go/blob/main/bigquery/nulls.go null 레퍼런스
type NaverData struct {
    ThumbnailURL           string              `bigquery:"thumbnail_url"`
    ProductID              string              `bigquery:"product_id"`
    NaverCatalogID1        bigquery.NullString `bigquery:"naver_catalog_id1"`
    // NaverCatalogId1Platform        string  `bigquery:"naver_catalog_id1_platform"`
    // NaverCatalogId1Price        int  `bigquery:"naver_catalog_id1_price"`
    NaverCatalogID2        bigquery.NullString `bigquery:"naver_catalog_id2"`
    // NaverCatalogId2Platform        string  `bigquery:"naver_catalog_id2_platform"`
    // NaverCatalogId2Price        int  `bigquery:"naver_catalog_id2_price"`
    NaverCatalogID3        bigquery.NullString `bigquery:"naver_catalog_id3"`
    // NaverCatalogId3Platform        string  `bigquery:"naver_catalog_id3_platform"`
    // NaverCatalogId3Price        int  `bigquery:"naver_catalog_id3_price"`
    NaverShoppingURL       string              `bigquery:"naver_shopping_url"`
    NaverMinPriceParttimer int                 `bigquery:"naver_min_price_parttimer"`
    ParttimerSearchedDt    civil.DateTime      `bigquery:"parttimer_searched_dt"`
    Proxy                  string

}
