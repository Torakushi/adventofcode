import sys

import numpy as np
import pandas as pd


def firstPart():
    df = pd.read_csv("../data.txt", sep=" ", names=["dir", "pow"], dtype="str")
    df["pow"] = df["pow"].astype(int)
    depth, direction = 0, 0
    for _, row in df.iterrows():
        if row["dir"] == "forward":
            direction += row["pow"]
        elif row["dir"] == "down":
            depth += row["pow"]
        elif row["dir"] == "up":
            depth -= row["pow"]
        else:
            raise Exception(f"wrong instruction {row[0]}")
    print(f"Depth {depth}, direction {direction}, multiplication {depth*direction}")


def secondPart():
    df = pd.read_csv("../data.txt", sep=" ", names=["dir", "pow"], dtype="str")
    df["pow"] = df["pow"].astype(int)
    depth, direction, aim = 0, 0, 0
    for _, row in df.iterrows():
        if row["dir"] == "forward":
            direction += row["pow"]
            depth += row["pow"] * aim
        elif row["dir"] == "down":
            aim += row["pow"]
        elif row["dir"] == "up":
            aim -= row["pow"]
        else:
            raise Exception(f"wrong instruction {row[0]}")
    print(f"Depth {depth}, direction {direction}, multiplication {depth*direction}")


def main():
    firstPart()
    secondPart()
    return 0


if __name__ == "__main__":
    sys.exit(main())
