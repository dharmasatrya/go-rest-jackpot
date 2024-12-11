-- Create the 'users' table
CREATE TABLE users (
  user_id SERIAL PRIMARY KEY,      -- Primary key
  name VARCHAR(255),               -- User's name
  email VARCHAR(255) UNIQUE,       -- Email, unique per user
  password VARCHAR(255),           -- Password (should be hashed)
  date_of_birth DATE,              -- User's birth date
  is_admin BOOLEAN DEFAULT FALSE,  -- Admin flag (default to FALSE)
  is_soft_banned BOOLEAN DEFAULT FALSE -- Soft ban flag
);

-- Create the 'wallets' table
CREATE TABLE wallets (
  wallet_id SERIAL PRIMARY KEY,    -- Primary key
  user_id INT,                     -- Foreign key to 'users'
  balance DECIMAL DEFAULT 0,       -- Wallet balance
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE -- On user delete, delete wallet
);

-- Create the 'topups' table
CREATE TABLE topups (
  topup_id SERIAL PRIMARY KEY,     -- Primary key
  amount DECIMAL NOT NULL          -- Topup amount
);

-- Create the 'deposits' table
CREATE TABLE deposits (
  topup_id INT,                    -- Foreign key to 'topups'
  user_id INT,                     -- Foreign key to 'users'
  PRIMARY KEY (topup_id, user_id), -- Composite primary key (topup_id, user_id)
  FOREIGN KEY (topup_id) REFERENCES topups(topup_id) ON DELETE CASCADE, -- On delete topup, delete deposit
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE -- On delete user, delete deposit
);

-- Create the 'games' table
CREATE TABLE games (
  game_id SERIAL PRIMARY KEY,      -- Primary key
  wallet_id INT,                   -- Foreign key to 'wallets'
  user_id INT,                     -- Foreign key to 'users'
  slot_outcome VARCHAR(255),       -- Slot outcome (e.g., "WIN", "LOSE")
  game_result VARCHAR(255),        -- Game result (e.g., "WIN", "LOSE")
  balance_shift DECIMAL,           -- Balance change after the game
  FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id) ON DELETE CASCADE, -- On wallet delete, delete game
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE -- On user delete, delete game
);

-- Create the 'modifiers' table
CREATE TABLE modifiers (
  modifier_id SERIAL PRIMARY KEY,  -- Primary key
  user_id INT,                     -- Foreign key to 'users'
  curse_applied BOOLEAN DEFAULT FALSE,  -- Curse applied flag
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE -- On user delete, delete modifier
);

-- Create the 'records' table
CREATE TABLE records (
  record_id SERIAL PRIMARY KEY,    -- Primary key
  user_id INT,                     -- Foreign key to 'users'
  games_played INT DEFAULT 0,      -- Total games played
  games_won INT DEFAULT 0,         -- Total games won
  games_lost INT DEFAULT 0,        -- Total games lost
  winrate DECIMAL(5,2),            -- Winrate as a percentage
  total_topup DECIMAL DEFAULT 0,   -- Total amount deposited by the user
  total_earning DECIMAL DEFAULT 0, -- Total earnings (win amount)
  total_lost DECIMAL DEFAULT 0,    -- Total lost (amount spent)
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE -- On user delete, delete record
);

-- Add any necessary indexes if needed (not strictly required, but could improve performance on certain queries)
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_wallets_user_id ON wallets (user_id);
CREATE INDEX idx_games_user_id ON games (user_id);
CREATE INDEX idx_deposits_user_id ON deposits (user_id);
CREATE INDEX idx_records_user_id ON records (user_id);


-- Insert 3 users
INSERT INTO users (name, email, password, date_of_birth, is_admin, is_soft_banned)
VALUES
('Alice', 'alice@example.com', 'password123', '1985-01-15', FALSE, FALSE),
('Bob', 'bob@example.com', 'password123', '1990-05-20', FALSE, FALSE),
('Charlie', 'charlie@example.com', 'password123', '1992-10-05', TRUE, FALSE);

-- Insert wallets for the 3 users
INSERT INTO wallets (user_id, balance)
VALUES
(1, 100000),  -- Alice's wallet with balance 500
(2, 200000),  -- Bob's wallet with balance 200
(3, 5000000); -- Charlie's wallet with balance 1000

-- Insert topups for the 3 users
INSERT INTO topups (amount)
VALUES
(100000),
(500000),
(1000000);

-- Insert deposits linking the users to topups
INSERT INTO deposits (topup_id, user_id)
VALUES
(1, 1),  -- Alice's deposit of 100
(2, 2),  -- Bob's deposit of 50
(3, 3);  -- Charlie's deposit of 300

-- Insert games for the 3 users (Example outcomes: "WIN", "LOSE")
INSERT INTO games (wallet_id, user_id, slot_outcome, game_result, balance_shift)
VALUES
(1, 1, 'WIN', 'WIN', 100.00),  -- Alice played and won, balance increased by 100
(2, 2, 'LOSE', 'LOSE', -50.00), -- Bob played and lost, balance decreased by 50
(3, 3, 'WIN', 'WIN', 200.00);  -- Charlie played and won, balance increased by 200

-- Insert modifiers for the 3 users (Example curse and blessing)
INSERT INTO modifiers (user_id, curse_applied)
VALUES
(1, FALSE),  -- Alice has no curse
(2, TRUE),   -- Bob has a curse
(3, FALSE);  -- Charlie has no curse

-- Insert records for the 3 users (Example games played, won, lost, etc.)
INSERT INTO records (user_id, games_played, games_won, games_lost, winrate, total_topup, total_earning, total_lost)
VALUES
(1, 5, 3, 2, 60.00, 200.00, 300.00, 100.00),  -- Alice has played 5 games, won 3, lost 2, winrate 60%
(2, 3, 1, 2, 33.33, 100.00, 150.00, 50.00),   -- Bob has played 3 games, won 1, lost 2, winrate 33.33%
(3, 7, 5, 2, 71.43, 500.00, 700.00, 100.00);  -- Charlie has played 7 games, won 5, lost 2, winrate 71.43%
