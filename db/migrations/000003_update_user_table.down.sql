-- Down migration for removing the 'total_comments' and 'total_likes' columns
ALTER TABLE "posts"
DROP COLUMN IF EXISTS total_comments,
DROP COLUMN IF EXISTS total_likes;
