CREATE TABLE articles(
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    search_vector tsvector GENERATED ALWAYS AS(
        to_tsvector('english', title || ' ' || body)
    ) STORED
);

CREATE INDEX article_search_gin_idx ON articles USING GIN(search_vector);

SELECT title
FROM articles
WHERE search_vector @@ websearch_to_tsquery('english', 'postgres indexes');