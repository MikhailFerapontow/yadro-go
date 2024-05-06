CREATE TABLE IF NOT EXISTS comic_keyword
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    comic_id INTEGER REFERENCES Comic (id),
    word_id  INTEGER REFERENCES Keyword (id),
    weight   INTEGER
);
