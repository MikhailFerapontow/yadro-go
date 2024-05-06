CREATE VIEW IF NOT EXISTS kw_comic
AS
SELECT
    ckw.comic_id as comic_id,
    c.url as url,
    k.id as keyword_id,
    ckw.weight as weight
FROM comic_keyword as ckw
JOIN Comic as c ON ckw.comic_id = c.id
JOIN Keyword as k ON ckw.word_id = k.id;

-- SELECT url, weight FROM kw_comic WHERE keyword_id = 99
