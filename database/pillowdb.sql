-- -------------------------------------------------------------
-- Zapra Gartiast (zapra@me.com)
-- 
-- Pillow database schema
--
-- Generation Time: 2025-09-13 11:24:50.2810 AM
-- -------------------------------------------------------------


-- Table Definition
CREATE TABLE "public"."Roles" (
    "id" uuid NOT NULL,
    "name" varchar(50) NOT NULL,
    "description" text,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Table Definition
CREATE TABLE "public"."Role_Permissions" (
    "role_id" uuid NOT NULL,
    "permission_id" uuid,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("role_id")
);

-- Table Definition
CREATE TABLE "public"."Users" (
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
CREATE TABLE "public"."Organizations" (
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
CREATE TABLE "public"."User_Roles" (
    "user_id" uuid,
    "role_id" uuid,
    "scope" varchar(50) DEFAULT 'org'::character varying,
    "parent_role_id" uuid,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now()
);

-- Table Definition
CREATE TABLE "public"."Permissions" (
    "id" uuid NOT NULL,
    "name" varchar(50) NOT NULL,
    "description" text,
    "scope_level" varchar(50) DEFAULT 'user'::character varying,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Table Definition
CREATE TABLE "public"."User_Organizations" (
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
CREATE TABLE "public"."Audit_Log" (
    "id" uuid NOT NULL,
    "user_id" uuid,
    "action" varchar(100),
    "details" text,
    "timestamp" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);



-- Indices
CREATE UNIQUE INDEX "Roles_name_key" ON public."Roles" USING btree (name);
ALTER TABLE "public"."Role_Permissions" ADD FOREIGN KEY ("role_id") REFERENCES "public"."Roles"("id");
ALTER TABLE "public"."Role_Permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "public"."Permissions"("id");


-- Indices
CREATE UNIQUE INDEX "Users_username_key" ON public."Users" USING btree (username);
CREATE UNIQUE INDEX "Users_email_key" ON public."Users" USING btree (email);
ALTER TABLE "public"."Organizations" ADD FOREIGN KEY ("parent_org_id") REFERENCES "public"."Organizations"("id");
ALTER TABLE "public"."Organizations" ADD FOREIGN KEY ("managed_by") REFERENCES "public"."Users"("id");


-- Indices
CREATE UNIQUE INDEX "Organizations_name_key" ON public."Organizations" USING btree (name);
ALTER TABLE "public"."User_Roles" ADD FOREIGN KEY ("user_id") REFERENCES "public"."Users"("id");
ALTER TABLE "public"."User_Roles" ADD FOREIGN KEY ("role_id") REFERENCES "public"."Roles"("id");


-- Indices
CREATE UNIQUE INDEX "Permissions_name_key" ON public."Permissions" USING btree (name);
ALTER TABLE "public"."User_Organizations" ADD FOREIGN KEY ("invited_by") REFERENCES "public"."Users"("id");
ALTER TABLE "public"."User_Organizations" ADD FOREIGN KEY ("role_id") REFERENCES "public"."Roles"("id");
ALTER TABLE "public"."User_Organizations" ADD FOREIGN KEY ("user_id") REFERENCES "public"."Users"("id");
ALTER TABLE "public"."User_Organizations" ADD FOREIGN KEY ("org_id") REFERENCES "public"."Organizations"("id");
ALTER TABLE "public"."Audit_Log" ADD FOREIGN KEY ("user_id") REFERENCES "public"."Users"("id");
