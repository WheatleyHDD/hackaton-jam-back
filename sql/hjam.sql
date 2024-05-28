CREATE TABLE IF NOT EXISTS "users" (
	"email" text NOT NULL UNIQUE,
	"username" text NOT NULL,
	"avatar" text,
	"first_name" text NOT NULL,
	"last_name" text NOT NULL,
	"middle_name" text,
	"password" text NOT NULL,
	"about" text,
	"work_place" text,
	"work_time" text,
	"loc" text,
	"perms" bigint NOT NULL DEFAULT '0',
	PRIMARY KEY ("email")
);

CREATE TABLE IF NOT EXISTS "skills" (
	"user_email" text NOT NULL UNIQUE,
	"skill" text NOT NULL,
	PRIMARY KEY ("user_email")
);

CREATE TABLE IF NOT EXISTS "contacts" (
	"user_email" text NOT NULL UNIQUE,
	"contact_link" text NOT NULL,
	PRIMARY KEY ("user_email")
);

CREATE TABLE IF NOT EXISTS "events" (
	"urid" text NOT NULL UNIQUE,
	"id" bigserial NOT NULL UNIQUE,
	"name" text NOT NULL,
	"start_time" timestamp with time zone NOT NULL,
	"end_time" timestamp with time zone NOT NULL,
	"prize" text,
	"location" text,
	"desc" text NOT NULL,
	"requirements" text,
	"partners" text,
	"icon" text,
	"is_irl" boolean NOT NULL DEFAULT false,
	"team_requirements_type" bigint NOT NULL,
	"team_requirements_value" bigint NOT NULL,
	PRIMARY KEY ("urid")
);

CREATE TABLE IF NOT EXISTS "event_members" (
	"event_uri" text NOT NULL,
	"member_email" text NOT NULL,
	PRIMARY KEY ("event_uri")
);

CREATE TABLE IF NOT EXISTS "event_orgs" (
	"event_uri" text NOT NULL,
	"organizator_email" text NOT NULL,
	PRIMARY KEY ("event_uri")
);

CREATE TABLE IF NOT EXISTS "event_blog" (
	"id" bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
	"event_uri" text NOT NULL,
	"title" text NOT NULL,
	"author" text NOT NULL,
	"post_date" date NOT NULL,
	"post_text" text NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "tokens" (
	"token" text NOT NULL UNIQUE,
	"user_email" text NOT NULL,
	PRIMARY KEY ("token")
);

CREATE TABLE IF NOT EXISTS "teams" (
	"id" bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
	"name" text NOT NULL DEFAULT 'Без названия',
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "teams_members" (
	"team_id" bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
	"member_email" text NOT NULL,
	"role" text NOT NULL,
	PRIMARY KEY ("team_id")
);


ALTER TABLE "skills" ADD CONSTRAINT "skills_fk0" FOREIGN KEY ("user_email") REFERENCES "users"("email");
ALTER TABLE "contacts" ADD CONSTRAINT "contacts_fk0" FOREIGN KEY ("user_email") REFERENCES "users"("email");

ALTER TABLE "event_members" ADD CONSTRAINT "event_members_fk0" FOREIGN KEY ("event_uri") REFERENCES "events"("urid");

ALTER TABLE "event_members" ADD CONSTRAINT "event_members_fk1" FOREIGN KEY ("member_email") REFERENCES "users"("email");
ALTER TABLE "event_orgs" ADD CONSTRAINT "event_orgs_fk0" FOREIGN KEY ("event_uri") REFERENCES "events"("urid");

ALTER TABLE "event_orgs" ADD CONSTRAINT "event_orgs_fk1" FOREIGN KEY ("organizator_email") REFERENCES "users"("email");
ALTER TABLE "event_blog" ADD CONSTRAINT "event_blog_fk1" FOREIGN KEY ("event_uri") REFERENCES "events"("urid");

ALTER TABLE "event_blog" ADD CONSTRAINT "event_blog_fk3" FOREIGN KEY ("author") REFERENCES "users"("email");
ALTER TABLE "tokens" ADD CONSTRAINT "tokens_fk1" FOREIGN KEY ("user_email") REFERENCES "users"("email");

ALTER TABLE "teams_members" ADD CONSTRAINT "teams_members_fk0" FOREIGN KEY ("team_id") REFERENCES "teams"("id");

ALTER TABLE "teams_members" ADD CONSTRAINT "teams_members_fk1" FOREIGN KEY ("member_email") REFERENCES "users"("email");