package rule

import (
	"fmt"
	"time"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// ResolutionPipeline — Orchestrator điều phối toàn bộ pipeline 12 bước
type ResolutionPipeline struct {
	config             *types.PipelineConfig
	legalNormalizer    *LegalStatusResolver
	semanticNorm       *LandUseCodeResolver
	conflictClassifier *ConflictClassifier
}

// NewResolutionPipeline tạo ResolutionPipeline
func NewResolutionPipeline(
	config *types.PipelineConfig,
	legalNormalizer *LegalStatusResolver,
	semanticNorm *LandUseCodeResolver,
	conflictClassifier *ConflictClassifier,
) *ResolutionPipeline {
	if config == nil {
		config = types.DefaultPipelineConfig()
	}
	return &ResolutionPipeline{
		config:             config,
		legalNormalizer:    legalNormalizer,
		semanticNorm:       semanticNorm,
		conflictClassifier: conflictClassifier,
	}
}

// Name trả về tên engine
func (p *ResolutionPipeline) Name() string { return "ResolutionPipeline" }

// Execute chạy toàn bộ pipeline, trả về kết quả + trace
func (p *ResolutionPipeline) Execute(
	candidates []*types.LayerCandidate,
	mode types.ResolveMode,
) (*types.ResolveResult, *types.ResolutionTrace) {
	trace := &types.ResolutionTrace{StartedAt: time.Now()}
	result := &types.ResolveResult{Mode: mode}

	if len(candidates) == 0 {
		p.addTraceStep(trace, "no_candidates", 0, 0, "skipped", nil)
		trace.CompletedAt = time.Now()
		return result, trace
	}

	// Step 1: Semantic Normalization
	stepStart := time.Now()
	for _, c := range candidates {
		p.semanticNorm.NormalizeCandidate(c)
	}
	p.addTraceStep(trace, "semantic_normalize",
		len(candidates), len(candidates), "ok",
		map[string]any{"engine": p.semanticNorm.Name()})

	// Step 2: Legal Normalization
	stepStart = time.Now()
	normalizedLegals := make(map[uint64]*types.NormalizedLegal, len(candidates))
	for _, c := range candidates {
		normalizedLegals[c.LayerID] = p.legalNormalizer.Normalize(c)
	}
	p.addTraceStep(trace, "legal_normalize",
		len(candidates), len(candidates), "ok",
		map[string]any{"engine": p.legalNormalizer.Name()})

	// Step 3-4: Score & Rank + Tie-break (delegated to ParcelEngine)
	// Step 5: Conflict Detection & Classification (delegated to ParcelEngine)
	// Step 6: Risk Assessment (delegated to ParcelEngine)
	// Step 7: Build Response (delegated to ParcelEngine)
	// Step 8: Fallback check
	// Step 9: Set ResolutionStatus

	_ = stepStart // used above

	trace.CompletedAt = time.Now()
	return result, trace
}

// addTraceStep thêm 1 bước vào trace
func (p *ResolutionPipeline) addTraceStep(
	trace *types.ResolutionTrace,
	stepName string, inputCount, outputCount int,
	status string, details map[string]any,
) {
	if !p.config.EnableTrace {
		return
	}
	trace.Steps = append(trace.Steps, types.TraceStep{
		StepName:    stepName,
		InputCount:  inputCount,
		OutputCount: outputCount,
		Status:      status,
		Details:     details,
	})
}

// ResolveStatus xác định ResolutionStatus từ kết quả
func (p *ResolutionPipeline) ResolveStatus(
	hasConflict bool, hasWarning bool, isFallback bool,
) types.ResolutionStatus {
	if isFallback {
		return types.StatusFallbackResolved
	}
	if hasConflict {
		return types.StatusResolvedWithConflict
	}
	if hasWarning {
		return types.StatusResolvedWithWarning
	}
	return types.StatusResolved
}

// Validate kiểm tra config
func (p *ResolutionPipeline) Validate() error {
	if p.legalNormalizer == nil {
		return fmt.Errorf("ResolutionPipeline: legalNormalizer is nil")
	}
	if p.conflictClassifier == nil {
		return fmt.Errorf("ResolutionPipeline: conflictClassifier is nil")
	}
	return nil
}
