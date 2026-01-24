-- Create users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  auth_id UUID NOT NULL UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
  gender TEXT NOT NULL CHECK (gender IN ('male', 'female')),
  preferred_gender TEXT CHECK (preferred_gender IN ('male', 'female', 'any')),
  is_online BOOLEAN DEFAULT FALSE,
  last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create matchmaking queue table
CREATE TABLE matchmaking_queue (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  gender TEXT NOT NULL,
  preferred_gender TEXT NOT NULL,
  joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create chat sessions table
CREATE TABLE chat_sessions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  status TEXT DEFAULT 'active' CHECK (status IN ('active', 'ended', 'reported')),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  ended_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT different_users CHECK (user1_id != user2_id)
);

-- Create messages table
CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  chat_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
  sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create reports/moderation table
CREATE TABLE reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  reported_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  chat_session_id UUID REFERENCES chat_sessions(id) ON DELETE SET NULL,
  reason TEXT NOT NULL,
  description TEXT,
  status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'resolved')),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  reviewed_at TIMESTAMP WITH TIME ZONE,
  reviewer_id UUID REFERENCES users(id) ON DELETE SET NULL
);

-- Enable Row Level Security (RLS)
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE matchmaking_queue ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE reports ENABLE ROW LEVEL SECURITY;

-- RLS Policies for users table
CREATE POLICY "Users can view their own profile" ON users
  FOR SELECT USING (auth.uid() = auth_id);

CREATE POLICY "Users can update their own profile" ON users
  FOR UPDATE USING (auth.uid() = auth_id);

-- RLS Policies for matchmaking_queue table
CREATE POLICY "Users can view their own queue entry" ON matchmaking_queue
  FOR SELECT USING (auth.uid() = (SELECT auth_id FROM users WHERE id = user_id));

CREATE POLICY "Users can insert their own queue entry" ON matchmaking_queue
  FOR INSERT WITH CHECK (auth.uid() = (SELECT auth_id FROM users WHERE id = user_id));

CREATE POLICY "Users can delete their own queue entry" ON matchmaking_queue
  FOR DELETE USING (auth.uid() = (SELECT auth_id FROM users WHERE id = user_id));

-- RLS Policies for chat_sessions table
CREATE POLICY "Users can view their own chat sessions" ON chat_sessions
  FOR SELECT USING (
    auth.uid() = (SELECT auth_id FROM users WHERE id = user1_id) OR
    auth.uid() = (SELECT auth_id FROM users WHERE id = user2_id)
  );

-- RLS Policies for messages table
CREATE POLICY "Users can view messages from their chats" ON messages
  FOR SELECT USING (
    auth.uid() = sender_id OR
    auth.uid() = (
      SELECT COALESCE(
        (SELECT auth_id FROM users WHERE id = user1_id),
        (SELECT auth_id FROM users WHERE id = user2_id)
      ) FROM chat_sessions WHERE id = chat_id
    )
  );

CREATE POLICY "Users can insert messages to their chats" ON messages
  FOR INSERT WITH CHECK (
    auth.uid() = (SELECT auth_id FROM users WHERE id = sender_id)
  );

-- RLS Policies for reports table
CREATE POLICY "Users can view their own reports" ON reports
  FOR SELECT USING (auth.uid() = (SELECT auth_id FROM users WHERE id = reporter_id));

CREATE POLICY "Users can create reports" ON reports
  FOR INSERT WITH CHECK (auth.uid() = (SELECT auth_id FROM users WHERE id = reporter_id));

-- Create indexes for performance
CREATE INDEX idx_users_auth_id ON users(auth_id);
CREATE INDEX idx_matchmaking_queue_gender ON matchmaking_queue(gender);
CREATE INDEX idx_matchmaking_queue_joined_at ON matchmaking_queue(joined_at);
CREATE INDEX idx_chat_sessions_user1 ON chat_sessions(user1_id);
CREATE INDEX idx_chat_sessions_user2 ON chat_sessions(user2_id);
CREATE INDEX idx_chat_sessions_status ON chat_sessions(status);
CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_reports_status ON reports(status);
