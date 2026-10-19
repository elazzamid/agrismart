package knowledge

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PublishedKnowledge struct {
	ID string `json:"id"`
	Slug string `json:"slug"`
	Title string `json:"title"`
	Summary string `json:"summary,omitempty"`
	VersionID string `json:"version_id"`
	VersionNo int `json:"version_no"`
	Content string `json:"content"`
	SourceID *string `json:"source_id,omitempty"`
	CropID string `json:"crop_id,omitempty"`
	GrowthStageID *string `json:"growth_stage_id,omitempty"`
}

type RetrievalFilter struct { CropID string; GrowthStageID string; Query string }

const publishedRetrievalQuery = `
SELECT d.id,d.slug,d.title,COALESCE(d.summary,''),v.id,v.version_no,v.content,v.source_id,
       COALESCE(k.crop_id::text,''),k.growth_stage_id
FROM knowledge_documents d
JOIN LATERAL (
  SELECT v1.* FROM knowledge_versions v1
  WHERE v1.document_id=d.id
    AND v1.version_no = (SELECT MAX(v2.version_no) FROM knowledge_versions v2 WHERE v2.document_id=d.id)
    AND NULLIF(BTRIM(v1.source_id), '') IS NOT NULL
    AND EXISTS (
      SELECT 1 FROM knowledge_validations kv
      WHERE kv.version_id=v1.id AND kv.decision='approved'
        AND NOT EXISTS (
          SELECT 1 FROM knowledge_validations kv2
          WHERE kv2.version_id=v1.id
            AND (kv2.validated_at, kv2.id) > (kv.validated_at, kv.id)
        )
    )
  ORDER BY v1.version_no DESC LIMIT 1
) v ON TRUE
LEFT JOIN LATERAL (
  SELECT k1.crop_id,k1.growth_stage_id FROM knowledge_crop_links k1
  WHERE k1.document_id=d.id
    AND ($1='' OR k1.crop_id::text=$1)
    AND ($2='' OR k1.growth_stage_id::text=$2)
  ORDER BY k1.crop_id,k1.growth_stage_id NULLS LAST LIMIT 1
) k ON TRUE
WHERE d.status='published'
  AND ($1='' OR k.crop_id IS NOT NULL)
  AND ($2='' OR k.growth_stage_id IS NOT NULL)
  AND ($3='' OR d.title ILIKE '%'||$3||'%' OR COALESCE(d.summary,'') ILIKE '%'||$3||'%' OR v.content ILIKE '%'||$3||'%')
ORDER BY d.updated_at DESC,v.version_no DESC
LIMIT $4`

func (s *Service) SearchPublished(ctx context.Context, filter RetrievalFilter, limit int) ([]PublishedKnowledge, error) {
	if limit <= 0 || limit > 50 { limit = 20 }
	cropID := strings.TrimSpace(filter.CropID)
	stageID := strings.TrimSpace(filter.GrowthStageID)
	q := strings.TrimSpace(filter.Query)
	rows, err := s.db.Query(ctx, publishedRetrievalQuery, cropID, stageID, q, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	result := make([]PublishedKnowledge, 0)
	for rows.Next() {
		var item PublishedKnowledge
		if err := rows.Scan(&item.ID,&item.Slug,&item.Title,&item.Summary,&item.VersionID,&item.VersionNo,&item.Content,&item.SourceID,&item.CropID,&item.GrowthStageID); err != nil { return nil, err }
		result = append(result,item)
	}
	if err := rows.Err(); err != nil { return nil, err }
	return result,nil
}

var _ *pgxpool.Pool
