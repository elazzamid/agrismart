# Agricultural Seed Data Policy

Seed data is split into two categories.

## Structural reference data

Safe for repository seed files when the values are stable and primarily identify entities, such as crop names, generic growth-stage labels, and catalog examples.

## Agronomic recommendations

Do not encode fertilizer rates, pesticide doses, spray intervals, PHI, tank mixes, or other actionable recommendations in generic seed files without a traceable authoritative source and validation metadata.

Recommendations will belong to the knowledge-base model and must support:

- source reference
- author
- validator
- validation date
- version
- publication status
- applicable crop/variety/growth stage

The AI layer must consume published/validated knowledge rather than treating seed data as expert advice.
