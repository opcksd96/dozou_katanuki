package wails

import (
	"context"
	app "dozou_katanuki/application/relation"
)

type RelationApp struct {
	ctx        context.Context
	getGraphUC *app.GetGraphUseCase
}

func NewRelationApp(getGraphUC *app.GetGraphUseCase) *RelationApp {
	return &RelationApp{getGraphUC: getGraphUC}
}

func (a *RelationApp) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *RelationApp) GetKnowledgeGraph(centerAccountID string) (*app.GraphData, error) {
	return a.getGraphUC.Execute(a.ctx, centerAccountID)
}
