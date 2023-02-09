# mapper
**mapper层即DAO层**, 使用gorm直接对数据库进行操作


## zxl, 2023.2.9
- 把实体定义从 mapper 包中移到 model 包中
- 将数据库地址从硬编码的方式改为从配置文件中读取
- connector.go中添加OSS服务初始化代码