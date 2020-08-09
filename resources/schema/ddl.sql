CREATE TABLE princess_article (
	article_id bigserial NOT NULL,
	author bigint NOT NULL,
	status bigint NULL,
	article_type bigint NULL,
	title varchar(1000) NOT NULL,
	article_content text NULL,
	config text NULL,
	time_created bigint NOT NULL,
	time_updated bigint NULL,
	CONSTRAINT article_pk PRIMARY KEY (article_id)
);
CREATE INDEX princess_article_created_status_idx ON princess_article (time_created DESC,status);

COMMENT ON COLUMN princess_article.status IS '0 - normal
1 - deleted';
COMMENT ON COLUMN princess_article.article_type IS '1 - article
';

CREATE TABLE princess_user (
	user_id bigserial NOT NULL,
	user_name varchar(300) NOT NULL,
	user_password varchar(500) NOT NULL,
	display_name varchar(300) NULL,
	status bigint NULL,
	config text NULL,
	user_type bigint NULL,
	time_created bigint NOT NULL,
	time_updated bigint NULL,
	CONSTRAINT princess_user_pk PRIMARY KEY (user_id)
);
CREATE UNIQUE INDEX princess_user_user_name_idx ON princess_user (user_name);

COMMENT ON COLUMN princess_user.status IS '0 - normal
1 - deleted';
COMMENT ON COLUMN princess_user.user_type IS '1 - user';
