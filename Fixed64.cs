using System;
using System.Collections;

namespace FixedNumber
{
    class Fixed64
    {
        private static readonly ulong maskS = 0x8000000000000000;
        private static readonly ulong maskE = 0x7FF0000000000000;
        private static readonly int sizeM = 52;
        private static readonly ulong maskM = 0x000FFFFFFFFFFFFF;

        private static ulong precision = 0;


        static void Main(string[] args)
        {
            DoubleToFixed64(23123.14);

            Int64 v = 213;
            Console.WriteLine( Transfer(v));    
        }

        public  static double Transfer(Int64 v){  
            double d = 0;
            unsafe
            {
                double* pd = (double*)&v;
                d = *pd;
            } 
            return d;
        }

        public static void DoubleToFixed64(double value) {
            byte[] bits=BitConverter.GetBytes(value);
            foreach (var item in bits)
            {
                Console.WriteLine(item);
            }
            for (int i = 0; i < 64; i++){
                if ((bits[0] & (1 << i)) > 0){
                    Console.Write(0);
                } else {
                    Console.Write(1);
                }
            }
        }
    }
}
