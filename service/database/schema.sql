-- Schema del database WASAText.
-- Riferimento per la relazione: lo schema effettivo vive nella costante `schema`
-- in database.go e viene applicato automaticamente a ogni avvio.

-- 1. Utenti
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    photo_url TEXT NOT NULL DEFAULT ''
);

-- 2. Sessioni (autenticazione Bearer)
CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 3. Conversazioni (chat dirette e gruppi)
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL DEFAULT '',      -- usato per i gruppi
    photo_url TEXT NOT NULL DEFAULT '', -- usato per i gruppi
    is_group BOOLEAN NOT NULL DEFAULT 0
);

-- 4. Membri di una conversazione
CREATE TABLE IF NOT EXISTS conversation_members (
    conversation_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY (conversation_id, user_id),
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 5. Messaggi (proprietà legata a sender_id, non allo username)
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    content_type TEXT NOT NULL,   -- 'text' oppure 'photo'
    content_value TEXT NOT NULL,  -- testo del messaggio o photoUrl
    reply_to_message_id TEXT,     -- opzionale, per le risposte
    is_forwarded BOOLEAN NOT NULL DEFAULT 0,  -- messaggio inoltrato
    data_sent DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_to_message_id) REFERENCES messages(id) ON DELETE SET NULL
);

-- 6. Stato di consegna/lettura per destinatario (spunte)
CREATE TABLE IF NOT EXISTS message_status (
    message_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    status TEXT NOT NULL,          -- 'delivered' oppure 'read'
    PRIMARY KEY (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 7. Reazioni (una per utente per messaggio)
CREATE TABLE IF NOT EXISTS reactions (
    id TEXT PRIMARY KEY,
    message_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    reaction_type TEXT NOT NULL,
    UNIQUE (message_id, user_id),
    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
