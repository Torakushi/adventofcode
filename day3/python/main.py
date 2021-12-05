import sys

import numpy as np
import pandas as pd


def firstPart():
    file = open("../data.txt", "r")
    data = []
    for line in file.readlines():
        if line[-1] == "\n":
            line = line[:-1]
        s = []
        s.extend(line)
        data.append(s)

    df = pd.DataFrame(data)
    i, gam, eps = 0, 0, 0
    for col in df:
        gam += int(df[col].value_counts(ascending=False).index.values[0]) * 2 ** (
            11 - i
        )
        eps += int(df[col].value_counts(ascending=True).index.values[0]) * 2 ** (11 - i)
        i += 1

    print(gam * eps)


# 5852595
def secondPart():
    file = open("../data.txt", "r")
    data = []
    for line in file.readlines():
        if line[-1] == "\n":
            line = line[:-1]
        s = []
        s.extend(line)
        data.append(s)

    df = pd.DataFrame(data).astype(int)
    dfBis = df
    for col in df:
        v = 1
        val_count = dfBis[col].value_counts(ascending=False)
        if val_count.drop_duplicates().shape[0] == 2:
            v = int(val_count.index.values[0])

        dfBis = dfBis[dfBis[col] == v]
        if dfBis.shape[0] == 1:
            break

    lineOx = dfBis.head()

    dfBis = df
    for col in df:
        v = 0
        val_count = dfBis[col].value_counts(ascending=True)
        if val_count.drop_duplicates().shape[0] == 2:
            v = int(val_count.index.values[0])

        dfBis = dfBis[dfBis[col] == v]

        if dfBis.shape[0] == 1:
            break

    lineCO2 = dfBis.head()

    Ox, Co2 = 0, 0
    for i in range(12):
        Ox += int(lineOx[i]) * 2 ** (11 - i)
        Co2 += int(lineCO2[i]) * 2 ** (11 - i)

    print(Ox * Co2)


def main():
    firstPart()
    secondPart()
    return 0


if __name__ == "__main__":
    sys.exit(main())
