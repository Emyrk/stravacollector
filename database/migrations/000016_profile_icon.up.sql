BEGIN;

ALTER TABLE athletes
	ADD COLUMN profile_pic_link TEXT NOT NULL DEFAULT '',
	ADD COLUMN profile_pic_link_medium TEXT NOT NULL DEFAULT '';

COMMIT;