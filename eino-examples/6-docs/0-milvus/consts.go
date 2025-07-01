package main

import "github.com/milvus-io/milvus/client/v2/entity"

const (
	typ                           = "Milvus"
	defaultCollection             = "eino_collection"
	defaultDescription            = "the collection for eino"
	defaultCollectionID           = "id"
	defaultCollectionIDDesc       = "the unique id of the document"
	defaultCollectionVector       = "vector"
	defaultCollectionVectorDesc   = "the vector of the document"
	defaultCollectionContent      = "content"
	defaultCollectionContentDesc  = "the content of the document"
	defaultCollectionMetadata     = "metadata"
	defaultCollectionMetadataDesc = "the metadata of the document"

	defaultDim = 81920

	defaultIndexField = "vector"

	defaultConsistencyLevel = ConsistencyLevelBounded
	defaultMetricType       = HAMMING
)

const (
	ConsistencyLevelStrong     ConsistencyLevel = 1
	ConsistencyLevelSession    ConsistencyLevel = 2
	ConsistencyLevelBounded    ConsistencyLevel = 3
	ConsistencyLevelEventually ConsistencyLevel = 4
	ConsistencyLevelCustomized ConsistencyLevel = 5

	L2             = MetricType(entity.L2)
	IP             = MetricType(entity.IP)
	CONSINE        = MetricType(entity.COSINE)
	HAMMING        = MetricType(entity.HAMMING)
	JACCARD        = MetricType(entity.JACCARD)
	TANIMOTO       = MetricType(entity.TANIMOTO)
	SUBSTRUCTURE   = MetricType(entity.SUBSTRUCTURE)
	SUPERSTRUCTURE = MetricType(entity.SUPERSTRUCTURE)
)

type ConsistencyLevel entity.ConsistencyLevel

func (c *ConsistencyLevel) getConsistencyLevel() entity.ConsistencyLevel {
	return entity.ConsistencyLevel(*c - 1)
}

// MetricType is the metric type for vector by eino
type MetricType entity.MetricType

// getMetricType returns the metric type
func (t *MetricType) getMetricType() entity.MetricType {
	return entity.MetricType(*t)
}
