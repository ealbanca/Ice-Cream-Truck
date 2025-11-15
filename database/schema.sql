-- Create the database
SET @OLD_UNIQUE_CHECKS = @@UNIQUE_CHECKS,
    UNIQUE_CHECKS = 0;
SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS,
    FOREIGN_KEY_CHECKS = 0;
SET @OLD_SQL_MODE = @@SQL_MODE,
    SQL_MODE = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `iceCreamTruck` DEFAULT CHARACTER SET utf8;
USE `iceCreamTruck`;
-- -----------------------------------------------------
-- Table `iceCreamTruck`.`Topping`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `iceCreamTruck`.`Topping` (
    `topping_id` INT NOT NULL AUTO_INCREMENT,
    `topping_name` VARCHAR(45) NOT NULL,
    `topping_description` VARCHAR(255) NULL,
    `topping_price` DECIMAL(10, 2) NOT NULL,
    `topping_quantity` INT NOT NULL,
    PRIMARY KEY (`topping_id`),
    UNIQUE INDEX `Ingredients_id_UNIQUE` (`topping_id` ASC) VISIBLE
) ENGINE = InnoDB;
-- -----------------------------------------------------
-- Table `iceCreamTruck`.`Flavor`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `iceCreamTruck`.`Flavor` (
    `flavor_id` INT NOT NULL,
    `flavor_name` VARCHAR(45) NOT NULL,
    `flavor_quantity` INT NOT NULL,
    `flavor_price` DECIMAL(10, 2) NULL,
    PRIMARY KEY (`flavor_id`)
) ENGINE = InnoDB;
-- -----------------------------------------------------
-- Table `iceCreamTruck`.`Cup`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `iceCreamTruck`.`Cup` (
    `cup_id` INT NOT NULL,
    `cup_name` VARCHAR(45) NOT NULL,
    `cup_price` DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (`cup_id`, `cup_name`)
) ENGINE = InnoDB;
-- -----------------------------------------------------
-- Table `iceCreamTruck`.`IceCream`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `iceCreamTruck`.`IceCream` (
    `iceCream_id` INT NOT NULL AUTO_INCREMENT,
    `IceCream_name` VARCHAR(45) NULL,
    `topping_id` INT NOT NULL,
    `flavor_id` INT NOT NULL,
    `cup_id` INT NOT NULL,
    `iceCream_price` DECIMAL(10, 2) NULL,
    PRIMARY KEY (`iceCream_id`),
    INDEX `fk_IceCream_Ingredients_idx` (`topping_id` ASC) VISIBLE,
    INDEX `fk_IceCream_Flavor1_idx` (`flavor_id` ASC) VISIBLE,
    INDEX `fk_IceCream_Cup1_idx` (`cup_id` ASC) VISIBLE,
    CONSTRAINT `fk_IceCream_Ingredients` FOREIGN KEY (`topping_id`) REFERENCES `iceCreamTruck`.`Topping` (`topping_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT `fk_IceCream_Flavor1` FOREIGN KEY (`flavor_id`) REFERENCES `iceCreamTruck`.`Flavor` (`flavor_id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT `fk_IceCream_Cup1` FOREIGN KEY (`cup_id`) REFERENCES `iceCreamTruck`.`Cup` (`cup_id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE = InnoDB;
SET SQL_MODE = @OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS = @OLD_UNIQUE_CHECKS;
-- Insert initial data into Topping table
INSERT INTO Topping (
        topping_name,
        topping_description,
        topping_price,
        topping_quantity
    )
VALUES (
        'Sprinkles',
        'Colorful sugar sprinkles',
        0.50,
        200
    ),
    (
        'Chocolate Chips',
        'Mini chocolate chip pieces',
        0.75,
        150
    ),
    (
        'Caramel Drizzle',
        'Sweet caramel sauce',
        0.60,
        100
    ),
    (
        'Whipped Cream',
        'Light and fluffy whipped cream',
        0.40,
        120
    );
-- Insert initial data into Flavor table
INSERT INTO Flavor (
        flavor_id,
        flavor_name,
        flavor_quantity,
        flavor_price
    )
VALUES (1, 'Vanilla', 100, 1.50),
    (2, 'Chocolate', 90, 1.60),
    (3, 'Strawberry', 80, 1.55),
    (4, 'Mint', 70, 1.65);
-- Insert initial data into Cup table
INSERT INTO Cup (cup_id, cup_name, cup_price)
VALUES (1, 'Small Cup', 0.50),
    (2, 'Medium Cup', 0.75),
    (3, 'Large Cup', 1.00),
    (4, 'Waffle Cone', 1.25);
-- Insert initial data into IceCream table
INSERT INTO IceCream (
        IceCream_name,
        topping_id,
        flavor_id,
        cup_id,
        iceCream_price
    )
VALUES ('Vanilla Delight', 1, 1, 1, 3.50),
    ('Chocolate Crunch', 2, 2, 2, 4.25),
    ('Strawberry Swirl', 3, 3, 3, 4.75),
    ('Mint Supreme', 4, 4, 4, 5.00);