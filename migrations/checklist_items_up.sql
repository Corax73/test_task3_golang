CREATE TABLE IF NOT EXISTS "checklist_items" (
    "id" serial PRIMARY KEY,
    "checklist_id" int NOT NULL,
    "description" varchar NOT NULL DEFAULT '',
    "created_at" date NOT NULL DEFAULT 'now()',
    FOREIGN KEY ("checklist_id") REFERENCES "checklists" ("id")
);
CREATE INDEX ON "checklist_items" ("checklist_id");