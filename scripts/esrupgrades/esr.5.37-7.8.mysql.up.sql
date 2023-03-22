/* ==> mysql/000041_create_upload_sessions.up.sql <== */
/* Release 5.37 was meant to contain the index idx_uploadsessions_type, but a bug prevented that.
   This part of the migration #41 adds such index */
/* ==> mysql/000075_alter_upload_sessions_index.up.sql <== */
/* ==> mysql/000090_create_enums.up.sql <== */
DELIMITER //
CREATE PROCEDURE MigrateUploadSessions ()
BEGIN
	-- 'CREATE INDEX idx_uploadsessions_type ON UploadSessions(Type);'
	DECLARE CreateIndex BOOLEAN;
	DECLARE CreateIndexQuery TEXT DEFAULT NULL;

	-- 'DROP INDEX idx_uploadsessions_user_id ON UploadSessions; CREATE INDEX idx_uploadsessions_user_id ON UploadSessions(UserId);'
	DECLARE AlterIndex BOOLEAN;
	DECLARE AlterIndexQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE UploadSessions MODIFY COLUMN Type ENUM("attachment", "import");'
	DECLARE AlterColumn BOOLEAN;
	DECLARE AlterColumnQuery TEXT DEFAULT NULL;

	SELECT COUNT(*) = 0 FROM INFORMATION_SCHEMA.STATISTICS
		WHERE table_name = 'UploadSessions'
		AND table_schema = DATABASE()
		AND index_name = 'idx_uploadsessions_type'
		INTO CreateIndex;

	SELECT IFNULL(GROUP_CONCAT(column_name ORDER BY seq_in_index), '') = 'Type' FROM information_schema.statistics
		WHERE table_name = 'UploadSessions'
		AND table_schema = DATABASE()
		AND index_name = 'idx_uploadsessions_user_id'
		GROUP BY index_name
		INTO AlterIndex;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'UploadSessions'
		AND table_schema = DATABASE()
		AND column_name = 'Type'
		AND REPLACE(LOWER(column_type), '"', "'") != "enum('attachment','import')"
		INTO AlterColumn;

	IF CreateIndex THEN
		SET CreateIndexQuery = 'ADD INDEX idx_uploadsessions_type (Type)';
	END IF;

	IF AlterIndex THEN
		SET AlterIndexQuery = 'DROP INDEX idx_uploadsessions_user_id, ADD INDEX idx_uploadsessions_user_id (UserId)';
	END IF;

	IF AlterColumn THEN
		SET AlterColumnQuery = 'MODIFY COLUMN Type ENUM("attachment", "import")';
	END IF;

	IF CreateIndex OR AlterIndex OR AlterColumn THEN
		SET @query = CONCAT('ALTER TABLE UploadSessions ', CONCAT_WS(', ', CreateIndexQuery, AlterIndexQuery, AlterColumnQuery));
		PREPARE stmt FROM @query;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
	END IF;
END//
DELIMITER ;
CALL MigrateUploadSessions ();
DROP PROCEDURE IF EXISTS MigrateUploadSessions;

/* ==> mysql/000054_create_crt_channelmembership_count.up.sql <== */
/* fixCRTChannelMembershipCounts fixes the channel counts, i.e. the total message count,
total root message count, mention count, and mention count in root messages for users
who have viewed the channel after the last post in the channel */

DELIMITER //
CREATE PROCEDURE MigrateCRTChannelMembershipCounts ()
BEGIN
	IF(
		SELECT
			EXISTS (
			SELECT
				* FROM Systems
			WHERE
				Name = 'CRTChannelMembershipCountsMigrationComplete') = 0) THEN
		UPDATE
			ChannelMembers
			INNER JOIN Channels ON Channels.Id = ChannelMembers.ChannelId SET
				MentionCount = 0, MentionCountRoot = 0, MsgCount = Channels.TotalMsgCount, MsgCountRoot = Channels.TotalMsgCountRoot, LastUpdateAt = (
				SELECT
					(SELECT ROUND(UNIX_TIMESTAMP(NOW(3))*1000)))
	WHERE
		ChannelMembers.LastViewedAt >= Channels.LastPostAt;
		INSERT INTO Systems
			VALUES('CRTChannelMembershipCountsMigrationComplete', 'true');
	END IF;
END//
DELIMITER ;
CALL MigrateCRTChannelMembershipCounts ();
DROP PROCEDURE IF EXISTS MigrateCRTChannelMembershipCounts;

/* ==> mysql/000055_create_crt_thread_count_and_unreads.up.sql <== */
/* fixCRTThreadCountsAndUnreads Marks threads as read for users where the last
reply time of the thread is earlier than the time the user viewed the channel.
Marking a thread means setting the mention count to zero and setting the
last viewed at time of the the thread as the last viewed at time
of the channel */

DELIMITER //
CREATE PROCEDURE MigrateCRTThreadCountsAndUnreads ()
BEGIN
	IF(SELECT EXISTS(SELECT * FROM Systems WHERE Name = 'CRTThreadCountsAndUnreadsMigrationComplete') = 0) THEN
		UPDATE
			ThreadMemberships
			INNER JOIN (
				SELECT
					PostId,
					UserId,
					ChannelMembers.LastViewedAt AS CM_LastViewedAt,
					Threads.LastReplyAt
				FROM
					Threads
					INNER JOIN ChannelMembers ON ChannelMembers.ChannelId = Threads.ChannelId
				WHERE
					Threads.LastReplyAt <= ChannelMembers.LastViewedAt) AS q ON ThreadMemberships.Postid = q.PostId
				AND ThreadMemberships.UserId = q.UserId SET LastViewed = q.CM_LastViewedAt + 1, UnreadMentions = 0, LastUpdated = (
				SELECT
					(SELECT ROUND(UNIX_TIMESTAMP(NOW(3))*1000)));
		INSERT INTO Systems
			VALUES('CRTThreadCountsAndUnreadsMigrationComplete', 'true');
	END IF;
END//
DELIMITER ;
CALL MigrateCRTThreadCountsAndUnreads ();
DROP PROCEDURE IF EXISTS MigrateCRTThreadCountsAndUnreads;

/* ==> mysql/000056_upgrade_channels_v6.0.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Channels'
        AND table_schema = DATABASE()
        AND index_name = 'idx_channels_team_id_display_name'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_channels_team_id_display_name ON Channels(TeamId, DisplayName);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Channels'
        AND table_schema = DATABASE()
        AND index_name = 'idx_channels_team_id_type'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_channels_team_id_type ON Channels(TeamId, Type);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Channels'
        AND table_schema = DATABASE()
        AND index_name = 'idx_channels_team_id'
    ) > 0,
    'DROP INDEX idx_channels_team_id ON Channels;',
    'SELECT 1'
));

PREPARE removeIndexIfExists FROM @preparedStatement;
EXECUTE removeIndexIfExists;
DEALLOCATE PREPARE removeIndexIfExists;

/* ==> mysql/000057_upgrade_command_webhooks_v6.0.up.sql <== */

DELIMITER //
CREATE PROCEDURE MigrateRootId_CommandWebhooks () BEGIN DECLARE ParentId_EXIST INT;
SELECT COUNT(*)
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_NAME = 'CommandWebhooks'
  AND table_schema = DATABASE()
  AND COLUMN_NAME = 'ParentId' INTO ParentId_EXIST;
IF(ParentId_EXIST > 0) THEN
    UPDATE CommandWebhooks SET RootId = ParentId WHERE RootId = '' AND RootId != ParentId;
END IF;
END//
DELIMITER ;
CALL MigrateRootId_CommandWebhooks ();
DROP PROCEDURE IF EXISTS MigrateRootId_CommandWebhooks;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'CommandWebhooks'
        AND table_schema = DATABASE()
        AND column_name = 'ParentId'
    ) > 0,
    'ALTER TABLE CommandWebhooks DROP COLUMN ParentId;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000058_upgrade_channelmembers_v6.0.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'ChannelMembers'
        AND table_schema = DATABASE()
        AND column_name = 'NotifyProps'
        AND column_type != 'JSON'
    ) > 0,
    'ALTER TABLE ChannelMembers MODIFY COLUMN NotifyProps JSON;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'ChannelMembers'
        AND table_schema = DATABASE()
        AND index_name = 'idx_channelmembers_user_id'
    ) > 0,
    'DROP INDEX idx_channelmembers_user_id ON ChannelMembers;',
    'SELECT 1'
));

PREPARE removeIndexIfExists FROM @preparedStatement;
EXECUTE removeIndexIfExists;
DEALLOCATE PREPARE removeIndexIfExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'ChannelMembers'
        AND table_schema = DATABASE()
        AND index_name = 'idx_channelmembers_user_id_channel_id_last_viewed_at'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_channelmembers_user_id_channel_id_last_viewed_at ON ChannelMembers(UserId, ChannelId, LastViewedAt);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'ChannelMembers'
        AND table_schema = DATABASE()
        AND index_name = 'idx_channelmembers_channel_id_scheme_guest_user_id'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_channelmembers_channel_id_scheme_guest_user_id ON ChannelMembers(ChannelId, SchemeGuest, UserId);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

/* ==> mysql/000059_upgrade_users_v6.0.up.sql <== */
/* ==> mysql/000074_upgrade_users_v6.3.up.sql <== */
/* ==> mysql/000077_upgrade_users_v6.5.up.sql <== */
/* ==> mysql/000088_remaining_migrations.up.sql <== */
DELIMITER //
CREATE PROCEDURE MigrateUsers ()
BEGIN
	-- 'ALTER TABLE Users MODIFY COLUMN Props JSON;',
	DECLARE ChangeProps BOOLEAN;
	DECLARE ChangePropsQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Users MODIFY COLUMN NotifyProps JSON;',
	DECLARE ChangeNotifyProps BOOLEAN;
	DECLARE ChangeNotifyPropsQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Users ALTER Timezone DROP DEFAULT;',
	DECLARE DropTimezoneDefault BOOLEAN;
	DECLARE DropTimezoneDefaultQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Users MODIFY COLUMN Timezone JSON;',
	DECLARE ChangeTimezone BOOLEAN;
	DECLARE ChangeTimezoneQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Users MODIFY COLUMN Roles text;',
	DECLARE ChangeRoles BOOLEAN;
	DECLARE ChangeRolesQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Users DROP COLUMN AcceptedTermsOfServiceId;',
	DECLARE DropTermsOfService BOOLEAN;
	DECLARE DropTermsOfServiceQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Users DROP COLUMN AcceptedServiceTermsId;',
	DECLARE DropServiceTerms BOOLEAN;
	DECLARE DropServiceTermsQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Users DROP COLUMN ThemeProps',
	DECLARE DropThemeProps BOOLEAN;
	DECLARE DropThemePropsQuery TEXT DEFAULT NULL;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Users'
		AND table_schema = DATABASE()
		AND column_name = 'Props'
		AND LOWER(column_type) != 'json'
		INTO ChangeProps;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Users'
		AND table_schema = DATABASE()
		AND column_name = 'NotifyProps'
		AND LOWER(column_type) != 'json'
		INTO ChangeNotifyProps;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Users'
		AND column_default IS NOT NULL
		INTO DropTimezoneDefault;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Users'
		AND table_schema = DATABASE()
		AND column_name = 'Timezone'
		AND LOWER(column_type) != 'json'
		INTO ChangeTimezone;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Users'
		AND table_schema = DATABASE()
		AND column_name = 'Roles'
		AND LOWER(column_type) != 'text'
		INTO ChangeRoles;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Users'
		AND table_schema = DATABASE()
		AND column_name = 'AcceptedTermsOfServiceId'
		INTO DropTermsOfService;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Users'
		AND table_schema = DATABASE()
		AND column_name = 'AcceptedServiceTermsId'
		INTO DropServiceTerms;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Users'
		AND table_schema = DATABASE()
		AND column_name = 'ThemeProps'
		INTO DropThemeProps;

	IF ChangeProps THEN
		SET ChangePropsQuery = 'MODIFY COLUMN Props JSON';
	END IF;

	IF ChangeNotifyProps THEN
		SET ChangeNotifyPropsQuery = 'MODIFY COLUMN NotifyProps JSON';
	END IF;

	IF DropTimezoneDefault THEN
		SET DropTimezoneDefaultQuery = 'ALTER Timezone DROP DEFAULT';
	END IF;

	IF ChangeTimezone THEN
		SET ChangeTimezoneQuery = 'MODIFY COLUMN Timezone JSON';
	END IF;

	IF ChangeRoles THEN
		SET ChangeRolesQuery = 'MODIFY COLUMN Roles text';
	END IF;

	IF DropTermsOfService THEN
		SET DropTermsOfServiceQuery = 'DROP COLUMN AcceptedTermsOfServiceId';
	END IF;

	IF DropServiceTerms THEN
		SET DropServiceTermsQuery = 'DROP COLUMN AcceptedServiceTermsId';
	END IF;

	IF DropThemeProps THEN
		INSERT INTO Preferences(UserId, Category, Name, Value) SELECT Id, '', '', ThemeProps FROM Users WHERE Users.ThemeProps != 'null';
		SET DropThemePropsQuery = 'DROP COLUMN ThemeProps';
	END IF;

	IF ChangeProps OR ChangeNotifyProps OR DropTimezoneDefault OR ChangeTimezone OR ChangeRoles OR DropTermsOfService OR DropServiceTerms OR DropThemeProps THEN
		SET @query = CONCAT('ALTER TABLE Users ', CONCAT_WS(', ', ChangePropsQuery, ChangeNotifyPropsQuery, DropTimezoneDefaultQuery, ChangeTimezoneQuery, ChangeRolesQuery, DropTermsOfServiceQuery, DropServiceTermsQuery, DropThemePropsQuery));
		PREPARE stmt FROM @query;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
	END IF;
END//
DELIMITER ;
CALL MigrateUsers ();
DROP PROCEDURE IF EXISTS MigrateUsers;

/* ==> mysql/000060_upgrade_jobs_v6.0.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Jobs'
        AND table_schema = DATABASE()
        AND column_name = 'Data'
        AND column_type != 'JSON'
    ) > 0,
    'ALTER TABLE Jobs MODIFY COLUMN Data JSON;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;


/* ==> mysql/000061_upgrade_link_metadata_v6.0.up.sql <== */

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'LinkMetadata'
        AND table_schema = DATABASE()
        AND column_name = 'Data'
        AND column_type != 'JSON'
    ) > 0,
    'ALTER TABLE LinkMetadata MODIFY COLUMN Data JSON;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000062_upgrade_sessions_v6.0.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Sessions'
        AND table_schema = DATABASE()
        AND column_name = 'Props'
        AND column_type != 'JSON'
    ) > 0,
    'ALTER TABLE Sessions MODIFY COLUMN Props JSON;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;


/* ==> mysql/000063_upgrade_threads_v6.0.up.sql <== */
/* ==> mysql/000083_threads_threaddeleteat.up.sql <== */
/* ==> mysql/000096_threads_threadteamid.up.sql <== */
DELIMITER //
CREATE PROCEDURE MigrateThreads ()
BEGIN
	-- 'ALTER TABLE Threads MODIFY COLUMN Participants JSON;'
	DECLARE ChangeParticipants BOOLEAN;
	DECLARE ChangeParticipantsQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Threads DROP COLUMN DeleteAt;'
	DECLARE DropDeleteAt BOOLEAN;
	DECLARE DropDeleteAtQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Threads ADD COLUMN ThreadDeleteAt bigint(20);'
	DECLARE CreateThreadDeleteAt BOOLEAN;
	DECLARE CreateThreadDeleteAtQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Threads DROP COLUMN TeamId;'
	DECLARE DropTeamId BOOLEAN;
	DECLARE DropTeamIdQuery TEXT DEFAULT NULL;

	-- 'ALTER TABLE Threads ADD COLUMN ThreadTeamId varchar(26) DEFAULT NULL;'
	DECLARE CreateThreadTeamId BOOLEAN;
	DECLARE CreateThreadTeamIdQuery TEXT DEFAULT NULL;

	-- CREATE INDEX idx_threads_channel_id_last_reply_at ON Threads(ChannelId, LastReplyAt);
	DECLARE CreateIndex BOOLEAN;
	DECLARE CreateIndexQuery TEXT DEFAULT NULL;

	-- DROP INDEX idx_threads_channel_id ON Threads;
	DECLARE DropIndex BOOLEAN;
	DECLARE DropIndexQuery TEXT DEFAULT NULL;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Threads'
		AND table_schema = DATABASE()
		AND column_name = 'Participants'
		AND LOWER(column_type) != 'json'
		INTO ChangeParticipants;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
		WHERE table_name = 'Threads'
		AND table_schema = DATABASE()
		AND column_name = 'DeleteAt'
		INTO DropDeleteAt;

	SELECT COUNT(*) = 0 FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Threads'
		AND table_schema = DATABASE()
		AND column_name = 'ThreadDeleteAt'
		INTO CreateThreadDeleteAt;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
		WHERE table_name = 'Threads'
		AND table_schema = DATABASE()
		AND column_name = 'TeamId'
		INTO DropTeamId;

	SELECT COUNT(*) = 0 FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Threads'
		AND table_schema = DATABASE()
		AND column_name = 'ThreadTeamId'
		INTO CreateThreadTeamId;

	SELECT COUNT(*) = 0 FROM INFORMATION_SCHEMA.STATISTICS
		WHERE table_name = 'Threads'
		AND table_schema = DATABASE()
		AND index_name = 'idx_threads_channel_id_last_reply_at'
		INTO CreateIndex;

	SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
		WHERE table_name = 'Threads'
		AND table_schema = DATABASE()
		AND index_name = 'idx_threads_channel_id'
		INTO DropIndex;

	IF ChangeParticipants THEN
		SET ChangeParticipantsQuery = 'MODIFY COLUMN Participants JSON';
	END IF;

	IF DropDeleteAt THEN
		SET DropDeleteAtQuery = 'DROP COLUMN DeleteAt';
	END IF;

	IF CreateThreadDeleteAt THEN
		SET CreateThreadDeleteAtQuery = 'ADD COLUMN ThreadDeleteAt bigint(20)';
	END IF;

	IF DropTeamId THEN
		SET DropTeamIdQuery = 'DROP COLUMN TeamId';
	END IF;

	IF CreateThreadTeamId THEN
		SET CreateThreadTeamIdQuery = 'ADD COLUMN ThreadTeamId varchar(26) DEFAULT NULL';
	END IF;

	IF CreateIndex THEN
		SET CreateIndexQuery = 'ADD INDEX idx_threads_channel_id_last_reply_at (ChannelId, LastReplyAt)';
	END IF;

	IF DropIndex THEN
		SET DropIndexQuery = 'DROP INDEX idx_threads_channel_id';
	END IF;

	IF ChangeParticipants OR DropDeleteAt OR CreateThreadDeleteAt OR DropTeamId OR CreateThreadTeamId OR CreateIndex OR DropIndex THEN
		SET @query = CONCAT('ALTER TABLE Threads ', CONCAT_WS(', ', ChangeParticipantsQuery, DropDeleteAtQuery, CreateThreadDeleteAtQuery, DropTeamIdQuery, CreateThreadTeamIdQuery, CreateIndexQuery, DropIndexQuery));
		PREPARE stmt FROM @query;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
	END IF;

	UPDATE Threads, Posts
		SET Threads.ThreadDeleteAt = Posts.DeleteAt
		WHERE Posts.Id = Threads.PostId
		AND Threads.ThreadDeleteAt IS NULL;

	UPDATE Threads, Channels
		SET Threads.ThreadTeamId = Channels.TeamId
		WHERE Channels.Id = Threads.ChannelId
		AND Threads.ThreadTeamId IS NULL;
END//
DELIMITER ;
CALL MigrateThreads ();
DROP PROCEDURE IF EXISTS MigrateThreads;

/* ==> mysql/000064_upgrade_status_v6.0.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Status'
        AND table_schema = DATABASE()
        AND index_name = 'idx_status_status_dndendtime'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_status_status_dndendtime ON Status(Status, DNDEndTime);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Status'
        AND table_schema = DATABASE()
        AND index_name = 'idx_status_status'
    ) > 0,
    'DROP INDEX idx_status_status ON Status;',
    'SELECT 1'
));

PREPARE removeIndexIfExists FROM @preparedStatement;
EXECUTE removeIndexIfExists;
DEALLOCATE PREPARE removeIndexIfExists;

/* ==> mysql/000065_upgrade_groupchannels_v6.0.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'GroupChannels'
        AND table_schema = DATABASE()
        AND index_name = 'idx_groupchannels_schemeadmin'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_groupchannels_schemeadmin ON GroupChannels(SchemeAdmin);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

/* ==> mysql/000066_upgrade_posts_v6.0.up.sql <== */
DELIMITER //
CREATE PROCEDURE MigrateRootId_Posts ()
BEGIN
DECLARE ParentId_EXIST INT;
DECLARE Alter_FileIds INT;
DECLARE Alter_Props INT;
SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_NAME = 'Posts'
  AND table_schema = DATABASE()
  AND COLUMN_NAME = 'ParentId' INTO ParentId_EXIST;
SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE table_name = 'Posts'
  AND table_schema = DATABASE()
  AND column_name = 'FileIds'
  AND column_type != 'text' INTO Alter_FileIds;
SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
  WHERE table_name = 'Posts'
  AND table_schema = DATABASE()
  AND column_name = 'Props'
  AND column_type != 'JSON' INTO Alter_Props;
IF (Alter_Props OR Alter_FileIds) THEN
	IF(ParentId_EXIST > 0) THEN
		UPDATE Posts SET RootId = ParentId WHERE RootId = '' AND RootId != ParentId;
		ALTER TABLE Posts MODIFY COLUMN FileIds text, MODIFY COLUMN Props JSON, DROP COLUMN ParentId;
	ELSE
		ALTER TABLE Posts MODIFY COLUMN FileIds text, MODIFY COLUMN Props JSON;
	END IF;
END IF;
END//
DELIMITER ;
CALL MigrateRootId_Posts ();
DROP PROCEDURE IF EXISTS MigrateRootId_Posts;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Posts'
        AND table_schema = DATABASE()
        AND index_name = 'idx_posts_root_id_delete_at'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_posts_root_id_delete_at ON Posts(RootId, DeleteAt);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Posts'
        AND table_schema = DATABASE()
        AND index_name = 'idx_posts_root_id'
    ) > 0,
    'DROP INDEX idx_posts_root_id ON Posts;',
    'SELECT 1'
));

PREPARE removeIndexIfExists FROM @preparedStatement;
EXECUTE removeIndexIfExists;
DEALLOCATE PREPARE removeIndexIfExists;

/* ==> mysql/000067_upgrade_channelmembers_v6.1.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'ChannelMembers'
        AND table_schema = DATABASE()
        AND column_name = 'Roles'
        AND column_type != 'text'
    ) > 0,
    'ALTER TABLE ChannelMembers MODIFY COLUMN Roles text;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000068_upgrade_teammembers_v6.1.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'TeamMembers'
        AND table_schema = DATABASE()
        AND column_name = 'Roles'
        AND column_type != 'text'
    ) > 0,
    'ALTER TABLE TeamMembers MODIFY COLUMN Roles text;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000069_upgrade_jobs_v6.1.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Jobs'
        AND table_schema = DATABASE()
        AND index_name = 'idx_jobs_status_type'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_jobs_status_type ON Jobs(Status, Type);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

/* ==> mysql/000070_upgrade_cte_v6.1.up.sql <== */
DELIMITER //
CREATE PROCEDURE Migrate_LastRootPostAt ()
BEGIN
DECLARE
	LastRootPostAt_EXIST INT;
	SELECT
		COUNT(*)
	FROM
		INFORMATION_SCHEMA.COLUMNS
	WHERE
		TABLE_NAME = 'Channels'
		AND table_schema = DATABASE()
		AND COLUMN_NAME = 'LastRootPostAt' INTO LastRootPostAt_EXIST;
	IF(LastRootPostAt_EXIST = 0) THEN
        ALTER TABLE Channels ADD COLUMN LastRootPostAt bigint DEFAULT 0;
		UPDATE
			Channels
			INNER JOIN (
				SELECT
					Channels.Id channelid,
					COALESCE(MAX(Posts.CreateAt), 0) AS lastrootpost
				FROM
					Channels
					LEFT JOIN Posts FORCE INDEX (idx_posts_channel_id_update_at) ON Channels.Id = Posts.ChannelId
				WHERE
					Posts.RootId = ''
				GROUP BY
					Channels.Id) AS q ON q.channelid = Channels.Id SET LastRootPostAt = lastrootpost;
	END IF;
END//
DELIMITER ;
CALL Migrate_LastRootPostAt ();
DROP PROCEDURE IF EXISTS Migrate_LastRootPostAt;

/* ==> mysql/000071_upgrade_sessions_v6.1.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Sessions'
        AND table_schema = DATABASE()
        AND column_name = 'Roles'
        AND column_type != 'text'
    ) > 0,
    'ALTER TABLE Sessions MODIFY COLUMN Roles text;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000072_upgrade_schemes_v6.3.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Schemes'
        AND table_schema = DATABASE()
        AND column_name = 'DefaultPlaybookAdminRole'
    ) > 0,
    'SELECT 1',
    'ALTER TABLE Schemes ADD COLUMN DefaultPlaybookAdminRole VARCHAR(64) DEFAULT "";'
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Schemes'
        AND table_schema = DATABASE()
        AND column_name = 'DefaultPlaybookMemberRole'
    ) > 0,
    'SELECT 1',
    'ALTER TABLE Schemes ADD COLUMN DefaultPlaybookMemberRole VARCHAR(64) DEFAULT "";'
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Schemes'
        AND table_schema = DATABASE()
        AND column_name = 'DefaultRunAdminRole'
    ) > 0,
    'SELECT 1',
    'ALTER TABLE Schemes ADD COLUMN DefaultRunAdminRole VARCHAR(64) DEFAULT "";'
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Schemes'
        AND table_schema = DATABASE()
        AND column_name = 'DefaultRunMemberRole'
    ) > 0,
    'SELECT 1',
    'ALTER TABLE Schemes ADD COLUMN DefaultRunMemberRole VARCHAR(64) DEFAULT "";'
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

/* ==> mysql/000073_upgrade_plugin_key_value_store_v6.3.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT Count(*) FROM Information_Schema.Columns
        WHERE table_name = 'PluginKeyValueStore'
        AND table_schema = DATABASE()
        AND column_name = 'PKey'
        AND column_type != 'varchar(150)'
    ) > 0,
    'ALTER TABLE PluginKeyValueStore MODIFY COLUMN PKey varchar(150);',
    'SELECT 1'
));

PREPARE alterTypeIfExists FROM @preparedStatement;
EXECUTE alterTypeIfExists;
DEALLOCATE PREPARE alterTypeIfExists;

/* ==> mysql/000076_upgrade_lastrootpostat.up.sql <== */
DELIMITER //
CREATE PROCEDURE Migrate_LastRootPostAt_Default ()
BEGIN
	IF (
			SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_NAME = 'Channels'
			AND TABLE_SCHEMA = DATABASE()
			AND COLUMN_NAME = 'LastRootPostAt'
			AND (COLUMN_DEFAULT IS NULL OR COLUMN_DEFAULT != 0)
		) = 1 THEN
		ALTER TABLE Channels ALTER COLUMN LastRootPostAt SET DEFAULT 0;
	END IF;
END//
DELIMITER ;
CALL Migrate_LastRootPostAt_Default ();
DROP PROCEDURE IF EXISTS Migrate_LastRootPostAt_Default;

DELIMITER //
CREATE PROCEDURE Migrate_LastRootPostAt_Fix ()
BEGIN
	IF (
		SELECT COUNT(*)
		FROM Channels
		WHERE LastRootPostAt IS NULL
	) > 0 THEN
	-- fixes migrate cte and sets the LastRootPostAt for channels that don't have it set
		UPDATE
			Channels
			INNER JOIN (
				SELECT
					Channels.Id channelid,
					COALESCE(MAX(Posts.CreateAt), 0) AS lastrootpost
				FROM
					Channels
					LEFT JOIN Posts FORCE INDEX (idx_posts_channel_id_update_at) ON Channels.Id = Posts.ChannelId
				WHERE
					Posts.RootId = ''
				GROUP BY
					Channels.Id) AS q ON q.channelid = Channels.Id
				SET
					LastRootPostAt = lastrootpost
				WHERE
					LastRootPostAt IS NULL;

		-- sets LastRootPostAt to 0, for channels with no posts
		UPDATE Channels SET LastRootPostAt=0 WHERE LastRootPostAt IS NULL;
	END IF;
END//
DELIMITER ;
CALL Migrate_LastRootPostAt_Fix ();
DROP PROCEDURE IF EXISTS Migrate_LastRootPostAt_Fix;

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000078_create_oauth_mattermost_app_id.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'OAuthApps'
        AND table_schema = DATABASE()
        AND column_name = 'MattermostAppID'
    ) > 0,
	'SELECT 1',
    'ALTER TABLE OAuthApps ADD COLUMN MattermostAppID varchar(32);'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000079_usergroups_displayname_index.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'UserGroups'
        AND table_schema = DATABASE()
        AND index_name = 'idx_usergroups_displayname'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_usergroups_displayname ON UserGroups(DisplayName);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

/* ==> mysql/000080_posts_createat_id.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Posts'
        AND table_schema = DATABASE()
        AND index_name = 'idx_posts_create_at_id'
    ) > 0,
    'SELECT 1;',
    'CREATE INDEX idx_posts_create_at_id on Posts(CreateAt, Id) LOCK=NONE;'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

/* ==> mysql/000081_threads_deleteat.up.sql <== */
-- Replaced by 000083_threads_threaddeleteat.up.sql

/* ==> mysql/000082_upgrade_oauth_mattermost_app_id.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'OAuthApps'
        AND table_schema = DATABASE()
        AND column_name = 'MattermostAppID'
    ) > 0,
    'UPDATE OAuthApps SET MattermostAppID = "" WHERE MattermostAppID IS NULL;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'OAuthApps'
        AND table_schema = DATABASE()
        AND column_name = 'MattermostAppID'
    ) > 0,
    'ALTER TABLE OAuthApps MODIFY MattermostAppID varchar(32) NOT NULL DEFAULT "";',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000084_recent_searches.up.sql <== */
CREATE TABLE IF NOT EXISTS RecentSearches (
    UserId CHAR(26),
    SearchPointer int,
    Query json,
    CreateAt bigint NOT NULL,
    PRIMARY KEY (UserId, SearchPointer)
);
/* ==> mysql/000085_fileinfo_add_archived_column.up.sql <== */

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'FileInfo'
        AND table_schema = DATABASE()
        AND column_name = 'Archived'
    ) > 0,
    'SELECT 1',
    'ALTER TABLE FileInfo ADD COLUMN Archived boolean NOT NULL DEFAULT false;'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000086_add_cloud_limits_archived.up.sql <== */
SET @preparedStatement = (SELECT IF(
	NOT EXISTS(
		SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS
		WHERE table_name = 'Teams'
		AND table_schema = DATABASE()
		AND column_name = 'CloudLimitsArchived'
	),
	'ALTER TABLE Teams ADD COLUMN CloudLimitsArchived BOOLEAN NOT NULL DEFAULT FALSE;',
	'SELECT 1'
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

/* ==> mysql/000087_sidebar_categories_index.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'SidebarCategories'
        AND table_schema = DATABASE()
        AND index_name = 'idx_sidebarcategories_userid_teamid'
    ) > 0,
    'SELECT 1;',
    'CREATE INDEX idx_sidebarcategories_userid_teamid on SidebarCategories(UserId, TeamId) LOCK=NONE;'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

/* ==> mysql/000088_remaining_migrations.up.sql <== */
DROP TABLE IF EXISTS JobStatuses;

DROP TABLE IF EXISTS PasswordRecovery;

/* ==> mysql/000089_add-channelid-to-reaction.up.sql <== */
SET @preparedStatement = (SELECT IF(
    NOT EXISTS(
        SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Reactions'
        AND table_schema = DATABASE()
        AND column_name = 'ChannelId'
    ),
    'ALTER TABLE Reactions ADD COLUMN ChannelId varchar(26) NOT NULL DEFAULT "";',
    'SELECT 1;'
));

PREPARE addColumnIfNotExists FROM @preparedStatement;
EXECUTE addColumnIfNotExists;
DEALLOCATE PREPARE addColumnIfNotExists;


UPDATE Reactions SET ChannelId = COALESCE((select ChannelId from Posts where Posts.Id = Reactions.PostId), '') WHERE ChannelId="";


SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'Reactions'
        AND table_schema = DATABASE()
        AND index_name = 'idx_reactions_channel_id'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_reactions_channel_id ON Reactions(ChannelId);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

/* ==> mysql/000090_create_enums.up.sql <== */
SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Channels'
        AND table_schema = DATABASE()
        AND column_name = 'Type'
        AND column_type != 'ENUM("D", "O", "G", "P")'
    ) > 0,
    'ALTER TABLE Channels MODIFY COLUMN Type ENUM("D", "O", "G", "P");',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Teams'
        AND table_schema = DATABASE()
        AND column_name = 'Type'
        AND column_type != 'ENUM("I", "O")'
    ) > 0,
    'ALTER TABLE Teams MODIFY COLUMN Type ENUM("I", "O");',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000091_create_post_reminder.up.sql <== */
CREATE TABLE IF NOT EXISTS PostReminders (
    PostId varchar(26) NOT NULL,
    UserId varchar(26) NOT NULL,
    TargetTime bigint,
    PRIMARY KEY (PostId, UserId)
);

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'PostReminders'
        AND table_schema = DATABASE()
        AND index_name = 'idx_postreminders_targettime'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_postreminders_targettime ON PostReminders(TargetTime);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;
/* ==> mysql/000092_add_createat_to_teammembers.up.sql <== */
SET @preparedStatement = (SELECT IF(
    NOT EXISTS(
        SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'TeamMembers'
        AND table_schema = DATABASE()
        AND column_name = 'CreateAt'
    ),
    'ALTER TABLE TeamMembers ADD COLUMN CreateAt bigint DEFAULT 0;',
    'SELECT 1;'
));

PREPARE addColumnIfNotExists FROM @preparedStatement;
EXECUTE addColumnIfNotExists;
DEALLOCATE PREPARE addColumnIfNotExists;

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.STATISTICS
        WHERE table_name = 'TeamMembers'
        AND table_schema = DATABASE()
        AND index_name = 'idx_teammembers_create_at'
    ) > 0,
    'SELECT 1',
    'CREATE INDEX idx_teammembers_createat ON TeamMembers(CreateAt);'
));

PREPARE createIndexIfNotExists FROM @preparedStatement;
EXECUTE createIndexIfNotExists;
DEALLOCATE PREPARE createIndexIfNotExists;

/* ==> mysql/000093_notify_admin.up.sql <== */
CREATE TABLE IF NOT EXISTS NotifyAdmin (
    UserId varchar(26) NOT NULL,
    CreateAt bigint(20) DEFAULT NULL,
    RequiredPlan varchar(26) NOT NULL,
    RequiredFeature varchar(100) NOT NULL,
    Trial BOOLEAN NOT NULL,
    PRIMARY KEY (UserId, RequiredFeature, RequiredPlan)
);

/* ==> mysql/000094_threads_teamid.up.sql <== */
-- Replaced by 000096_threads_threadteamid.up.sql

/* ==> mysql/000095_remove_posts_parentid.up.sql <== */
-- While upgrading from 5.x to 6.x with manual queries, there is a chance that this
-- migration is skipped. In that case, we need to make sure that the column is dropped.

SET @preparedStatement = (SELECT IF(
    (
        SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'Posts'
        AND table_schema = DATABASE()
        AND column_name = 'ParentId'
    ) > 0,
    'ALTER TABLE Posts DROP COLUMN ParentId;',
    'SELECT 1'
));

PREPARE alterIfExists FROM @preparedStatement;
EXECUTE alterIfExists;
DEALLOCATE PREPARE alterIfExists;

/* ==> mysql/000097_create_posts_priority.up.sql <== */
CREATE TABLE IF NOT EXISTS PostsPriority (
    PostId varchar(26) NOT NULL,
    ChannelId varchar(26) NOT NULL,
    Priority varchar(32) NOT NULL,
    RequestedAck tinyint(1),
    PersistentNotifications tinyint(1),
    PRIMARY KEY (PostId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET @preparedStatement = (SELECT IF(
    NOT EXISTS(
        SELECT 1 FROM INFORMATION_SCHEMA.COLUMNS
        WHERE table_name = 'ChannelMembers'
        AND table_schema = DATABASE()
        AND column_name = 'UrgentMentionCount'
    ),
    'ALTER TABLE ChannelMembers ADD COLUMN UrgentMentionCount bigint(20);',
    'SELECT 1;'
));

PREPARE alterIfNotExists FROM @preparedStatement;
EXECUTE alterIfNotExists;
DEALLOCATE PREPARE alterIfNotExists;

/* ==> mysql/000098_create_post_acknowledgements.up.sql <== */
CREATE TABLE IF NOT EXISTS PostAcknowledgements (
    PostId varchar(26) NOT NULL,
    UserId varchar(26) NOT NULL,
    AcknowledgedAt bigint(20) DEFAULT NULL,
    PRIMARY KEY (PostId, UserId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/* ==> mysql/000099_create_drafts.up.sql <== */
/* ==> mysql/000100_add_draft_priority_column.up.sql <== */
CREATE TABLE IF NOT EXISTS Drafts (
    CreateAt bigint(20) DEFAULT NULL,
    UpdateAt bigint(20) DEFAULT NULL,
    DeleteAt bigint(20) DEFAULT NULL,
    UserId varchar(26) NOT NULL,
    ChannelId varchar(26) NOT NULL,
    RootId varchar(26) DEFAULT '',
    Message text,
    Props text,
    FileIds text,
	Priority text,
    PRIMARY KEY (UserId, ChannelId, RootId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

/* ==> mysql/000101_create_true_up_review_history.up.sql <== */
CREATE TABLE IF NOT EXISTS TrueUpReviewHistory (
	DueDate bigint(20),
	Completed boolean,
    PRIMARY KEY (DueDate)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
