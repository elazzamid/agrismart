# AgriSmart — Master Plan

Version: 1.0
Status: MVP foundation

## 1. Vision
AgriSmart is a mobile-first digital assistant for Indonesian farmers, helping them make better cultivation decisions from land preparation through harvest.

## 2. Product principle
The product must answer: **“Apa yang harus saya lakukan terhadap tanaman saya sekarang?”**

The platform is not only an article portal. It connects verified agricultural knowledge to the farmer’s current crop, growth stage, activities, observations, and records.

## 3. MVP crops
- Cabai
- Padi
- Jagung

The initial knowledge base should be deep and validated for these crops rather than broad and shallow across many crops.

## 4. MVP modules
1. Authentication and farmer profile
2. Farm and plot management
3. Crop, variety, and growth-stage knowledge
4. Fertilizer knowledge and calculation foundation
5. Pest, disease, and weed knowledge
6. Pesticide information with safety/label guardrails
7. Farm activity records
8. Farm expenses and harvest records
9. Knowledge/article management
10. Admin panel
11. Knowledge-grounded AI assistant foundation

## 5. Post-MVP
- Image-based diagnosis
- Weather integration
- Smart notifications
- Commodity prices
- Expert/adviser workflows
- Farmer groups
- Marketplace
- Supplier/distributor ecosystem
- Advanced production analytics

## 6. Architecture
Start with a modular monolith.

Suggested stack:
- Backend: Go
- Database: PostgreSQL
- Frontend: mobile-first web application
- Object storage: for plant images and documents
- AI: retrieval/knowledge-grounded assistant behind a safety layer

Avoid microservices until there is a proven need.

## 7. Knowledge base
Knowledge is structured around:

Crop → Variety → Growth Stage → Nutrient/Fertilization → Pest/Disease/Weed → Management → Harvest

Important knowledge records should have:
- source
- author
- validator
- validation date
- version
- status

Statuses:
DRAFT → REVIEW → VALIDATED → PUBLISHED

## 8. AI principles
AI is an assistance layer, not the source of truth.

Flow:
Farmer → Intent → Knowledge Retrieval → Verified Knowledge → AI Reasoning → Safety Validation → Answer

The AI must be able to say that available data is insufficient rather than inventing an answer.

## 9. Pesticide safety
Pesticide-related information must prioritize integrated pest management, verified label information, and safe-use guidance.

The system must not invent or infer:
- dose
- interval
- PHI
- tank-mix compatibility
- off-label usage

## 10. Core domain entities
users
farmer_profiles
farms
farm_plots
crops
crop_varieties
crop_growth_stages
fertilizers
fertilizer_nutrients
pesticides
active_ingredients
pesticide_targets
pests
diseases
weeds
crop_pests
crop_diseases
crop_weeds
cultivation_guides
fertilization_guides
management_guides
farm_activities
farm_inputs
farm_expenses
harvests
articles
knowledge_documents
ai_conversations
ai_messages
ai_diagnoses
ai_recommendations

## 11. Development rules
1. Requirements before implementation.
2. Database design before major business logic.
3. No large feature expansion inside a milestone.
4. Prefer simple, testable modules.
5. Every important knowledge claim has provenance.
6. AI cannot bypass knowledge and safety layers.
7. Every milestone must be independently verifiable.
8. Important changes are committed to Git.
9. Avoid large refactors without a concrete need.
10. Priority order: Correctness → Safety → Simplicity → Performance → Features.

## 12. Definition of done
A feature is done only when its requirement is implemented, tested, manually verified where applicable, documented, and committed.

## 13. MVP success criteria
A farmer can:
1. create an account
2. create a farm/plot
3. select a crop and variety
4. track crop age/growth stage
5. read stage-relevant guidance
6. record farm activities
7. record fertilizer/input usage
8. inspect pest/disease information
9. record expenses and harvest
10. ask the AI assistant questions grounded in verified knowledge

## 14. Milestones
### M001 — Technical Foundation
Repository structure, backend/frontend skeleton, PostgreSQL, configuration, health endpoint, initial schema, authentication foundation, Docker development environment, automated verification, Git workflow.

### M002 — Core Farm Domain
Farm, plot, crop, variety, growth-stage models and APIs.

### M003 — Knowledge Base
Crop, fertilizer, pest, disease, weed, pesticide and guide management.

### M004 — Farm Operations
Activities, inputs, expenses, harvest.

### M005 — Farmer Experience
Mobile-first dashboard and workflows.

### M006 — AI Foundation
Knowledge retrieval, chat, provenance and safety controls.

### M007 — Smart Features
Image diagnosis, weather, notifications.

### M008 — Ecosystem
Experts, groups, suppliers, marketplace.
