package knowledge

import (
    "context"
    "strings"
)

type DiagnosisCandidate struct {
    DocumentID string `json:"document_id"`
    Slug string `json:"slug"`
    Title string `json:"title"`
    ProblemType string `json:"problem_type"`
    ProblemID string `json:"problem_id"`
    Score int `json:"score"`
    Evidence string `json:"evidence,omitempty"`
}

// DiagnosePublished performs conservative keyword matching against published
// knowledge. It is intentionally a candidate finder, not a definitive diagnosis.
func (s *Service) DiagnosePublished(ctx context.Context, cropID string, symptoms []string, limit int) ([]DiagnosisCandidate, error) {
    if limit <= 0 || limit > 20 { limit = 10 }
    cleaned := make([]string, 0, len(symptoms))
    for _, symptom := range symptoms {
        if v := strings.TrimSpace(symptom); v != "" { cleaned = append(cleaned, v) }
    }
    if len(cleaned) == 0 { return []DiagnosisCandidate{}, nil }

    rows, err := s.db.Query(ctx, `
        SELECT DISTINCT d.id, d.slug, d.title,
            CASE WHEN kp.pest_id IS NOT NULL THEN 'pest' ELSE 'disease' END,
            COALESCE(kp.pest_id, kd.disease_id)::text,
            GREATEST(
                (SELECT COUNT(*) FROM unnest($2::text[]) q WHERE v.content ILIKE '%' || q || '%'),
                (SELECT COUNT(*) FROM unnest($2::text[]) q WHERE d.title ILIKE '%' || q || '%')
            )::int AS score,
            v.content
        FROM knowledge_documents d
        JOIN knowledge_versions v ON v.document_id=d.id
        LEFT JOIN knowledge_pest_links kp ON kp.document_id=d.id
        LEFT JOIN knowledge_disease_links kd ON kd.document_id=d.id
        WHERE d.status='published'
          AND ($1='' OR EXISTS (SELECT 1 FROM knowledge_crop_links k WHERE k.document_id=d.id AND k.crop_id::text=$1))
          AND EXISTS (SELECT 1 FROM knowledge_validations kv WHERE kv.version_id=v.id AND kv.decision='approved')
          AND EXISTS (SELECT 1 FROM unnest($2::text[]) q WHERE v.content ILIKE '%' || q || '%' OR d.title ILIKE '%' || q || '%')
        ORDER BY score DESC, d.updated_at DESC
        LIMIT $3`, strings.TrimSpace(cropID), cleaned, limit)
    if err != nil { return nil, err }
    defer rows.Close()
    result := make([]DiagnosisCandidate, 0)
    for rows.Next() {
        var c DiagnosisCandidate
        if err := rows.Scan(&c.DocumentID,&c.Slug,&c.Title,&c.ProblemType,&c.ProblemID,&c.Score,&c.Evidence); err != nil { return nil, err }
        result = append(result,c)
    }
    if err := rows.Err(); err != nil { return nil, err }
    return result,nil
}
