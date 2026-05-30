-- MySQL dump 10.13  Distrib 8.0.43, for Win64 (x86_64)
--
-- Host: localhost    Database: ticketek_db
-- ------------------------------------------------------
-- Server version	9.4.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `sectores`
--

DROP TABLE IF EXISTS `sectores`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sectores` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `evento_id` int unsigned NOT NULL,
  `nombre` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Ej: Campo, Platea Baja, VIP',
  `capacidad_maxima` int unsigned NOT NULL,
  `capacidad_disponible` int unsigned NOT NULL,
  `precio` decimal(10,2) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_sectores_evento` (`evento_id`),
  CONSTRAINT `fk_sectores_evento` FOREIGN KEY (`evento_id`) REFERENCES `eventos` (`id`) ON DELETE CASCADE,
  CONSTRAINT `chk_sectores_capacidad` CHECK ((`capacidad_disponible` <= `capacidad_maxima`)),
  CONSTRAINT `chk_sectores_precio` CHECK ((`precio` >= 0))
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sectores`
--

LOCK TABLES `sectores` WRITE;
/*!40000 ALTER TABLE `sectores` DISABLE KEYS */;
INSERT INTO `sectores` VALUES (1,1,'Campo',2000,1999,15000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(2,1,'Platea Baja',1500,1499,25000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(3,1,'Platea Alta',1000,1000,18000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(4,1,'VIP',500,500,45000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(5,2,'Popular',1500,1499,8000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(6,2,'Platea',1200,1200,20000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(7,2,'Platea Preferencial',300,300,35000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(8,3,'General',1500,1500,5000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(9,3,'VIP Pass',500,500,12000.00,'2026-05-29 18:44:26','2026-05-29 18:44:26');
/*!40000 ALTER TABLE `sectores` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-05-30 16:37:59
