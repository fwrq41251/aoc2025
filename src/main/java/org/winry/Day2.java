package org.winry;

import java.util.ArrayList;
import java.util.List;

public class Day2 {

    public static void main(String[] args) {
        var input = """
                119-210,907313-1048019,7272640820-7272795557,6315717352-6315818282,42-65,2234869-2439411,1474-2883,33023-53147,1-15,6151-14081,3068-5955,65808-128089,518490556-518593948,3535333552-3535383752,7340190-7548414,547-889,34283147-34389716,44361695-44519217,607492-669180,7071078-7183353,67-115,969-1469,3636264326-3636424525,762682710-762831570,827113-906870,205757-331544,290-523,86343460-86510016,5536957-5589517,132876-197570,676083-793651,23-41,17920-31734,440069-593347
                """;
        var lines = input.strip().split(",");
        var repeatNumberSum = 0L;
        for (var line : lines) {
            var parts = line.split("-");
            var start = Long.parseLong(parts[0]);
            var end = Long.parseLong(parts[1]);
            for (long i = start; i <= end; i++) {
                if (isRepeatNumber(i)) {
                    repeatNumberSum += i;
                }
            }
        }
        System.out.println("Sum of repeat numbers: " + repeatNumberSum);
    }

    private static List<Integer> getDivisors(int n) {
        var divisors = new ArrayList<Integer>();
        for (int i = 1; i <= Math.sqrt(n); i++) {
            if (n % i == 0) {
                if (i != n) { // 模式长度不能等于字符串总长度
                    divisors.add(i);
                }
                int otherDivisor = n / i;
                if (otherDivisor != n && otherDivisor != i) {
                    divisors.add(otherDivisor);
                }
            }
        }
        return divisors;
    }

    private static boolean isRepeatNumber(long number) {
        var strNum = Long.toString(number);
        var length = strNum.length();

        if (length == 1) {
            return true;
        }

        var divisors = getDivisors(length);
        for (var divisor : divisors) {
            var numRepeats = length / divisor;
            var pattern = strNum.substring(0, divisor);
            var repeatedString = pattern.repeat(numRepeats);

            if (repeatedString.equals(strNum)) {
                return true;
            }
        }
        return false;
    }
}
