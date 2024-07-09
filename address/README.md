# 实体地址
**业务介绍:**
地址在电商、店铺等业务中经常出现，主要有收货地址、退货地址、注册地址、店铺地址等，其主要信息为省市区、完整地址、联系人、联系电话等，同时有分类如：默认地址、公司、家等，有的类型只能有一条，有的可以有多条
**技术方案:**
业务简单，复用率高，增加租户概念
**功能介绍:**
1. 基本的 增、改、删、查 接口


|名称|标题|必填|类型|格式|可空|默认值|案例|描述|
|:--|:--|:--|:--|:--|:--|:--|:--|:--|
|Fbusiness_id|租户ID|true|string||false||||
|label|标签|false|string||false||||
|is_default|默认|false|string||false||||
|contact_phone|联系手机号|false|string||false||||
|Fcontact_name|联系人|false|string||false||||
|address|详细地址|false|string||false||||
|provice|省|false|string|string|false||||
|provice_id|省ID|false|string|string|false||||
|city|城市|false|string|string|false||||
|city_id|城市ID|false|string|string|false||||
|area|区|false|string|string|false||||
|area_id|区ID|false|string|string|false||||