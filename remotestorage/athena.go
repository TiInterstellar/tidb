package remotestorage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/glue"
)

type Athena struct {
	glue   *glue.Glue
	client *athena.Athena
}

// TODO: just a quick reference for now, should be implemented.
func (a *Athena) Setup() {
	// 1. Create a glue database
	_, _ = a.glue.CreateDatabase(nil)

	// 2. Create a glue table with pre-defined schema
	_, _ = a.glue.CreateTable(nil)

	// 3. Create a glue crawler to crawl the table data from s3
	_, _ = a.glue.CreateCrawler(nil)

	// 4. Create a glue job to run the crawler
	_, _ = a.glue.CreateJob(nil)

	// 5. Create a athena catalog with glue
	_, _ = a.client.CreateDataCatalog(&athena.CreateDataCatalogInput{
		Name: aws.String("implement-me"),
		Type: aws.String("glue"),
		Parameters: map[string]*string{
			"catalog_id": aws.String("implement-me"),
		},
	})
}
