CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_songs_group_name ON songs(group_name);
CREATE INDEX idx_songs_song_name ON songs(song_name);
