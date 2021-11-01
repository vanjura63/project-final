DROP TABLE IF EXISTS comment;
CREATE TABLE comment (
    id BIGSERIAL PRIMARY KEY,
    news_id BIGINT NOT NULL,
    text_comment TEXT NOT NULL,
    parent_id BIGINT NOT NULL,
    pub_time BIGINT NOT NULL,
    profane BOOLEAN NOT NULL,
    );