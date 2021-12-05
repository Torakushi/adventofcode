import sys
import numpy as np


def firstPart():
    n = np.loadtxt("../data.txt")
    # First Method
    print(f"First part: With numpy {np.count_nonzero(n[:-1]-n[1:] < 0)}")
    # Second Method
    print(f"First part: With List comprehension {len([i for i in range(1, len(n)) if n[i] > n[i-1]])}")


def secondPart():
    n = np.loadtxt("../data.txt")
    # First Method
    sum = [n[i]+n[i+1]+n[i+2] for i in range(len(n)-2)]

    print(f"Second part: With numpy {np.count_nonzero(np.delete(sum,-1)-sum[1:] < 0)}")


def main():
    firstPart()
    secondPart()
    return 0


if __name__ == '__main__':
    sys.exit(main())
