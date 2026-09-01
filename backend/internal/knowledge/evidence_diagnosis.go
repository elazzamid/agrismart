package knowledge

import (
    "context"
    "strings"
)

type EvidenceDiagnosis struct {
    ProblemType string `json:"problem_type"`
    ProblemID string `json:"problem_id"`
    MatchedSymptoms int `json:"matched_symptoms"`
    Score int `json:"score"`
    Evidence []string `json:"evidence,omitempty"`
}

// DiagnoseBySymptoms ranks published problems using explicit symptom mappings.
// It is a candidate ranking mechanism, not a definitive diagnosis.
func (s *Service) DiagnoseBySymptoms(ctx context.Context, symptomIDs []string, cropID string, limit int) ([]EvidenceDiagnosis, error) {
    if limit <= 0 || limit > 20 { limit = 10 }
    cleaned := make([]string, 0, len(symptomIDs))
    for _, id := range symptomIDs { if v := strings.TrimSpace(id); v != "" { cleaned = append(cleaned, v) } }
    if len(cleaned) == 0 { return []EvidenceDiagnosis{}, nil }
    rows, err := s.db.Query(ctx, `
        SELECT CASE WHEN ps.pest_id IS NOT NULL THEN 'pest' ELSE 'disease' END,
               COALESCE(ps.pest_id, ps.disease_id)::text,
               COUNT(DISTINCT ps.symptom_id)::int,
               SUM(ps.weight)::int,
               ARRAY_AGG(DISTINCT COALESCE(s.name, ps.evidence_note))
        FROM problem_symptoms ps
        JOIN symptoms s ON s.id=ps.symptom_id
        WHERE ps.symptom_id::text = ANY($1::text[])
          AND ($2='' OR EXISTS (
              SELECT 1 FROM problem_symptom_sources pss
              JOIN knowledge_documents d ON d.id=pss.document_id
              JOIN knowledge_versions v ON v.id=pss.version_id AND v.document_id=d.id
              WHERE pss.symptom_id=ps.symptom_id
                AND ((pss.pest_id IS NOT NULL AND pss.pest_id=ps.pest_id) OR (pss.disease_id IS NOT NULL AND pss.disease_id=ps.disease_id))
                AND d.status='published'
                AND EXISTS (SELECT 1 FROM knowledge_crop_links k WHERE k.document_id=d.id AND k.crop_id::text=$2)
                AND v.version_no = (SELECT MAX(v2.version_no) FROM knowledge_versions v2 WHERE v2.document_id=d.id)
                AND EXISTS (
                    SELECT 1 FROM knowledge_validations kv
                    WHERE kv.version_id=v.id AND kv.decision='approved'
                      AND NOT EXISTS (
                          SELECT 1 FROM knowledge_validations kv2
                          WHERE kv2.version_id=v.id AND (kv2.validated_at, kv2.id) > (kv.validated_at, kv.id)
                      )
                )
          ))
          AND (ps.pest_id IS NOT NULL) <> (ps.disease_id IS NOT NULL)
        GROUP BY ps.pest_id, ps.disease_id
        ORDER BY score DESC, matched_symptoms DESC
        LIMIT $3`, cleaned, strings.TrimSpace(cropID), limit)
    if err != nil { return nil, err }
    defer rows.Close()
    result := make([]EvidenceDiagnosis, 0)
    for rows.Next() {
        var d EvidenceDiagnosis
        if err := rows.Scan(&d.ProblemType,&d.ProblemID,&d.MatchedSymptoms,&d.Score,&d.Evidence); err != nil { return nil, err }
        result = append(result,d)
    }
    if err := rows.Err(); err != nil { return nil, err }
    return result,nil
}
