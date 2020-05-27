package elasticsearch

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/willqiang/bookstore_utils-go/logger"
	"time"
)

var (
	Client esClientInterface = &esClient{}
)

type esClientInterface interface {
	setClient(*elastic.Client)
	Index(string, string, interface{}) (*elastic.IndexResponse, error)
	Get(string, string, string) (*elastic.GetResult, error)
	Search(string, elastic.Query) (*elastic.SearchResult, error)
}

type esClient struct {
	client *elastic.Client
}

func Init() {
	// Obtain a client for an Elasticsearch cluster of two nodes,
	// running on 10.0.1.1 and 10.0.1.2. Do not run the sniffer.
	// Set the healthcheck interval to 10s. When requests fail,
	// retry 5 times. Print error messages to os.Stderr and informational
	// messages to os.Stdout.
	log := logger.GetLogger()
	client, err := elastic.NewClient(
		elastic.SetURL("http://192.168.56.4:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log),
		elastic.SetInfoLog(log),
	)
	if err != nil {
		panic(err)
	}

	Client.setClient(client)
}

func (c *esClient) setClient(client *elastic.Client) {
	c.client = client
}

func (c *esClient) Index(index string, docType string, doc interface{}) (*elastic.IndexResponse, error) {
	fmt.Println(c.client)
	ctx := context.Background()
	result, err := c.client.Index().
		Index(index).
		Type(docType).
		BodyJson(doc).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to index document in index %s", index), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) Get(index string, docType string, id string) (*elastic.GetResult, error) {
	ctx := context.Background()
	result, err := c.client.Get().
		Index(index).
		Type(docType).
		Id(id).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when tring to get id %s", id), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) Search(index string, query elastic.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	//result, err := c.client.Search(index).Query(query).RestTotalHitsAsInt(true).Do(ctx)
	result, err := c.client.Search(index).Query(query).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return result, nil
}
