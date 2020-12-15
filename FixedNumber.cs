using System.Transactions;
using UnityEngine;
using UnityEditor;
using System.IO;
using System.Collections.Generic;
using System;
using System.Collections;

public class Fixed64 :  IComparable<Fixed64>, IEquatable<Fixed64>
{

	[MenuItem("Test/FixedNumber")]
	public static void StartRun() {
		var filse = new string[] { "TestData/AddData.txt", "TestData/SubData.txt", "TestData/MulData.txt", "TestData/DivData.txt" };
		List<Fixed64> t = new List<Fixed64>();
		foreach (var file in filse)
        {
			var op = new List<string>();
			var v0 = new List<ulong>();
			var v1 = new List<ulong>();
			var v2 = new List<ulong>();
			Debug.Log("Start Test");
			if (ReadData(file,ref op,ref v0,ref v1,ref v2))
            {
				for (int i = 0; i < op.Count; i++)
                {
					var f0 = new Fixed64(v0[i]);
					var f1 = new Fixed64(v1[i]);
					var f2 = new Fixed64(v2[i]);
					switch (op[i])
                    {
						case "Add":
                            if (f0 + f1 != f2)
                            {
								Debug.Log("Shoudle Equal:" + (f0 + f1).ToString() + " " + f2.ToString());
                            }
							break;
						case "Sub":
							if (f0 - f1 != f2)
							{
								Debug.Log("Shoudle Equal:" + (f0 - f1).ToString() + " " + f2.ToString());
							}
							break;
						case "Mul":
							if (f0 * f1 != f2)
							{
								Debug.Log("Shoudle Equal:" + (f0 * f1).ToString() + " " + f2.ToString());
							}
							break;

						case "Div":
							if (f0 / f1 != f2)
							{
								Debug.Log("Shoudle Equal:" + (f0 / f1).ToString() + " " + f2.ToString());
							}
							break;

					}
				}
            }
			Debug.Log("End Test");
		}
	}

	public Fixed64() { }

	public Fixed64(ulong value,bool transform = false)
    {
		if (transform)
        {
			v = value << precisionBitsNum;
        } else
        {
			v = value;
        }
    }

	public Fixed64(double value)
	{
		ulong valueBits = 0;
		unsafe
		{
			ulong* pl = (ulong*)&value;
			valueBits = *pl;
		}

		var s = valueBits & mask64S;
		var e = valueBits & mask64E;
		var m = valueBits & mask64M | ((ulong)1 << size64M);
		var realE = ((Int64)e >> size64M) - bias64;
		ulong fixedD = 0;
		ulong fixedP = 0;

		if (realE >= 0) {
			var pBitsNum = size64M - (int)realE;
			if (pBitsNum < 0)
			{
				if ((long)precision - pBitsNum > 63 - 53) {
					throw new System.InvalidOperationException("Fixed number: part digital overflow");
				}

				fixedD = m << -pBitsNum << precisionBitsNum;
			}
			else
			{
				fixedD = m >> pBitsNum << precisionBitsNum;
				var pBitsFlowNum = pBitsNum - precisionBitsNum;
				var pMask = (((ulong)1 << pBitsNum) - 1);


				if (pBitsFlowNum <= 0) {
					fixedP = (m & pMask) << (-1 * pBitsFlowNum);
				}
				else
				{
					fixedP = roundOdd(m & pMask, (pBitsFlowNum)) >> pBitsFlowNum;
				}
			}
		}
		else
		{
			var allMBitsNum = 52 - realE;
			var pFlowBitsNum = allMBitsNum - precisionBitsNum;
			if (pFlowBitsNum > 53) {
				v = 0;
				return;
			}
			if (pFlowBitsNum <= 0) {
				fixedP = m << ((int)pFlowBitsNum * -1);
			}
			else
			{
				fixedP = roundOdd(m, (int)pFlowBitsNum) >> (int)pFlowBitsNum;
			}
		}

		v |= s;
		v |= fixedD;
		v |= fixedP;
	}

	public static bool ReadData(string path, ref List<string> op,ref List<ulong> v1, ref List<ulong> v2, ref List<ulong> v3)
    {
        using (StringReader sr = new StringReader(File.ReadAllText(path)))
        {
            while (true)
            {
				var line = sr.ReadLine();
				if (line == null)
                {
					return true; 
                }
				var opAndData = line.Split(':');
				if (opAndData.Length != 2)
				{
					return false;
				}
				var nums = opAndData[1].Split(' ');
				if (nums.Length != 4)
				{
					return false;
				}
				var num1 = System.Convert.ToUInt64(nums[0]);
				var num2 = System.Convert.ToUInt64(nums[1]);
				var num3 = System.Convert.ToUInt64(nums[2]);

				op.Add(opAndData[0]);
				v1.Add(num1);
				v2.Add(num2);
				v3.Add(num3);
			}
        }
    }


	public static Fixed64 operator +(Fixed64 f1, Fixed64 f2)
	{
		var ret = new Fixed64();
		var fS = f1.v & mask64S;
		var oS = f2.v & mask64S;
		f1.v = bitClear(f1.v, mask64S);
		f2.v = bitClear(f2.v, mask64S);
		if (fS == oS) {
			ret.v = (f1.v + f2.v) | fS;
			return ret;
		}
		else
		{
			if (f1.v > f2.v)
			{
				ret.v = (f1.v - f2.v) | fS;
				return ret;
			} else if (f1.v < f2.v) {
				ret.v = (f2.v - f1.v) | oS;
				return ret;
			} else
			{
				return ret;
			}
		}
	}

	public static Fixed64 operator -(Fixed64 f1, Fixed64 f2)
	{
		f2.v = f2.v ^ mask64S;
		return f1 + f2;
	}

	public static Fixed64 operator *(Fixed64 f1, Fixed64 f2)
	{
		var fS = f1.v & mask64S;
		var oS = f2.v & mask64S;
		f1.v = bitClear(f1.v, mask64S);
		f2.v = bitClear(f2.v, mask64S);
		var r = Mul64(f1.v, f2.v);
		var hi = r.Item1;
		var lo = r.Item2;
		lo = roundOdd(lo, precisionBitsNum) >> precisionBitsNum;
		hi = (hi & decimalBitsMask) << (64 - precisionBitsNum);
		var ret = new Fixed64();
		ret.v = hi | lo;
		if (ret.v == 0)
		{
			return ret;
		}
		ret.v |= fS ^ oS;
		return ret;
	}

	public static Fixed64 operator /(Fixed64 f1, Fixed64 f2)
	{
		var fS = f1.v & mask64S;
		var oS = f2.v & mask64S;
		f1.v = bitClear(f1.v, mask64S);
		f2.v = bitClear(f2.v, mask64S);
		var quo = Div64(f1.v >> 64 - precisionBitsNum, f1.v << precisionBitsNum, f2.v).Item1;
		var ret = new Fixed64();
		ret.v = quo;
		return ret;
	}

	public static  bool operator >(Fixed64 f1, Fixed64 f2)
    {
		return (f1.v ^ mask64S) > (f2.v ^ mask64S);
    }
	public static bool operator <(Fixed64 f1, Fixed64 f2)
	{
		return (f1.v ^ mask64S) < (f2.v ^ mask64S);
	}
	public static bool operator ==(Fixed64 f1, Fixed64 f2)
	{
		return f1.v == f2.v;
	}
	public static bool operator !=(Fixed64 f1, Fixed64 f2)
	{
		return f1.v != f2.v;
	}

	public Fixed64 Abs()
	{
		var ret = new Fixed64();
		ret.v = bitClear(v, mask64S);
		return ret;
	}

	public Int64 Round()
	{
		var d = v & decimalBitsMask;
		var roundUp = (Int64)0;
		if (precisionBitsNum > 0 && (d >= ((ulong)1 << (precisionBitsNum - 1)))) {
			roundUp = 1;
		}
		var ret = ((Int64)(bitClear(v, mask64S) >> precisionBitsNum) + roundUp);
		return (v & mask64S) > 0 ? -1 * ret : ret;
	}

	public static Tuple<ulong, ulong> Div64(ulong hi, ulong lo, ulong y)
	{
		if (y == 0)
		{
			throw new System.InvalidOperationException("Fixed number: Div Zero");
		}
		if (y <= hi)
		{
			throw new System.InvalidOperationException("Fixed number: Div Result overflow");
		}

		var s = 64 - Len64(y);
		y <<= s;

		var yn1 = y >> 32;
		var yn0 = y & mask32;
		var un32 = hi << s | lo >> (64 - s);
		var un10 = lo << s;
		var un1 = un10 >> 32;
		var un0 = un10 & mask32;
		var q1 = un32 / yn1;
		var rhat = un32 - q1 * yn1;

		while (q1 >= two32 || q1 * yn0 > two32 * rhat + un1)
		{
			q1--;
			rhat += yn1;
			if (rhat >= two32)
			{
				break;
			}
		}

		var un21 = un32 * two32 + un1 - q1 * y;
		var q0 = un21 / yn1;
		rhat = un21 - q0 * yn1;

		while (q0 >= two32 || q0 * yn0 > two32 * rhat + un0)
		{
			q0--;
			rhat += yn1;
			if (rhat >= two32)
			{
				break;
			}
		}

		return new Tuple<ulong, ulong>(q1 * two32 + q0, (un21 * two32 + un0 - q0 * y) >> s);
	}

	public static Tuple<ulong, ulong> Mul64(ulong x, ulong y)
	{
		var x0 = x & mask32;
		var x1 = x >> 32;
		var y0 = y & mask32;
		var y1 = y >> 32;
		var w0 = x0 * y0;
		var t = x1 * y0 + (w0 >> 32);
		var w1 = t & mask32;
		var w2 = t >> 32;
		w1 += x0 * y1;
		var hi = x1 * y1 + w2 + (w1 >> 32);
		var lo = x * y;
		return new Tuple<ulong, ulong>(hi, lo);
	}

	public Int64 ToInt64()
    {
		var ret = (Int64)bitClear(v,mask64S) >> precisionBitsNum;
		return (v & mask64S) > 0 ? -ret :ret ;
    }

	public double ToDouble()
	{
		var number = bitClear(v, mask64S);
		var idx = Len64(number);
		if (idx != 0)
		{
			var e = idx - precisionBitsNum - 1;
			number = (((ulong)1 << idx) - 1) & number;
			if (idx > size64M)
			{
				number >>= idx - size64M - 1;
			}
			else
			{
				number <<= size64M - idx + 1;
			}
			number = bitClear(number, mask64E);
			number |= v & mask64S;
			number |= (ulong)(e + bias64) << size64M;
		}

		unsafe
		{
			double* dp = (double*)&number;
			return *dp;
		}
	}

	public override string ToString()
	{
		return ToBase10N(5);
	}

    public string ToBase10N(int n)
    {
        if (n > 18)
        {
            throw new System.InvalidOperationException("Fixed number:BaseN greater 18");
        }

        var floatStr = "";
        if ((v & mask64S) > 0)
        {
            floatStr += "-";
        }

        if (n == 0)
        {
            return floatStr + insertToFloatStrBase10((ulong)(Abs().Round()), 0).ToString();
        }

        n = n + 1;// To round to the end
        var number = bitClear(v, mask64S);
        var d = number >> precisionBitsNum;
        var p = number & decimalBitsMask;

        var mulResult = Mul64(p, (ulong)(Math.Pow(10, (double)n)));
        var hi = mulResult.Item1;
        var lo = mulResult.Item2;
        var divResult = Div64(hi, lo, (ulong)1 << precisionBitsNum);
        var quo = divResult.Item1;

        var partD = insertToFloatStrBase10(d, 0);
        var partP = insertToFloatStrBase10(quo, n);
        roundEndBase10(ref partP).RemoveAt(partP.Count - 1);
        var upToTop = carrayUpBase10(ref partP);

        if (upToTop)
        {
            partD[partD.Count - 1]++;
            upToTop = carrayUpBase10(ref partD);
        }

        if (upToTop)
        {
            floatStr += '1';
        }

        floatStr += new string(partD.ToArray());
        floatStr += '.';
        floatStr += new string(partP.ToArray());
        return floatStr;
    }

    private static List<char> roundEndBase10(ref List<char> p)
    {
		if (p[p.Count - 1] > '4') {
			p[p.Count - 1] = '0';
			p[p.Count - 2]++;
		}
		return p;
    }

    private static bool carrayUpBase10(ref List<char> p)
    {
        var i = p.Count - 1;

        for (; i >= 1; i--)
        {
            if (p[i] > '9')
            {
                p[i] = '0';
                p[i - 1]++;
            }
            else
            {
                return false;
            }
        }

        var upToTop = p[0] == ('9' + 1);
        if (upToTop)
        {
            p[0] = '0';
        }

        return upToTop;
    }

    public static List<char> insertToFloatStrBase10(ulong v, int n )
    {
		List<char> ret = new List<char>();
		while(v > 0) {
			ret.Add((char)((v % 10) + '0'));
			v /= 10;
		}
		while (n > ret.Count) {
			ret.Add('0');
		}
		if (ret.Count > 0) {
			ret.Reverse();
		}
		else
		{
			ret.Add('0');
  		}

		return ret;
	}

	private static ulong roundOdd(ulong v, int precisionBitsNum)
	{
		var precisionBitsMask = ((ulong)1 << precisionBitsNum) - 1;
		var flow = v & precisionBitsMask;
		var cond = ((ulong)1 << (precisionBitsNum - 1));
		ulong endBit = 0;
		if (flow > cond || flow == cond && (v & ((ulong)1 << precisionBitsNum)) > 0) {
			endBit = 1;
		}

		return bitClear(v,precisionBitsMask) + (endBit << precisionBitsNum);
	}

	private static ulong bitClear(ulong v,ulong mask)
    {
		return v & (0xFFFFFFFFFFFFFFFF ^ mask);//&^
	}

	private static string ulong2BitsStr(ulong v)
    {
		string retStr = "";
		for (int i = 0; i < 64; i++)
        {
			if ((v & mask64S) > 0)
            {
				retStr += "1";
            } else
            {
				retStr += "0";
            }
			v = v << 1;
		}
		return retStr;
    }

	public static int Len64(ulong x)
    {
		int n = 0;
		if (x >= (ulong)1 << 32) {
			x >>= 32;
			n = 32;
		}
		if (x >= 1 << 16)
		{
			x >>= 16;
			n += 16;
		}
		if (x >= 1 << 8) {
			x >>= 8;
			n += 8;
		}
        return n + (int)len8tab[x];
    }

	public bool Equals(Fixed64 other)
	{
		return this == other;
	}

	public override int GetHashCode()
    {
		return v.GetHashCode();
    }

    public int CompareTo(Fixed64 obj)
    {

		if (this > obj)
        {
			return 1;
        } 
		if (this < obj)
        {
			return -1;
        }
		return 0;
    }

	public override bool Equals(object obj)
	{
		if (obj.GetType() != GetType())
		{
			return false;
		}
		return this == (Fixed64)obj;
	}

	private static readonly byte[] len8tab = new byte[]{
		0x00,0x01, 0x02, 0x02, 0x03, 0x03, 0x03, 0x03, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04,
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05,
		0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
		0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
		0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
		0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
		0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
		0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	};

    private static readonly ulong mask64S = (ulong)1 << 63;
    private static readonly int size64M = 52;
	private static readonly ulong mask64E = (((ulong)1 << 11) - 1) << size64M;
	private static readonly ulong mask64M = ((ulong)1 << size64M) - 1;
    private static readonly int bias64 = (1 << (11 - 1)) - 1;
	private static readonly ulong mask32 = ((ulong)1 << 32) - 1;
	private static readonly ulong two32 = (ulong)1 << 32;

	private static ulong precision = 20;
    private static ulong decimalBitsMask = ((ulong)1 << 20) - 1;
	private static int precisionBitsNum = (int)precision;

	private ulong v = 0;
}
