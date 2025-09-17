-- -------------------------------------------------------------
-- Zapra Gartiast (zapra@me.com)
-- 
-- Pillow database schema
--
-- Generation Time: 2025-09-13 11:24:50.2810 AM
-- -------------------------------------------------------------


-- Table Definition
CREATE TABLE "public"."roles" (
    "id" uuid NOT NULL,
    "name" varchar(50) NOT NULL,
    "description" text,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Table Definition
CREATE TABLE "public"."role_permissions" (
    "role_id" uuid NOT NULL,
    "permission_id" uuid,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now()
);

-- Table Definition
CREATE TABLE "public"."users" (
    "id" uuid NOT NULL,
    "username" varchar(50) NOT NULL,
    "password_hash" text NOT NULL,
    "email" varchar(100),
    "is_active" bool DEFAULT true,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Table Definition
CREATE TABLE "public"."organizations" (
    "id" uuid NOT NULL,
    "name" varchar(100) NOT NULL,
    "description" text,
    "domain" varchar(100),
    "managed_by" uuid,
    "parent_org_id" uuid,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Table Definition
CREATE TABLE "public"."user_roles" (
    "user_id" uuid,
    "role_id" uuid,
    "scope" varchar(50) DEFAULT 'org'::character varying,
    "parent_role_id" uuid,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now()
);

-- Table Definition
CREATE TABLE "public"."permissions" (
    "id" uuid NOT NULL,
    "name" varchar(50) NOT NULL,
    "description" text,
    "scope_level" varchar(50) DEFAULT 'user'::character varying,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Table Definition
CREATE TABLE "public"."user_organizations" (
    "id" uuid NOT NULL,
    "user_id" uuid,
    "org_id" uuid,
    "role_id" uuid,
    "invited_by" uuid,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Table Definition
CREATE TABLE "public"."audit_log" (
    "id" uuid NOT NULL,
    "user_id" uuid,
    "action" varchar(100),
    "details" text,
    "timestamp" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Create custom_fields table for global field definitions
CREATE TABLE IF NOT EXISTS "public"."custom_fields" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "name" varchar(100) NOT NULL,
    "label" varchar(255) NOT NULL,
    "type" varchar(50) NOT NULL,
    "required" boolean DEFAULT false,
    "options" jsonb, -- For select/multiselect fields
    "validation" jsonb, -- Validation rules
    "order" integer DEFAULT 0,
    "is_active" boolean DEFAULT true,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Create user_custom_field_values table for storing user-specific values
CREATE TABLE IF NOT EXISTS "public"."user_custom_field_values" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),
    "user_id" uuid NOT NULL,
    "field_id" uuid NOT NULL,
    "value" text, -- Store as text, parse based on field type
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id"),
    UNIQUE ("user_id", "field_id") -- One value per user per field
);

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS "custom_fields_name_idx" ON "public"."custom_fields" ("name");
CREATE INDEX IF NOT EXISTS "custom_fields_active_idx" ON "public"."custom_fields" ("is_active");
CREATE INDEX IF NOT EXISTS "user_custom_field_values_user_id_idx" ON "public"."user_custom_field_values" ("user_id");
CREATE INDEX IF NOT EXISTS "user_custom_field_values_field_id_idx" ON "public"."user_custom_field_values" ("field_id");

-- Add foreign key constraints
ALTER TABLE "public"."user_custom_field_values"
ADD CONSTRAINT "fk_user_custom_field_values_user_id"
FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE;

ALTER TABLE "public"."user_custom_field_values"
ADD CONSTRAINT "fk_user_custom_field_values_field_id"
FOREIGN KEY ("field_id") REFERENCES "public"."custom_fields"("id") ON DELETE CASCADE;

-- Add comments
COMMENT ON TABLE "public"."custom_fields" IS 'Global custom field definitions that apply to all users';
COMMENT ON TABLE "public"."user_custom_field_values" IS 'User-specific values for global custom fields';
COMMENT ON COLUMN "public"."custom_fields"."type" IS 'Field type: text, textarea, number, email, phone, date, boolean, select, multiselect';
COMMENT ON COLUMN "public"."custom_fields"."options" IS 'JSON array of options for select/multiselect fields';
COMMENT ON COLUMN "public"."custom_fields"."validation" IS 'JSON object with validation rules (min_length, max_length, min, max, pattern)';

-- Insert some sample global fields (optional)
-- Uncomment these lines to add sample fields
/*
INSERT INTO "public"."custom_fields" ("name", "label", "type", "required", "order") VALUES
('full_name', 'Full Name', 'text', false, 1),
('phone_number', 'Phone Number', 'phone', false, 2),
('date_of_birth', 'Date of Birth', 'date', false, 3),
('department', 'Department', 'select', false, 4),
('bio', 'Bio', 'textarea', false, 5);

-- Add options for department field
UPDATE "public"."custom_fields"
SET "options" = '["Engineering", "Marketing", "Sales", "HR", "Finance"]'::jsonb
WHERE "name" = 'department';
*/

-- Indices
CREATE INDEX "roles_name_key" ON public."roles" USING btree (name);
ALTER TABLE "public"."role_permissions" ADD FOREIGN KEY ("role_id") REFERENCES "public"."roles"("id");
ALTER TABLE "public"."role_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "public"."permissions"("id");


-- Indices
CREATE UNIQUE INDEX "users_username_key" ON public."users" USING btree (username);
CREATE UNIQUE INDEX "users_email_key" ON public."users" USING btree (email);
ALTER TABLE "public"."organizations" ADD FOREIGN KEY ("parent_org_id") REFERENCES "public"."organizations"("id");
ALTER TABLE "public"."organizations" ADD FOREIGN KEY ("managed_by") REFERENCES "public"."users"("id");


-- Indices
CREATE UNIQUE INDEX "organizations_name_key" ON public."organizations" USING btree (name);
ALTER TABLE "public"."user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");
ALTER TABLE "public"."user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "public"."roles"("id");


-- Indices
CREATE UNIQUE INDEX "permissions_name_key" ON public."permissions" USING btree (name);
ALTER TABLE "public"."user_organizations" ADD FOREIGN KEY ("invited_by") REFERENCES "public"."users"("id");
ALTER TABLE "public"."user_organizations" ADD FOREIGN KEY ("role_id") REFERENCES "public"."roles"("id");
ALTER TABLE "public"."user_organizations" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");
ALTER TABLE "public"."user_organizations" ADD FOREIGN KEY ("org_id") REFERENCES "public"."organizations"("id");
ALTER TABLE "public"."audit_log" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");
