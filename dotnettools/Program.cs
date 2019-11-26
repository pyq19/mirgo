using System;

namespace dotnettools
{
    class Program
    {
        static void Main(string[] args)
        {
            Manager.ConnectionString = "server=127.0.0.1;uid=root;pwd=root;database=mir";
            Manager manager = new Manager();

            string path = "/opt/gopath/src/github.com/yenkeia/mir-go/dotnettools/database/Server.MirDB";
            manager.loadFromFile(path);

            manager.saveDataToDatabase();

            manager.loadFromDatabase();

            manager.test();
        }
    }
}