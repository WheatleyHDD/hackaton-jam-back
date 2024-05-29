CREATE TABLE "users" (
	"email" varchar(255) NOT NULL,
	"username" varchar(255) NOT NULL,
	"avatar" varchar(255) NOT NULL DEFAULT 'https://i.imgur.com/b0zqmkj.jpeg',
	"first_name" varchar(255) NOT NULL,
	"last_name" varchar(255) NOT NULL,
	"middle_name" varchar(255),
	"password" varchar(255) NOT NULL,
	"about" TEXT,
	"work_place" varchar(255),
	"work_time" varchar(255),
	"loc" varchar(255),
	"perms" int NOT NULL DEFAULT 0,
	CONSTRAINT "users_pk" PRIMARY KEY ("email")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "tokens" (
	"token" varchar(255) NOT NULL,
	"user_email" varchar(255) NOT NULL,
	CONSTRAINT "tokens_pk" PRIMARY KEY ("token")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "skills" (
	"user_email" varchar(255) NOT NULL,
	"skill" varchar(255) NOT NULL,
	CONSTRAINT "skills_pk" PRIMARY KEY ("user_email")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "contacts" (
	"user_email" varchar(255) NOT NULL,
	"contact_link" varchar(255) NOT NULL,
	CONSTRAINT "contacts_pk" PRIMARY KEY ("user_email")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "teams" (
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL UNIQUE DEFAULT 'Без названия',
	CONSTRAINT "teams_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "teams_members" (
	"team_id" bigint NOT NULL,
	"member_email" varchar(255) NOT NULL,
	"role" varchar(255) NOT NULL,
	CONSTRAINT "teams_members_pk" PRIMARY KEY ("team_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "events" (
	"urid" varchar(255) NOT NULL,
	"id" serial NOT NULL,
	"name" varchar(255) NOT NULL,
	"start_time" timestamp with time zone NOT NULL,
	"end_time" timestamp with time zone NOT NULL,
	"prize" varchar(255),
	"location" varchar(255),
	"desc" TEXT NOT NULL,
	"requirements" TEXT,
	"icon" varchar(255),
	"is_irl" bool NOT NULL DEFAULT 'false',
	"team_requirements_type" int NOT NULL DEFAULT '0',
	"team_requirements_value" int NOT NULL DEFAULT '5',
	CONSTRAINT "events_pk" PRIMARY KEY ("urid")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "event_orgs" (
	"event_uri" varchar(255) NOT NULL,
	"organizator_email" varchar(255) NOT NULL,
	CONSTRAINT "event_orgs_pk" PRIMARY KEY ("event_uri")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "event_members" (
	"event_uri" varchar(255) NOT NULL,
	"member_email" varchar(255) NOT NULL,
	CONSTRAINT "event_members_pk" PRIMARY KEY ("event_uri")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "event_blog" (
	"id" serial NOT NULL,
	"event_uri" varchar(255) NOT NULL,
	"title" varchar(255) NOT NULL,
	"author" varchar(255) NOT NULL,
	"post_date" DATE NOT NULL,
	"post_text" TEXT NOT NULL,
	CONSTRAINT "event_blog_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "event_tags" (
	"event_uri" varchar(255) NOT NULL,
	"tag" varchar(255) NOT NULL,
	CONSTRAINT "event_tags_pk" PRIMARY KEY ("event_uri")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "event_partners" (
	"event_uri" varchar(255) NOT NULL,
	"logo_url" varchar(255) NOT NULL,
	CONSTRAINT "event_partners_pk" PRIMARY KEY ("event_uri")
) WITH (
  OIDS=FALSE
);




ALTER TABLE "tokens" ADD CONSTRAINT "tokens_fk0" FOREIGN KEY ("user_email") REFERENCES "users"("email");

ALTER TABLE "skills" ADD CONSTRAINT "skills_fk0" FOREIGN KEY ("user_email") REFERENCES "users"("email");

ALTER TABLE "contacts" ADD CONSTRAINT "contacts_fk0" FOREIGN KEY ("user_email") REFERENCES "users"("email");


ALTER TABLE "teams_members" ADD CONSTRAINT "teams_members_fk0" FOREIGN KEY ("team_id") REFERENCES "teams"("id");
ALTER TABLE "teams_members" ADD CONSTRAINT "teams_members_fk1" FOREIGN KEY ("member_email") REFERENCES "users"("email");


ALTER TABLE "event_orgs" ADD CONSTRAINT "event_orgs_fk0" FOREIGN KEY ("event_uri") REFERENCES "events"("urid");
ALTER TABLE "event_orgs" ADD CONSTRAINT "event_orgs_fk1" FOREIGN KEY ("organizator_email") REFERENCES "users"("email");

ALTER TABLE "event_members" ADD CONSTRAINT "event_members_fk0" FOREIGN KEY ("event_uri") REFERENCES "events"("urid");
ALTER TABLE "event_members" ADD CONSTRAINT "event_members_fk1" FOREIGN KEY ("member_email") REFERENCES "users"("email");

ALTER TABLE "event_blog" ADD CONSTRAINT "event_blog_fk0" FOREIGN KEY ("event_uri") REFERENCES "events"("urid");
ALTER TABLE "event_blog" ADD CONSTRAINT "event_blog_fk1" FOREIGN KEY ("author") REFERENCES "users"("email");

ALTER TABLE "event_tags" ADD CONSTRAINT "event_tags_fk0" FOREIGN KEY ("event_uri") REFERENCES "events"("urid");

ALTER TABLE "event_partners" ADD CONSTRAINT "event_partners_fk0" FOREIGN KEY ("event_uri") REFERENCES "events"("urid");


-- Пробные данные
INSERT INTO "users" ("email", "username", "first_name", "last_name", "password", "perms") VALUES ('thatmaidguy1@ya.ru', 'admin', 'Админ', 'Админов', '$2a$10$DmTlEGzS/Ix0JFfTT3hmH.ZLliSvSMRlkTBVoo2F6uBZiQwXP1YVy', 10);
INSERT INTO "users" ("email", "username", "first_name", "last_name", "password", "perms") VALUES ('thatmaidguy2@ya.ru', 'organizator', 'Организатор', 'Организаторов', '$2a$10$DmTlEGzS/Ix0JFfTT3hmH.ZLliSvSMRlkTBVoo2F6uBZiQwXP1YVy', 1);
INSERT INTO "users" ("email", "username", "first_name", "last_name", "password", "perms") VALUES ('thatmaidguy3@ya.ru', 'user', 'Иван', 'Иванов', '$2a$10$DmTlEGzS/Ix0JFfTT3hmH.ZLliSvSMRlkTBVoo2F6uBZiQwXP1YVy', 0);

--INSERT INTO "events" ("urid", "name", "start_time", "end_time", "prize", "location", "desc", "requirements", "icon", "is_irl", "team_requirements_type", "team_requirements_value") VALUES ("example_event", "Example Event1", "29-05-2024", "30-05-2024", "200 рублей выплот", "Екатеринбург", "Тестовое описание", "тест", "", false, 0, 5);
--INSERT INTO "event_orgs" ("event_uri", "organizator_email") VALUES ("example_event", "thatmaidguy2@ya.ru")
--INSERT INTO "event_orgs" ("event_uri", "organizator_email") VALUES ("example_event", "thatmaidguy2@ya.ru")