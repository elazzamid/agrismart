package knowledge

import (
    "context"
    "strings"
)

type NutrientRequirement struct {
    CropID string `json:"crop_id"`
    GrowthStageID *string `json:"growth_stage_id,omitempty"`
    NutrientCode string `json:"nutrient_code"`
    NutrientName string `json:"nutrient_name"`
    RequirementMin *float64 `json:"requirement_min,omitempty"`
    RequirementMax *float64 `json:"requirement_max,omitempty"`
    Unit string `json:"unit"`
    SourceDocumentID *string `json:"source_document_id,omitempty"`
    SourceVersionID *string `json:"source_version_id,omitempty"`
    Notes string `json:"notes,omitempty"`
}

// ListNutrientRequirements returns only requirements tied to a published
// knowledge document whose referenced version has an approved validation.
func (s *Service) ListNutrientRequirements(ctx context.Context, cropID, growthStageID string, limit int) ([]NutrientRequirement, error) {
    if limit <= 0 || limit > 100 { limit = 50 }
    cropID, growthStageID = strings.TrimSpace(cropID), strings.TrimSpace(growthStageID)
    rows, err := s.db.Query(ctx, `
SELECT r.crop_id::text,r.growth_stage_id::text,n.code,n.name,r.requirement_min,r.requirement_max,r.unit,
       r.source_document_id::text,r.source_version_id::text,COALESCE(r.notes,'')
FROM crop_nutrient_requirements r
JOIN nutrients n ON n.id=r.nutrient_id
LEFT JOIN knowledge_documents d ON d.id=r.source_document_id AND d.status='published'
LEFT JOIN knowledge_versions v ON v.id=r.source_version_id AND v.document_id=d.id
WHERE r.crop_id::text=$1
  AND ($2='' OR r.growth_stage_id::text=$2)
  AND r.source_document_id IS NOT NULL
  AND r.source_version_id IS NOT NULL
  AND EXISTS (SELECT 1 FROM knowledge_validations kv WHERE kv.version_id=r.source_version_id AND kv.decision='approved')
ORDER BY n.code
LIMIT $3`, cropID, growthStageID, limit)
    if err != nil { return nil, err }
    defer rows.Close()
    result := make([]NutrientRequirement,0)
    for rows.Next() {
        var x NutrientRequirement
        if err := rows.Scan(&x.CropID,&x.GrowthStageID,&x.NutrientCode,&x.NutrientName,&x.RequirementMin,&x.RequirementMax,&x.Unit,&x.SourceDocumentID,&x.SourceVersionID,&x.Notes); err != nil { return nil, err }
        result=append(result,x)
    }
    if err := rows.Err(); err != nil { return nil, err }
    return result,nil
}
