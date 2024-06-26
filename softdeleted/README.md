# 软删除
**业务介绍:**
处理业务删除逻辑
**技术方案:**
1. 通过记录删除时间实现
2. 通过某个字段(如status) 特殊值标记 

**注意:**

1. 由于方案1,2在条件筛选时,其本质是 字段使用的固定值含义不同，方案1中 空字符""表示正常记录值，方案2中固定的值(如status=0) 表示删除的记录，因此在 
```go
GetDeletedAtField() (valueType ValueType, softDeletedField SoftDeletedField) 
``` 
函数签名中，增加valueType 返回，其值有2个枚举值
```go
ValueType_Delete
ValueType_OK 
```

2. 由于Update 方法需要用于常规的含义，方便和其它单元组件集成，所以若删除记录，需要调用本包下的```Delete```方法

**功能介绍:**
1. 删业务数据
2. 更新、查询增加删除字段条件过滤



