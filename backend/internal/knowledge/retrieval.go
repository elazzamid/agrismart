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

type RetrievalFilter struct {
    CropID string
    GrowthStageID string
    Query string
}

func (s *Service) SearchPublished(ctx context.Context, filter RetrievalFilter, limit int) ([]PublishedKnowledge, error) {
    if limit <= 0 || limit > 50 { limit = 20 }
    query := `
        SELECT d.id, d.slug, d.title, COALESCE(d.summary,''), v.id, v.version_no, v.content,
               v.source_id, kcl.crop_id, kcl.growth_stage_id
        FROM knowledge_documents d
        JOIN LATERAL (
            SELECT v1.* FROM knowledge_versions v1
            WHERE v1.document_id = d.id
              AND EXISTS (
                  SELECT 1 FROM knowledge_validations kv
                  WHERE kv.version_id = v1.id AND kv.decision = 'approved'
              )
            ORDER BY v1.version_no DESC
            LIMIT 1
        ) v ON TRUE
        LEFT JOIN LATERAL (
            SELECT DISTINCT k.document_id, k.crop_id, k.growth_stage_id
            FROM knowledge_crop_links k
            WHERE k.document_id = d.id
              AND ($1 = '' OR k.crop_id::text = $1)
              AND ($2 = '' OR k.growth_stage_id::text = $2)
            LIMIT 1
        ) kcl ON TRUE
        WHERE d.status = 'published'
          AND ($1 = '' OR kcl.crop_id IS NOT NULL)
          AND ($2 = '' OR kcl.growth_stage_id IS NOT NULL)
          AND ($3 = '' OR d.title ILIKE '%'||$3||'%' OR COALESCE(d.summary,'') ILIKE '%'||$3||'%' OR v.content ILIKE '%'||$3||'%')
        ORDER BY d.updated_at DESC, v.version_no DESC
        LIMIT $4`
    rows, err := s.db.Query(ctx, query, strings.TrimSpace(filter.CropID), strings.TrimSpace(filter.GrowthStageID), strings.TrimSpace(filter.Query), limit)
    if err != nil { return nil, err }
    defer rows.Close()
    result := make([]PublishedKnowledge, 0)
    for rows.Next() {
        var item PublishedKnowledge
        if err := rows.Scan(&item.ID,&item.Slug,&item.Title,&item.Summary,&item.VersionID,&item.VersionNo,&item.Content,&item.SourceID,&item.CropID,&item.GrowthStageID); err != nil { return nil, err }
        result = append(result, item)
    }
    if err := rows.Err(); err != nil { return nil, err }
    return result, nil
}

var _ *pgxpool.Pool
