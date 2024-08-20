
-- DROP SCHEMA `xl_direct`;
-- CREATE SCHEMA `xl_direct`;

CREATE TABLE IF NOT EXISTS "services" (
  "id" SERIAL PRIMARY KEY,
  "category" varchar(20) NOT NULL,
  "code" varchar(25) UNIQUE NOT NULL,
  "name" varchar(50) NOT NULL,
  "price" float(5) DEFAULT 0,
  "product_id"  varchar(25),
  "sid_optin" varchar(35) NOT NULL,
  "sid_mt" varchar(35) NOT NULL,
  "renewal_day" int DEFAULT 0,
  "trial_day" int DEFAULT 0,
  "url_telco" varchar(85),
  "url_portal" varchar(85),
  "url_callback" varchar(85),
  "url_notif_sub" varchar(85),
  "url_notif_unsub" varchar(85),
  "url_notif_renewal" varchar(85),
  "url_postback" varchar(85),
  "url_postback_billable" varchar(85)
);

CREATE TABLE IF NOT EXISTS "contents" (
  "id" SERIAL PRIMARY KEY,
  "service_id" int NOT NULL,
  "name" varchar(20) NOT NULL,
  "value" varchar(400) NOT NULL,
  "tid" varchar(5) NOT NULL,
  "sequence" int NOT NULL DEFAULT 0,
  FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);


CREATE TABLE IF NOT EXISTS "schedules" (
  "id" int,
  "name" varchar(20) NOT NULL,
  "publish_at" timestamp,
  "unlocked_at" timestamp,
  "is_unlocked" bool DEFAULT false,
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "subscriptions" (
  "id" SERIAL PRIMARY KEY,
  "service_id" int NOT NULL,
  "category" varchar(20) NOT NULL,
  "msisdn" varchar(60) NOT NULL,
  "sub_id" int NOT NULL,
  "channel" varchar(20),
  "camp_keyword" varchar(55),
  "camp_sub_keyword" varchar(55),
  "adnet" varchar(55),
  "pub_id" varchar(55),
  "aff_sub" varchar(100),
  "latest_trxid" varchar(100),
  "latest_keyword" varchar(180),
  "latest_subject" varchar(20),
  "latest_status" varchar(20),
  "latest_pin" varchar(10),
  "latest_payload" varchar(100),
  "amount" float(8) DEFAULT 0,
  "trial_at" timestamp,
  "renewal_at" timestamp,
  "unsub_at" timestamp,
  "charge_at" timestamp,
  "retry_at" timestamp,
  "purge_at" timestamp,
  "success" int DEFAULT 0,
  "failed" int DEFAULT 0,
  "ip_address" varchar(50),
  "purge_reason" varchar(100),
  "is_trial" bool DEFAULT false,
  "content_sequence" int NOT NULL DEFAULT 0,
  "is_retry" bool DEFAULT false,
  "is_confirm" bool DEFAULT false,
  "is_purge" bool DEFAULT false,
  "is_active" bool DEFAULT false,
  "total_firstpush" int DEFAULT 0,
  "total_renewal" int DEFAULT 0,
  "total_amount_firstpush" float(8) DEFAULT 0,
  "total_amount_renewal" float(8) DEFAULT 0,
  "created_at" timestamp,
  "updated_at" timestamp,
  FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);

CREATE TABLE IF NOT EXISTS "transactions" (
  "id" SERIAL PRIMARY KEY,
  "tx_id" varchar(100),
  "service_id" int NOT NULL,
  "msisdn" varchar(60) NOT NULL,
  "sub_id" int NOT NULL,
  "channel" varchar(20) NOT NULL,
  "camp_keyword" varchar(55),
  "camp_sub_keyword" varchar(55),
  "adnet" varchar(55),
  "pub_id" varchar(55),
  "aff_sub" varchar(100),
  "keyword" varchar(180),
  "amount" float(8) DEFAULT 0,
  "pin" varchar(10),
  "status" varchar(45),
  "status_code" varchar(45),
  "status_detail" varchar(100),
  "subject" varchar(45),
  "ip_address" varchar(45),
  "payload" text,
  "created_at" timestamp,
  "updated_at" timestamp,
  FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);

CREATE TABLE IF NOT EXISTS "histories" (
  "id" SERIAL PRIMARY KEY,
  "service_id" int NOT NULL,
  "msisdn" varchar(60) NOT NULL,
  "sub_id" int NOT NULL,
  "channel" varchar(20),
  "adnet" varchar(20),
  "keyword" varchar(180),
  "subject" varchar(20),
  "status" varchar(45),
  "ip_address" varchar(45),
  "created_at" timestamp,
  FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);

CREATE TABLE IF NOT EXISTS "blacklists" (
  "id" SERIAL PRIMARY KEY,
  "msisdn" varchar(60) UNIQUE NOT NULL,
  "created_at" timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "traffics_campaign" (
  "id" SERIAL PRIMARY KEY,
  "tx_id" varchar(100) UNIQUE NOT NULL,
  "service_id" int NOT NULL,
  "camp_keyword" varchar(55),
  "camp_sub_keyword" varchar(55),
  "adnet" varchar(55),
  "pub_id" varchar(55),
  "aff_sub" varchar(100),
  "browser" varchar(200),
  "os" varchar(100),
  "device" varchar(200),
  "referer" varchar(300),
  "ip_address" varchar(45),
  "created_at" timestamp,
  FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);

CREATE TABLE IF NOT EXISTS "traffics_mo" (
  "id" SERIAL PRIMARY KEY,
  "tx_id" varchar(100) UNIQUE NOT NULL,
  "service_id" int NOT NULL,
  "msisdn" varchar(60) NOT NULL,
  "channel" varchar(20),
  "camp_keyword" varchar(55),
  "camp_sub_keyword" varchar(55),
  "subject" varchar(55),
  "adnet" varchar(55),
  "pub_id" varchar(55),
  "aff_sub" varchar(100),
  "is_charge" boolean DEFAULT false,
  "ip_address" varchar(45),
  "created_at" timestamp,
  FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);

CREATE TABLE IF NOT EXISTS "dailypushes" (
  "id" SERIAL PRIMARY KEY,
  "tx_id" varchar(100) UNIQUE NOT NULL,
  "subscription_id" int NOT NULL,
  "service_id" int NOT NULL,
  "msisdn" varchar(60) NOT NULL,
  "channel" varchar(20),
  "camp_keyword" varchar(55),
  "camp_sub_keyword" varchar(55),
  "subject" varchar(55),
  "adnet" varchar(55),
  "pub_id" varchar(55),
  "aff_sub" varchar(100),
  "status_code" varchar(45),
  "status_detail" varchar(100),
  "is_charge" boolean DEFAULT false,
  "ip_address" varchar(45),
  "created_at" timestamp,
  "updated_at" timestamp,
  FOREIGN KEY ("subscription_id") REFERENCES "subscriptions" ("id"),
  FOREIGN KEY ("service_id") REFERENCES "services" ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "uidx_msisdn" ON "blacklists" ("msisdn");
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_service_msisdn" ON "subscriptions" ("service_id", "msisdn");
CREATE INDEX IF NOT EXISTS "idx_latest_pin" ON "subscriptions" ("latest_pin");
CREATE INDEX IF NOT EXISTS "idx_category_msisdn" ON "subscriptions" ("category", "msisdn");
CREATE INDEX IF NOT EXISTS "idx_service_msisdn" ON "transactions" ("service_id", "msisdn");
CREATE INDEX IF NOT EXISTS "idx_service_name" ON "contents" ("service_id", "name");
CREATE INDEX IF NOT EXISTS "idx_name_publish_at" ON "schedules" ("name", "publish_at");
CREATE INDEX IF NOT EXISTS "idx_traffic_service_msisdn" ON "traffics_mo" ("service_id", "msisdn");
CREATE INDEX IF NOT EXISTS "idx_dp_service_msisdn" ON "dailypushes" ("service_id", "msisdn");


ALTER TABLE "services" ADD "is_content_sequence" bool DEFAULT false;
ALTER TABLE "contents" ADD "sequence" int NOT NULL DEFAULT 0;
ALTER TABLE "subscriptions" ADD "content_sequence" int NOT NULL DEFAULT 0;
ALTER TABLE "traffics_campaign" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
ALTER TABLE "traffics_mo" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
