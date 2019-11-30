using System;
using SqlSugar;

namespace dotnettools
{
    class Program
    {
        static void Main(string[] args)
        {
            // string connectionString = "server=127.0.0.1;uid=root;pwd=root;database=mir";
            string projectDirPath = System.Environment.CurrentDirectory;
            string sqlitePath = "/src/github.com/yenkeia/mir-go/dotnettools/mir.sqlite";
            string connectionString = "Data Source=" + projectDirPath + sqlitePath;
            SqlSugarClient DB = new SqlSugarClient(new ConnectionConfig()
            {
                ConnectionString = connectionString,
                // DbType = DbType.MySql,//设置数据库类型
                DbType = DbType.Sqlite,//设置数据库类型
                IsAutoCloseConnection = true,//自动释放数据务，如果存在事务，在事务结束后释放
                InitKeyType = InitKeyType.Attribute //从实体特性中读取主键自增列信息
            });
            Manager.DB = DB;

            Manager manager = new Manager();

            string path = projectDirPath + "/src/github.com/yenkeia/mir-go/dotnettools/database/Server.MirDB";
            manager.loadFromFile(path);

            manager.saveDataToDatabase();

            manager.loadFromDatabase();
        }
    }
}