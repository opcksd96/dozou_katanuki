package relation

import (
	"context"
	"fmt"
	"dozou_katanuki/domain/relation"
)

type BuildGraphUseCase struct {
	repo relation.RelationRepository
}

func NewBuildGraphUseCase(repo relation.RelationRepository) *BuildGraphUseCase {
	return &BuildGraphUseCase{repo: repo}
}

func (uc *BuildGraphUseCase) Execute(ctx context.Context, articleID string) error {
	// TODO: Phase 2
	// ツイート解析 (本文, メンション, 引用, メディアハッシュ) を行い、エッジを自動生成する。
	fmt.Println("[BuildGraphUseCase] 自動抽出ロジックは Phase 2 で実装されます。")
	return nil
}

func (uc *BuildGraphUseCase) AddManualRelation(ctx context.Context, rel *relation.AccountRelation) error {
	return uc.repo.SaveRelation(ctx, rel)
}
