
// 演示数据 
db.checkitem.insert({
    "name":"检查DB当前登录用户信息",
    "target_ip":"server",
    "target_ip2":"",
    "options":{ },
    "cmd": {
      "category":"mssql",
      "res_type":"multi_rows",  
      "command":"select login_name,host_name, host_process_id,login_time from sys.dm_exec_sessions where host_name!=''",
      "args":[ 
          "server=192.168.230.113\SD;user id=readonly;password=readonly;" 
      ]
    },
    "checkways":[
        {
            "way":"shouldBeGreaterThan",
            "params":[0],
            "level":"warn",
            "options":[]
        } 
    ]
});

db.checkitem.insert({
    "name":"检查Server运行目录文件",
    "target_ip":"server",
    "target_ip2":"",
    "options":{ },
    "cmd": {
      "category":"cmd",
      "res_type":"string",  
      "command":"dir",
      "args":[ 
          ".",
          "/B" 
      ]
    },
  "checkways":[
      { 
        "way":"ShouldContain",
        "params":[".exe"],
        "level":"warn",
        "options":[]
      } 
    ]
});