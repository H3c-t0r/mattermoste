-- MySQL dump 10.13  Distrib 8.0.23, for Linux (x86_64)
--
-- Host: localhost    Database: dbwidnrtyyj7nhxnj5nkq5s7te7c
-- ------------------------------------------------------
-- Server version	8.0.23

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `Systems`
--

LOCK TABLES `Systems` WRITE;
/*!40000 ALTER TABLE `Systems` DISABLE KEYS */;
INSERT INTO `Systems` VALUES ('about_subsection_permissions','true');
INSERT INTO `Systems` VALUES ('add_billing_permissions','true');
INSERT INTO `Systems` VALUES ('add_bot_permissions','true');
INSERT INTO `Systems` VALUES ('add_convert_channel_permissions','true');
INSERT INTO `Systems` VALUES ('add_manage_guests_permissions','true');
INSERT INTO `Systems` VALUES ('add_system_console_permissions','true');
INSERT INTO `Systems` VALUES ('add_system_roles_permissions','true');
INSERT INTO `Systems` VALUES ('add_use_group_mentions_permission','true');
INSERT INTO `Systems` VALUES ('AdvancedPermissionsMigrationComplete','true');
INSERT INTO `Systems` VALUES ('apply_channel_manage_delete_to_channel_user','true');
INSERT INTO `Systems` VALUES ('authentication_subsection_permissions','true');
INSERT INTO `Systems` VALUES ('channel_moderations_permissions','true');
INSERT INTO `Systems` VALUES ('compliance_subsection_permissions','true');
INSERT INTO `Systems` VALUES ('ContentExtractionConfigDefaultTrueMigrationComplete','true');
INSERT INTO `Systems` VALUES ('custom_groups_permissions','true');
INSERT INTO `Systems` VALUES ('CustomGroupAdminRoleCreationMigrationComplete','true');
INSERT INTO `Systems` VALUES ('download_compliance_export_results','true');
INSERT INTO `Systems` VALUES ('emoji_permissions_split','true');
INSERT INTO `Systems` VALUES ('EmojisPermissionsMigrationComplete','true');
INSERT INTO `Systems` VALUES ('environment_subsection_permissions','true');
INSERT INTO `Systems` VALUES ('experimental_subsection_permissions','true');
INSERT INTO `Systems` VALUES ('GuestRolesCreationMigrationComplete','true');
INSERT INTO `Systems` VALUES ('integrations_subsection_permissions','true');
INSERT INTO `Systems` VALUES ('list_join_public_private_teams','true');
INSERT INTO `Systems` VALUES ('manage_secure_connections_permissions','true');
INSERT INTO `Systems` VALUES ('manage_shared_channel_permissions','true');
INSERT INTO `Systems` VALUES ('PlaybookRolesCreationMigrationComplete','true');
INSERT INTO `Systems` VALUES ('playbooks_manage_roles','true');
INSERT INTO `Systems` VALUES ('playbooks_permissions','true');
INSERT INTO `Systems` VALUES ('remove_channel_manage_delete_from_team_user','true');
INSERT INTO `Systems` VALUES ('remove_permanent_delete_user','true');
INSERT INTO `Systems` VALUES ('reporting_subsection_permissions','true');
INSERT INTO `Systems` VALUES ('site_subsection_permissions','true');
INSERT INTO `Systems` VALUES ('SystemConsoleRolesCreationMigrationComplete','true');
INSERT INTO `Systems` VALUES ('test_email_ancillary_permission','true');
INSERT INTO `Systems` VALUES ('Version','5.31.0');
INSERT INTO `Systems` VALUES ('view_members_new_permission','true');
INSERT INTO `Systems` VALUES ('webhook_permissions_split','true');
/*!40000 ALTER TABLE `Systems` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping data for table `Roles`
--

LOCK TABLES `Roles` WRITE;
/*!40000 ALTER TABLE `Roles` DISABLE KEYS */;
INSERT INTO `Roles` VALUES ('3ndsqn4sbbyjxpzccrzmzejstw','team_guest','authentication.roles.team_guest.name','authentication.roles.team_guest.description',1605167829008,1662271986894,0,' view_team',1,1);
INSERT INTO `Roles` VALUES ('44bq9f9s93b7f811ex5r1b4s1w','system_custom_group_admin','authentication.roles.system_custom_group_admin.name','authentication.roles.system_custom_group_admin.description',1662271985879,1662271986897,0,' manage_custom_group_members create_custom_group edit_custom_group delete_custom_group',0,1);
INSERT INTO `Roles` VALUES ('6jaz4y4nmjnxunkmogjf95fiha','system_user_manager','authentication.roles.system_user_manager.name','authentication.roles.system_user_manager.description',0,1662271986902,0,' delete_public_channel sysconsole_write_user_management_channels list_public_teams sysconsole_read_authentication_ldap sysconsole_write_user_management_groups manage_private_channel_members sysconsole_read_user_management_permissions sysconsole_read_authentication_password read_channel join_private_teams manage_team sysconsole_read_user_management_teams sysconsole_read_authentication_email sysconsole_write_user_management_teams read_public_channel_groups list_private_teams convert_private_channel_to_public manage_team_roles convert_public_channel_to_private sysconsole_read_authentication_openid view_team add_user_to_team read_ldap_sync_job read_public_channel test_ldap manage_private_channel_properties delete_private_channel manage_channel_roles sysconsole_read_authentication_guest_access sysconsole_read_user_management_channels sysconsole_read_user_management_groups manage_public_channel_members sysconsole_read_authentication_saml remove_user_from_team join_public_teams manage_public_channel_properties sysconsole_read_authentication_mfa sysconsole_read_authentication_signup read_private_channel_groups',0,1);
INSERT INTO `Roles` VALUES ('6pahsh5hg7rpjfhz4f5c1wsbfw','team_admin','authentication.roles.team_admin.name','authentication.roles.team_admin.description',0,1662271986906,0,' remove_user_from_team manage_slash_commands manage_team_roles delete_others_posts manage_others_slash_commands import_team manage_others_outgoing_webhooks convert_private_channel_to_public delete_post playbook_private_manage_roles convert_public_channel_to_private manage_channel_roles playbook_public_manage_roles manage_outgoing_webhooks manage_others_incoming_webhooks manage_incoming_webhooks manage_team',1,1);
INSERT INTO `Roles` VALUES ('c7oo8yeiojfu8xjyuyxn3fhxpc','team_post_all','authentication.roles.team_post_all.name','authentication.roles.team_post_all.description',0,1662271986910,0,' create_post use_channel_mentions',0,1);
INSERT INTO `Roles` VALUES ('cmqctq1egt877y9ua9pdsknoiw','team_post_all_public','authentication.roles.team_post_all_public.name','authentication.roles.team_post_all_public.description',0,1662271986914,0,' create_post_public use_channel_mentions',0,1);
INSERT INTO `Roles` VALUES ('dmmkxxmi3b8pdgcq9pjtf6mfao','playbook_member','authentication.roles.playbook_member.name','authentication.roles.playbook_member.description',1662271985811,1662271986918,0,' playbook_private_manage_members playbook_private_manage_properties run_create playbook_public_view playbook_public_manage_members playbook_public_manage_properties playbook_private_view',1,1);
INSERT INTO `Roles` VALUES ('dwjkqj9bj7r4xr8zb18z3tg53y','playbook_admin','authentication.roles.playbook_admin.name','authentication.roles.playbook_admin.description',1662271985841,1662271986921,0,' playbook_private_manage_roles playbook_private_manage_properties playbook_public_make_private playbook_public_manage_members playbook_public_manage_roles playbook_public_manage_properties playbook_private_manage_members',1,1);
INSERT INTO `Roles` VALUES ('hh56iy3patffuc3h76soondcga','channel_admin','authentication.roles.channel_admin.name','authentication.roles.channel_admin.description',0,1662271986924,0,' manage_channel_roles use_group_mentions',1,1);
INSERT INTO `Roles` VALUES ('hkcrew7wttb5fbuw3ime6g7nzc','system_read_only_admin','authentication.roles.system_read_only_admin.name','authentication.roles.system_read_only_admin.description',0,1662271986928,0,' sysconsole_read_environment_database sysconsole_read_experimental_features sysconsole_read_compliance_compliance_export sysconsole_read_environment_performance_monitoring sysconsole_read_environment_file_storage sysconsole_read_user_management_channels read_public_channel_groups read_elasticsearch_post_aggregation_job sysconsole_read_integrations_integration_management sysconsole_read_environment_push_notification_server read_compliance_export_job sysconsole_read_user_management_teams sysconsole_read_environment_logging sysconsole_read_about_edition_and_license sysconsole_read_site_customization sysconsole_read_reporting_site_statistics sysconsole_read_site_emoji sysconsole_read_authentication_guest_access test_ldap read_audits sysconsole_read_site_posts download_compliance_export_result sysconsole_read_compliance_compliance_monitoring sysconsole_read_site_announcement_banner sysconsole_read_integrations_gif sysconsole_read_authentication_email sysconsole_read_site_file_sharing_and_downloads sysconsole_read_experimental_bleve sysconsole_read_compliance_data_retention_policy read_channel sysconsole_read_experimental_feature_flags sysconsole_read_environment_image_proxy view_team sysconsole_read_authentication_openid sysconsole_read_environment_web_server sysconsole_read_integrations_cors read_ldap_sync_job sysconsole_read_authentication_saml get_analytics read_private_channel_groups sysconsole_read_reporting_team_statistics sysconsole_read_compliance_custom_terms_of_service sysconsole_read_authentication_ldap sysconsole_read_environment_smtp read_other_users_teams sysconsole_read_user_management_permissions sysconsole_read_environment_session_lengths read_public_channel read_data_retention_job sysconsole_read_user_management_groups sysconsole_read_environment_high_availability sysconsole_read_site_public_links sysconsole_read_authentication_password sysconsole_read_environment_rate_limiting list_public_teams sysconsole_read_site_users_and_teams sysconsole_read_authentication_signup get_logs read_license_information sysconsole_read_site_notices list_private_teams read_elasticsearch_post_indexing_job sysconsole_read_site_notifications sysconsole_read_authentication_mfa sysconsole_read_integrations_bot_accounts sysconsole_read_reporting_server_logs sysconsole_read_site_localization sysconsole_read_environment_elasticsearch sysconsole_read_user_management_users sysconsole_read_plugins sysconsole_read_environment_developer',0,1);
INSERT INTO `Roles` VALUES ('iiwt9pt6wiyb9e1enixtxs5yme','run_admin','authentication.roles.run_admin.name','authentication.roles.run_admin.description',1662271985864,1662271986932,0,' run_manage_properties run_manage_members',1,1);
INSERT INTO `Roles` VALUES ('jg1f1xfh3bb73pua938orwg9ie','system_guest','authentication.roles.global_guest.name','authentication.roles.global_guest.description',1605167829015,1662271986937,0,' create_direct_channel create_group_channel',1,1);
INSERT INTO `Roles` VALUES ('k891n5tpd3n9peue79azejjocy','system_post_all_public','authentication.roles.system_post_all_public.name','authentication.roles.system_post_all_public.description',0,1662271986941,0,' use_channel_mentions create_post_public',0,1);
INSERT INTO `Roles` VALUES ('kb6r9i58x7dxdb3srfohd66sse','system_admin','authentication.roles.global_admin.name','authentication.roles.global_admin.description',0,1662271986948,0,' list_public_teams edit_brand manage_private_channel_properties sysconsole_read_user_management_teams playbook_public_create manage_others_bots invalidate_caches manage_shared_channels sysconsole_write_environment_logging manage_others_outgoing_webhooks sysconsole_read_reporting_team_statistics sysconsole_read_plugins list_team_channels use_group_mentions sysconsole_read_site_users_and_teams sysconsole_write_site_localization get_analytics sysconsole_read_experimental_bleve manage_team_roles sysconsole_read_site_localization use_slash_commands edit_post sysconsole_write_user_management_channels test_elasticsearch list_private_teams add_ldap_public_cert join_public_teams manage_slash_commands manage_others_incoming_webhooks manage_public_channel_members sysconsole_read_environment_elasticsearch sysconsole_write_site_customization delete_others_emojis run_manage_members create_emojis sysconsole_write_authentication_email sysconsole_write_compliance_compliance_export add_saml_private_cert create_bot sysconsole_write_environment_rate_limiting add_saml_public_cert edit_other_users sysconsole_write_integrations_integration_management read_user_access_token create_elasticsearch_post_indexing_job sysconsole_write_user_management_users assign_system_admin_role sysconsole_write_user_management_groups sysconsole_read_authentication_guest_access sysconsole_write_about_edition_and_license sysconsole_read_authentication_ldap sysconsole_read_experimental_feature_flags sysconsole_read_integrations_cors sysconsole_read_user_management_groups join_public_channels sysconsole_read_experimental_features test_ldap sysconsole_write_environment_elasticsearch sysconsole_write_reporting_server_logs sysconsole_read_environment_image_proxy sysconsole_read_site_announcement_banner sysconsole_read_reporting_site_statistics sysconsole_write_authentication_mfa sysconsole_read_authentication_openid purge_bleve_indexes playbook_public_manage_members delete_emojis sysconsole_write_environment_file_storage sysconsole_write_reporting_site_statistics playbook_private_manage_members import_team sysconsole_write_environment_web_server sysconsole_write_authentication_password read_public_channel_groups create_compliance_export_job sysconsole_read_authentication_password list_users_without_team sysconsole_read_authentication_mfa add_ldap_private_cert create_data_retention_job read_license_information sysconsole_write_authentication_signup sysconsole_read_environment_push_notification_server edit_others_posts download_compliance_export_result create_ldap_sync_job sysconsole_write_authentication_ldap sysconsole_write_plugins read_data_retention_job sysconsole_write_compliance_data_retention_policy sysconsole_read_site_public_links manage_bots manage_system sysconsole_write_compliance_custom_terms_of_service playbook_public_manage_roles playbook_public_manage_properties playbook_private_create sysconsole_write_experimental_bleve sysconsole_read_authentication_email promote_guest get_saml_cert_status add_user_to_team sysconsole_write_site_users_and_teams create_custom_group manage_private_channel_members read_jobs sysconsole_write_experimental_features read_other_users_teams sysconsole_write_reporting_team_statistics sysconsole_read_environment_file_storage create_post_bleve_indexes_job sysconsole_read_site_file_sharing_and_downloads playbook_private_make_public playbook_public_view create_user_access_token create_public_channel read_channel sysconsole_read_user_management_channels sysconsole_read_user_management_permissions read_public_channel sysconsole_read_compliance_custom_terms_of_service sysconsole_write_site_emoji sysconsole_read_integrations_gif sysconsole_read_site_customization sysconsole_write_integrations_cors invite_user create_direct_channel sysconsole_write_user_management_teams run_create manage_custom_group_members read_ldap_sync_job sysconsole_read_site_notifications playbook_private_manage_properties sysconsole_read_integrations_bot_accounts convert_public_channel_to_private invalidate_email_invite reload_config get_saml_metadata_from_idp manage_secure_connections delete_private_channel sysconsole_read_about_edition_and_license convert_private_channel_to_public sysconsole_read_environment_developer recycle_database_connections remove_saml_private_cert manage_oauth sysconsole_write_environment_database sysconsole_write_site_notifications sysconsole_write_authentication_guest_access sysconsole_write_compliance_compliance_monitoring sysconsole_write_environment_image_proxy create_post_public manage_jobs remove_user_from_team delete_others_posts create_post_ephemeral playbook_private_view create_elasticsearch_post_aggregation_job remove_reaction add_reaction sysconsole_write_environment_high_availability sysconsole_write_authentication_openid sysconsole_write_user_management_permissions add_saml_idp_cert sysconsole_read_site_posts view_members sysconsole_write_environment_smtp sysconsole_read_authentication_saml create_post use_channel_mentions create_team playbook_private_manage_roles get_public_link sysconsole_write_billing manage_system_wide_oauth sysconsole_read_environment_database sysconsole_write_environment_session_lengths run_manage_properties sysconsole_write_authentication_saml sysconsole_read_environment_web_server sysconsole_read_environment_rate_limiting manage_public_channel_properties create_group_channel sysconsole_read_compliance_data_retention_policy sysconsole_read_environment_high_availability manage_others_slash_commands sysconsole_read_compliance_compliance_export delete_custom_group sysconsole_read_user_management_system_roles purge_elasticsearch_indexes view_team sysconsole_read_environment_performance_monitoring manage_channel_roles playbook_public_make_private remove_saml_public_cert demote_to_guest sysconsole_write_environment_performance_monitoring read_audits sysconsole_write_site_announcement_banner upload_file revoke_user_access_token read_others_bots test_email read_elasticsearch_post_aggregation_job sysconsole_read_compliance_compliance_monitoring join_private_teams delete_post sysconsole_write_site_public_links manage_team edit_custom_group sysconsole_write_experimental_feature_flags sysconsole_write_user_management_system_roles remove_others_reactions manage_license_information sysconsole_read_authentication_signup read_compliance_export_job sysconsole_write_environment_developer remove_saml_idp_cert manage_incoming_webhooks sysconsole_read_site_emoji assign_bot sysconsole_write_integrations_gif sysconsole_read_user_management_users delete_public_channel manage_outgoing_webhooks sysconsole_write_site_posts remove_ldap_private_cert sysconsole_write_site_file_sharing_and_downloads sysconsole_read_integrations_integration_management sysconsole_read_environment_logging test_site_url sysconsole_read_environment_session_lengths read_elasticsearch_post_indexing_job sysconsole_read_billing sysconsole_read_site_notices sysconsole_read_reporting_server_logs sysconsole_write_integrations_bot_accounts sysconsole_write_site_notices create_private_channel read_private_channel_groups run_view read_bots manage_roles test_s3 sysconsole_write_environment_push_notification_server get_logs invite_guest remove_ldap_public_cert sysconsole_read_environment_smtp',1,1);
INSERT INTO `Roles` VALUES ('km7kijhdtjbajquwu36uqneyoc','system_post_all','authentication.roles.system_post_all.name','authentication.roles.system_post_all.description',0,1662271986953,0,' create_post use_channel_mentions',0,1);
INSERT INTO `Roles` VALUES ('no7s4436sjbzzqjpupg85mszty','custom_group_user','authentication.roles.custom_group_user.name','authentication.roles.custom_group_user.description',1662271985801,1662271986956,0,'',0,0);
INSERT INTO `Roles` VALUES ('qo7e17c1m3rezyjqx5iq9dpmxe','system_manager','authentication.roles.system_manager.name','authentication.roles.system_manager.description',0,1662271986960,0,' sysconsole_write_environment_image_proxy sysconsole_read_environment_developer read_ldap_sync_job sysconsole_read_reporting_team_statistics recycle_database_connections get_logs read_private_channel_groups test_elasticsearch sysconsole_read_environment_logging purge_elasticsearch_indexes sysconsole_write_site_posts sysconsole_read_environment_database sysconsole_read_environment_performance_monitoring manage_team sysconsole_read_authentication_password sysconsole_write_site_users_and_teams sysconsole_read_user_management_channels sysconsole_write_environment_rate_limiting sysconsole_write_site_notifications read_license_information edit_brand sysconsole_read_plugins sysconsole_read_environment_high_availability sysconsole_read_environment_file_storage sysconsole_read_environment_elasticsearch sysconsole_write_environment_web_server sysconsole_write_environment_smtp sysconsole_write_environment_performance_monitoring sysconsole_write_environment_session_lengths sysconsole_write_user_management_groups convert_private_channel_to_public manage_private_channel_properties sysconsole_read_site_posts list_private_teams sysconsole_read_authentication_ldap sysconsole_read_authentication_guest_access sysconsole_read_site_emoji sysconsole_write_integrations_integration_management convert_public_channel_to_private manage_private_channel_members read_elasticsearch_post_aggregation_job manage_team_roles sysconsole_write_site_file_sharing_and_downloads read_channel read_public_channel sysconsole_read_authentication_openid add_user_to_team sysconsole_write_environment_developer sysconsole_write_site_localization sysconsole_read_about_edition_and_license test_s3 reload_config sysconsole_write_environment_elasticsearch test_site_url sysconsole_write_site_announcement_banner get_analytics sysconsole_read_environment_push_notification_server sysconsole_read_authentication_signup test_email sysconsole_write_integrations_bot_accounts sysconsole_write_integrations_cors view_team sysconsole_write_integrations_gif sysconsole_read_site_notices sysconsole_read_environment_image_proxy sysconsole_read_integrations_cors sysconsole_write_environment_push_notification_server join_public_teams test_ldap create_elasticsearch_post_aggregation_job sysconsole_read_environment_session_lengths sysconsole_write_environment_file_storage manage_public_channel_members sysconsole_write_site_customization sysconsole_read_site_announcement_banner sysconsole_read_environment_smtp sysconsole_write_user_management_teams delete_public_channel sysconsole_write_environment_logging read_public_channel_groups sysconsole_read_site_users_and_teams sysconsole_read_reporting_site_statistics sysconsole_read_site_localization sysconsole_read_site_customization sysconsole_read_environment_rate_limiting sysconsole_read_environment_web_server sysconsole_write_user_management_permissions sysconsole_read_site_file_sharing_and_downloads sysconsole_write_site_public_links sysconsole_read_site_public_links sysconsole_read_authentication_email read_elasticsearch_post_indexing_job sysconsole_read_authentication_saml remove_user_from_team delete_private_channel sysconsole_write_user_management_channels sysconsole_read_reporting_server_logs sysconsole_read_integrations_bot_accounts sysconsole_read_user_management_teams list_public_teams create_elasticsearch_post_indexing_job sysconsole_write_site_emoji invalidate_caches sysconsole_read_integrations_integration_management sysconsole_write_environment_high_availability sysconsole_read_user_management_permissions join_private_teams manage_channel_roles sysconsole_write_site_notices manage_public_channel_properties sysconsole_write_environment_database sysconsole_read_site_notifications sysconsole_read_user_management_groups sysconsole_read_integrations_gif sysconsole_read_authentication_mfa',0,1);
INSERT INTO `Roles` VALUES ('rkr97ikkh7fixy86qsoo5rqm4c','system_user_access_token','authentication.roles.system_user_access_token.name','authentication.roles.system_user_access_token.description',0,1662271986965,0,' create_user_access_token read_user_access_token revoke_user_access_token',0,1);
INSERT INTO `Roles` VALUES ('rxzdk5irm7rcffcfej9e33kqeo','team_user','authentication.roles.team_user.name','authentication.roles.team_user.description',0,1662271986968,0,' invite_user view_team read_public_channel playbook_public_create add_user_to_team playbook_private_create create_private_channel list_team_channels create_public_channel join_public_channels',1,1);
INSERT INTO `Roles` VALUES ('x768jnyzw3rkfx7xb66ehcac6o','channel_user','authentication.roles.channel_user.name','authentication.roles.channel_user.description',0,1662271986972,0,' manage_public_channel_properties create_post manage_private_channel_properties delete_public_channel manage_private_channel_members get_public_link delete_post delete_private_channel upload_file edit_post remove_reaction use_channel_mentions add_reaction read_channel use_slash_commands manage_public_channel_members',1,1);
INSERT INTO `Roles` VALUES ('ynn8aynsn7n1trtbuq6p4cyzhe','channel_guest','authentication.roles.channel_guest.name','authentication.roles.channel_guest.description',1605167829001,1662271986975,0,' read_channel add_reaction remove_reaction upload_file edit_post create_post use_channel_mentions use_slash_commands',1,1);
INSERT INTO `Roles` VALUES ('yqyby79r9jggxg7a9dnenuawmo','run_member','authentication.roles.run_member.name','authentication.roles.run_member.description',1662271985813,1662271986979,0,' run_view',1,1);
INSERT INTO `Roles` VALUES ('zzehkfnp67bg5g1owh6eptdcxc','system_user','authentication.roles.global_user.name','authentication.roles.global_user.description',0,1662271986983,0,' create_emojis join_public_teams list_public_teams edit_custom_group delete_emojis create_team create_group_channel manage_custom_group_members view_members delete_custom_group create_custom_group create_direct_channel',1,1);
/*!40000 ALTER TABLE `Roles` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-09-04  6:13:45
