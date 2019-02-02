package rolles

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/olivere/elastic"
)

func FindLeastFullNode(es *elastic.Client, ctx context.Context) (string, error) {

	res, err := es.NodesStats().CompletionFields("fs").Do(ctx)
	if err != nil {
		return "", fmt.Errorf("Failed to get node stats -- %v", err)
	}

	maxBytes := int64(0)
	maxNodes := make([]string, 0)
	for node, stats := range res.Nodes {
		if stats.FS != nil {
			freeBytes := stats.FS.Total.AvailableInBytes
			if freeBytes > maxBytes {
				maxBytes = freeBytes
				maxNodes = []string{node}
			} else if freeBytes == maxBytes {
				maxNodes = append(maxNodes, node)
			}
		}
	}

	if len(maxNodes) == 0 {
		return "", fmt.Errorf("all nodes returned 0 byes available")
	} else if len(maxNodes) == 1 {
		return maxNodes[0], nil
	}
	return maxNodes[rand.Intn(len(maxNodes))], nil
}
