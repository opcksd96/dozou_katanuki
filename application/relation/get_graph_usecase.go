package relation

import (
	"context"
	"dozou_katanuki/domain/relation"
)

type GraphNode struct {
	Data map[string]interface{} `json:"data"`
}

type GraphEdge struct {
	Data map[string]interface{} `json:"data"`
}

type GraphData struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

type GetGraphUseCase struct {
	repo relation.RelationRepository
}

func NewGetGraphUseCase(repo relation.RelationRepository) *GetGraphUseCase {
	return &GetGraphUseCase{repo: repo}
}

func (uc *GetGraphUseCase) Execute(ctx context.Context, centerAccountID string) (*GraphData, error) {
	// 簡略化のため、すべてのRelationを取得してグラフ化するロジックとする
	// 実際には centerAccountID を起点に探索する
	relations, err := uc.repo.GetAllRelations(ctx)
	if err != nil {
		return nil, err
	}

	nodesMap := make(map[string]bool)
	var nodes []GraphNode
	var edges []GraphEdge

	for _, rel := range relations {
		if !nodesMap[rel.SourceAccountID] {
			nodes = append(nodes, GraphNode{Data: map[string]interface{}{"id": rel.SourceAccountID}})
			nodesMap[rel.SourceAccountID] = true
		}
		if !nodesMap[rel.TargetAccountID] {
			nodes = append(nodes, GraphNode{Data: map[string]interface{}{"id": rel.TargetAccountID, "name": rel.TargetHandle}})
			nodesMap[rel.TargetAccountID] = true
		}

		edges = append(edges, GraphEdge{
			Data: map[string]interface{}{
				"id":     rel.ID,
				"source": rel.SourceAccountID,
				"target": rel.TargetAccountID,
				"type":   rel.Type,
				"weight": rel.Weight,
			},
		})
	}

	if nodes == nil {
		nodes = []GraphNode{}
	}
	if edges == nil {
		edges = []GraphEdge{}
	}

	return &GraphData{Nodes: nodes, Edges: edges}, nil
}
