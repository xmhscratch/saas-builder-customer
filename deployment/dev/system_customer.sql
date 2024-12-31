-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: mariadb_master:3306
-- Generation Time: Dec 31, 2024 at 10:02 AM
-- Server version: 10.8.8-MariaDB-log
-- PHP Version: 8.2.8

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `claimh_customer_1ef19f81-a3a1-45a2-9203-b792abcddc52`
--
CREATE DATABASE IF NOT EXISTS `claimh_customer_1ef19f81-a3a1-45a2-9203-b792abcddc52` DEFAULT CHARACTER SET ascii COLLATE ascii_general_ci;
USE `claimh_customer_1ef19f81-a3a1-45a2-9203-b792abcddc52`;

DELIMITER $$
--
-- Procedures
--
DROP PROCEDURE IF EXISTS `GetIndexableCustomerAttributes`$$
CREATE DEFINER=`root`@`%` PROCEDURE `GetIndexableCustomerAttributes` (IN `customerId` CHAR(27))   BEGIN
    SET @sql = NULL;
    SET @customerId = customerId;

    SELECT
        GROUP_CONCAT(DISTINCT
            CONCAT(
                'MAX(IF(attributes.code_name = ''',
                code_name,
                ''', attributes.value, NULL)) AS ',
                CONCAT(code_name, '_', entity_type)
            )
        ) INTO @sql
    FROM (
        (
            SELECT e.code_name, e.sort_order, 'varchar' AS entity_type, v.value
            FROM customer_attribute_varchar AS v
            JOIN attributes AS e ON e.id = v.attribute_id
            WHERE (v.customer_id = @customerId)
        ) UNION (
            SELECT e.code_name, e.sort_order, 'datetime' AS entity_type, v.value
            FROM customer_attribute_datetime AS v
            JOIN attributes AS e ON e.id = v.attribute_id
            WHERE (v.customer_id = @customerId)
        ) UNION (
            SELECT e.code_name, e.sort_order, 'decimal' AS entity_type, v.value
            FROM customer_attribute_decimal AS v
            JOIN attributes AS e ON e.id = v.attribute_id
            WHERE (v.customer_id = @customerId)
        ) UNION (
            SELECT e.code_name, e.sort_order, 'int' AS entity_type, v.value
            FROM customer_attribute_int AS v
            JOIN attributes AS e ON e.id = v.attribute_id
            WHERE (v.customer_id = @customerId)
        ) UNION (
            SELECT e.code_name, e.sort_order, 'boolean' AS entity_type, v.value
            FROM customer_attribute_boolean AS v
            JOIN attributes AS e ON e.id = v.attribute_id
            WHERE (v.customer_id = @customerId)
        ) UNION (
            SELECT e.code_name, e.sort_order, 'text' AS entity_type, v.value
            FROM customer_attribute_text AS v
            JOIN attributes AS e ON e.id = v.attribute_id
            WHERE (v.customer_id = @customerId)
        ) ORDER BY sort_order ASC
    ) AS attributes;

    SET @sql = CONCAT("SELECT '", @customerId, "' AS customer_id, ", "customers.deleted_at AS deleted_at, ", @sql, " FROM  (\r\n        (\r\n            SELECT e.code_name, e.sort_order, 'varchar' AS entity_type, v.value\r\n            FROM customer_attribute_varchar AS v\r\n            JOIN attributes AS e ON e.id = v.attribute_id\r\n            WHERE (v.customer_id = @customerId)\r\n        ) UNION (\r\n            SELECT e.code_name, e.sort_order, 'datetime' AS entity_type, v.value\r\n            FROM customer_attribute_datetime AS v\r\n            JOIN attributes AS e ON e.id = v.attribute_id\r\n            WHERE (v.customer_id = @customerId)\r\n        ) UNION (\r\n            SELECT e.code_name, e.sort_order, 'decimal' AS entity_type, v.value\r\n            FROM customer_attribute_decimal AS v\r\n            JOIN attributes AS e ON e.id = v.attribute_id\r\n            WHERE (v.customer_id = @customerId)\r\n        ) UNION (\r\n            SELECT e.code_name, e.sort_order, 'int' AS entity_type, v.value\r\n            FROM customer_attribute_int AS v\r\n            JOIN attributes AS e ON e.id = v.attribute_id\r\n            WHERE (v.customer_id = @customerId)\r\n        ) UNION (\r\n            SELECT e.code_name, e.sort_order, 'boolean' AS entity_type, v.value\r\n            FROM customer_attribute_boolean AS v\r\n            JOIN attributes AS e ON e.id = v.attribute_id\r\n            WHERE (v.customer_id = @customerId)\r\n        ) UNION (\r\n            SELECT e.code_name, e.sort_order, 'text' AS entity_type, v.value\r\n            FROM customer_attribute_text AS v\r\n            JOIN attributes AS e ON e.id = v.attribute_id\r\n            WHERE (v.customer_id = @customerId)\r\n        ) ORDER BY sort_order ASC\r\n    ) AS attributes, (SELECT deleted_at AS deleted_at FROM customers WHERE (id = @customerId)) AS customers GROUP BY customer_id");

IF(!ISNULL(@sql) && LENGTH(trim(@sql)) > 0) THEN
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;     
END IF;
END$$

DROP PROCEDURE IF EXISTS `UpdateCustomerSyncedAt`$$
CREATE DEFINER=`root`@`%` PROCEDURE `UpdateCustomerSyncedAt` (IN `customerId` CHAR(27))   BEGIN
    DECLARE EXIT HANDLER FOR SQLEXCEPTION 
        BEGIN
            ROLLBACK;
            RESIGNAL;
        END;

    START TRANSACTION;
        SET @customerId = customerId;
        SET @dateNow = CURRENT_TIMESTAMP();

        UPDATE `customers`
        SET `synced_at` = @dateNow
        WHERE `customers`.`id` = @customerId;
    COMMIT;

    SELECT DATE_FORMAT(CONVERT_TZ(@dateNow, @@session.time_zone, '+00:00'),'%Y-%m-%dT%TZ') AS `synced_at`;
END$$

--
-- Functions
--
DROP FUNCTION IF EXISTS `StringRandomize`$$
CREATE DEFINER=`root`@`%` FUNCTION `StringRandomize` (`size` SMALLINT(6)) RETURNS VARCHAR(255) CHARSET ascii COLLATE ascii_general_ci DETERMINISTIC BEGIN
    SET @tokenString = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    SET @result = "";
    SET @c = 0;

    WHILE @c < size DO
        SET @result = CONCAT(@result, SUBSTRING(@tokenString, ROUND(RAND() * LENGTH(@tokenString)), 1));
        SET @c = @c + 1;
    END WHILE;

    RETURN @result;
END$$

DELIMITER ;

-- --------------------------------------------------------

--
-- Table structure for table `attributes`
--

DROP TABLE IF EXISTS `attributes`;
CREATE TABLE IF NOT EXISTS `attributes` (
  `id` varchar(27) NOT NULL,
  `entity_type` enum('blob','boolean','datetime','decimal','int','text','varchar') NOT NULL,
  `code_name` varchar(100) NOT NULL,
  `metadata` text DEFAULT NULL,
  `label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `default_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `source_data` varchar(255) DEFAULT NULL,
  `input_renderer` varchar(255) DEFAULT NULL,
  `list_renderer` varchar(255) DEFAULT NULL,
  `is_filterable` tinyint(1) NOT NULL DEFAULT 0,
  `is_visible_on_front` tinyint(1) NOT NULL DEFAULT 1,
  `is_visible_in_list` tinyint(1) NOT NULL DEFAULT 0,
  `is_configurable` tinyint(1) NOT NULL DEFAULT 0,
  `is_user_defined` tinyint(1) NOT NULL DEFAULT 1,
  `is_read_only` tinyint(1) NOT NULL DEFAULT 0,
  `is_required` tinyint(1) NOT NULL DEFAULT 0,
  `is_unique` tinyint(1) NOT NULL DEFAULT 0,
  `list_column_size` smallint(6) NOT NULL DEFAULT 3,
  `display_format` varchar(12) DEFAULT NULL,
  `sort_order` int(11) NOT NULL,
  `note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp(),
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code_name` (`code_name`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `attributes`
--

TRUNCATE TABLE `attributes`;
--
-- Dumping data for table `attributes`
--

INSERT IGNORE INTO `attributes` (`id`, `entity_type`, `code_name`, `metadata`, `label`, `description`, `default_value`, `source_data`, `input_renderer`, `list_renderer`, `is_filterable`, `is_visible_on_front`, `is_visible_in_list`, `is_configurable`, `is_user_defined`, `is_read_only`, `is_required`, `is_unique`, `list_column_size`, `display_format`, `sort_order`, `note`, `created_at`, `updated_at`, `deleted_at`) VALUES
('2Ldv6o4KKKKN8T4fFGUab2F7c7Y', 'int', 'country', '{}', 'Country', NULL, '229', 'enumeration/ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 'Dropdown.Selection', 'Flag', 1, 1, 0, 1, 1, 0, 0, 0, 3, NULL, 14, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6o5iOHw5Qx7RoXMcZEQAm3p', 'datetime', 'birthday', '{\"format\":\"MMM Do YY\",\"showTodayButton\":true}', 'Birthday', NULL, NULL, NULL, 'Picker.Datetime', 'Datetime', 1, 1, 0, 0, 0, 0, 0, 0, 3, NULL, 9, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6p2ZI9efINxMEEwoxDUiFWq', 'datetime', 'appointment', '{\"format\":\"MMM Do YY\",\"showTodayButton\":true,\"hasRange\":true,\"type\":\"date\"}', 'Appointment', NULL, NULL, NULL, 'Picker.Datetime', 'Datetime', 0, 1, 1, 1, 1, 0, 0, 0, 2, NULL, 11, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6pdily8utFiGS32OSWG2M1N', 'text', 'address1', NULL, 'Address 1', NULL, NULL, NULL, 'Input.Text', 'Text', 1, 1, 0, 0, 0, 0, 0, 0, 3, NULL, 10, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6pfkqn4vzFlD9OPxQWGBKDq', 'blob', 'avatar', '{}', 'Avatar', NULL, NULL, NULL, 'Picker.Uploader', 'File', 0, 1, 1, 1, 1, 0, 0, 0, 1, NULL, 8, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6rp8yxIy0uAy7ABIfoNhqFd', 'varchar', 'firstName', NULL, 'First Name', NULL, NULL, NULL, 'Input.Text', NULL, 1, 1, 1, 0, 0, 0, 1, 0, 2, NULL, 1, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6rqTruc1nUgTh9fTEhGfJVU', 'varchar', 'phone', '{\"uriScheme\":\"tel\"}', 'Phone', NULL, NULL, NULL, 'Input.Text', 'Deeplink', 1, 1, 0, 0, 0, 0, 0, 0, 3, NULL, 5, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6s933tvH8OrDICzGARY6s7f', 'int', 'gender', NULL, 'Gender', NULL, '1', 'enumeration/b60a1289-51a9-552f-a946-fdd8622d0f9c', 'Dropdown.Selection', 'Icon', 1, 1, 1, 0, 0, 0, 0, 0, 1, NULL, 3, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6tnbVwdT3ElnU3Pe9XfVGOm', 'varchar', 'city', NULL, 'City', NULL, NULL, NULL, 'Input.Text', NULL, 1, 1, 0, 0, 0, 0, 0, 0, 3, NULL, 4, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6ucScAJ4HhbHDALcxwwlpYk', 'varchar', 'password', NULL, 'Password', NULL, NULL, NULL, NULL, NULL, 0, 0, 0, 0, 0, 1, 0, 0, 3, NULL, 0, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6uI49DDiRXlyRo0i54wzzLo', 'varchar', 'lastName', NULL, 'Last Name', NULL, NULL, NULL, 'Input.Text', NULL, 1, 1, 1, 0, 0, 0, 0, 0, 2, NULL, 2, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL),
('2Ldv6uQB6rntTZmcvAaQSRrlVnM', 'varchar', 'website', '{}', 'Website', NULL, NULL, NULL, 'Input.Text', 'Deeplink', 0, 1, 1, 1, 1, 0, 0, 0, 1, NULL, 6, NULL, '2022-09-02 17:13:18', '2022-09-02 17:13:18', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `customers`
--

DROP TABLE IF EXISTS `customers`;
CREATE TABLE IF NOT EXISTS `customers` (
  `id` char(27) NOT NULL,
  `email_address` varchar(512) NOT NULL,
  `initial_points` double DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `synced_at` datetime DEFAULT NULL,
  `current_points` double DEFAULT 0,
  `earned_points` double DEFAULT 0,
  `spent_points` double DEFAULT 0,
  PRIMARY KEY (`id`),
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customers`
--

TRUNCATE TABLE `customers`;
--
-- Dumping data for table `customers`
--

INSERT IGNORE INTO `customers` (`id`, `email_address`, `initial_points`, `created_at`, `updated_at`, `deleted_at`, `synced_at`, `current_points`, `earned_points`, `spent_points`) VALUES
('2ETVohmxehXDJZS8d4iwKJGZKuO', 'admin@localhost.com', '', 0, '2022-09-08 06:48:57', '2022-09-08 06:48:57', NULL, '2023-01-28 16:06:30', 0, 0, 0);

-- --------------------------------------------------------

--
-- Table structure for table `customer_attribute_blob`
--

DROP TABLE IF EXISTS `customer_attribute_blob`;
CREATE TABLE IF NOT EXISTS `customer_attribute_blob` (
  `customer_id` varchar(27) NOT NULL,
  `attribute_id` varchar(27) NOT NULL,
  `value` mediumblob DEFAULT NULL,
  `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`customer_id`,`attribute_id`),
  KEY `attribute_id` (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_attribute_blob`
--

TRUNCATE TABLE `customer_attribute_blob`;
-- --------------------------------------------------------

--
-- Table structure for table `customer_attribute_boolean`
--

DROP TABLE IF EXISTS `customer_attribute_boolean`;
CREATE TABLE IF NOT EXISTS `customer_attribute_boolean` (
  `customer_id` varchar(27) NOT NULL,
  `attribute_id` varchar(27) NOT NULL,
  `value` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`customer_id`,`attribute_id`),
  KEY `attribute_id` (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_attribute_boolean`
--

TRUNCATE TABLE `customer_attribute_boolean`;
-- --------------------------------------------------------

--
-- Table structure for table `customer_attribute_datetime`
--

DROP TABLE IF EXISTS `customer_attribute_datetime`;
CREATE TABLE IF NOT EXISTS `customer_attribute_datetime` (
  `customer_id` varchar(27) NOT NULL,
  `attribute_id` varchar(27) NOT NULL,
  `value` datetime DEFAULT NULL,
  `value2` datetime DEFAULT NULL,
  PRIMARY KEY (`customer_id`,`attribute_id`),
  KEY `attribute_id` (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_attribute_datetime`
--

TRUNCATE TABLE `customer_attribute_datetime`;
-- --------------------------------------------------------

--
-- Table structure for table `customer_attribute_decimal`
--

DROP TABLE IF EXISTS `customer_attribute_decimal`;
CREATE TABLE IF NOT EXISTS `customer_attribute_decimal` (
  `customer_id` varchar(27) NOT NULL,
  `attribute_id` varchar(27) NOT NULL,
  `value` decimal(12,4) DEFAULT NULL,
  PRIMARY KEY (`customer_id`,`attribute_id`),
  KEY `attribute_id` (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_attribute_decimal`
--

TRUNCATE TABLE `customer_attribute_decimal`;
-- --------------------------------------------------------

--
-- Table structure for table `customer_attribute_int`
--

DROP TABLE IF EXISTS `customer_attribute_int`;
CREATE TABLE IF NOT EXISTS `customer_attribute_int` (
  `customer_id` varchar(27) NOT NULL,
  `attribute_id` varchar(27) NOT NULL,
  `value` int(11) DEFAULT NULL,
  PRIMARY KEY (`customer_id`,`attribute_id`),
  KEY `attribute_id` (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_attribute_int`
--

TRUNCATE TABLE `customer_attribute_int`;
--
-- Dumping data for table `customer_attribute_int`
--

INSERT IGNORE INTO `customer_attribute_int` (`customer_id`, `attribute_id`, `value`) VALUES
('2ETVohmxehXDJZS8d4iwKJGZKuO', '2Ldv6o4KKKKN8T4fFGUab2F7c7Y', 0),
('2ETVohmxehXDJZS8d4iwKJGZKuO', '2Ldv6s933tvH8OrDICzGARY6s7f', 0);

-- --------------------------------------------------------

--
-- Table structure for table `customer_attribute_text`
--

DROP TABLE IF EXISTS `customer_attribute_text`;
CREATE TABLE IF NOT EXISTS `customer_attribute_text` (
  `customer_id` varchar(27) NOT NULL,
  `attribute_id` varchar(27) NOT NULL,
  `value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`customer_id`,`attribute_id`),
  KEY `attribute_id` (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_attribute_text`
--

TRUNCATE TABLE `customer_attribute_text`;
-- --------------------------------------------------------

--
-- Table structure for table `customer_attribute_varchar`
--

DROP TABLE IF EXISTS `customer_attribute_varchar`;
CREATE TABLE IF NOT EXISTS `customer_attribute_varchar` (
  `customer_id` varchar(27) NOT NULL,
  `attribute_id` varchar(27) NOT NULL,
  `value` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`customer_id`,`attribute_id`),
  KEY `attribute_id` (`attribute_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_attribute_varchar`
--

TRUNCATE TABLE `customer_attribute_varchar`;
--
-- Dumping data for table `customer_attribute_varchar`
--

INSERT IGNORE INTO `customer_attribute_varchar` (`customer_id`, `attribute_id`, `value`) VALUES
('2ETVohmxehXDJZS8d4iwKJGZKuO', '2Ldv6rp8yxIy0uAy7ABIfoNhqFd', 'sfdgs'),
('2ETVohmxehXDJZS8d4iwKJGZKuO', '2Ldv6uI49DDiRXlyRo0i54wzzLo', 'dfgdfg');

-- --------------------------------------------------------

--
-- Table structure for table `customer_events`
--

DROP TABLE IF EXISTS `customer_events`;
CREATE TABLE IF NOT EXISTS `customer_events` (
  `event_name` varchar(60) NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `aggregate_method` enum('mean','median','count','sum','first','last','min','max','stddev') NOT NULL DEFAULT 'count',
  `synchronization_interval` int(11) NOT NULL DEFAULT 5000,
  PRIMARY KEY (`event_name`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_events`
--

TRUNCATE TABLE `customer_events`;
-- --------------------------------------------------------

--
-- Table structure for table `customer_event_series`
--

DROP TABLE IF EXISTS `customer_event_series`;
CREATE TABLE IF NOT EXISTS `customer_event_series` (
  `event_name` varchar(60) NOT NULL,
  `customer_id` char(27) NOT NULL,
  `snapshot_date` date NOT NULL,
  `event_occurrence` int(11) NOT NULL,
  `event_value` decimal(12,4) NOT NULL,
  `aggregate_method` varchar(256) NOT NULL DEFAULT 'count',
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`event_name`,`customer_id`,`snapshot_date`),
  KEY `customer_id` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `customer_event_series`
--

TRUNCATE TABLE `customer_event_series`;
-- --------------------------------------------------------

--
-- Table structure for table `enumerations`
--

DROP TABLE IF EXISTS `enumerations`;
CREATE TABLE IF NOT EXISTS `enumerations` (
  `id` varchar(36) NOT NULL,
  `label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `default_value` int(11) DEFAULT NULL,
  `list_renderer` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `enumerations`
--

TRUNCATE TABLE `enumerations`;
--
-- Dumping data for table `enumerations`
--

INSERT IGNORE INTO `enumerations` (`id`, `label`, `description`, `default_value`, `list_renderer`) VALUES
('b60a1289-51a9-552f-a946-fdd8622d0f9c', 'Gender', NULL, 1, 'Icon'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 'Country Flag', NULL, 229, 'Flag'),
('e6145e3e-e7aa-50f6-bd86-7c32189ed0cf', 'Yes / no', NULL, 2, 'Icon');

-- --------------------------------------------------------

--
-- Table structure for table `enumeration_values`
--

DROP TABLE IF EXISTS `enumeration_values`;
CREATE TABLE IF NOT EXISTS `enumeration_values` (
  `enum_id` varchar(36) NOT NULL,
  `value` int(11) NOT NULL,
  `label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`enum_id`,`value`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `enumeration_values`
--

TRUNCATE TABLE `enumeration_values`;
--
-- Dumping data for table `enumeration_values`
--

INSERT IGNORE INTO `enumeration_values` (`enum_id`, `value`, `label`, `data`) VALUES
('b60a1289-51a9-552f-a946-fdd8622d0f9c', 1, 'Male', 'mars'),
('b60a1289-51a9-552f-a946-fdd8622d0f9c', 2, 'Female', 'venus'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 1, 'Afghanistan', 'af'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 2, 'Aland Islands', 'ax'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 3, 'Albania', 'al'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 4, 'Algeria', 'dz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 5, 'American Samoa', 'as'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 6, 'Andorra', 'ad'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 7, 'Angola', 'ao'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 8, 'Anguilla', 'ai'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 9, 'Antigua', 'ag'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 10, 'Argentina', 'ar'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 11, 'Armenia', 'am'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 12, 'Aruba', 'aw'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 13, 'Australia', 'au'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 14, 'Austria', 'at'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 15, 'Azerbaijan', 'az'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 16, 'Bahamas', 'bs'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 17, 'Bahrain', 'bh'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 18, 'Bangladesh', 'bd'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 19, 'Barbados', 'bb'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 20, 'Belarus', 'by'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 21, 'Belgium', 'be'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 22, 'Belize', 'bz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 23, 'Benin', 'bj'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 24, 'Bermuda', 'bm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 25, 'Bhutan', 'bt'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 26, 'Bolivia', 'bo'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 27, 'Bosnia', 'ba'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 28, 'Botswana', 'bw'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 29, 'Bouvet Island', 'bv'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 30, 'Brazil', 'br'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 31, 'British Virgin Islands', 'vg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 32, 'Brunei', 'bn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 33, 'Bulgaria', 'bg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 34, 'Burkina Faso', 'bf'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 35, 'Burma', 'mm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 36, 'Burundi', 'bi'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 37, 'Caicos Islands', 'tc'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 38, 'Cambodia', 'kh'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 39, 'Cameroon', 'cm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 40, 'Canada', 'ca'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 41, 'Cape Verde', 'cv'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 42, 'Cayman Islands', 'ky'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 43, 'Central African Republic', 'cf'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 44, 'Chad', 'td'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 45, 'Chile', 'cl'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 46, 'China', 'cn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 47, 'Christmas Island', 'cx'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 48, 'Cocos Islands', 'cc'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 49, 'Colombia', 'co'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 50, 'Comoros', 'km'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 51, 'Congo Brazzaville', 'cg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 52, 'Congo', 'cd'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 53, 'Cook Islands', 'ck'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 54, 'Costa Rica', 'cr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 55, 'Cote Divoire', 'ci'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 56, 'Croatia', 'hr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 57, 'Cuba', 'cu'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 58, 'Cyprus', 'cy'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 59, 'Czech Republic', 'cz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 60, 'Denmark', 'dk'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 61, 'Djibouti', 'dj'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 62, 'Dominica', 'dm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 63, 'Dominican Republic', 'do'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 64, 'Ecuador', 'ec'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 65, 'Egypt', 'eg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 66, 'El Salvador', 'sv'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 67, 'England', 'gb'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 68, 'Equatorial Guinea', 'gq'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 69, 'Eritrea', 'er'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 70, 'Estonia', 'ee'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 71, 'Ethiopia', 'et'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 72, 'European Union', 'eu'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 73, 'Falkland Islands', 'fk'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 74, 'Faroe Islands', 'fo'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 75, 'Fiji', 'fj'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 76, 'Finland', 'fi'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 77, 'France', 'fr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 78, 'French Guiana', 'gf'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 79, 'French Polynesia', 'pf'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 80, 'French Territories', 'tf'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 81, 'Gabon', 'ga'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 82, 'Gambia', 'gm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 83, 'Georgia', 'ge'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 84, 'Germany', 'de'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 85, 'Ghana', 'gh'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 86, 'Gibraltar', 'gi'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 87, 'Greece', 'gr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 88, 'Greenland', 'gl'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 89, 'Grenada', 'gd'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 90, 'Guadeloupe', 'gp'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 91, 'Guam', 'gu'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 92, 'Guatemala', 'gt'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 93, 'Guinea-Bissau', 'gw'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 94, 'Guinea', 'gn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 95, 'Guyana', 'gy'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 96, 'Haiti', 'ht'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 97, 'Heard Island', 'hm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 98, 'Honduras', 'hn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 99, 'Hong Kong', 'hk'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 100, 'Hungary', 'hu'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 101, 'Iceland', 'is'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 102, 'India', 'in'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 103, 'Indian Ocean Territory', 'io'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 104, 'Indonesia', 'id'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 105, 'Iran', 'ir'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 106, 'Iraq', 'iq'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 107, 'Ireland', 'ie'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 108, 'Israel', 'il'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 109, 'Italy', 'it'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 110, 'Jamaica', 'jm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 111, 'Japan', 'jp'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 112, 'Jordan', 'jo'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 113, 'Kazakhstan', 'kz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 114, 'Kenya', 'ke'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 115, 'Kiribati', 'ki'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 116, 'Kuwait', 'kw'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 117, 'Kyrgyzstan', 'kg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 118, 'Laos', 'la'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 119, 'Latvia', 'lv'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 120, 'Lebanon', 'lb'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 121, 'Lesotho', 'ls'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 122, 'Liberia', 'lr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 123, 'Libya', 'ly'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 124, 'Liechtenstein', 'li'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 125, 'Lithuania', 'lt'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 126, 'Luxembourg', 'lu'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 127, 'Macau', 'mo'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 128, 'Macedonia', 'mk'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 129, 'Madagascar', 'mg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 130, 'Malawi', 'mw'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 131, 'Malaysia', 'my'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 132, 'Maldives', 'mv'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 133, 'Mali', 'ml'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 134, 'Malta', 'mt'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 135, 'Marshall Islands', 'mh'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 136, 'Martinique', 'mq'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 137, 'Mauritania', 'mr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 138, 'Mauritius', 'mu'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 139, 'Mayotte', 'yt'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 140, 'Mexico', 'mx'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 141, 'Micronesia', 'fm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 142, 'Moldova', 'md'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 143, 'Monaco', 'mc'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 144, 'Mongolia', 'mn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 145, 'Montenegro', 'me'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 146, 'Montserrat', 'ms'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 147, 'Morocco', 'ma'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 148, 'Mozambique', 'mz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 149, 'Namibia', 'na'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 150, 'Nauru', 'nr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 151, 'Nepal', 'np'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 152, 'Netherlands Antilles', 'an'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 153, 'Netherlands', 'nl'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 154, 'New Caledonia', 'nc'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 155, 'New Guinea', 'pg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 156, 'New Zealand', 'nz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 157, 'Nicaragua', 'ni'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 158, 'Niger', 'ne'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 159, 'Nigeria', 'ng'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 160, 'Niue', 'nu'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 161, 'Norfolk Island', 'nf'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 162, 'North Korea', 'kp'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 163, 'Northern Mariana Islands', 'mp'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 164, 'Norway', 'no'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 165, 'Oman', 'om'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 166, 'Pakistan', 'pk'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 167, 'Palau', 'pw'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 168, 'Palestine', 'ps'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 169, 'Panama', 'pa'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 170, 'Paraguay', 'py'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 171, 'Peru', 'pe'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 172, 'Philippines', 'ph'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 173, 'Pitcairn Islands', 'pn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 174, 'Poland', 'pl'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 175, 'Portugal', 'pt'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 176, 'Puerto Rico', 'pr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 177, 'Qatar', 'qa'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 178, 'Reunion', 're'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 179, 'Romania', 'ro'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 180, 'Russia', 'ru'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 181, 'Rwanda', 'rw'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 182, 'Saint Helena', 'sh'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 183, 'Saint Kitts and Nevis', 'kn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 184, 'Saint Lucia', 'lc'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 185, 'Saint Pierre', 'pm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 186, 'Saint Vincent', 'vc'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 187, 'Samoa', 'ws'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 188, 'San Marino', 'sm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 189, 'Sandwich Islands', 'gs'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 190, 'Sao Tome', 'st'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 191, 'Saudi Arabia', 'sa'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 192, 'Senegal', 'sn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 193, 'Serbia', 'cs'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 194, 'Serbia', 'rs'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 195, 'Seychelles', 'sc'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 196, 'Sierra Leone', 'sl'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 197, 'Singapore', 'sg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 198, 'Slovakia', 'sk'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 199, 'Slovenia', 'si'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 200, 'Solomon Islands', 'sb'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 201, 'Somalia', 'so'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 202, 'South Africa', 'za'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 203, 'South Korea', 'kr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 204, 'Spain', 'es'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 205, 'Sri Lanka', 'lk'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 206, 'Sudan', 'sd'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 207, 'Suriname', 'sr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 208, 'Svalbard', 'sj'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 209, 'Swaziland', 'sz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 210, 'Sweden', 'se'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 211, 'Switzerland', 'ch'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 212, 'Syria', 'sy'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 213, 'Taiwan', 'tw'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 214, 'Tajikistan', 'tj'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 215, 'Tanzania', 'tz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 216, 'Thailand', 'th'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 217, 'Timorleste', 'tl'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 218, 'Togo', 'tg'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 219, 'Tokelau', 'tk'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 220, 'Tonga', 'to'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 221, 'Trinidad', 'tt'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 222, 'Tunisia', 'tn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 223, 'Turkey', 'tr'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 224, 'Turkmenistan', 'tm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 225, 'Tuvalu', 'tv'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 226, 'Uganda', 'ug'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 227, 'Ukraine', 'ua'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 228, 'United Arab Emirates', 'ae'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 229, 'United States', 'us'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 230, 'Uruguay', 'uy'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 231, 'Us Minor Islands', 'um'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 232, 'Us Virgin Islands', 'vi'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 233, 'Uzbekistan', 'uz'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 234, 'Vanuatu', 'vu'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 235, 'Vatican City', 'va'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 236, 'Venezuela', 've'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 237, 'Vietnam', 'vn'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 238, 'Wallis and Futuna', 'wf'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 239, 'Western Sahara', 'eh'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 240, 'Yemen', 'ye'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 241, 'Zambia', 'zm'),
('ca8f216f-cd1c-5ed4-8fc3-1a3fc66f4c98', 242, 'Zimbabwe', 'zw'),
('e6145e3e-e7aa-50f6-bd86-7c32189ed0cf', 1, 'Yes', 'green toggle on'),
('e6145e3e-e7aa-50f6-bd86-7c32189ed0cf', 2, 'No', 'grey toggle off');

-- --------------------------------------------------------

--
-- Table structure for table `groups`
--

DROP TABLE IF EXISTS `groups`;
CREATE TABLE IF NOT EXISTS `groups` (
  `id` char(27) NOT NULL,
  `code_name` varchar(100) NOT NULL,
  `category_name` varchar(100) NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `sort_order` smallint(6) NOT NULL,
  `parent_id` char(27) DEFAULT NULL,
  `node_left` smallint(6) NOT NULL,
  `node_right` smallint(6) NOT NULL,
  `node_level` smallint(6) NOT NULL,
  `node_depth` smallint(6) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code_name` (`code_name`),
  KEY `category_name` (`category_name`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `groups`
--

TRUNCATE TABLE `groups`;
-- --------------------------------------------------------

--
-- Table structure for table `group_categories`
--

DROP TABLE IF EXISTS `group_categories`;
CREATE TABLE IF NOT EXISTS `group_categories` (
  `category_name` varchar(100) NOT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`category_name`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `group_categories`
--

TRUNCATE TABLE `group_categories`;
-- --------------------------------------------------------

--
-- Table structure for table `group_customer_linker`
--

DROP TABLE IF EXISTS `group_customer_linker`;
CREATE TABLE IF NOT EXISTS `group_customer_linker` (
  `group_id` char(27) NOT NULL,
  `customer_id` char(27) NOT NULL,
  PRIMARY KEY (`group_id`,`customer_id`),
  KEY `customer_id` (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=ascii COLLATE=ascii_general_ci;

--
-- Truncate table before insert `group_customer_linker`
--

TRUNCATE TABLE `group_customer_linker`;
-- --------------------------------------------------------

--
-- Constraints for table `customer_attribute_blob`
--
ALTER TABLE `customer_attribute_blob`
  ADD CONSTRAINT `customer_attribute_blob_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `customer_attribute_blob_ibfk_2` FOREIGN KEY (`attribute_id`) REFERENCES `attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `customer_attribute_boolean`
--
ALTER TABLE `customer_attribute_boolean`
  ADD CONSTRAINT `customer_attribute_boolean_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `customer_attribute_boolean_ibfk_2` FOREIGN KEY (`attribute_id`) REFERENCES `attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `customer_attribute_datetime`
--
ALTER TABLE `customer_attribute_datetime`
  ADD CONSTRAINT `customer_attribute_datetime_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `customer_attribute_datetime_ibfk_2` FOREIGN KEY (`attribute_id`) REFERENCES `attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `customer_attribute_decimal`
--
ALTER TABLE `customer_attribute_decimal`
  ADD CONSTRAINT `customer_attribute_decimal_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `customer_attribute_decimal_ibfk_2` FOREIGN KEY (`attribute_id`) REFERENCES `attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `customer_attribute_int`
--
ALTER TABLE `customer_attribute_int`
  ADD CONSTRAINT `customer_attribute_int_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `customer_attribute_int_ibfk_2` FOREIGN KEY (`attribute_id`) REFERENCES `attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `customer_attribute_text`
--
ALTER TABLE `customer_attribute_text`
  ADD CONSTRAINT `customer_attribute_text_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `customer_attribute_text_ibfk_2` FOREIGN KEY (`attribute_id`) REFERENCES `attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `customer_attribute_varchar`
--
ALTER TABLE `customer_attribute_varchar`
  ADD CONSTRAINT `customer_attribute_varchar_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `customer_attribute_varchar_ibfk_2` FOREIGN KEY (`attribute_id`) REFERENCES `attributes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `customer_event_series`
--
ALTER TABLE `customer_event_series`
  ADD CONSTRAINT `customer_event_series_ibfk_1` FOREIGN KEY (`event_name`) REFERENCES `customer_events` (`event_name`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `customer_event_series_ibfk_2` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `enumeration_values`
--
ALTER TABLE `enumeration_values`
  ADD CONSTRAINT `enumeration_values_ibfk_1` FOREIGN KEY (`enum_id`) REFERENCES `enumerations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `groups`
--
ALTER TABLE `groups`
  ADD CONSTRAINT `groups_ibfk_1` FOREIGN KEY (`category_name`) REFERENCES `group_categories` (`category_name`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `group_customer_linker`
--
ALTER TABLE `group_customer_linker`
  ADD CONSTRAINT `group_customer_linker_ibfk_1` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `group_customer_linker_ibfk_2` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
