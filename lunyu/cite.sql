-- MySQL dump 10.13  Distrib 5.7.21, for macos10.13 (x86_64)
--
-- Host: localhost    Database: world
-- ------------------------------------------------------
-- Server version	5.7.21

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `cite`
--

DROP TABLE IF EXISTS `cite`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cite` (
  `ID` bigint(20) NOT NULL AUTO_INCREMENT,
  `SEQ` varchar(50) DEFAULT NULL COMMENT '编号',
  `CTENT` varchar(500) DEFAULT NULL COMMENT '内容',
  `TAGS` varchar(200) DEFAULT NULL COMMENT '标签',
  `REFS` varchar(200) DEFAULT NULL COMMENT '参考',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cite`
--

LOCK TABLES `cite` WRITE;
/*!40000 ALTER TABLE `cite` DISABLE KEYS */;
INSERT INTO `cite` VALUES (1,'1.1','有朋自远方来\r\n不亦说乎 不亦说乎','学习 朋友','3.2 5.3 13.26'),(2,'1.2','学而时习之','学习','1.3 9.9'),(3,'1.3','人不知而不愠','君子','9.9 1.2 13.26'),(6,'2.12','君子不器','君子','13.25 10.18'),(12,'9.9','父母唯其疾之忧','父母','1.2 1.4 1.3'),(13,'13.25','君子易事而难说也。说之不以道，不说也；及其使人也，器之。小人难事而易说也。说之虽不以道，说也；及其使人也，求备也。','君子','2.14 18.10 10.18'),(14,'10.18','君子不施其亲，不使大臣怨乎不以。故旧无大故，则不弃也。无求备于一人。','君子','13.25 2.12'),(15,'12.1','克己复礼为仁。一日克己复礼，天下归仁焉。为仁由己，而由人乎哉？\r\n非礼勿视，非礼勿听，非礼无言，非礼勿动','仁 君子',''),(16,'12.4','君子不忧不惧。内省不疚，夫何忧何惧？','君子','9.29'),(17,'12.3','仁者，其言也讱。为之难，言之得无讱乎！','君子','13.27'),(18,'12.2','出门如见大宾，使民如承大祭。己所不欲，勿施于人。在邦无怨，在家无怨。','仁 君子','15.24'),(19,'13.27','刚、毅、木、讷，近仁','仁 君子','12.3'),(20,'13.26','君子泰而不骄，小人骄而不泰','君子','1.1 1.3');
/*!40000 ALTER TABLE `cite` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-03-27 12:10:17
