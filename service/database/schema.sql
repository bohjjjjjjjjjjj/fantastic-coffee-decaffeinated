-- 1. Tabella Utenti
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    photo_url TEXT DEFAULT ''
);

-- 2. Tabella Sessioni (per l'autenticazione Bearer)
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. Tabella Conversazioni (gestisce sia chat 1-to-1 che gruppi)
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    name TEXT DEFAULT '',         -- Usato se è un gruppo
    photo_url TEXT DEFAULT '',     -- Usato se è un gruppo
    is_group BOOLEAN DEFAULT FALSE
);

-- 4. Tabella Membri del Gruppo / Conversazione
CREATE TABLE IF NOT EXISTS conversation_members (
    conversation_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY (conversation_id, user_id),
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 5. Tabella Messaggi
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    sender_username TEXT NOT NULL,
    content_type TEXT NOT NULL,   -- 'text' oppure 'photo'
    content_value TEXT NOT NULL,  -- Il testo del messaggio O la photoUrl
    reply_to_message_id TEXT,    -- Facoltativo, per le risposte
    data_sent DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'sent',   -- 'sent', 'delivered', 'read', 'failed'
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_to_message_id) REFERENCES messages(id) ON DELETE SET NULL
);

-- 6. Tabella Reazioni
CREATE TABLE IF NOT EXISTS reactions (
    id TEXT PRIMARY KEY,
    message_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    reaction_type TEXT NOT NULL,  -- Es. "like", "heart", ecc.
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);