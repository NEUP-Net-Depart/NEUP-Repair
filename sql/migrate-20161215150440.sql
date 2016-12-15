ALTER TABLE `orders` ADD COLUMN `area` VARCHAR(255);
ALTER TABLE `orders` ADD COLUMN `contact` LONGTEXT;

UPDATE `orders` SET `area` = "无" where `area` is NULL;
UPDATE `orders` SET `contact` = "无" where `contact` is NULL;
