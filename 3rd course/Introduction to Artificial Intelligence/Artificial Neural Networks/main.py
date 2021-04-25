import matplotlib.pyplot as plt
import numpy as np
#from sklearn.linear_model import LinearRegression


def read_file():
    dataMatrix = [[], []]
    f = open("sunspot.txt", "r")
    lines = f.readlines()

    for line in lines:
        splitLine = line.split('\t')
        dataMatrix[0].append(int(splitLine[0]))
        dataMatrix[1].append(int(splitLine[1]))
    f.close()
    return dataMatrix

data = read_file()
print(data)