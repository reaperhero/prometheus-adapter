# api

- 获取所有规则

```
path: /adapter/rules
method: GET

response:
{
    "count": 2,
    "rules": {
        "http://10.10.12.48:30090": [
            {
                "id": 149,
                "status": true,
                "sql": "kube_pod_container_status_restarts_total{}"
            }
        ],
        "http://10.10.12.49:30090": [
            {
                "id": 150,
                "status": true,
                "sql": "kube_pod_container_status_restarts_total{}"
            }
        ]
    }
}
```

- 获取指定实例规则
```
path: /adapter/rules?instance=http://10.10.12.48:30090
method: GET



response:
{
    "count": 1,
    "rules": {
        "http://10.10.12.48:30090": [
            {
                "id": 149,
                "status": true,
                "sql": "kube_pod_container_status_restarts_total{}"
            }
        ]
    }
}

```

- 创建规则

```
path: /adapter/rule
method: POST

request:
{
    "cap_name": "自定义指标名称",
    "cap_sql": "kube_pod_container_status_restarts_total{}",
    "instance":"http://10.10.12.48:30090"
}

response:
{
    "msg": "ok"
    "data": ""
} 
```


- 删除规则

```
path: /adapter/rule/:id
method: DELETE

response:
{
    "code": ""
    "msg": ""
} 
```


- 更新规则

```
path: /adapter/rule/:id
method: PUT

request:
{
    "cap_name": ""
    "instance": "",
    "sql": ""
}

response:
{
    "code": ""
    "msg": ""
} 
```

- 更新状态

```
path: /adapter/rulestatus/:id?status=true
method: PUT

response:
{
    "code": ""
    "msg": ""
} 
```