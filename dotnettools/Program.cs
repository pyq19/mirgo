using System;

namespace dotnettools
{
    class Program
    {
        static void Main(string[] args)
        {
            Loader loader = new Loader();
            string path = "/opt/gopath/src/github.com/yenkeia/mir-go/dotnettools/database/Server.MirDB";
            loader.loadFromFile(path);
            Console.WriteLine("加载完成");
        }
    }
}