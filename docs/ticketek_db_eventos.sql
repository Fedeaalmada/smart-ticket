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
-- Table structure for table `eventos`
--

DROP TABLE IF EXISTS `eventos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `eventos` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `titulo` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `descripcion` text COLLATE utf8mb4_unicode_ci,
  `foto_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `fecha` date NOT NULL,
  `horario` time NOT NULL,
  `duracion_minutos` int unsigned NOT NULL DEFAULT '120',
  `categoria` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Ej: Concierto, Teatro, Deporte',
  `capacidad_total` int unsigned NOT NULL,
  `estado` enum('activo','cancelado','agotado') COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'activo',
  `creado_por` int unsigned NOT NULL COMMENT 'FK al admin que lo creó',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_eventos_fecha` (`fecha`),
  KEY `idx_eventos_estado` (`estado`),
  KEY `idx_eventos_categoria` (`categoria`),
  KEY `fk_eventos_creado_por` (`creado_por`),
  CONSTRAINT `fk_eventos_creado_por` FOREIGN KEY (`creado_por`) REFERENCES `usuarios` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `eventos`
--

LOCK TABLES `eventos` WRITE;
/*!40000 ALTER TABLE `eventos` DISABLE KEYS */;
INSERT INTO `eventos` VALUES (1,'Tan Biónica - La Última Noche Mágica Tour','La banda más convocante del rock nacional regresa con su gira de despedida. Una noche única e irrepetible.','https://example.com/tanbiónica.jpg','2026-08-15','21:00:00',150,'Concierto',5000,'activo',1,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(2,'River vs Boca - Superclásico','El partido más esperado del año en el Monumental.','https://example.com/superclasico.jpg','2026-09-07','17:00:00',120,'Deporte',3000,'activo',1,'2026-05-29 18:44:26','2026-05-29 18:44:26'),(3,'Comic Con Argentina 2026','La convención de cultura pop más grande del país. Celebridades, cosplay, gaming y mucho más.','https://example.com/comiccon.jpg','2026-10-10','10:00:00',480,'Feria',2000,'activo',1,'2026-05-29 18:44:26','2026-05-29 18:44:26');
/*!40000 ALTER TABLE `eventos` ENABLE KEYS */;
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
