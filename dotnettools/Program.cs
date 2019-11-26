using System;

namespace dotnettools
{
    class Program
    {
        static void Main(string[] args)
        {
            Manager manager = new Manager();

            string path = "/opt/gopath/src/github.com/yenkeia/mir-go/dotnettools/database/Server.MirDB";
            manager.loadFromFile(path);

            // manager.saveDataToDatabase();

            // manager.loadFromDatabase("SQLite", "db.sqlite");

            manager.test();
        }
    }
}